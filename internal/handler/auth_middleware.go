package handler

import (
	"banksystem/internal/service"
	"context"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	jwtService *service.JWTService
	logger     Logger
}

func NewAuthMiddleware(jwtService *service.JWTService, logger Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		logger:     logger,
	}
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.logger.Printf("No Authorization header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.logger.Printf("Invalid Authorization header format")
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		userID, err := m.jwtService.ValidateToken(parts[1])
		if err != nil {
			m.logger.Printf("Token validation failed: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
