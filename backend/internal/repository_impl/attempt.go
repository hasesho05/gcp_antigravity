package repository_impl

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/cockroachdb/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gcp_antigravity/backend/internal/domain"
	"gcp_antigravity/backend/internal/repository"
)

type attemptRepository struct {
	client *firestore.Client
}

func NewAttemptRepository(client *firestore.Client) repository.AttemptRepository {
	return &attemptRepository{client: client}
}

func (r *attemptRepository) Save(ctx context.Context, attempt domain.Attempt) error {
	if attempt.UserID == "" || attempt.ID == "" {
		return errors.New("UserIDとAttemptIDは必須です")
	}

	docRef := r.client.Collection("users").Doc(attempt.UserID).Collection("attempts").Doc(attempt.ID)
	
	if tx, ok := GetTransaction(ctx); ok {
		return tx.Set(docRef, attempt)
	}

	_, err := docRef.Set(ctx, attempt)
	if err != nil {
		return errors.Wrap(err, "firestore: attemptの保存に失敗しました")
	}

	return nil
}

func (r *attemptRepository) Find(ctx context.Context, attemptID string, userID string) (*domain.Attempt, error) {
	if userID == "" || attemptID == "" {
		return nil, errors.New("UserIDとAttemptIDは必須です")
	}
	docRef := r.client.Collection("users").Doc(userID).Collection("attempts").Doc(attemptID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, errors.Wrap(domain.ErrNotFound, "attemptが見つかりませんでした")
		}
		return nil, errors.Wrap(err, "firestore: attemptの取得に失敗しました")
	}

	var attempt domain.Attempt
	if err := doc.DataTo(&attempt); err != nil {
		return nil, errors.Wrap(err, "firestore: attemptのデータマッピングに失敗しました")
	}

	return &attempt, nil
}
