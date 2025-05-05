package models

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a role in the system
type Role struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	Permissions []Permission `json:"permissions,omitempty" db:"-"`
}

// Permission represents a permission in the system
type Permission struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// RolePermission represents a many-to-many relationship between roles and permissions
type RolePermission struct {
	RoleID       uuid.UUID `json:"role_id" db:"role_id"`
	PermissionID uuid.UUID `json:"permission_id" db:"permission_id"`
}

// UserRole represents a many-to-many relationship between users and roles
type UserRole struct {
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	RoleID         uuid.UUID `json:"role_id" db:"role_id"`
	OrganizationID uuid.UUID `json:"organization_id" db:"organization_id"`
}

// RoleRequest represents the data needed to create or update a role
type RoleRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions" validate:"required"`
}

// PermissionRequest represents the data needed to create or update a permission
type PermissionRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// AssignRoleRequest represents the data needed to assign a role to a user
type AssignRoleRequest struct {
	UserID         uuid.UUID `json:"user_id" validate:"required"`
	RoleID         uuid.UUID `json:"role_id" validate:"required"`
	OrganizationID uuid.UUID `json:"organization_id" validate:"required"`
}

// RoleResponse represents the role data returned to clients
type RoleResponse struct {
	ID          uuid.UUID          `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// PermissionResponse represents the permission data returned to clients
type PermissionResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts a Role to a RoleResponse
func (r *Role) ToResponse() RoleResponse {
	permissionResponses := make([]PermissionResponse, 0, len(r.Permissions))
	for _, p := range r.Permissions {
		permissionResponses = append(permissionResponses, p.ToResponse())
	}

	return RoleResponse{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Permissions: permissionResponses,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// ToResponse converts a Permission to a PermissionResponse
func (p *Permission) ToResponse() PermissionResponse {
	return PermissionResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
