package models

import (
	"time"

	"github.com/google/uuid"
)

// ProjectStatus represents the status of a project
type ProjectStatus string

// Project statuses
const (
	ProjectStatusPlanning  ProjectStatus = "planning"
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusOnHold    ProjectStatus = "on_hold"
	ProjectStatusCompleted ProjectStatus = "completed"
	ProjectStatusCancelled ProjectStatus = "cancelled"
)

// Project represents a project in the system
type Project struct {
	ID             uuid.UUID     `json:"id" db:"id"`
	OrganizationID uuid.UUID     `json:"organization_id" db:"organization_id"`
	Name           string        `json:"name" db:"name"`
	Description    string        `json:"description" db:"description"`
	Status         ProjectStatus `json:"status" db:"status"`
	StartDate      *time.Time    `json:"start_date,omitempty" db:"start_date"`
	EndDate        *time.Time    `json:"end_date,omitempty" db:"end_date"`
	CreatedBy      uuid.UUID     `json:"created_by" db:"created_by"`
	CreatedAt      time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" db:"updated_at"`
	
	// Related entities
	Creator        *User         `json:"creator,omitempty" db:"-"`
}

// ProjectRequest represents the data needed to create or update a project
type ProjectRequest struct {
	OrganizationID uuid.UUID     `json:"organization_id" validate:"required,uuid4"`
	Name           string        `json:"name" validate:"required,min=3,max=255"`
	Description    string        `json:"description" validate:"max=5000"`
	Status         ProjectStatus `json:"status" validate:"required,oneof=planning active on_hold completed cancelled"`
	StartDate      *time.Time    `json:"start_date,omitempty"`
	EndDate        *time.Time    `json:"end_date,omitempty"`
}

// ProjectResponse represents the project data returned to clients
type ProjectResponse struct {
	ID             uuid.UUID     `json:"id"`
	OrganizationID uuid.UUID     `json:"organization_id"`
	Name           string        `json:"name"`
	Description    string        `json:"description"`
	Status         ProjectStatus `json:"status"`
	StartDate      *time.Time    `json:"start_date,omitempty"`
	EndDate        *time.Time    `json:"end_date,omitempty"`
	CreatedBy      uuid.UUID     `json:"created_by"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	
	// Related entities
	Creator        *UserResponse `json:"creator,omitempty"`
}

// ToResponse converts a Project to a ProjectResponse
func (p *Project) ToResponse() ProjectResponse {
	response := ProjectResponse{
		ID:             p.ID,
		OrganizationID: p.OrganizationID,
		Name:           p.Name,
		Description:    p.Description,
		Status:         p.Status,
		StartDate:      p.StartDate,
		EndDate:        p.EndDate,
		CreatedBy:      p.CreatedBy,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
	
	if p.Creator != nil {
		creatorResponse := p.Creator.ToResponse()
		response.Creator = &creatorResponse
	}
	
	return response
}

// NewProject creates a new Project from a ProjectRequest
func NewProject(req ProjectRequest, createdBy uuid.UUID) *Project {
	now := time.Now()
	return &Project{
		ID:             uuid.New(),
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		Description:    req.Description,
		Status:         req.Status,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		CreatedBy:      createdBy,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// ProjectListParams represents the parameters for listing projects
type ProjectListParams struct {
	OrganizationID uuid.UUID      `query:"organization_id" validate:"required,uuid4"`
	Status         *ProjectStatus `query:"status"`
	CreatedBy      *uuid.UUID     `query:"created_by"`
	SearchTerm     *string        `query:"search"`
	SortBy         string         `query:"sort_by" default:"name"`
	SortOrder      string         `query:"sort_order" default:"asc"`
	Page           int            `query:"page" default:"1"`
	PageSize       int            `query:"page_size" default:"20"`
}
