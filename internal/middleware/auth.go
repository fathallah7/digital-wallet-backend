package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/fathallah7/wallet-service/internal/handler"
)

type AuthMiddleware struct {
	jwtSecret []byte
}

func NewAuthMiddleware(jwtSecret []byte) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: jwtSecret}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			handler.WriteError(w, http.StatusUnauthorized, nil, "missing token")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handler.WriteError(w, http.StatusUnauthorized, nil, "invalid token format")
			return
		}
		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return m.jwtSecret, nil
		})
		if err != nil || !token.Valid {
			handler.WriteError(w, http.StatusUnauthorized, nil, "invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			handler.WriteError(w, http.StatusUnauthorized, nil, "invalid token claims")
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			handler.WriteError(w, http.StatusUnauthorized, nil, "invalid user id in token")
			return
		}

		ctx := context.WithValue(r.Context(), handler.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
