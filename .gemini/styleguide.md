GCP認定資格プラットフォーム バックエンドスタイルガイド本ドキュメントは、GCP認定資格模擬試験プラットフォームのGoバックエンドにおける設計原則、ディレクトリ構造、コーディング規則、および確定されたドメインモデルを定義します。1. プロジェクト概要と技術スタック1.1. プロジェクトの目的個人開発プロジェクトとして、低コスト（FirestoreのRead/Write最適化）で運用可能なGCP認定資格模擬試験プラットフォームのバックエンドを構築します。将来的なネイティブアプリ（JSON API）との連携を見据え、クリーンで拡張性の高い設計を採用します。1.2. 技術スタックCategoryTechnologyNoteBackendGo (1.22+)Cloud Run, Standard net/httpDatabaseFirestoreNoSQL, コスト最適化されたスキーマAuthFirebase Authenticationユーザー認証とUID管理Errorgithub.com/cockroachdb/errorsスタックトレース付きのエラー管理ToolQuicktypeGo StructからTypeScript型定義を自動生成2. アーキテクチャとディレクトリ構成2.1. アーキテクチャ原則Clean Architecture + Domain Driven Design (Lightweight) を採用し、依存関係を外側（infra）から内側（domain）へと一方向にする原則を厳守します。Domain: アプリケーションの核となるデータ構造とビジネスルールを定義。他の層に依存しない。Usecase: アプリケーション固有の操作（採点、進捗保存、アクセス制御）を定義。Repository: データアクセスのインターフェース（抽象）を定義。Repository Impl: Repositoryインターフェースの具象実装。Infra: 外部システム（Firestore）との接続、共通ヘルパーなど低レイヤーの処理を定義。Handler: 外部からのI/F（HTTP/JSON）に特化し、Usecaseを呼び出す。2.2. ディレクトリ構造.
├── Makefile                 # ビルド、テスト、型定義生成コマンド
├── cmd
│   └── api
│       └── main.go          # エントリーポイント (DIとルーティング)
├── internal
│   ├── domain               # 1. ドメイン層 (純粋なエンティティ)
│   │   ├── user.go
│   │   ├── exam.go
│   │   ├── exam_set.go
│   │   ├── question.go
│   │   ├── attempt.go
│   │   ├── stats.go
│   │   └── error.go
│   ├── handler              # 2. プレゼンテーション層
│   │   ├── admin
│   │   └── client
│   ├── usecase              # 3. ユースケース層 (Interactor)
│   │   ├── exam.go
│   │   ├── input
│   │   └── output
│   ├── repository           # 4. リポジトリインターフェース
│   │   └── repository.go
│   ├── repository_impl      # 5. リポジトリ実装層
│   │   └── exam.go
│   └── infra                # 6. インフラ層 (Firestoreドライバ)
│       └── firestore
│           └── client.go
└── scripts
    └── dump_json.go         # Quicktype用JSON生成スクリプト
3. Goコーディング規則3.1. ドメインモデルの命名規則Struct名は単数形 (User, Question, Attempt) を使用する。フィールド名はキャメルケースを使用する。JSON TagとFirestore Tagの厳守:json タグは キャメルケース (examId)firestore タグは スネークケース (exam_id)タグが一つでも欠けるとQuicktype連携またはDB操作に支障をきたすため、必ず両方記述する。3.2. エラーハンドリング標準エラー: github.com/cockroachdb/errors を利用する。(backend/internal/domain/error.go)


Wrapの利用: repository や infra 層からエラーが返される際は、必ず errors.Wrap(err, "...") を使用し、呼び出し元のコンテキスト情報を付加する。これによりスタックトレースが保持され、デバッグが容易になる。ハンドラでの処理: handler 層でエラーがキャッチされた場合、fmt.Printf("%+v\n", err) を使用してスタックトレースをログに出力する。3.3. DTOとレスポンス処理DTO（Data Transfer Object）は usecase/input および usecase/output に配置する。レスポンス用の構造体は output 内部に定義し、ToResponseData() のようなポインタレシーバメソッドを通じてドメインエンティティをレスポンス形式に変換する。4. データモデルと型定義4.1. Quicktype連携ワークフローinternal/domain および internal/usecase/output のGo構造体が、フロントエンドのTypeScript型定義のSingle Source of Truth (SSOT) となります。Go Structを更新。make generate-sample でJSONサンプルを生成。QuicktypeでJSONをTSインターフェースに変換。4.2. ドメインモデル (Go Structs)Exam Entity (internal/domain/exam.go)GCP認定試験そのものを表します。// Exam は認定試験を表します（例: "Google Cloud Certified - Professional Cloud Developer"）。
type Exam struct {
	ID          string    `json:"id" firestore:"id"`                   // 例: "professional_cloud_developer"
	Code        string    `json:"code" firestore:"code"`               // 例: "PCD"
	Name        string    `json:"name" firestore:"name"`               // 例: "Professional Cloud Developer"
	Description string    `json:"description" firestore:"description"` // 例: "あなたの能力を評価します..."
	ImageURL    string    `json:"imageUrl" firestore:"image_url"`      // 試験のロゴ/アイコンのURL
	CreatedAt   time.Time `json:"createdAt" firestore:"created_at"`
}
ExamSet Entity (internal/domain/exam_set.go)1つの資格試験に含まれる、模擬試験の単位です。// ExamSet は模擬試験のセットを表します（例: "Practice Exam 1"）。
// Firestore Path: exams/{examID}/sets/{id}
type ExamSet struct {
	ID          string    `json:"id" firestore:"id"`                   // 例: "practice_exam_1"
	ExamID      string    `json:"examId" firestore:"exam_id"`          // 親のExam ID
	Name        string    `json:"name" firestore:"name"`               // 例: "Practice Exam 1"
	Description string    `json:"description" firestore:"description"` // 例: "50 questions covering all domains"
	QuestionIDs []string  `json:"questionIds" firestore:"question_ids"` // 含まれる問題IDのリスト (冗長化)
	CreatedAt   time.Time `json:"createdAt" firestore:"created_at"`
}
User Entity (internal/domain/user.go)サブスクリプションとアクセス制御の基盤となる情報です。// User はFirebase AuthのUIDをベースとした、アプリケーション独自のユーザー情報を保持します。
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
Question Entity (internal/domain/question.go)問題のマスターデータ。リッチテキスト(QuestionText, OverallExplanation)を想定しています。// Question は1つの問題を表すマスターデータです。
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
Attempt Entity (internal/domain/attempt.go)ユーザーの受験履歴。中断・再開に必要な情報を全て含みます。// Attempt はユーザーの1回の受験データを表します。
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

// AttemptStatus は受験の進捗状態を定義します。
type AttemptStatus string

const (
	StatusInProgress AttemptStatus = "in_progress" // 進行中
	StatusPaused     AttemptStatus = "paused"      // 中断中
	StatusCompleted  AttemptStatus = "completed"   // 完了
)
Stats Entity (internal/domain/stats.go)累積成績と分野別弱点分析のためのデータ構造です。// UserExamStats は「資格ごと(ExamID)」の累積成績です。
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
	AccuracyRate int    `json:"accuracyRate" firestore:"accuracy_rate"` // パーセンテージ (0-100)
}
5. Firestore設計とアクセス制御5.1. コレクション設計（コスト最適化）サブコレクションとMap構造を多用し、Read/Write回数を削減します。Collection NameDoc ID Pattern目的exams/{examID}/sets/{setID}資格試験(Exam)、模擬試験セット(ExamSet)、問題(Question)のマスターデータ。`sets`や`questions`はサブコレクション。users/{User.ID} (Firebase UID)ユーザーのロールとサブスクリプション管理。users/{uid}/attemptsAuto IDトランザクションデータ。中断/完了時のみWriteし、コストを最小化。users/{uid}/stats{ExamID}集計データ。試験完了時にBackendでトランザクション更新。5.2. アクセス制御ロジック (Usecase Layer)試験開始リクエスト時（POST /users/me/attempts）に、以下のルールでアクセスを制御します。User RoleSubscription StatusTarget ExamIDAccess ResultRoleProSubActiveAll Examsアクセス許可RoleFreeN/Acloud_digital_leader (CDL)アクセス許可 (無料ユーザー特典)RoleFreeN/AOther Examsアクセス拒否RoleAdminN/AAll Examsアクセス許可RoleProSubExpired / SubCanceledAll Examsアクセス拒否 (無料ユーザーと同等に扱う)