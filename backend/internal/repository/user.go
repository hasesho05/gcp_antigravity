package repository

import (
	"context"

	"nearline/backend/internal/domain"
)

// UserRepository はユーザーエンティティの永続化を管理します。
type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	Find(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}
