package repository

import (
	"context"

	"nearline/backend/internal/domain"
)

// ExamRepository はExamエンティティの永続化を管理します。
type ExamRepository interface {
	FindAll(ctx context.Context) ([]domain.Exam, error)
	Find(ctx context.Context, id string) (*domain.Exam, error)
	FindSets(ctx context.Context, examID string) ([]domain.ExamSet, error)
}
