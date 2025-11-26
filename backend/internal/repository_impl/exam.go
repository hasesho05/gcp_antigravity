package repository_impl

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/cockroachdb/errors"
	"google.golang.org/api/iterator"

	"nearline/backend/internal/domain"
	"nearline/backend/internal/repository"
)

type examRepository struct {
	client *firestore.Client
}

func NewExamRepository(client *firestore.Client) repository.ExamRepository {
	return &examRepository{client: client}
}

func (r *examRepository) FindAll(ctx context.Context) ([]domain.Exam, error) {
	var exams []domain.Exam
	iter := r.client.Collection("exams").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "試験の取得に失敗しました")
		}

		var exam domain.Exam
		if err := doc.DataTo(&exam); err != nil {
			return nil, errors.Wrap(err, "ドキュメントの変換に失敗しました")
		}
		exams = append(exams, exam)
	}
	return exams, nil
}
