package usecase

import (
	"context"

	"nearline/backend/internal/domain"
	"nearline/backend/internal/repository"
)

type ExamUsecase interface {
	ListExams(ctx context.Context) ([]domain.Exam, error)
	GetExam(ctx context.Context, id string) (*domain.Exam, error)
	ListExamSets(ctx context.Context, examID string) ([]domain.ExamSet, error)
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

func (u *examUsecase) GetExam(ctx context.Context, id string) (*domain.Exam, error) {
	return u.examRepo.Find(ctx, id)
}

func (u *examUsecase) ListExamSets(ctx context.Context, examID string) ([]domain.ExamSet, error) {
	return u.examRepo.FindSets(ctx, examID)
}
