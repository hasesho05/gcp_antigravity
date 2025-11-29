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
	userUsecase     usecase.UserUsecase
}

func NewClientHandler(qu usecase.QuestionUsecase, au usecase.AttemptUsecase, su usecase.StatsUsecase, eu usecase.ExamUsecase, uu usecase.UserUsecase) *ClientHandler {
	return &ClientHandler{
		questionUsecase: qu,
		attemptUsecase:  au,
		statsUsecase:    su,
		examUsecase:     eu,
		userUsecase:     uu,
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

func (h *ClientHandler) GetExam(w http.ResponseWriter, r *http.Request) {
	examID := chi.URLParam(r, "examID")
	exam, err := h.examUsecase.GetExam(r.Context(), examID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "試験が見つかりませんでした", http.StatusNotFound)
			return
		}
		fmt.Printf("internal server error: %+v\n", err)
		http.Error(w, "サーバー内部エラーが発生しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exam)
}

func (h *ClientHandler) GetExamSets(w http.ResponseWriter, r *http.Request) {
	examID := chi.URLParam(r, "examID")
	examSets, err := h.examUsecase.ListExamSets(r.Context(), examID)
	if err != nil {
		fmt.Printf("internal server error: %+v\n", err)
		http.Error(w, "サーバー内部エラーが発生しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(examSets)
}

func (h *ClientHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "認証されていません", http.StatusUnauthorized)
		return
	}

	// Firebase AuthからEmailを取得するために、クライアントから送られてくることを期待するか、
	// あるいはサーバー側でFirebase Admin SDKを使って取得するか。
	// ここでは簡単のため、リクエストボディからEmailを受け取るか、あるいはトークンに含まれる情報を使うのが一般的だが、
	// middlewareでトークンを検証した際にClaimsからEmailを取得してContextに入れるのがベスト。
	// 現状のmiddlewareはUIDのみなので、一旦リクエストボディからEmailを受け取る形にするか、
	// またはUserUsecase内でAdmin SDKを使って取得するように変更するか。
	// ここではシンプルにリクエストボディから受け取る形を実装する。

	var req struct {
		Email    string              `json:"email"`
		Provider domain.AuthProvider `json:"provider"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "リクエストボディが無効です", http.StatusBadRequest)
		return
	}

	// プロバイダーが指定されていない場合はデフォルトでパスワード認証とする（あるいはエラーにする）
	if req.Provider == "" {
		req.Provider = domain.ProviderPassword
	}

	user, err := h.userUsecase.CreateUser(r.Context(), userID, req.Email, req.Provider)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			// 既に存在する場合は、そのユーザー情報を取得して返す（冪等性）
			// Note: ここでは簡易的にGetUserを呼ぶか、エラーから復帰させる。
			// Usecase側で既に存在チェックをしているが、Race Condition等でここに来る可能性もある。
			// または、UsecaseがErrAlreadyExistsを返すのは「Emailが重複しているがIDが違う」場合のみにするべきか検討が必要。
			// 現状のUsecase実装では:
			// 1. IDで検索 -> あればそれを返す (nil error)
			// 2. Emailで検索 -> あれば ErrAlreadyExists
			// なので、ここに来るのは「IDは新しいがEmailが既存」の場合。
			// つまり、別のUIDで同じEmailを使おうとしている -> 本当にConflict。
			// しかし、Firebase AuthのGoogleログイン等では、同じEmailなら同じUIDになるはず（設定による）。
			// もし「User not found in backend」と言われているのにここに来るなら、
			// 「IDはDBにない」かつ「EmailはDBにある」状態。
			// これは「以前別のプロバイダで登録したEmailで、今回は別のプロバイダ（UIDが違う）でログインした」場合に起こりうる。
			// この場合、アカウントをリンクするか、エラーにするか。
			// 今回はシンプルにConflictを返すままでよいが、フロントエンドのループの原因はこれではない可能性が高い。
			// フロントエンドは 404 -> Create -> 409 -> ... となっているのか？
			// いや、ログを見ると "User not found in backend, creating new user..." が出続けているということは、
			// CreateUserが成功していない（あるいは成功してもstateが更新されていない）？
			
			// とりあえず、ログを出力してデバッグしやすくする
			fmt.Printf("CreateUser conflict: %v\n", err)
			http.Error(w, "このメールアドレスは既に登録されています", http.StatusConflict)
			return
		}
		fmt.Printf("internal server error: %+v\n", err)
		http.Error(w, "サーバー内部エラーが発生しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *ClientHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "認証されていません", http.StatusUnauthorized)
		return
	}

	user, err := h.userUsecase.GetUser(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "ユーザーが見つかりませんでした", http.StatusNotFound)
			return
		}
		fmt.Printf("internal server error: %+v\n", err)
		http.Error(w, "サーバー内部エラーが発生しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

