package middleware

import (
	"net/http"

	"github.com/Jerinji2016/halooid/backend/internal/auth"
	"github.com/Jerinji2016/halooid/backend/internal/rbac"
	"github.com/google/uuid"
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
			// Get token claims from context
			claims, err := auth.GetTokenClaimsFromContext(r.Context())
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// First check if the permission is in the token claims
			// This is more efficient than making a database call
			for _, p := range claims.Permissions {
				if p == permission {
					// User has the permission in their token
					next.ServeHTTP(w, r)
					return
				}
			}

			// Get organization ID from token or use default
			orgID := claims.OrgID
			if orgID == uuid.Nil {
				orgID = DefaultOrganizationID
			}

			// If not found in token, check the database as a fallback
			// This handles cases where permissions might have been updated since the token was issued
			hasPermission, err := m.rbacService.HasPermission(r.Context(), claims.UserID, orgID, permission)
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
			// Get token claims from context
			claims, err := auth.GetTokenClaimsFromContext(r.Context())
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// First check if any of the permissions are in the token claims
			// This is more efficient than making a database call
			for _, requiredPermission := range permissions {
				for _, userPermission := range claims.Permissions {
					if userPermission == requiredPermission {
						// User has at least one of the required permissions in their token
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			// Get organization ID from token or use default
			orgID := claims.OrgID
			if orgID == uuid.Nil {
				orgID = DefaultOrganizationID
			}

			// If not found in token, check the database as a fallback
			// This handles cases where permissions might have been updated since the token was issued
			for _, permission := range permissions {
				hasPermission, err := m.rbacService.HasPermission(r.Context(), claims.UserID, orgID, permission)
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
			// Get token claims from context
			claims, err := auth.GetTokenClaimsFromContext(r.Context())
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// First check if all of the permissions are in the token claims
			// This is more efficient than making a database call
			allPermissionsInToken := true
			permissionsToCheck := make([]string, 0)

			for _, requiredPermission := range permissions {
				found := false
				for _, userPermission := range claims.Permissions {
					if userPermission == requiredPermission {
						found = true
						break
					}
				}

				if !found {
					allPermissionsInToken = false
					permissionsToCheck = append(permissionsToCheck, requiredPermission)
				}
			}

			// If all permissions are in the token, we can proceed
			if allPermissionsInToken {
				next.ServeHTTP(w, r)
				return
			}

			// Get organization ID from token or use default
			orgID := claims.OrgID
			if orgID == uuid.Nil {
				orgID = DefaultOrganizationID
			}

			// Check the database for the permissions not found in the token
			for _, permission := range permissionsToCheck {
				hasPermission, err := m.rbacService.HasPermission(r.Context(), claims.UserID, orgID, permission)
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
