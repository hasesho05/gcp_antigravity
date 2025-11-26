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

type userStatsRepository struct {
	client *firestore.Client
}

func NewUserStatsRepository(client *firestore.Client) repository.UserStatsRepository {
	return &userStatsRepository{client: client}
}

func (r *userStatsRepository) Save(ctx context.Context, stats domain.UserExamStats) error {
	if stats.UserID == "" || stats.ExamID == "" {
		return errors.New("UserIDとExamIDは必須です")
	}

	docRef := r.client.Collection("users").Doc(stats.UserID).Collection("stats").Doc(stats.ExamID)
	
	if tx, ok := GetTransaction(ctx); ok {
		return tx.Set(docRef, stats)
	}

	_, err := docRef.Set(ctx, stats)
	if err != nil {
		return errors.Wrap(err, "firestore: statsの保存に失敗しました")
	}
	return nil
}

func (r *userStatsRepository) Find(ctx context.Context, userID, examID string) (*domain.UserExamStats, error) {
	if userID == "" || examID == "" {
		return nil, errors.New("UserIDとExamIDは必須です")
	}
	docRef := r.client.Collection("users").Doc(userID).Collection("stats").Doc(examID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil // Not found is not an error, usecase will handle it
		}
		return nil, errors.Wrap(err, "firestore: statsの取得に失敗しました")
	}

	var stats domain.UserExamStats
	if err := doc.DataTo(&stats); err != nil {
		return nil, errors.Wrap(err, "firestore: statsのデータマッピングに失敗しました")
	}

	return &stats, nil
}
