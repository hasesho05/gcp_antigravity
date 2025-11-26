package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"nearline/backend/internal/handler/admin"
	client_handler "nearline/backend/internal/handler/client"
	"nearline/backend/internal/infra/auth"
	"nearline/backend/internal/infra/firestore"
	internal_middleware "nearline/backend/internal/middleware"
	"nearline/backend/internal/repository_impl"
	"nearline/backend/internal/usecase"

	"github.com/joho/godotenv"
)

func main() {
	// サーバーの起動
	if err := run(); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}

func run() error {
	// .envファイルを読み込む (開発環境用)
	// 本番環境では環境変数が直接設定されるため、ファイルがなくてもエラーにしない
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Firestoreの初期化
	client := firestore.NewClient(ctx)
	defer client.Close()

	// Firebase Authの初期化
	authClient, err := auth.NewClient(ctx)
	if err != nil {
		return errors.Wrap(err, "Firebase Authクライアントの初期化に失敗しました")
	}

	// 依存性の注入 (Dependency Injection)
	qRepo := repository_impl.NewQuestionRepository(client)
	aRepo := repository_impl.NewAttemptRepository(client)
	sRepo := repository_impl.NewUserStatsRepository(client)
	txRepo := repository_impl.NewTransactionRepository(client)
	examRepo := repository_impl.NewExamRepository(client) // 追加

	questionUsecase := usecase.NewQuestionUsecase(qRepo)
	attemptUsecase := usecase.NewAttemptUsecase(qRepo, aRepo, sRepo, txRepo)
	statsUsecase := usecase.NewStatsUsecase(sRepo)
	examUsecase := usecase.NewExamUsecase(examRepo, qRepo, aRepo, sRepo, txRepo)

	adminHandler := admin.NewAdminHandler(questionUsecase)
	clientHandler := client_handler.NewClientHandler(questionUsecase, attemptUsecase, statsUsecase, examUsecase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()

	// ミドルウェア
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware) // Chi用に適応したカスタムCORSミドルウェア

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintf(w, "nearline Backend is running!")	})

	// 認証ミドルウェア
	authMiddleware := internal_middleware.AuthMiddleware(authClient)

	// 管理者用ルート
	r.Route("/admin", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/exams/{examID}/sets/{examSetID}/questions", adminHandler.UploadQuestions)
	})

	// Exams (Public & Protected mixed)
	r.Route("/exams", func(r chi.Router) {
		// Public: List Exams
		r.Get("/", clientHandler.ListExams)

		// Protected: Exam Details & Questions
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Route("/{examID}", func(r chi.Router) {
				r.Get("/sets/{examSetID}/questions", clientHandler.GetQuestions)
			})
		})
	})

	// クライアント用ルート (Authenticated)
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		// Users
		r.Route("/users/me", func(r chi.Router) {
			r.Post("/attempts", clientHandler.StartAttempt)
			r.Put("/attempts/{attemptID}", clientHandler.UpdateAttempt)
			r.Post("/attempts/{attemptID}/complete", clientHandler.CompleteAttempt)
			r.Get("/stats/{examID}", clientHandler.GetStats)
		})
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Server listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
	return nil
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // In production, replace * with specific origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

