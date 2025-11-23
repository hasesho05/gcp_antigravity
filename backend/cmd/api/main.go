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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"gcp_antigravity/backend/internal/handler/admin"
	client_handler "gcp_antigravity/backend/internal/handler/client"
	"gcp_antigravity/backend/internal/infra/auth"
	"gcp_antigravity/backend/internal/infra/firestore"
	internal_middleware "gcp_antigravity/backend/internal/middleware"
	"gcp_antigravity/backend/internal/repository_impl"
	"gcp_antigravity/backend/internal/usecase"
)

func main() {
	// サーバーの起動
	if err := run(); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Firestoreの初期化
	client := firestore.NewClient(ctx)
	defer client.Close()

	// Firebase Authの初期化
	authClient, err := auth.NewClient(ctx)
	panic(err)
	}

	// 依存性の注入 (Dependency Injection)
	qRepo := repository_impl.NewQuestionRepository(client)
	aRepo := repository_impl.NewAttemptRepository(client)
	sRepo := repository_impl.NewUserStatsRepository(client)
	examUsecase := usecase.NewExamUsecase(qRepo, aRepo, sRepo)
	adminHandler := admin.NewAdminHandler(examUsecase)
	clientHandler := client_handler.NewClientHandler(examUsecase)

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
		fmt.Fprintf(w, "GCP Antigravity Backend is running!")
	})

	// 認証ミドルウェア
	authMiddleware := internal_middleware.AuthMiddleware(authClient)

	// 管理者用ルート
	r.Route("/admin", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/exams/{examID}/sets/{examSetID}/questions", adminHandler.UploadQuestions)
	})

	// クライアント用ルート
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		
		// Exams
		r.Route("/exams/{examID}", func(r chi.Router) {
			r.Get("/sets/{examSetID}/questions", clientHandler.GetQuestions)
		})

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
