package client

import (
	"encoding/json"
	"net/http"

	"github.com/cockroachdb/errors"

	"gcp_antigravity/backend/internal/domain"
	"gcp_antigravity/backend/internal/usecase"
	"gcp_antigravity/backend/internal/usecase/input"
)

type ClientHandler struct {
	usecase usecase.ExamUsecase
}

func NewClientHandler(u usecase.ExamUsecase) *ClientHandler {
	return &ClientHandler{usecase: u}
}

func (h *ClientHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	// In Go 1.22, we can get path values.
	// pattern: /exams/{examID}/sets/{setID}/questions
	setID := r.PathValue("setID")
	if setID == "" {
		http.Error(w, "setID is required", http.StatusBadRequest)
		return
	}

	questions, err := h.usecase.GetExamQuestions(r.Context(), setID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "questions not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

func (h *ClientHandler) StartAttempt(w http.ResponseWriter, r *http.Request) {
	// Mock UserID for now (Middleware should handle auth)
	// In a real scenario, we extract UID from context (set by Auth middleware)
	userID := "test_user_id" 

	var req input.CreateAttemptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	attempt, err := h.usecase.StartAttempt(r.Context(), userID, req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidArgument) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(attempt)
}
