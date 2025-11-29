package repository_impl

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/cockroachdb/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"nearline/backend/internal/domain"
	"nearline/backend/internal/repository"
)

type userRepository struct {
	client *firestore.Client
}

func NewUserRepository(client *firestore.Client) repository.UserRepository {
	return &userRepository{client: client}
}

func (r *userRepository) Create(ctx context.Context, user domain.User) error {
	_, err := r.client.Collection("users").Doc(user.ID).Set(ctx, user)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	return nil
}

func (r *userRepository) Find(ctx context.Context, id string) (*domain.User, error) {
	doc, err := r.client.Collection("users").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to find user")
	}

	var user domain.User
	if err := doc.DataTo(&user); err != nil {
		return nil, errors.Wrap(err, "failed to decode user")
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	iter := r.client.Collection("users").Where("email", "==", email).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		// iterator.Done is returned when there are no more items
		// We treat this as NotFound
		if status.Code(err) == codes.NotFound || err.Error() == "no more items in iterator" {
			return nil, domain.ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to find user by email")
	}

	var user domain.User
	if err := doc.DataTo(&user); err != nil {
		return nil, errors.Wrap(err, "failed to decode user")
	}

	return &user, nil
}
