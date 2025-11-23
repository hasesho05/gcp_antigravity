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
