package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"

	"gcp_antigravity/backend/internal/domain"
	"gcp_antigravity/backend/internal/repository"
	"gcp_antigravity/backend/internal/usecase/input"
)

type ExamUsecase interface {
	UploadQuestions(ctx context.Context, req input.UploadQuestionsRequest) error
	GetExamQuestions(ctx context.Context, examSetID string) ([]domain.Question, error)
	StartAttempt(ctx context.Context, userID string, req input.CreateAttemptRequest) (*domain.Attempt, error)
	UpdateAttempt(ctx context.Context, userID, attemptID string, req input.UpdateAttemptRequest) error
	CompleteAttempt(ctx context.Context, userID, attemptID string, req input.CompleteAttemptRequest) (*domain.Attempt, error)
	GetUserExamStats(ctx context.Context, userID, examID string) (*domain.UserExamStats, error)
}

type examUsecase struct {
	repo repository.ExamRepository
}

func NewExamUsecase(repo repository.ExamRepository) ExamUsecase {
	return &examUsecase{repo: repo}
}

func (u *examUsecase) UploadQuestions(ctx context.Context, req input.UploadQuestionsRequest) error {
	if len(req.Questions) == 0 {
		return errors.Wrap(domain.ErrInvalidArgument, "問題が提供されていません")
	}

	var domainQuestions []domain.Question
	now := time.Now()

	for _, qInput := range req.Questions {
		// Generate ID: {ExamCode}_{SetID}_{Index}
		// e.g. PCD_SET1_001
		id := fmt.Sprintf("%s_%s_%03d", req.ExamCode, req.ExamSetID, qInput.Index)

		var options []domain.AnswerOption
		for _, o := range qInput.Options {
			options = append(options, domain.AnswerOption{
				ID:          o.ID,
				Text:        o.Text,
				Explanation: o.Explanation,
			})
		}

		q := domain.Question{
			ID:                 id,
			ExamID:             req.ExamID,
			ExamSetID:          req.ExamSetID,
			ExamCode:           req.ExamCode,
			QuestionText:       qInput.QuestionText,
			QuestionType:       qInput.QuestionType,
			Options:            options,
			CorrectAnswers:     qInput.CorrectAnswers,
			OverallExplanation: qInput.OverallExplanation,
			Domain:             qInput.Domain,
			ImageURL:           qInput.ImageURL,
			ReferenceURLs:      qInput.ReferenceURLs,
			CreatedAt:          now,
		}
		domainQuestions = append(domainQuestions, q)
	}

	if err := u.repo.BulkCreateQuestions(ctx, domainQuestions); err != nil {
		return err
	}

	return nil
}

func (u *examUsecase) GetExamQuestions(ctx context.Context, examSetID string) ([]domain.Question, error) {
	if examSetID == "" {
		return nil, errors.Wrap(domain.ErrInvalidArgument, "examSetIDは必須です")
	}
	return u.repo.GetQuestionsByExamSetID(ctx, examSetID)
}

func (u *examUsecase) StartAttempt(ctx context.Context, userID string, req input.CreateAttemptRequest) (*domain.Attempt, error) {
	if userID == "" {
		return nil, errors.Wrap(domain.ErrUnauthenticated, "userIDは必須です")
	}
	if req.ExamID == "" || req.ExamSetID == "" {
		return nil, errors.Wrap(domain.ErrInvalidArgument, "examIDとexamSetIDは必須です")
	}

	// 実際のアプリでは、既に進行中のAttemptがあるかどうかを確認するかもしれません。
	// 現状は、単に新しいAttemptを作成します。

	// Attempt IDの生成 (例: UUID または AutoID)。
	// UUIDライブラリをまだインポートしていないため、単純な時間ベースのものを使用するか、AutoIDの場合はリポジトリに生成を任せます。
	// しかし、リポジトリは構造体内にIDがあることを期待しています。
	// とりあえず単純な文字列を使用するか、必要であれば google/uuid をインポートします。
	// リクエスト通り依存関係を少なく保つため、擬似ランダム文字列または時間を使用します。
	// 実際にはFirestore AutoIDが最適ですが、ドメインでIDを定義しています。
	// このMVPでは単純な時間ベースのIDを使用します。
	attemptID := fmt.Sprintf("%s_%d", userID, time.Now().UnixNano())

	now := time.Now()
	attempt := domain.Attempt{
		ID:             attemptID,
		UserID:         userID,
		ExamID:         req.ExamID,
		ExamSetID:      req.ExamSetID,
		Status:         domain.StatusInProgress,
		Score:          0,
		TotalQuestions: 0, // ExamSetのメタデータから取得するか、問題をカウントして設定すべき？
		                   // 現状は0、または問題をフェッチしてカウントします。
		CurrentIndex:   0,
		Answers:        make(map[string][]string),
		StartedAt:      now,
		UpdatedAt:      now,
	}

	if err := u.repo.SaveAttempt(ctx, attempt); err != nil {
		return nil, err
	}

	return &attempt, nil
}

func (u *examUsecase) UpdateAttempt(ctx context.Context, userID, attemptID string, req input.UpdateAttemptRequest) error {
	attempt, err := u.repo.GetAttempt(ctx, attemptID, userID)
	if err != nil {
		return err
	}

	if attempt.Status == domain.StatusCompleted {
		return errors.Wrap(domain.ErrFailedPrecondition, "試験は既に完了しています")
	}

	attempt.CurrentIndex = req.CurrentIndex
	attempt.Answers = req.Answers
	attempt.UpdatedAt = time.Now()

	return u.repo.SaveAttempt(ctx, *attempt)
}

func (u *examUsecase) CompleteAttempt(ctx context.Context, userID, attemptID string, req input.CompleteAttemptRequest) (*domain.Attempt, error) {
	attempt, err := u.repo.GetAttempt(ctx, attemptID, userID)
	if err != nil {
		return nil, err
	}

	if attempt.Status == domain.StatusCompleted {
		return nil, errors.Wrap(domain.ErrFailedPrecondition, "試験は既に完了しています")
	}

	// スコア計算のために問題を取得
	questions, err := u.repo.GetQuestionsByExamSetID(ctx, attempt.ExamSetID)
	if err != nil {
		return nil, err
	}

	// スコア計算
	score := 0
	totalQuestions := len(questions)
	
	// アクセスしやすいように問題をIDでマップ化
	qMap := make(map[string]domain.Question)
	for _, q := range questions {
		qMap[q.ID] = q
	}

	// 分野別成績の集計
	domainCorrect := make(map[string]int)
	domainTotal := make(map[string]int)

	for qID, userAnswers := range req.Answers {
		q, ok := qMap[qID]
		if !ok {
			continue
		}

		domainTotal[q.Domain]++

		// 正誤判定 (現状は完全一致)
		// 複数選択の場合、userAnswersが全てのcorrectAnswersを含み、かつそれ以外を含まないことを確認する必要があります。
		if isCorrect(userAnswers, q.CorrectAnswers) {
			score++
			domainCorrect[q.Domain]++
		}
	}

	// Attemptの更新
	now := time.Now()
	attempt.Answers = req.Answers
	attempt.Status = domain.StatusCompleted
	attempt.Score = score
	attempt.TotalQuestions = totalQuestions
	attempt.CompletedAt = &now
	attempt.UpdatedAt = now

	if err := u.repo.SaveAttempt(ctx, *attempt); err != nil {
		return nil, err
	}

	// 統計情報の更新
	// 既存の統計情報を取得
	stats, err := u.repo.GetUserExamStats(ctx, userID, attempt.ExamID)
	if err != nil {
		return nil, err
	}

	if stats == nil {
		stats = &domain.UserExamStats{
			ExamID:      attempt.ExamID,
			UserID:      userID,
			DomainStats: make(map[string]domain.DomainScore),
		}
	}

	stats.TotalAttempts++
	// 平均スコアの計算: (OldAvg * (N-1) + NewScore) / N ? 
	// または合計スコアを保持して割る？
	// 単純化のため、近似計算を行うか、TotalScoreが保存されていればより良いです。
	// 現在の構造体にはAverageScoreがあります。
	// これを更新すると仮定します。
	// NewAvg = ((OldAvg * (N-1)) + NewScore) / N
	// 注意: Scoreはint（正解数）、AverageScoreはfloat（パーセンテージ？またはカウント？）。
	// 通常、AverageScoreはパーセンテージまたはスケーリングされた値を意味します。
	// ここではAverageScoreをパーセンテージ(0-100)と仮定します。
	// attempt.Scoreは生のカウントです。
	
	attemptPercentage := float64(attempt.Score) / float64(attempt.TotalQuestions) * 100
	if attempt.TotalQuestions == 0 {
		attemptPercentage = 0
	}

	currentTotalAvg := stats.AverageScore * float64(stats.TotalAttempts-1)
	stats.AverageScore = (currentTotalAvg + attemptPercentage) / float64(stats.TotalAttempts)

	stats.LastTakenAt = now

	// 分野別成績の更新
	for dName, total := range domainTotal {
		correct := domainCorrect[dName]
		
		dScore, ok := stats.DomainStats[dName]
		if !ok {
			dScore = domain.DomainScore{DomainName: dName}
		}
		dScore.TotalCount += total
		dScore.CorrectCount += correct
		if dScore.TotalCount > 0 {
			dScore.AccuracyRate = int(float64(dScore.CorrectCount) / float64(dScore.TotalCount) * 100)
		}
		stats.DomainStats[dName] = dScore
	}

	if err := u.repo.UpdateStats(ctx, *stats); err != nil {
		// エラーをログに出力するが、Attempt完了自体は失敗させない？
		// 理想的にはリトライするか、失敗させるべきです。
		return nil, err
	}

	return attempt, nil
}

func (u *examUsecase) GetUserExamStats(ctx context.Context, userID, examID string) (*domain.UserExamStats, error) {
	if userID == "" || examID == "" {
		return nil, errors.Wrap(domain.ErrInvalidArgument, "userIDとexamIDは必須です")
	}
	return u.repo.GetUserExamStats(ctx, userID, examID)
}

// 正誤判定のヘルパー関数
func isCorrect(userAns, correctAns []string) bool {
	if len(userAns) != len(correctAns) {
		return false
	}
	// ソートまたはマップチェック。通常は数が少ないため、単純なループで十分です。
	// またはマップに変換します。
	cMap := make(map[string]bool)
	for _, c := range correctAns {
		cMap[c] = true
	}
	for _, u := range userAns {
		if !cMap[u] {
			return false
		}
	}
	return true
}
