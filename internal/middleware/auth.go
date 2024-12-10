package middleware

import (
	"context"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/handler/errrender"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/user/project/internal/service"
)

const userIDContextKey = "user_id"

func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				errrender.RenderError(w, r, contract.ErrUnauthorized, "missing auth header")
				return
			}

			// Bearer token validation
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				errrender.RenderError(w, r, contract.ErrUnauthorized, "invalid auth header")
				return
			}

			tokenString := parts[1]

			// Parse the JWT token
			token, err := jwt.ParseWithClaims(tokenString, &service.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})
			if err != nil || !token.Valid {
				errrender.RenderError(w, r, contract.ErrUnauthorized, "invalid token")
				return
			}

			// Extract claims
			claims, ok := token.Claims.(*service.JWTClaims)
			if !ok {
				errrender.RenderError(w, r, contract.ErrUnauthorized, "invalid token")
				return
			}

			if claims.UID <= 0 {
				errrender.RenderError(w, r, contract.ErrUnauthorized, "invalid token")
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), userIDContextKey, claims.UID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
