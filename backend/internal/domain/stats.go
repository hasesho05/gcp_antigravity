package domain

import "time"

// UserExamStats は「資格ごと(ExamID)」の累積成績です。
type UserExamStats struct {
	ExamID        string                 `json:"examId" firestore:"exam_id"`
	UserID        string                 `json:"userId" firestore:"user_id"`
	TotalAttempts int                    `json:"totalAttempts" firestore:"total_attempts"`
	AverageScore  float64                `json:"averageScore" firestore:"average_score"`
	DomainStats   map[string]DomainScore `json:"domainStats" firestore:"domain_stats"` // 分野ごとの集計
	LastTakenAt   time.Time              `firestore:"last_taken_at"`
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
