package domain

import "time"

// Exam は認定試験を表します（例: "Google Cloud Certified - Professional Cloud Developer"）。
type Exam struct {
	ID          string    `json:"id" firestore:"id"`                   // 例: "professional_cloud_developer"
	Code        string    `json:"code" firestore:"code"`               // 例: "PCD"
	Name        string    `json:"name" firestore:"name"`               // 例: "Professional Cloud Developer"
	Description string    `json:"description" firestore:"description"` // 例: "あなたの能力を評価します..."
	ImageURL    string    `json:"imageUrl" firestore:"image_url"`      // 試験のロゴ/アイコンのURL
	CreatedAt   time.Time `json:"createdAt" firestore:"created_at"`
}
