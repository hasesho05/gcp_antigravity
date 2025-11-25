package usecase

import (
	"context"

	"gcp_antigravity/backend/internal/domain"
	"gcp_antigravity/backend/internal/repository"
)

type ExamUsecase interface {
	ListExams(ctx context.Context) ([]domain.Exam, error)
}

type examUsecase struct {
	examRepo repository.ExamRepository
	qRepo    repository.QuestionRepository
	aRepo    repository.AttemptRepository
	sRepo    repository.UserStatsRepository
	txRepo   repository.TransactionRepository
}

func NewExamUsecase(
	examRepo repository.ExamRepository,
	qRepo repository.QuestionRepository,
	aRepo repository.AttemptRepository,
	sRepo repository.UserStatsRepository,
	txRepo repository.TransactionRepository,
) ExamUsecase {
	return &examUsecase{
		examRepo: examRepo,
		qRepo:    qRepo,
		aRepo:    aRepo,
		sRepo:    sRepo,
		txRepo:   txRepo,
	
}

}

func (u *examUsecase) ListExams(ctx context.Context) ([]domain.Exam, error) {
	return u.examRepo.FindAll(ctx)
}
