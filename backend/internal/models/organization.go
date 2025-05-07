package models

import (
	"time"

	"github.com/google/uuid"
)

// Organization represents an organization in the system
type Organization struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Users       []User    `json:"users,omitempty" db:"-"`
}

// OrganizationUser represents a many-to-many relationship between organizations and users
type OrganizationUser struct {
	OrganizationID uuid.UUID `json:"organization_id" db:"organization_id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
}

// OrganizationRequest represents the data needed to create or update an organization
type OrganizationRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// OrganizationResponse represents the organization data returned to clients
type OrganizationResponse struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	IsActive    bool            `json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Users       []UserResponse  `json:"users,omitempty"`
}

// ToResponse converts an Organization to an OrganizationResponse
func (o *Organization) ToResponse() OrganizationResponse {
	userResponses := make([]UserResponse, 0, len(o.Users))
	for _, u := range o.Users {
		userResponses = append(userResponses, u.ToResponse())
	}

	return OrganizationResponse{
		ID:          o.ID,
		Name:        o.Name,
		Description: o.Description,
		IsActive:    o.IsActive,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
		Users:       userResponses,
	}
}

// NewOrganization creates a new Organization from an OrganizationRequest
func NewOrganization(req OrganizationRequest) *Organization {
	now := time.Now()
	return &Organization{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// AddUserToOrganizationRequest represents the data needed to add a user to an organization
type AddUserToOrganizationRequest struct {
	UserID         uuid.UUID `json:"user_id" validate:"required"`
	OrganizationID uuid.UUID `json:"organization_id" validate:"required"`
}

// OrganizationUserResponse represents the organization-user relationship data returned to clients
type OrganizationUserResponse struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	UserID         uuid.UUID `json:"user_id"`
}
