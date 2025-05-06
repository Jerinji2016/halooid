package rbac

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Handlers provides HTTP handlers for RBAC
type Handlers struct {
	service Service
}

// NewHandlers creates a new Handlers
func NewHandlers(service Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

// CreateRole handles role creation
func (h *Handlers) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req models.RoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Name == "" {
		http.Error(w, "Role name is required", http.StatusBadRequest)
		return
	}
	if len(req.Permissions) == 0 {
		http.Error(w, "At least one permission is required", http.StatusBadRequest)
		return
	}

	// Create role
	role, err := h.service.CreateRole(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrRoleNameExists) {
			http.Error(w, "Role name already exists", http.StatusConflict)
			return
		}
		if errors.Is(err, ErrPermissionNotFound) {
			http.Error(w, "One or more permissions not found", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}

// GetRole handles retrieving a role
func (h *Handlers) GetRole(w http.ResponseWriter, r *http.Request) {
	// Get role ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	// Get role
	role, err := h.service.GetRoleByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			http.Error(w, "Role not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

// ListRoles handles retrieving all roles
func (h *Handlers) ListRoles(w http.ResponseWriter, r *http.Request) {
	// Get roles
	roles, err := h.service.ListRoles(r.Context())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

// UpdateRole handles updating a role
func (h *Handlers) UpdateRole(w http.ResponseWriter, r *http.Request) {
	// Get role ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	var req models.RoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Name == "" {
		http.Error(w, "Role name is required", http.StatusBadRequest)
		return
	}
	if len(req.Permissions) == 0 {
		http.Error(w, "At least one permission is required", http.StatusBadRequest)
		return
	}

	// Update role
	role, err := h.service.UpdateRole(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			http.Error(w, "Role not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrRoleNameExists) {
			http.Error(w, "Role name already exists", http.StatusConflict)
			return
		}
		if errors.Is(err, ErrPermissionNotFound) {
			http.Error(w, "One or more permissions not found", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

// DeleteRole handles deleting a role
func (h *Handlers) DeleteRole(w http.ResponseWriter, r *http.Request) {
	// Get role ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	// Delete role
	err = h.service.DeleteRole(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			http.Error(w, "Role not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.WriteHeader(http.StatusNoContent)
}

// CreatePermission handles permission creation
func (h *Handlers) CreatePermission(w http.ResponseWriter, r *http.Request) {
	var req models.PermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Name == "" {
		http.Error(w, "Permission name is required", http.StatusBadRequest)
		return
	}

	// Create permission
	permission, err := h.service.CreatePermission(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrPermissionNameExists) {
			http.Error(w, "Permission name already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(permission)
}

// GetPermission handles retrieving a permission
func (h *Handlers) GetPermission(w http.ResponseWriter, r *http.Request) {
	// Get permission ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	// Get permission
	permission, err := h.service.GetPermissionByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrPermissionNotFound) {
			http.Error(w, "Permission not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permission)
}

// ListPermissions handles retrieving all permissions
func (h *Handlers) ListPermissions(w http.ResponseWriter, r *http.Request) {
	// Get permissions
	permissions, err := h.service.ListPermissions(r.Context())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permissions)
}

// UpdatePermission handles updating a permission
func (h *Handlers) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	// Get permission ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	var req models.PermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Name == "" {
		http.Error(w, "Permission name is required", http.StatusBadRequest)
		return
	}

	// Update permission
	permission, err := h.service.UpdatePermission(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrPermissionNotFound) {
			http.Error(w, "Permission not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrPermissionNameExists) {
			http.Error(w, "Permission name already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permission)
}

// DeletePermission handles deleting a permission
func (h *Handlers) DeletePermission(w http.ResponseWriter, r *http.Request) {
	// Get permission ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	// Delete permission
	err = h.service.DeletePermission(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrPermissionNotFound) {
			http.Error(w, "Permission not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.WriteHeader(http.StatusNoContent)
}

// AssignRoleToUser handles assigning a role to a user
func (h *Handlers) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	var req models.AssignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Assign role to user
	err := h.service.AssignRoleToUser(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrRoleNotFound) {
			http.Error(w, "Role not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.WriteHeader(http.StatusNoContent)
}

// RemoveRoleFromUser handles removing a role from a user
func (h *Handlers) RemoveRoleFromUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID, role ID, and organization ID from URL
	vars := mux.Vars(r)
	userIDStr := vars["user_id"]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	roleIDStr := vars["role_id"]
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	orgIDStr := vars["org_id"]
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	// Remove role from user
	err = h.service.RemoveRoleFromUser(r.Context(), userID, roleID, orgID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrRoleNotFound) {
			http.Error(w, "Role not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.WriteHeader(http.StatusNoContent)
}

// GetUserRoles handles retrieving all roles for a user
func (h *Handlers) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	// Get user ID and organization ID from URL
	vars := mux.Vars(r)
	userIDStr := vars["user_id"]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	orgIDStr := vars["org_id"]
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	// Get user roles
	roles, err := h.service.GetUserRoles(r.Context(), userID, orgID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

// CheckPermission handles checking if a user has a specific permission
func (h *Handlers) CheckPermission(w http.ResponseWriter, r *http.Request) {
	// Get user ID, organization ID, and permission name from URL
	vars := mux.Vars(r)
	userIDStr := vars["user_id"]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	orgIDStr := vars["org_id"]
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	permissionName := vars["permission"]
	if permissionName == "" {
		http.Error(w, "Permission name is required", http.StatusBadRequest)
		return
	}

	// Check if user has permission
	hasPermission, err := h.service.HasPermission(r.Context(), userID, orgID, permissionName)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrInvalidPermission) {
			http.Error(w, "Invalid permission", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"has_permission": hasPermission})
}
