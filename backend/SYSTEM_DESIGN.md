System Design Document (Finalized: V2)1. Design Principles & GoalsCost Optimization: FirestoreのRead/Write回数を最小限に抑える（特にRead）。Schema Sharing (Quicktype): GoのStructをSSOT（Single Source of Truth）とし、フロントエンドと型安全に通信する。Access Control: User の Role および SubscriptionStatus に基づき、試験へのアクセスを厳密に制御する。2. Access Control Logic (Usecase Layer)試験開始リクエスト(POST /users/me/attempts)時に、ユーザーの情報を参照し、以下のルールでアクセスを制御します。User RoleSubscription StatusTarget ExamIDAccess ResultRoleProSubActiveAll ExamsAccess Granted (プロユーザーは全てアクセス可)RoleFreeN/Acloud_digital_leader (CDL)Access Granted (CDLのみ無料アクセス可)RoleFreeN/AOther Exams (PCD, ACE etc.)Access DeniedRoleAdminN/AAll ExamsAccess Granted (管理目的のため)RoleProSubExpired / SubCanceledAll ExamsAccess Denied (有料期間終了ユーザーは無料ユーザーと同等)3. Domain Model Definitions (Go Structs)すべての構造体は internal/domain パッケージ内にあり、json タグはAPI通信用、firestore タグはDB操作用に使用されます。3.1 User Entity (internal/domain/user.go)ユーザーの認証・権限・サブスクリプション情報を管理します。// User はFirebase AuthのUIDをベースとした、アプリケーション独自のユーザー情報を保持します。
type User struct {
	ID                 string             `json:"id" firestore:"id"` // Firebase Auth UID
	Email              string             `json:"email" firestore:"email"`
	Role               UserRole           `json:"role" firestore:"role"`
	SubscriptionStatus SubscriptionStatus `json:"subscriptionStatus" firestore:"subscription_status"`
	CreatedAt          time.Time          `json:"createdAt" firestore:"created_at"`
}

type UserRole string
const (
	RoleFree  UserRole = "free"
	RolePro   UserRole = "pro"
	RoleAdmin UserRole = "admin"
)

type SubscriptionStatus string
const (
	SubActive   SubscriptionStatus = "active"
	SubExpired  SubscriptionStatus = "expired"
	SubCanceled SubscriptionStatus = "canceled"
)
3.2 Question Entity (internal/domain/question.go)問題のマスターデータです。// Question は1つの問題を表すマスターデータです。
type Question struct {
	ID                 string         `json:"id" firestore:"id"`
	ExamID             string         `json:"examId" firestore:"exam_id"`
	ExamSetID          string         `json:"examSetId" firestore:"exam_set_id"`
	ExamCode           string         `json:"examCode" firestore:"exam_code"`
	QuestionText       string         `json:"question" firestore:"question_text"`
	QuestionType       string         `json:"questionType" firestore:"question_type"`
	Options            []AnswerOption `json:"answerOptions" firestore:"options"`
	CorrectAnswers     []string       `json:"correctAnswers" firestore:"correct_answers"`
	OverallExplanation string         `json:"overallExplanation" firestore:"overall_explanation"`
	Domain             string         `json:"domain" firestore:"domain"`
	ImageURL           string         `json:"imageUrl,omitempty" firestore:"image_url,omitempty"`
	ReferenceURLs      []string       `json:"referenceUrls,omitempty" firestore:"reference_urls"`
	CreatedAt          time.Time      `json:"createdAt" firestore:"created_at"`
}

// AnswerOption は問題の個々の選択肢です。
type AnswerOption struct {
	ID          string `json:"id" firestore:"id"`
	Text        string `json:"answer" firestore:"text"`
	Explanation string `json:"explanation" firestore:"explanation"`
}
3.3 Attempt Entity (internal/domain/attempt.go)ユーザーの1回の受験データを表します。// Attempt はユーザーの1回の受験データを表します。
type Attempt struct {
	ID             string              `json:"id" firestore:"id"`
	UserID         string              `json:"userId" firestore:"user_id"`
	ExamID         string              `json:"examId" firestore:"exam_id"`
	ExamSetID      string              `json:"examSetId" firestore:"exam_set_id"`
	Status         AttemptStatus       `json:"status" firestore:"status"`
	Score          int                 `json:"score" firestore:"score"`
	TotalQuestions int                 `json:"totalQuestions" firestore:"total_questions"`
	CurrentIndex   int                 `json:"currentIndex" firestore:"current_index"`
	Answers        map[string][]string `json:"answers" firestore:"answers"` // Key: QuestionID
	StartedAt      time.Time           `json:"startedAt" firestore:"started_at"`
	UpdatedAt      time.Time           `json:"updatedAt" firestore:"updated_at"`
	CompletedAt    *time.Time          `json:"completedAt,omitempty" firestore:"completed_at,omitempty"`
}

type AttemptStatus string
const (
	StatusInProgress AttemptStatus = "in_progress"
	StatusPaused     AttemptStatus = "paused"
	StatusCompleted  AttemptStatus = "completed"
)
3.4 Stats Entity (internal/domain/stats.go)資格ごとの累積成績集計データです。// UserExamStats は「資格ごと(ExamID)」の累積成績です。
type UserExamStats struct {
	ExamID        string                 `json:"examId" firestore:"exam_id"`
	UserID        string                 `json:"userId" firestore:"user_id"`
	TotalAttempts int                    `json:"totalAttempts" firestore:"total_attempts"`
	AverageScore  float64                `json:"averageScore" firestore:"average_score"`
	DomainStats   map[string]DomainScore `json:"domainStats" firestore:"domain_stats"`
	LastTakenAt   time.Time              `json:"lastTakenAt" firestore:"last_taken_at"`
}

// DomainScore は特定分野ごとの成績集計です。
type DomainScore struct {
	DomainName   string `json:"domainName" firestore:"domain_name"`
	CorrectCount int    `json:"correctCount" firestore:"correct_count"`
	TotalCount   int    `json:"totalCount" firestore:"total_count"`
	AccuracyRate int    `json:"accuracyRate" firestore:"accuracy_rate"`
}
4. Database Design (Firestore)4.1 CollectionsCollection NameDoc ID PatternDomain EntityDescriptionusers{User.ID} (Firebase UID)Userユーザーのロール、サブスクリプションを保存。questions{ExamCode}_{SetID}_{Index}Questionマスターデータ。試験セット単位で一括Read。users/{uid}/attemptsAuto IDAttemptトランザクションデータ。中断/完了時に一括保存。users/{uid}/stats{ExamID}UserExamStats集計データ。試験完了時にBackendで差分更新（Increment）。


package domain

import "time"

// User はFirebase AuthのUIDをベースとした、アプリケーション独自のユーザー情報を保持します。
// ロールとサブスクリプションの状態を管理し、アクセス制御に使用されます。
type User struct {
	ID                 string             `json:"id" firestore:"id"` // Firebase Auth UID
	Email              string             `json:"email" firestore:"email"`
	Role               UserRole           `json:"role" firestore:"role"` // free, pro, admin
	SubscriptionStatus SubscriptionStatus `json:"subscriptionStatus" firestore:"subscription_status"` // active, expired, canceled
	CreatedAt          time.Time          `json:"createdAt" firestore:"created_at"`
}

// UserRole はユーザーの権限レベルを定義します。
type UserRole string

const (
	RoleFree  UserRole = "free"  // 無料ユーザー
	RolePro   UserRole = "pro"   // 有料サブスクリプションユーザー
	RoleAdmin UserRole = "admin" // 管理者
)

// SubscriptionStatus はサブスクリプションの状態を定義します。
type SubscriptionStatus string

const (
	SubActive   SubscriptionStatus = "active"
	SubExpired  SubscriptionStatus = "expired"
	SubCanceled SubscriptionStatus = "canceled"
)




package domain

import "time"

// Question は1つの問題を表すマスターデータです。
type Question struct {
	ID                 string         `json:"id" firestore:"id"`                                   // Document ID (e.g., "PCD_SET1_001")
	ExamID             string         `json:"examId" firestore:"exam_id"`                          // 資格ID (e.g., "professional_cloud_developer")
	ExamSetID          string         `json:"examSetId" firestore:"exam_set_id"`                   // 模擬試験セットID (e.g., "practice_exam_1")
	ExamCode           string         `json:"examCode" firestore:"exam_code"`                      // 資格コード (e.g., "PCD")
	QuestionText       string         `json:"question" firestore:"question_text"`                  // HTML string
	QuestionType       string         `json:"questionType" firestore:"question_type"`              // "multiple-choice" or "multi-select"
	Options            []AnswerOption `json:"answerOptions" firestore:"options"`                   // 選択肢リスト
	CorrectAnswers     []string       `json:"correctAnswers" firestore:"correct_answers"`          // 正解のOption IDリスト
	OverallExplanation string         `json:"overallExplanation" firestore:"overall_explanation"`  // 全体の解説 (HTML)
	Domain             string         `json:"domain" firestore:"domain"`                           // 分野 (e.g. "Compute")
	ImageURL           string         `json:"imageUrl,omitempty" firestore:"image_url,omitempty"`  // 解説図などのURL
	ReferenceURLs      []string       `json:"referenceUrls,omitempty" firestore:"reference_urls"`  // 参考リンク
	CreatedAt          time.Time      `json:"createdAt" firestore:"created_at"`
}

// AnswerOption は問題の個々の選択肢です。
type AnswerOption struct {
	ID          string `json:"id" firestore:"id"`                   // "a", "b", "c", "d" or UUID
	Text        string `json:"answer" firestore:"text"`             // 選択肢の文言
	Explanation string `json:"explanation" firestore:"explanation"` // この選択肢ごとの解説
}



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

package domain

import "time"

// UserExamStats は「資格ごと(ExamID)」の累積成績です。
type UserExamStats struct {
	ExamID        string                 `json:"examId" firestore:"exam_id"`
	UserID        string                 `json:"userId" firestore:"user_id"`
	TotalAttempts int                    `json:"totalAttempts" firestore:"total_attempts"`
	AverageScore  float64                `json:"averageScore" firestore:"average_score"`
	DomainStats   map[string]DomainScore `json:"domainStats" firestore:"domain_stats"` // 分野ごとの集計
	LastTakenAt   time.Time              `json:"lastTakenAt" firestore:"last_taken_at"`
}

// DomainScore は特定分野ごとの成績集計です。
type DomainScore struct {
	DomainName   string `json:"domainName" firestore:"domain_name"`
	CorrectCount int    `json:"correctCount" firestore:"correct_count"`
	TotalCount   int    `json:"totalCount" firestore:"total_count"`
	AccuracyRate int    `json:"accuracyRate" firestore:"accuracy_rate"` // パーセンテージ (0-100)
}