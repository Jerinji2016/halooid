package models

import (
	"time"

	"github.com/google/uuid"
)

// TaskStatus represents the status of a task
type TaskStatus string

// Task statuses
const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusReview     TaskStatus = "review"
	TaskStatusDone       TaskStatus = "done"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

// TaskPriority represents the priority of a task
type TaskPriority string

// Task priorities
const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
	TaskPriorityCritical TaskPriority = "critical"
)

// Task represents a task in the system
type Task struct {
	ID             uuid.UUID    `json:"id" db:"id"`
	ProjectID      *uuid.UUID   `json:"project_id,omitempty" db:"project_id"`
	Title          string       `json:"title" db:"title"`
	Description    string       `json:"description" db:"description"`
	Status         TaskStatus   `json:"status" db:"status"`
	Priority       TaskPriority `json:"priority" db:"priority"`
	DueDate        *time.Time   `json:"due_date,omitempty" db:"due_date"`
	CreatedBy      uuid.UUID    `json:"created_by" db:"created_by"`
	AssignedTo     *uuid.UUID   `json:"assigned_to,omitempty" db:"assigned_to"`
	EstimatedHours *float64     `json:"estimated_hours,omitempty" db:"estimated_hours"`
	ActualHours    *float64     `json:"actual_hours,omitempty" db:"actual_hours"`
	Tags           []string     `json:"tags" db:"-"`
	CreatedAt      time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at" db:"updated_at"`
	
	// Related entities
	Project        *Project     `json:"project,omitempty" db:"-"`
	Creator        *User        `json:"creator,omitempty" db:"-"`
	Assignee       *User        `json:"assignee,omitempty" db:"-"`
}

// TaskRequest represents the data needed to create or update a task
type TaskRequest struct {
	ProjectID      *uuid.UUID   `json:"project_id,omitempty" validate:"omitempty,uuid4"`
	Title          string       `json:"title" validate:"required,min=3,max=255"`
	Description    string       `json:"description" validate:"max=5000"`
	Status         TaskStatus   `json:"status" validate:"required,oneof=todo in_progress review done cancelled"`
	Priority       TaskPriority `json:"priority" validate:"required,oneof=low medium high critical"`
	DueDate        *time.Time   `json:"due_date,omitempty"`
	AssignedTo     *uuid.UUID   `json:"assigned_to,omitempty" validate:"omitempty,uuid4"`
	EstimatedHours *float64     `json:"estimated_hours,omitempty" validate:"omitempty,min=0"`
	Tags           []string     `json:"tags,omitempty" validate:"dive,max=50"`
}

// TaskResponse represents the task data returned to clients
type TaskResponse struct {
	ID             uuid.UUID    `json:"id"`
	ProjectID      *uuid.UUID   `json:"project_id,omitempty"`
	Title          string       `json:"title"`
	Description    string       `json:"description"`
	Status         TaskStatus   `json:"status"`
	Priority       TaskPriority `json:"priority"`
	DueDate        *time.Time   `json:"due_date,omitempty"`
	CreatedBy      uuid.UUID    `json:"created_by"`
	AssignedTo     *uuid.UUID   `json:"assigned_to,omitempty"`
	EstimatedHours *float64     `json:"estimated_hours,omitempty"`
	ActualHours    *float64     `json:"actual_hours,omitempty"`
	Tags           []string     `json:"tags"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	
	// Related entities
	Project        *ProjectResponse `json:"project,omitempty"`
	Creator        *UserResponse    `json:"creator,omitempty"`
	Assignee       *UserResponse    `json:"assignee,omitempty"`
}

// ToResponse converts a Task to a TaskResponse
func (t *Task) ToResponse() TaskResponse {
	response := TaskResponse{
		ID:             t.ID,
		ProjectID:      t.ProjectID,
		Title:          t.Title,
		Description:    t.Description,
		Status:         t.Status,
		Priority:       t.Priority,
		DueDate:        t.DueDate,
		CreatedBy:      t.CreatedBy,
		AssignedTo:     t.AssignedTo,
		EstimatedHours: t.EstimatedHours,
		ActualHours:    t.ActualHours,
		Tags:           t.Tags,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}
	
	if t.Project != nil {
		projectResponse := t.Project.ToResponse()
		response.Project = &projectResponse
	}
	
	if t.Creator != nil {
		creatorResponse := t.Creator.ToResponse()
		response.Creator = &creatorResponse
	}
	
	if t.Assignee != nil {
		assigneeResponse := t.Assignee.ToResponse()
		response.Assignee = &assigneeResponse
	}
	
	return response
}

// NewTask creates a new Task from a TaskRequest
func NewTask(req TaskRequest, createdBy uuid.UUID) *Task {
	now := time.Now()
	return &Task{
		ID:             uuid.New(),
		ProjectID:      req.ProjectID,
		Title:          req.Title,
		Description:    req.Description,
		Status:         req.Status,
		Priority:       req.Priority,
		DueDate:        req.DueDate,
		CreatedBy:      createdBy,
		AssignedTo:     req.AssignedTo,
		EstimatedHours: req.EstimatedHours,
		Tags:           req.Tags,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// TaskListParams represents the parameters for listing tasks
type TaskListParams struct {
	ProjectID  *uuid.UUID   `query:"project_id"`
	Status     *TaskStatus  `query:"status"`
	Priority   *TaskPriority `query:"priority"`
	CreatedBy  *uuid.UUID   `query:"created_by"`
	AssignedTo *uuid.UUID   `query:"assigned_to"`
	DueBefore  *time.Time   `query:"due_before"`
	DueAfter   *time.Time   `query:"due_after"`
	SearchTerm *string      `query:"search"`
	SortBy     string       `query:"sort_by" default:"created_at"`
	SortOrder  string       `query:"sort_order" default:"desc"`
	Page       int          `query:"page" default:"1"`
	PageSize   int          `query:"page_size" default:"20"`
}
