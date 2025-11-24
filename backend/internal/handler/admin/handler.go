package admin

import (
	"encoding/json"
	"net/http"
	"fmt"

	"github.com/cockroachdb/errors"

	"gcp_antigravity/backend/internal/domain"
	"gcp_antigravity/backend/internal/usecase"
	"gcp_antigravity/backend/internal/usecase/input"
)

type AdminHandler struct {
	usecase usecase.QuestionUsecase
}

func NewAdminHandler(u usecase.QuestionUsecase) *AdminHandler {
	return &AdminHandler{usecase: u}
}

func (h *AdminHandler) UploadQuestions(w http.ResponseWriter, r *http.Request) {
	// Path params are handled by the router (or we parse them manually if using std lib without pattern matching in older Go, 
	// but Go 1.22+ supports path values).
	// However, for simplicity and since we accept JSON body that includes IDs, we'll use the body primarily.
	// Ideally we should validate path params match body params.
	
	var req input.UploadQuestionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "リクエストボディが無効です", http.StatusBadRequest)
		return
	}

	// Basic validation of path params vs body could go here if we extracted path values.
	// examID := r.PathValue("examID")
	// setID := r.PathValue("setID")

	if err := h.usecase.UploadQuestions(r.Context(), req); err != nil {
		// Error handling with cockroachdb/errors
		if errors.Is(err, domain.ErrInvalidArgument) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Log the full error with stack trace for internal errors
		fmt.Printf("internal server error: %+v\n", err)
		http.Error(w, "サーバー内部エラーが発生しました", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"ok"}`))
}
