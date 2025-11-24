package usecase

import (
	"gcp_antigravity/backend/internal/repository"
)

type ExamUsecase interface {
}

type examUsecase struct {
	qRepo    repository.QuestionRepository
	aRepo    repository.AttemptRepository
	sRepo    repository.UserStatsRepository
	txRepo   repository.TransactionRepository
}

func NewExamUsecase(
	qRepo repository.QuestionRepository,
	aRepo repository.AttemptRepository,
	sRepo repository.UserStatsRepository,
	txRepo repository.TransactionRepository,
) ExamUsecase {
	return &examUsecase{
		qRepo:    qRepo,
		aRepo:    aRepo,
		sRepo:    sRepo,
		txRepo:   txRepo,
	}
}

