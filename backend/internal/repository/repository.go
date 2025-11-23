package repository

import (
	"context"

	"gcp_antigravity/backend/internal/domain"
)

type ExamRepository interface {
	BulkCreateQuestions(ctx context.Context, questions []domain.Question) error
	GetQuestionsByExamSetID(ctx context.Context, examSetID string) ([]domain.Question, error)
	SaveAttempt(ctx context.Context, attempt domain.Attempt) error
	GetAttempt(ctx context.Context, attemptID string, userID string) (*domain.Attempt, error)
	UpdateStats(ctx context.Context, stats domain.UserExamStats) error
}
