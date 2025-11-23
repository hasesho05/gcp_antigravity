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
	// バッチ書き込みを使用してコストを最適化
	batch := r.client.Batch()
	count := 0
	batchSize := 500 // Firestoreのバッチ制限は500オペレーションです。
	// 質問の数が単一のリクエストに対して妥当な範囲（例：50-60）であると仮定しています。
	// もしこれより多い場合は、チャンクに分割する必要があります。現時点では500未満と仮定しています。

	for _, q := range questions {
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
	// exam_set_id == examSetID の問題をクエリ
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
	// users/{uid}/attempts/{attemptID} に保存
	// 注意: 設計では users/{uid}/attempts となっています。
	// 正しいパスを使用していることを確認する必要があります。
	if attempt.UserID == "" || attempt.ID == "" {
		return errors.New("userIDとattemptIDは必須です")
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
	// 現時点ではプレースホルダー。
	// C-305では、必要であればトランザクションロジックを実装するか、
	// メモリ内で計算して単純なSet/Mergeを行います。
	// 設計ではバックエンドでの「インクリメント」となっているため、トランザクションまたはMergeが適しています。
	// 現状は単純なSetとします。
	
	if stats.UserID == "" || stats.ExamID == "" {
		return errors.New("userIDとexamIDは必須です")
	}

	docRef := r.client.Collection("users").Doc(stats.UserID).Collection("stats").Doc(stats.ExamID)
	_, err := docRef.Set(ctx, stats)
	if err != nil {
		return errors.Wrap(err, "firestore: failed to update stats")
	}
	return nil
}

func (r *examRepository) GetUserExamStats(ctx context.Context, userID, examID string) (*domain.UserExamStats, error) {
	docRef := r.client.Collection("users").Doc(userID).Collection("stats").Doc(examID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		if !doc.Exists() {
			// 見つからない場合は空の統計情報を返すか、エラーを返すか？
			// 通常はnilまたは特定のエラーを返すのが良い。
			// ここではnilを返し、ユースケース側で「見つからない」を「新規統計情報」として扱うようにします。
			return nil, nil 
		}
		return nil, errors.Wrap(err, "firestore: failed to get stats")
	}

	var stats domain.UserExamStats
	if err := doc.DataTo(&stats); err != nil {
		return nil, errors.Wrap(err, "firestore: failed to map stats data")
	}

	return &stats, nil
}
