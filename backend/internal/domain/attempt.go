package domain

import "time"

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
	CompletedAt    *time.Time          `json:"completedAt,omitempty" firestore:"completed_at,omitempty"`
}

// AttemptStatus は受験の進捗状態を定義します。
type AttemptStatus string

const (
	StatusInProgress AttemptStatus = "in_progress" // 進行中
	StatusPaused     AttemptStatus = "paused"      // 中断中
	StatusCompleted  AttemptStatus = "completed"   // 完了
)
