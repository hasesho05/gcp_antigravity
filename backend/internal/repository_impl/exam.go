package repository_impl

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/cockroachdb/errors"

	"gcp_antigravity/backend/internal/domain"
	"gcp_antigravity/backend/internal/repository"
)

type examRepository struct {
	client *firestore.Client
}

func NewExamRepository(client *firestore.Client) repository.ExamRepository {
	return &examRepository{client: client}
}

func (r *examRepository) BulkCreateQuestions(ctx context.Context, questions []domain.Question) error {
	batch := r.client.Batch()
	
	// Firestore Batch limit is 500 operations. 
	// Assuming questions length is reasonable for a single request (e.g. 50-60).
	// If larger, we should chunk it. For now, we assume < 500.

	for _, q := range questions {
		// Doc ID Pattern: {ExamCode}_{SetID}_{Index}
		// We assume ID is already generated in Usecase or we use the one provided.
		// According to design: "PCD_SET1_001"
		if q.ID == "" {
			return errors.New("question ID is required")
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

func (r *examRepository) GetQuestionsByExamSetID(ctx context.Context, examSetID string) ([]domain.Question, error) {
	// Query questions where exam_set_id == examSetID
	iter := r.client.Collection("questions").Where("exam_set_id", "==", examSetID).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "firestore: failed to get questions")
	}

	var questions []domain.Question
	for _, doc := range docs {
		var q domain.Question
		if err := doc.DataTo(&q); err != nil {
			return nil, errors.Wrap(err, "firestore: failed to map question data")
		}
		questions = append(questions, q)
	}

	return questions, nil
}

func (r *examRepository) SaveAttempt(ctx context.Context, attempt domain.Attempt) error {
	// Save to users/{uid}/attempts/{attemptID}
	// Note: The design says users/{uid}/attempts.
	// We need to ensure we are using the correct path.
	if attempt.UserID == "" || attempt.ID == "" {
		return errors.New("userID and attemptID are required")
	}

	docRef := r.client.Collection("users").Doc(attempt.UserID).Collection("attempts").Doc(attempt.ID)
	_, err := docRef.Set(ctx, attempt)
	if err != nil {
		return errors.Wrap(err, "firestore: failed to save attempt")
	}

	return nil
}

func (r *examRepository) GetAttempt(ctx context.Context, attemptID string, userID string) (*domain.Attempt, error) {
	docRef := r.client.Collection("users").Doc(userID).Collection("attempts").Doc(attemptID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		if !doc.Exists() {
			return nil, errors.Wrap(domain.ErrNotFound, "attempt not found")
		}
		return nil, errors.Wrap(err, "firestore: failed to get attempt")
	}

	var attempt domain.Attempt
	if err := doc.DataTo(&attempt); err != nil {
		return nil, errors.Wrap(err, "firestore: failed to map attempt data")
	}

	return &attempt, nil
}

func (r *examRepository) UpdateStats(ctx context.Context, stats domain.UserExamStats) error {
	// Placeholder for now. 
	// In C-305, we will implement this with transaction logic if needed, 
	// or simple set/merge if we calculate in memory.
	// Design says "Increment" in backend, so likely a transaction or merge.
	// For now, simple Set.
	
	if stats.UserID == "" || stats.ExamID == "" {
		return errors.New("userID and examID are required")
	}

	docRef := r.client.Collection("users").Doc(stats.UserID).Collection("stats").Doc(stats.ExamID)
	_, err := docRef.Set(ctx, stats)
	if err != nil {
		return errors.Wrap(err, "firestore: failed to update stats")
	}
	return nil
}
