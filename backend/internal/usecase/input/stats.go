package input

import (
	"nearline/backend/internal/domain"

	"github.com/cockroachdb/errors"
)

type GetUserExamStats struct {
	UserID string
	ExamID string
}

func NewGetUserExamStats(userID, examID string) (*GetUserExamStats, error) {
	if userID == "" {
		return nil, errors.Wrap(domain.ErrUnauthenticated, "userID is required")
	}
	if examID == "" {
		return nil, errors.Wrap(domain.ErrInvalidArgument, "examID is required")
	}

	return &GetUserExamStats{
		UserID: userID,
		ExamID: examID,
	}, nil
}
