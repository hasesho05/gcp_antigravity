package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/cockroachdb/errors"

	"gcp_antigravity/backend/internal/domain"
	"gcp_antigravity/backend/internal/repository"
	"gcp_antigravity/backend/internal/usecase/input"
	"github.com/google/uuid"
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
	qRepo    repository.QuestionRepository
	aRepo    repository.AttemptRepository
	sRepo    repository.UserStatsRepository
}

func NewExamUsecase(
	qRepo repository.QuestionRepository,
	aRepo repository.AttemptRepository,
	sRepo repository.UserStatsRepository,
) ExamUsecase {
	return &examUsecase{
		qRepo:    qRepo,
		aRepo:    aRepo,
		sRepo:    sRepo,
	}
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

		q, err := domain.NewQuestion(
			id,
			req.ExamID,
			req.ExamSetID,
			req.ExamCode,
			qInput.QuestionText,
			qInput.QuestionType,
			qInput.OverallExplanation,
			qInput.Domain,
			qInput.ImageURL,
			options,
			qInput.CorrectAnswers,
			qInput.ReferenceURLs,
			now,
		)
		if err != nil {
			return err
		}
		domainQuestions = append(domainQuestions, *q)
	}

	if err := u.qRepo.BulkCreate(ctx, domainQuestions); err != nil {
		return err
	}

	return nil
}

func (u *examUsecase) GetExamQuestions(ctx context.Context, examSetID string) ([]domain.Question, error) {
	if examSetID == "" {
		return nil, errors.Wrap(domain.ErrInvalidArgument, "examSetIDは必須です")
	}
	return u.qRepo.FindByExamSetID(ctx, examSetID)
}

func (u *examUsecase) StartAttempt(ctx context.Context, userID string, req input.CreateAttemptRequest) (*domain.Attempt, error) {
	if userID == "" {
		return nil, errors.Wrap(domain.ErrUnauthenticated, "userIDは必須です")
	}
	if req.ExamID == "" || req.ExamSetID == "" {
		return nil, errors.Wrap(domain.ErrInvalidArgument, "examIDとexamSetIDは必須です")
	}

	// 問題を取得して合計数を設定
	questions, err := u.qRepo.FindByExamSetID(ctx, req.ExamSetID)
	if err != nil {
		return nil, errors.Wrap(err, "attempt開始時の問題取得に失敗しました")
	}
	if len(questions) == 0 {
		return nil, errors.Wrap(domain.ErrNotFound, "指定された試験セットに問題が見つかりません")
	}

	attemptID := uuid.NewString()
	now := time.Now()

	attempt, err := domain.NewAttempt(attemptID, userID, req.ExamID, req.ExamSetID, len(questions), now)
	if err != nil {
		return nil, err
	}

	if err := u.aRepo.Save(ctx, *attempt); err != nil {
		return nil, err
	}

	return attempt, nil
}

func (u *examUsecase) UpdateAttempt(ctx context.Context, userID, attemptID string, req input.UpdateAttemptRequest) error {
	attempt, err := u.aRepo.Find(ctx, attemptID, userID)
	if err != nil {
		return err
	}

	if attempt.Status == domain.StatusCompleted {
		return errors.Wrap(domain.ErrFailedPrecondition, "試験は既に完了しています")
	}

	attempt.CurrentIndex = req.CurrentIndex
	attempt.Answers = req.Answers
	attempt.UpdatedAt = time.Now()

	return u.aRepo.Save(ctx, *attempt)
}

func (u *examUsecase) CompleteAttempt(ctx context.Context, userID, attemptID string, req input.CompleteAttemptRequest) (*domain.Attempt, error) {
	var completedAttempt *domain.Attempt

	err := u.sRepo.RunTransaction(ctx, func(ctx context.Context) error {
		attempt, err := u.aRepo.Find(ctx, attemptID, userID)
		if err != nil {
			return err
		}

		if attempt.Status == domain.StatusCompleted {
			return errors.Wrap(domain.ErrFailedPrecondition, "試験は既に完了しています")
		}

		questions, err := u.qRepo.FindByExamSetID(ctx, attempt.ExamSetID)
		if err != nil {
			return err
		}

		score := 0
		qMap := make(map[string]domain.Question)
		for _, q := range questions {
			qMap[q.ID] = q
		}

		domainCorrect := make(map[string]int)
		domainTotal := make(map[string]int)

		for qID, userAnswers := range req.Answers {
			q, ok := qMap[qID]
			if !ok {
				continue
			}
			domainTotal[q.Domain]++
			if isCorrect(userAnswers, q.CorrectAnswers) {
				score++
				domainCorrect[q.Domain]++
			}
		}

		now := time.Now()
		attempt.Answers = req.Answers
		attempt.Status = domain.StatusCompleted
		attempt.Score = score
		attempt.CompletedAt = &now
		attempt.UpdatedAt = now

		stats, err := u.sRepo.Find(ctx, userID, attempt.ExamID)
		if err != nil {
			return err
		}
		if stats == nil {
			stats, err = domain.NewUserExamStats(userID, attempt.ExamID)
			if err != nil {
				return err
			}
		}

		stats.TotalAttempts++
		stats.TotalScore += score
		stats.TotalQuestionsAnswered += attempt.TotalQuestions
		stats.LastTakenAt = now

		for dName, total := range domainTotal {
			correct := domainCorrect[dName]
			dScore, ok := stats.DomainStats[dName]
			if !ok {
				dScore = domain.DomainScore{DomainName: dName}
			}
			dScore.TotalCount += total
			dScore.CorrectCount += correct
			if dScore.TotalCount > 0 {
				dScore.AccuracyRate = int(math.Round(float64(dScore.CorrectCount) / float64(dScore.TotalCount) * 100))
			}
			stats.DomainStats[dName] = dScore
		}

		if err := u.aRepo.Save(ctx, *attempt); err != nil {
			return err
		}
		if err := u.sRepo.Save(ctx, *stats); err != nil {
			return err
		}

		completedAttempt = attempt
		return nil
	})

	if err != nil {
		return nil, err
	}

	return completedAttempt, nil
}

func (u *examUsecase) GetUserExamStats(ctx context.Context, userID, examID string) (*domain.UserExamStats, error) {
	if userID == "" || examID == "" {
		return nil, errors.Wrap(domain.ErrInvalidArgument, "userIDとexamIDは必須です")
	}
	return u.sRepo.Find(ctx, userID, examID)
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
