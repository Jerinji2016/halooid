package models

import (
	"time"

	"github.com/google/uuid"
)

// Employee represents an employee in the system
type Employee struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	OrganizationID uuid.UUID  `json:"organization_id" db:"organization_id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	EmployeeID     string     `json:"employee_id" db:"employee_id"`
	Department     string     `json:"department" db:"department"`
	Position       string     `json:"position" db:"position"`
	HireDate       time.Time  `json:"hire_date" db:"hire_date"`
	ManagerID      *uuid.UUID `json:"manager_id,omitempty" db:"manager_id"`
	Salary         float64    `json:"salary" db:"salary"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	User           *User      `json:"user,omitempty" db:"-"`
	Manager        *Employee  `json:"manager,omitempty" db:"-"`
}

// EmployeeRequest represents the data needed to create or update an employee
type EmployeeRequest struct {
	UserID     uuid.UUID  `json:"user_id" validate:"required"`
	EmployeeID string     `json:"employee_id" validate:"required"`
	Department string     `json:"department" validate:"required"`
	Position   string     `json:"position" validate:"required"`
	HireDate   time.Time  `json:"hire_date" validate:"required"`
	ManagerID  *uuid.UUID `json:"manager_id,omitempty"`
	Salary     float64    `json:"salary" validate:"required,min=0"`
}

// EmployeeResponse represents the employee data returned to clients
type EmployeeResponse struct {
	ID             uuid.UUID        `json:"id"`
	OrganizationID uuid.UUID        `json:"organization_id"`
	EmployeeID     string           `json:"employee_id"`
	Department     string           `json:"department"`
	Position       string           `json:"position"`
	HireDate       time.Time        `json:"hire_date"`
	ManagerID      *uuid.UUID       `json:"manager_id,omitempty"`
	Salary         float64          `json:"salary"`
	IsActive       bool             `json:"is_active"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	User           *UserResponse    `json:"user,omitempty"`
	Manager        *EmployeeResponse `json:"manager,omitempty"`
}

// ToResponse converts an Employee to an EmployeeResponse
func (e *Employee) ToResponse() EmployeeResponse {
	response := EmployeeResponse{
		ID:             e.ID,
		OrganizationID: e.OrganizationID,
		EmployeeID:     e.EmployeeID,
		Department:     e.Department,
		Position:       e.Position,
		HireDate:       e.HireDate,
		ManagerID:      e.ManagerID,
		Salary:         e.Salary,
		IsActive:       e.IsActive,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}

	if e.User != nil {
		userResponse := e.User.ToResponse()
		response.User = &userResponse
	}

	if e.Manager != nil {
		managerResponse := e.Manager.ToResponse()
		response.Manager = &managerResponse
	}

	return response
}

// NewEmployee creates a new Employee from an EmployeeRequest
func NewEmployee(req EmployeeRequest, organizationID uuid.UUID) *Employee {
	now := time.Now()
	return &Employee{
		ID:             uuid.New(),
		OrganizationID: organizationID,
		UserID:         req.UserID,
		EmployeeID:     req.EmployeeID,
		Department:     req.Department,
		Position:       req.Position,
		HireDate:       req.HireDate,
		ManagerID:      req.ManagerID,
		Salary:         req.Salary,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// EmployeeListParams represents the parameters for listing employees
type EmployeeListParams struct {
	OrganizationID uuid.UUID `query:"organization_id" validate:"required"`
	Department     string    `query:"department"`
	Position       string    `query:"position"`
	IsActive       *bool     `query:"is_active"`
	SortBy         string    `query:"sort_by" default:"employee_id"`
	SortOrder      string    `query:"sort_order" default:"asc"`
	Page           int       `query:"page" default:"1"`
	PageSize       int       `query:"page_size" default:"10"`
}
