package repository_impl

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/cockroachdb/errors"

	"gcp_antigravity/backend/internal/domain"
	"gcp_antigravity/backend/internal/repository"
)

type questionRepository struct {
	client *firestore.Client
}

func NewQuestionRepository(client *firestore.Client) repository.QuestionRepository {
	return &questionRepository{client: client}
}

func (r *questionRepository) BulkCreate(ctx context.Context, questions []domain.Question) error {
	batch := r.client.Batch()
	for _, q := range questions {
		if q.ID == "" {
			return errors.New("質問IDは必須です")
		}
		docRef := r.client.Collection("questions").Doc(q.ID)
		batch.Set(docRef, q)
	}

	_, err := batch.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "firestore: failed to bulk create questions")
	}

	return nil
}

func (r *questionRepository) FindByExamSetID(ctx context.Context, examSetID string) ([]domain.Question, error) {
	iter := r.client.Collection("questions").Where("exam_set_id", "==", examSetID).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "firestore: failed to get questions")
	}

	var questions []domain.Question
	for _, doc := range docs {
		var q domain.Question
		if err := doc.DataTo(&q); err != nil {
			// 1つのドキュメントのマッピングに失敗しても、全体を失敗させる
			return nil, errors.Wrap(err, "firestore: failed to map question data")
		}
		questions = append(questions, q)
	}

	return questions, nil
}
