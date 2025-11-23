package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"

	"gcp_antigravity/backend/internal/domain"
	"gcp_antigravity/backend/internal/repository"
	"gcp_antigravity/backend/internal/usecase/input"
)

type ExamUsecase interface {
	UploadQuestions(ctx context.Context, req input.UploadQuestionsRequest) error
	GetExamQuestions(ctx context.Context, examSetID string) ([]domain.Question, error)
	StartAttempt(ctx context.Context, userID string, req input.CreateAttemptRequest) (*domain.Attempt, error)
}

type examUsecase struct {
	repo repository.ExamRepository
}

func NewExamUsecase(repo repository.ExamRepository) ExamUsecase {
	return &examUsecase{repo: repo}
}

func (u *examUsecase) UploadQuestions(ctx context.Context, req input.UploadQuestionsRequest) error {
	if len(req.Questions) == 0 {
		return errors.Wrap(domain.ErrInvalidArgument, "no questions provided")
	}

	var domainQuestions []domain.Question
	now := time.Now()

	for _, qInput := range req.Questions {
		// Generate ID: {ExamCode}_{SetID}_{Index}
		// e.g. PCD_SET1_001
		id := fmt.Sprintf("%s_%s_%03d", req.ExamCode, req.ExamSetID, qInput.Index)

		var options []domain.AnswerOption
		for _, o := range qInput.Options {
			options = append(options, domain.AnswerOption{
				ID:          o.ID,
				Text:        o.Text,
				Explanation: o.Explanation,
			})
		}

		q := domain.Question{
			ID:                 id,
			ExamID:             req.ExamID,
			ExamSetID:          req.ExamSetID,
			ExamCode:           req.ExamCode,
			QuestionText:       qInput.QuestionText,
			QuestionType:       qInput.QuestionType,
			Options:            options,
			CorrectAnswers:     qInput.CorrectAnswers,
			OverallExplanation: qInput.OverallExplanation,
			Domain:             qInput.Domain,
			ImageURL:           qInput.ImageURL,
			ReferenceURLs:      qInput.ReferenceURLs,
			CreatedAt:          now,
		}
		domainQuestions = append(domainQuestions, q)
	}

	if err := u.repo.BulkCreateQuestions(ctx, domainQuestions); err != nil {
		return err
	}

	return nil
}

func (u *examUsecase) GetExamQuestions(ctx context.Context, examSetID string) ([]domain.Question, error) {
	if examSetID == "" {
		return nil, errors.Wrap(domain.ErrInvalidArgument, "examSetID is required")
	}
	return u.repo.GetQuestionsByExamSetID(ctx, examSetID)
}

func (u *examUsecase) StartAttempt(ctx context.Context, userID string, req input.CreateAttemptRequest) (*domain.Attempt, error) {
	if userID == "" {
		return nil, errors.Wrap(domain.ErrUnauthenticated, "userID is required")
	}
	if req.ExamID == "" || req.ExamSetID == "" {
		return nil, errors.Wrap(domain.ErrInvalidArgument, "examID and examSetID are required")
	}

	// In a real app, we might check if there is already an in-progress attempt.
	// For now, we just create a new one.

	// Generate Attempt ID (e.g. UUID or AutoID). 
	// Since we don't have a UUID lib imported yet, we'll use a simple time-based one or rely on Repository to generate if it was AutoID.
	// But our repository expects an ID in the struct.
	// Let's assume we can use a simple string for now or import google/uuid if needed.
	// To keep dependencies low as per request, I'll use a pseudo-random string or just time.
	// Actually, Firestore AutoID is best, but we are defining ID in domain.
	// Let's use a simple time-based ID for this MVP.
	attemptID := fmt.Sprintf("%s_%d", userID, time.Now().UnixNano())

	now := time.Now()
	attempt := domain.Attempt{
		ID:             attemptID,
		UserID:         userID,
		ExamID:         req.ExamID,
		ExamSetID:      req.ExamSetID,
		Status:         domain.StatusInProgress,
		Score:          0,
		TotalQuestions: 0, // Should be populated from ExamSet metadata or by counting questions? 
		                   // For now 0, or we fetch questions to count them.
		CurrentIndex:   0,
		Answers:        make(map[string][]string),
		StartedAt:      now,
		UpdatedAt:      now,
	}

	if err := u.repo.SaveAttempt(ctx, attempt); err != nil {
		return nil, err
	}

	return &attempt, nil
}
