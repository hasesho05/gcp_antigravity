package repository

import (
	"context"

	"gcp_antigravity/backend/internal/domain"
)

// QuestionRepository は問題エンティティの永続化を管理します。
type QuestionRepository interface {
	BulkCreate(ctx context.Context, questions []domain.Question) error
	FindByExamSetID(ctx context.Context, examSetID string) ([]domain.Question, error)
}

// AttemptRepository は受験記録エンティティの永続化を管理します。
type AttemptRepository interface {
	Save(ctx context.Context, attempt domain.Attempt) error
	Find(ctx context.Context, attemptID string, userID string) (*domain.Attempt, error)
}

// UserStatsRepository はユーザーの試験統計エンティティの永続化を管理します。
type UserStatsRepository interface {
	Save(ctx context.Context, stats domain.UserExamStats) error
	Find(ctx context.Context, userID, examID string) (*domain.UserExamStats, error)
}

// TransactionRepositoryはトランザクション操作を管理します。
type TransactionRepository interface {
	Run(ctx context.Context, f func(txCtx context.Context) error) error
}
