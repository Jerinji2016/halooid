package auth

import (
	"context"
	"errors"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/google/uuid"
)

// Context keys
type contextKey string

const (
	userIDKey    contextKey = "user_id"
	tokenKey     contextKey = "token"
	tokenClaimsKey contextKey = "token_claims"
)

// Common errors
var (
	ErrUserIDNotFound    = errors.New("user ID not found in context")
	ErrTokenNotFound     = errors.New("token not found in context")
	ErrTokenClaimsNotFound = errors.New("token claims not found in context")
)

// SetUserIDInContext sets the user ID in the context
func SetUserIDInContext(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserIDFromContext gets the user ID from the context
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrUserIDNotFound
	}
	return userID, nil
}

// SetTokenInContext sets the token in the context
func SetTokenInContext(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

// GetTokenFromContext gets the token from the context
func GetTokenFromContext(ctx context.Context) (string, error) {
	token, ok := ctx.Value(tokenKey).(string)
	if !ok {
		return "", ErrTokenNotFound
	}
	return token, nil
}

// SetTokenClaimsInContext sets the token claims in the context
func SetTokenClaimsInContext(ctx context.Context, claims *models.TokenClaims) context.Context {
	return context.WithValue(ctx, tokenClaimsKey, claims)
}

// GetTokenClaimsFromContext gets the token claims from the context
func GetTokenClaimsFromContext(ctx context.Context) (*models.TokenClaims, error) {
	claims, ok := ctx.Value(tokenClaimsKey).(*models.TokenClaims)
	if !ok {
		return nil, ErrTokenClaimsNotFound
	}
	return claims, nil
}
