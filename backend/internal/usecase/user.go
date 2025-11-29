package usecase

import (
	"context"

	"github.com/cockroachdb/errors"

	"nearline/backend/internal/domain"
	"nearline/backend/internal/repository"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, id string, email string, provider domain.AuthProvider) (*domain.User, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) CreateUser(ctx context.Context, id string, email string, provider domain.AuthProvider) (*domain.User, error) {
	// 既に存在するか確認 (ID)
	existingUser, err := u.userRepo.Find(ctx, id)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, errors.Wrap(err, "failed to check existing user by ID")
	}
	if existingUser != nil {
		return existingUser, nil // 既に存在する場合はそのユーザーを返す（冪等性）
	}

	// Emailの重複チェック
	existingEmailUser, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, errors.Wrap(err, "failed to check existing user by email")
	}
	if existingEmailUser != nil {
		return nil, domain.ErrAlreadyExists
	}

	newUser := domain.NewUser(id, email, provider)

	if err := u.userRepo.Create(ctx, *newUser); err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	return newUser, nil
}

func (u *userUsecase) GetUser(ctx context.Context, id string) (*domain.User, error) {
	user, err := u.userRepo.Find(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}
	return user, nil
}
