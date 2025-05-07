package models

import (
	"time"

	"github.com/google/uuid"
)

// TokenType represents the type of token
type TokenType string

const (
	// AccessToken is used for API access
	AccessToken TokenType = "access"
	// RefreshToken is used to get a new access token
	RefreshToken TokenType = "refresh"
)

// TokenPair represents an access token and refresh token pair
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // Access token expiration in seconds
}

// TokenClaims represents the claims in a JWT token
type TokenClaims struct {
	UserID      uuid.UUID `json:"user_id"`
	Email       string    `json:"email"`
	TokenType   TokenType `json:"token_type"`
	Roles       []string  `json:"roles,omitempty"`
	Permissions []string  `json:"permissions,omitempty"`
	OrgID       uuid.UUID `json:"org_id,omitempty"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// RefreshTokenRequest represents the request to refresh an access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
