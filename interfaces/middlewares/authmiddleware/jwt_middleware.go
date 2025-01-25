// interfaces/middlewares/authmiddleware/jwt_middleware.go
package authmiddleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Abiti0233/memo-app/configs"
	"github.com/Abiti0233/memo-app/utils"
	"github.com/golang-jwt/jwt"
)

type key string

const (
	userIDKey key = "userID"
)

// JWTAuthentication はJWTトークンを検証するためのミドルウェア
func JWTAuthentication(config *configs.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "認可ヘッダーがありません", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "無効な認可ヘッダー形式", http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]
			claims, err := utils.ValidateJWT(tokenStr, config.JWTSecret)
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					http.Error(w, "無効なトークン署名", http.StatusUnauthorized)
					return
				}
				http.Error(w, "無効なトークン", http.StatusUnauthorized)
				return
			}

			// ExpiresAt に直接アクセスするのではなく、StandardClaims を介してアクセス
			if claims.StandardClaims.ExpiresAt < time.Now().Unix() {
				http.Error(w, "トークンの有効期限が切れました", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
