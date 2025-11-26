package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"

	"nearline/backend/internal/domain"
	"nearline/backend/internal/repository"
	"nearline/backend/internal/usecase/input"
	"nearline/backend/internal/usecase/output"
)

type QuestionUsecase interface {
	UploadQuestions(ctx context.Context, req input.UploadQuestionsRequest) error
	GetExamQuestions(ctx context.Context, input *input.GetExamQuestions) ([]domain.Question, error)
}

type questionUsecase struct {
	qRepo repository.QuestionRepository
}

func NewQuestionUsecase(qRepo repository.QuestionRepository) QuestionUsecase {
	return &questionUsecase{qRepo: qRepo}
}

func (u *questionUsecase) UploadQuestions(ctx context.Context, req input.UploadQuestionsRequest) error {
	if len(req.Questions) == 0 {
		return errors.Wrap(domain.ErrInvalidArgument, "問題が提供されていません")
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

		q, err := domain.NewQuestion(
			id,
			req.ExamID,
			req.ExamSetID,
			req.ExamCode,
			qInput.QuestionText,
			qInput.QuestionType,
			qInput.OverallExplanation,
			qInput.Domain,
			qInput.ImageURL,
			options,
			qInput.CorrectAnswers,
			qInput.ReferenceURLs,
			now,
		)
		if err != nil {
			return err
		}
		domainQuestions = append(domainQuestions, *q)
	}

	if err := u.qRepo.BulkCreate(ctx, domainQuestions); err != nil {
		return err
	}

	return nil
}

func (u *questionUsecase) GetExamQuestions(ctx context.Context, input *input.GetExamQuestions) ([]domain.Question, error) {
	questions, err := u.qRepo.FindByExamSetID(ctx, input.ExamSetID)
	if err != nil {
		return nil, err
	}
	return output.NewQuestions(questions), nil
}
