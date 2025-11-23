package middleware

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(authClient *auth.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "認証トークンが必要です", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == authHeader {
				http.Error(w, "認証トークンの形式が無効です", http.StatusUnauthorized)
				return
			}

			// Verify ID Token
			decodedToken, err := authClient.VerifyIDToken(r.Context(), token)
			if err != nil {
				http.Error(w, "認証トークンが無効です", http.StatusUnauthorized)
				return
			}

			// Set UID in context
			ctx := context.WithValue(r.Context(), UserIDKey, decodedToken.UID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
