package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/go-chi/chi/v5"

	"nearline/backend/internal/domain"
	"nearline/backend/internal/middleware"
	"nearline/backend/internal/usecase"
	"nearline/backend/internal/usecase/input"
)

type ClientHandler struct {
	questionUsecase usecase.QuestionUsecase
	attemptUsecase  usecase.AttemptUsecase
	statsUsecase    usecase.StatsUsecase
	examUsecase     usecase.ExamUsecase
}

func NewClientHandler(qu usecase.QuestionUsecase, au usecase.AttemptUsecase, su usecase.StatsUsecase, eu usecase.ExamUsecase) *ClientHandler {
	return &ClientHandler{
		questionUsecase: qu,
		attemptUsecase:  au,
		statsUsecase:    su,
		examUsecase:     eu,
	}
}

func (h *ClientHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	// Go 1.22ではパス値を取得できますが、ここではChiを使用しています。
	// パターン: /exams/{examID}/sets/{examSetID}/questions
	examID := chi.URLParam(r, "examID")
	examSetID := chi.URLParam(r, "examSetID")

	input, err := input.NewGetExamQuestions(examSetID)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidArgument) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("internal server error: %+v\n", err)
		http.Error(w, "サーバー内部エラーが発生しました", http.StatusInternalServerError)
		return
	}

	questions, err := h.questionUsecase.GetExamQuestions(r.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "問題が見つかりませんでした", http.StatusNotFound)
			return
		}
		fmt.Printf("internal server error: %+v\n", err)
		http.Error(w, "サーバー内部エラーが発生しました", http.StatusInternalServerError)
		return
	}

	// ルートのexamIDと取得した問題のexamIDが一致するか検証
	if len(questions) > 0 && questions[0].ExamID != examID {
		http.Error(w, "URLのexamIDと問題のexamIDが一致しません", http.StatusBadRequest)
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

	attempt, err := h.attemptUsecase.StartAttempt(r.Context(), userID, req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidArgument) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("internal server error: %+v\n", err)
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

	if err := h.attemptUsecase.UpdateAttempt(r.Context(), userID, attemptID, req); err != nil {
		if errors.Is(err, domain.ErrFailedPrecondition) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		fmt.Printf("internal server error: %+v\n", err)
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

	var reqBody input.CompleteAttemptRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "リクエストボディが無効です", http.StatusBadRequest)
		return
	}

	input, err := input.NewCompleteAttempt(userID, attemptID, reqBody.Answers)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthenticated) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if errors.Is(err, domain.ErrInvalidArgument) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("internal server error: %+v\n", err)
		http.Error(w, "サーバー内部エラーが発生しました", http.StatusInternalServerError)
		return
	}

	attempt, err := h.attemptUsecase.CompleteAttempt(r.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrFailedPrecondition) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		fmt.Printf("internal server error: %+v\n", err)
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

	input, err := input.NewGetUserExamStats(userID, examID)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthenticated) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if errors.Is(err, domain.ErrInvalidArgument) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("internal server error: %+v\n", err)
		http.Error(w, "サーバー内部エラーが発生しました", http.StatusInternalServerError)
		return
	}

	stats, err := h.statsUsecase.GetUserExamStats(r.Context(), input)
	if err != nil {
		fmt.Printf("internal server error: %+v\n", err)
		http.Error(w, "サーバー内部エラーが発生しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *ClientHandler) ListExams(w http.ResponseWriter, r *http.Request) {
	exams, err := h.examUsecase.ListExams(r.Context())
	if err != nil {
		fmt.Printf("internal server error: %+v\n", err)
		http.Error(w, "サーバー内部エラーが発生しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exams)
}
