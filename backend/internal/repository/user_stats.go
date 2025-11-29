package repository

import (
	"context"

	"nearline/backend/internal/domain"
)

// UserStatsRepository はユーザーの試験統計エンティティの永続化を管理します。
type UserStatsRepository interface {
	Save(ctx context.Context, stats domain.UserExamStats) error
	Find(ctx context.Context, userID, examID string) (*domain.UserExamStats, error)
}
