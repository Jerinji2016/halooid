package middleware

import (
	"errors"

	"github.com/Jerinji2016/halooid/backend/internal/auth"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ErrUserNotFound is returned when the user ID is not found in the context
var ErrUserNotFound = errors.New("user ID not found in context")

// GetUserIDFromContext retrieves the user ID from the Echo context
func GetUserIDFromContext(c echo.Context) (uuid.UUID, error) {
	// Get token claims from context
	claims, ok := c.Get("user").(*auth.TokenClaims)
	if !ok || claims == nil {
		return uuid.Nil, ErrUserNotFound
	}
	
	return claims.UserID, nil
}
