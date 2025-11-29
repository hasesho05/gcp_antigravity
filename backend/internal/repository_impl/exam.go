package repository_impl

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/cockroachdb/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func (r *examRepository) Find(ctx context.Context, id string) (*domain.Exam, error) {
	doc, err := r.client.Collection("exams").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrNotFound
		}
		return nil, errors.Wrap(err, "試験の取得に失敗しました")
	}

	var exam domain.Exam
	if err := doc.DataTo(&exam); err != nil {
		return nil, errors.Wrap(err, "ドキュメントの変換に失敗しました")
	}
	return &exam, nil
}

func (r *examRepository) FindSets(ctx context.Context, examID string) ([]domain.ExamSet, error) {
	var examSets []domain.ExamSet
	iter := r.client.Collection("exams").Doc(examID).Collection("sets").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "試験セットの取得に失敗しました")
		}

		var examSet domain.ExamSet
		if err := doc.DataTo(&examSet); err != nil {
			return nil, errors.Wrap(err, "ドキュメントの変換に失敗しました")
		}
		examSets = append(examSets, examSet)
	}
	return examSets, nil
}
