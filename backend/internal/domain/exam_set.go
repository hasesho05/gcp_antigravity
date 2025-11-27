package domain

import "time"

// ExamSet は模擬試験のセットを表します（例: "Practice Exam 1"）。
// Firestore Path: exams/{examID}/sets/{id}
type ExamSet struct {
	ID          string    `json:"id" firestore:"id"`                   // 例: "practice_exam_1"
	ExamID      string    `json:"examId" firestore:"exam_id"`          // 親のExam ID
	Name        string    `json:"name" firestore:"name"`               // 例: "Practice Exam 1"
	Description string    `json:"description" firestore:"description"` // 例: "50 questions covering all domains"
	QuestionIDs []string  `json:"questionIds" firestore:"question_ids"` // 含まれる問題IDのリスト (冗長化)
	CreatedAt   time.Time `json:"createdAt" firestore:"created_at"`
}
