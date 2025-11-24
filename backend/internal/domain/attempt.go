package domain

import (
	"time"

	"github.com/cockroachdb/errors"
)

// Attempt はユーザーの1回の受験データを表します。
type Attempt struct {
	ID             string              `json:"id" firestore:"id"`
	UserID         string              `json:"userId" firestore:"user_id"`
	ExamID         string              `json:"examId" firestore:"exam_id"`         // 資格ID
	ExamSetID      string              `json:"examSetId" firestore:"exam_set_id"`  // 模擬試験セットID
	Status         AttemptStatus       `json:"status" firestore:"status"`
	Score          int                 `json:"score" firestore:"score"`
	TotalQuestions int                 `json:"totalQuestions" firestore:"total_questions"`
	CurrentIndex   int                 `json:"currentIndex" firestore:"current_index"`
	Answers        map[string][]string `json:"answers" firestore:"answers"` // Key: QuestionID, Value: Selected Option IDs
	StartedAt      time.Time           `json:"startedAt" firestore:"started_at"`
	UpdatedAt      time.Time           `json:"updatedAt" firestore:"updated_at"`
	CompletedAt    *time.Time          `firestore:"completed_at,omitempty"`
}

// NewAttempt は新しいAttemptドメインオブジェクトを生成します。
func NewAttempt(id, userID, examID, examSetID string, totalQuestions int, now time.Time) (*Attempt, error) {
	if id == "" || userID == "" || examID == "" || examSetID == "" {
		return nil, errors.New("AttemptのID, UserID, ExamID, ExamSetIDは必須です")
	}

	return &Attempt{
		ID:             id,
		UserID:         userID,
		ExamID:         examID,
		ExamSetID:      examSetID,
		Status:         StatusInProgress,
		Score:          0,
		TotalQuestions: totalQuestions,
		CurrentIndex:   0,
		Answers:        make(map[string][]string),
		StartedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// AttemptStatus は受験の進捗状態を定義します。
type AttemptStatus string

const (
	StatusInProgress AttemptStatus = "in_progress" // 進行中
	StatusPaused     AttemptStatus = "paused"      // 中断中
	StatusCompleted  AttemptStatus = "completed"   // 完了
)

