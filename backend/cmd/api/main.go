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

	"gcp_antigravity/backend/internal/handler/admin"
	client_handler "gcp_antigravity/backend/internal/handler/client"
	"gcp_antigravity/backend/internal/infra/firestore"
	"gcp_antigravity/backend/internal/repository_impl"
	"gcp_antigravity/backend/internal/usecase"
)

func main() {
	ctx := context.Background()

	// Initialize Firestore
	client := firestore.NewClient(ctx)
	defer client.Close()

	// Dependency Injection
	examRepo := repository_impl.NewExamRepository(client)
	examUsecase := usecase.NewExamUsecase(examRepo)
	adminHandler := admin.NewAdminHandler(examUsecase)
	clientHandler := client_handler.NewClientHandler(examUsecase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GCP Antigravity Backend is running!")
	})

	// Admin Routes
	mux.HandleFunc("POST /admin/exams/{examID}/sets/{setID}/questions", adminHandler.UploadQuestions)

	// Client Routes
	mux.HandleFunc("GET /exams/{examID}/sets/{setID}/questions", clientHandler.GetQuestions)
	mux.HandleFunc("POST /users/me/attempts", clientHandler.StartAttempt)

	// CORS Middleware
	handler := corsMiddleware(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
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
