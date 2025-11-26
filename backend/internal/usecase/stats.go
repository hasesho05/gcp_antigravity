package usecase

import (
	"context"

	"nearline/backend/internal/domain"
	"nearline/backend/internal/repository"
	"nearline/backend/internal/usecase/input"
	"nearline/backend/internal/usecase/output"
)

type StatsUsecase interface {
	GetUserExamStats(ctx context.Context, input *input.GetUserExamStats) (*domain.UserExamStats, error)
}

type statsUsecase struct {
	sRepo repository.UserStatsRepository
}

func NewStatsUsecase(sRepo repository.UserStatsRepository) StatsUsecase {
	return &statsUsecase{sRepo: sRepo}
}

func (u *statsUsecase) GetUserExamStats(ctx context.Context, input *input.GetUserExamStats) (*domain.UserExamStats, error) {
	stats, err := u.sRepo.Find(ctx, input.UserID, input.ExamID)
	if err != nil {
		return nil, err // DBエラーなどの場合はそのまま返す
	}

	if stats == nil {
		// 統計情報が存在しない場合は、新しい空の統計オブジェクトを生成して返す
		stats, err = domain.NewUserExamStats(input.UserID, input.ExamID)
		if err != nil {
			// userID, examIDは検証済みなので基本的には発生しない
			return nil, err
		}
	}

	return output.NewUserExamStats(stats), nil
}
