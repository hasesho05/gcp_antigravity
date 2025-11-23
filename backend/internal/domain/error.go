package domain

import "github.com/cockroachdb/errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrInvalidArgument   = errors.New("invalid argument")
	ErrInternal          = errors.New("internal error")
	ErrUnauthenticated   = errors.New("unauthenticated")
	ErrPermissionDenied  = errors.New("permission denied")
	ErrAlreadyExists     = errors.New("already exists")
	ErrFailedPrecondition = errors.New("failed precondition")
)
