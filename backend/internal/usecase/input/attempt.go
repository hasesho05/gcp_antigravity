package input

import (
	"nearline/backend/internal/domain"

	"github.com/cockroachdb/errors"
)

type CompleteAttempt struct {
	UserID    string
	AttemptID string
	Answers   map[string][]string
}

func NewCompleteAttempt(userID, attemptID string, answers map[string][]string) (*CompleteAttempt, error) {
	if userID == "" {
		return nil, errors.Wrap(domain.ErrUnauthenticated, "userID is required")
	}
	if attemptID == "" {
		return nil, errors.Wrap(domain.ErrInvalidArgument, "attemptID is required")
	}

	return &CompleteAttempt{
		UserID:    userID,
		AttemptID: attemptID,
		Answers:   answers,
	}, nil
}
