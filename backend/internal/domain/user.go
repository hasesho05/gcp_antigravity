package domain

import "time"

// User はFirebase AuthのUIDをベースとした、アプリケーション独自のユーザー情報を保持します。
// ロールとサブスクリプションの状態を管理し、アクセス制御に使用されます。
type User struct {
	ID                 string             `json:"id" firestore:"id"` // Firebase Auth UID
	Email              string             `json:"email" firestore:"email"`
	Provider           AuthProvider       `json:"provider" firestore:"provider"` // google, password, etc.
	Role               UserRole           `json:"role" firestore:"role"` // free, pro, admin
	SubscriptionStatus SubscriptionStatus `json:"subscriptionStatus" firestore:"subscription_status"` // active, expired, canceled
	CreatedAt          time.Time          `json:"createdAt" firestore:"created_at"`
}

func NewUser(id, email string, provider AuthProvider) *User {
	return &User{
		ID:                 id,
		Email:              email,
		Provider:           provider,
		Role:               RoleFree,
		SubscriptionStatus: SubActive,
		CreatedAt:          time.Now(),
	}
}

// AuthProvider は認証プロバイダーを定義します。
type AuthProvider string

// tygo:enum
const (
	ProviderGoogle   AuthProvider = "google.com"
	ProviderPassword AuthProvider = "password"
	ProviderGithub   AuthProvider = "github.com"
)

// UserRole はユーザーの権限レベルを定義します。
type UserRole string

// tygo:enum
const (
	RoleFree  UserRole = "free"  // 無料ユーザー
	RolePro   UserRole = "pro"   // 有料サブスクリプションユーザー
	RoleAdmin UserRole = "admin" // 管理者
)

// SubscriptionStatus はサブスクリプションの状態を定義します。
//
// tygo:enum
type SubscriptionStatus string

const (
	SubActive   SubscriptionStatus = "active"
	SubExpired  SubscriptionStatus = "expired"
	SubCanceled SubscriptionStatus = "canceled"
)
