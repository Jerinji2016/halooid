package middleware

import (
	"net/http"

	"github.com/Jerinji2016/halooid/backend/internal/auth"
	"github.com/Jerinji2016/halooid/backend/internal/rbac"
)

// RBACMiddleware provides RBAC middleware
type RBACMiddleware struct {
	rbacService rbac.Service
}

// NewRBACMiddleware creates a new RBACMiddleware
func NewRBACMiddleware(rbacService rbac.Service) *RBACMiddleware {
	return &RBACMiddleware{
		rbacService: rbacService,
	}
}

// RequirePermission creates a middleware that requires a specific permission
func (m *RBACMiddleware) RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user ID from context
			userID, err := auth.GetUserIDFromContext(r.Context())
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get organization ID from request
			// In a real application, this would be determined based on the request
			// For now, we'll use a default organization ID
			orgID := DefaultOrganizationID

			// Check if user has permission
			hasPermission, err := m.rbacService.HasPermission(r.Context(), userID, orgID, permission)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if !hasPermission {
				http.Error(w, "Permission denied", http.StatusForbidden)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyPermission creates a middleware that requires any of the specified permissions
func (m *RBACMiddleware) RequireAnyPermission(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user ID from context
			userID, err := auth.GetUserIDFromContext(r.Context())
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get organization ID from request
			// In a real application, this would be determined based on the request
			// For now, we'll use a default organization ID
			orgID := DefaultOrganizationID

			// Check if user has any of the permissions
			for _, permission := range permissions {
				hasPermission, err := m.rbacService.HasPermission(r.Context(), userID, orgID, permission)
				if err != nil {
					continue
				}

				if hasPermission {
					// User has at least one of the required permissions
					next.ServeHTTP(w, r)
					return
				}
			}

			// User doesn't have any of the required permissions
			http.Error(w, "Permission denied", http.StatusForbidden)
		})
	}
}

// RequireAllPermissions creates a middleware that requires all of the specified permissions
func (m *RBACMiddleware) RequireAllPermissions(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user ID from context
			userID, err := auth.GetUserIDFromContext(r.Context())
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get organization ID from request
			// In a real application, this would be determined based on the request
			// For now, we'll use a default organization ID
			orgID := DefaultOrganizationID

			// Check if user has all of the permissions
			for _, permission := range permissions {
				hasPermission, err := m.rbacService.HasPermission(r.Context(), userID, orgID, permission)
				if err != nil || !hasPermission {
					// User doesn't have one of the required permissions
					http.Error(w, "Permission denied", http.StatusForbidden)
					return
				}
			}

			// User has all of the required permissions
			next.ServeHTTP(w, r)
		})
	}
}
