package input

import (
	"gcp_antigravity/backend/internal/domain"

	"github.com/cockroachdb/errors"
)

type GetExamQuestions struct {
	ExamSetID string
}

func NewGetExamQuestions(examSetID string) (*GetExamQuestions, error) {
	if examSetID == "" {
		return nil, errors.Wrap(domain.ErrInvalidArgument, "examSetID is required")
	}

	return &GetExamQuestions{
		ExamSetID: examSetID,
	}, nil
}
