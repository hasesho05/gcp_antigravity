package domain

import "time"

// UserExamStats は「資格ごと(ExamID)」の累積成績です。
type UserExamStats struct {
	UserID                 string                 `firestore:"user_id"`
	ExamID                 string                 `firestore:"exam_id"`
	TotalAttempts          int                    `firestore:"total_attempts"`
	TotalScore             int                    `firestore:"total_score"` // 全てのAttemptでの合計正解数
	TotalQuestionsAnswered int                    `firestore:"total_questions_answered"` // 全てのAttemptでの合計問題数
	DomainStats            map[string]DomainScore `firestore:"domain_stats"`
	LastTakenAt            time.Time              `firestore:"last_taken_at"`
}

// NewUserExamStats は新しいUserExamStatsドメインオブジェクトを生成します。
func NewUserExamStats(userID, examID string) (*UserExamStats, error) {
	if userID == "" || examID == "" {
		return nil, errors.New("統計情報のUserIDとExamIDは必須です")
	}
	return &UserExamStats{
		UserID:      userID,
		ExamID:      examID,
		DomainStats: make(map[string]DomainScore),
	}, nil
}

// DomainScore は特定分野ごとの成績集計です。
type DomainScore struct {
	DomainName   string `json:"domainName" firestore:"domain_name"`
	CorrectCount int    `json:"correctCount" firestore:"correct_count"`
	TotalCount   int    `json:"totalCount" firestore:"total_count"`
	AccuracyRate int    `json:"accuracyRate" firestore:"accuracy_rate"` // パーセンテージ (0-100)
}
