package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Jerinji2016/halooid/backend/internal/auth"
	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/google/uuid"
)

// Using auth package's context keys

// AuthMiddleware provides authentication middleware
type AuthMiddleware struct {
	authService auth.Service
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(authService auth.Service) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// Authenticate authenticates a request using JWT
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		tokenString, err := extractTokenFromHeader(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate token
		claims, err := m.authService.ValidateToken(r.Context(), tokenString, models.AccessToken)
		if err != nil {
			if errors.Is(err, auth.ErrExpiredToken) {
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user information to context
		ctx := auth.SetUserIDInContext(r.Context(), claims.UserID)
		ctx = auth.SetTokenInContext(ctx, tokenString)
		ctx = auth.SetTokenClaimsInContext(ctx, claims)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extracts the user ID from the context
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	return auth.GetUserIDFromContext(ctx)
}

// GetEmail extracts the email from the context
func GetEmail(ctx context.Context) (string, error) {
	claims, err := auth.GetTokenClaimsFromContext(ctx)
	if err != nil {
		return "", err
	}
	return claims.Email, nil
}

// extractTokenFromHeader extracts the token from the Authorization header
func extractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	// Check if the header has the Bearer prefix
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}
