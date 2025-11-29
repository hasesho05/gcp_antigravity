package repository

import (
	"context"

	"nearline/backend/internal/domain"
)

// AttemptRepository は受験記録エンティティの永続化を管理します。
type AttemptRepository interface {
	Save(ctx context.Context, attempt domain.Attempt) error
	Find(ctx context.Context, attemptID string, userID string) (*domain.Attempt, error)
}
