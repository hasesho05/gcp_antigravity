package repository

import (
	"context"

	"nearline/backend/internal/domain"
)

// QuestionRepository は問題エンティティの永続化を管理します。
type QuestionRepository interface {
	BulkCreate(ctx context.Context, questions []domain.Question) error
	FindByExamSetID(ctx context.Context, examSetID string) ([]domain.Question, error)
}
