package usecase

import (
	"context"
	"math"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"

	"gcp_antigravity/backend/internal/domain"
	"gcp_antigravity/backend/internal/repository"
	"gcp_antigravity/backend/internal/usecase/input"
	"gcp_antigravity/backend/internal/usecase/output"
)

type AttemptUsecase interface {
	StartAttempt(ctx context.Context, userID string, req input.CreateAttemptRequest) (*domain.Attempt, error)
	UpdateAttempt(ctx context.Context, userID, attemptID string, req input.UpdateAttemptRequest) error
	CompleteAttempt(ctx context.Context, input *input.CompleteAttempt) (*domain.Attempt, error)
}


type attemptUsecase struct {
	qRepo  repository.QuestionRepository
	aRepo  repository.AttemptRepository
	sRepo  repository.UserStatsRepository
	txRepo repository.TransactionRepository
}

func NewAttemptUsecase(
	qRepo repository.QuestionRepository,
	aRepo repository.AttemptRepository,
	sRepo repository.UserStatsRepository,
	txRepo repository.TransactionRepository,
) AttemptUsecase {
	return &attemptUsecase{
		qRepo:  qRepo,
		aRepo:  aRepo,
		sRepo:  sRepo,
		txRepo: txRepo,
	}
}

func (u *attemptUsecase) StartAttempt(ctx context.Context, userID string, req input.CreateAttemptRequest) (*domain.Attempt, error) {
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

func (u *attemptUsecase) UpdateAttempt(ctx context.Context, userID, attemptID string, req input.UpdateAttemptRequest) error {
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

func (u *attemptUsecase) CompleteAttempt(ctx context.Context, input *input.CompleteAttempt) (*domain.Attempt, error) {
	var completedAttempt *domain.Attempt

	err := u.txRepo.Run(ctx, func(txCtx context.Context) error {
		attempt, err := u.aRepo.Find(txCtx, input.AttemptID, input.UserID)
		if err != nil {
			return err
		}

		if attempt.Status == domain.StatusCompleted {
			return errors.Wrap(domain.ErrFailedPrecondition, "試験は既に完了しています")
		}

		questions, err := u.qRepo.FindByExamSetID(txCtx, attempt.ExamSetID)
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

		for qID, userAnswers := range input.Answers {
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
		attempt.Answers = input.Answers
		attempt.Status = domain.StatusCompleted
		attempt.Score = score
		attempt.CompletedAt = &now
		attempt.UpdatedAt = now

		stats, err := u.sRepo.Find(txCtx, input.UserID, attempt.ExamID)
		if err != nil {
			return err
		}
		if stats == nil {
			stats, err = domain.NewUserExamStats(input.UserID, attempt.ExamID)
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

		if err := u.aRepo.Save(txCtx, *attempt); err != nil {
			return err
		}
		if err := u.sRepo.Save(txCtx, *stats); err != nil {
			return err
		}

		completedAttempt = attempt
		return nil
	})

	if err != nil {
		return nil, err
	}

	return output.NewAttemptOutput(completedAttempt), nil
}

// isCorrect is needed by CompleteAttempt
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
