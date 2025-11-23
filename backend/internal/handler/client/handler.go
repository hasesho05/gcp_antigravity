package client

import (
	"encoding/json"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/go-chi/chi/v5"

	"gcp_antigravity/backend/internal/domain"
	"gcp_antigravity/backend/internal/middleware"
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
	// Go 1.22ではパス値を取得できますが、ここではChiを使用しています。
	// パターン: /exams/{examID}/sets/{examSetID}/questions
	examID := chi.URLParam(r, "examID")
	examSetID := chi.URLParam(r, "examSetID")

	if examID == "" || examSetID == "" {
		http.Error(w, "examIDとexamSetIDは必須です", http.StatusBadRequest)
		return
	}

	questions, err := h.usecase.GetExamQuestions(r.Context(), examSetID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "問題が見つかりませんでした", http.StatusNotFound)
			return
		}
		http.Error(w, "サーバー内部エラーが発生しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

func (h *ClientHandler) StartAttempt(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "認証されていません", http.StatusUnauthorized)
		return
	} 

	var req input.CreateAttemptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "リクエストボディが無効です", http.StatusBadRequest)
		return
	}

	attempt, err := h.usecase.StartAttempt(r.Context(), userID, req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidArgument) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "サーバー内部エラーが発生しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(attempt)
}

func (h *ClientHandler) UpdateAttempt(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "認証されていません", http.StatusUnauthorized)
		return
	}
	attemptID := chi.URLParam(r, "attemptID")

	var req input.UpdateAttemptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "リクエストボディが無効です", http.StatusBadRequest)
		return
	}

	if err := h.usecase.UpdateAttempt(r.Context(), userID, attemptID, req); err != nil {
		if errors.Is(err, domain.ErrFailedPrecondition) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "サーバー内部エラーが発生しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *ClientHandler) CompleteAttempt(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "認証されていません", http.StatusUnauthorized)
		return
	}
	attemptID := chi.URLParam(r, "attemptID")

	var req input.CompleteAttemptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "リクエストボディが無効です", http.StatusBadRequest)
		return
	}

	attempt, err := h.usecase.CompleteAttempt(r.Context(), userID, attemptID, req)
	if err != nil {
		if errors.Is(err, domain.ErrFailedPrecondition) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "サーバー内部エラーが発生しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attempt)
}

func (h *ClientHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "認証されていません", http.StatusUnauthorized)
		return
	}
	examID := chi.URLParam(r, "examID")

	stats, err := h.usecase.GetUserExamStats(r.Context(), userID, examID)
	if err != nil {
		http.Error(w, "サーバー内部エラーが発生しました", http.StatusInternalServerError)
		return
	}

	if stats == nil {
		// 統計情報が存在しない場合は、空のオブジェクトを返す
		stats = &domain.UserExamStats{
			ExamID:      examID,
			UserID:      userID,
			DomainStats: make(map[string]domain.DomainScore),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
