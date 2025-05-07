package models

import (
	"time"

	"github.com/google/uuid"
)

// TimeEntry represents a time entry for a task
type TimeEntry struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	TaskID          uuid.UUID  `json:"task_id" db:"task_id"`
	UserID          uuid.UUID  `json:"user_id" db:"user_id"`
	StartTime       time.Time  `json:"start_time" db:"start_time"`
	EndTime         *time.Time `json:"end_time,omitempty" db:"end_time"`
	DurationMinutes *int       `json:"duration_minutes,omitempty" db:"duration_minutes"`
	Description     string     `json:"description" db:"description"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	
	// Related entities
	Task            *Task      `json:"task,omitempty" db:"-"`
	User            *User      `json:"user,omitempty" db:"-"`
}

// TimeEntryRequest represents the data needed to create or update a time entry
type TimeEntryRequest struct {
	TaskID          uuid.UUID  `json:"task_id" validate:"required,uuid4"`
	StartTime       time.Time  `json:"start_time" validate:"required"`
	EndTime         *time.Time `json:"end_time,omitempty"`
	DurationMinutes *int       `json:"duration_minutes,omitempty"`
	Description     string     `json:"description" validate:"max=1000"`
}

// TimeEntryResponse represents the time entry data returned to clients
type TimeEntryResponse struct {
	ID              uuid.UUID        `json:"id"`
	TaskID          uuid.UUID        `json:"task_id"`
	UserID          uuid.UUID        `json:"user_id"`
	StartTime       time.Time        `json:"start_time"`
	EndTime         *time.Time       `json:"end_time,omitempty"`
	DurationMinutes *int             `json:"duration_minutes,omitempty"`
	Description     string           `json:"description"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	
	// Related entities
	Task            *TaskResponse    `json:"task,omitempty"`
	User            *UserResponse    `json:"user,omitempty"`
}

// ToResponse converts a TimeEntry to a TimeEntryResponse
func (t *TimeEntry) ToResponse() TimeEntryResponse {
	response := TimeEntryResponse{
		ID:              t.ID,
		TaskID:          t.TaskID,
		UserID:          t.UserID,
		StartTime:       t.StartTime,
		EndTime:         t.EndTime,
		DurationMinutes: t.DurationMinutes,
		Description:     t.Description,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
	
	if t.Task != nil {
		taskResponse := t.Task.ToResponse()
		response.Task = &taskResponse
	}
	
	if t.User != nil {
		userResponse := t.User.ToResponse()
		response.User = &userResponse
	}
	
	return response
}

// NewTimeEntry creates a new TimeEntry from a TimeEntryRequest
func NewTimeEntry(req TimeEntryRequest, userID uuid.UUID) *TimeEntry {
	now := time.Now()
	
	// Calculate duration if end time is provided
	var durationMinutes *int
	if req.EndTime != nil && req.DurationMinutes == nil {
		duration := int(req.EndTime.Sub(req.StartTime).Minutes())
		durationMinutes = &duration
	} else {
		durationMinutes = req.DurationMinutes
	}
	
	return &TimeEntry{
		ID:              uuid.New(),
		TaskID:          req.TaskID,
		UserID:          userID,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		DurationMinutes: durationMinutes,
		Description:     req.Description,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// CalculateDuration calculates the duration of the time entry in minutes
func (t *TimeEntry) CalculateDuration() int {
	if t.DurationMinutes != nil {
		return *t.DurationMinutes
	}
	
	if t.EndTime != nil {
		return int(t.EndTime.Sub(t.StartTime).Minutes())
	}
	
	// If no end time, calculate duration from start time to now
	return int(time.Now().Sub(t.StartTime).Minutes())
}

// IsRunning returns true if the time entry is currently running (no end time)
func (t *TimeEntry) IsRunning() bool {
	return t.EndTime == nil
}

// Stop stops the time entry by setting the end time to the current time
func (t *TimeEntry) Stop() {
	now := time.Now()
	t.EndTime = &now
	t.UpdatedAt = now
	
	// Calculate duration
	duration := int(now.Sub(t.StartTime).Minutes())
	t.DurationMinutes = &duration
}

// TimeEntryListParams represents the parameters for listing time entries
type TimeEntryListParams struct {
	TaskID     *uuid.UUID `query:"task_id"`
	UserID     *uuid.UUID `query:"user_id"`
	StartAfter *time.Time `query:"start_after"`
	StartBefore *time.Time `query:"start_before"`
	IsRunning  *bool      `query:"is_running"`
	SortBy     string     `query:"sort_by" default:"start_time"`
	SortOrder  string     `query:"sort_order" default:"desc"`
	Page       int        `query:"page" default:"1"`
	PageSize   int        `query:"page_size" default:"20"`
}

// TimeEntryAggregation represents aggregated time entry data
type TimeEntryAggregation struct {
	TotalDurationMinutes int       `json:"total_duration_minutes"`
	TaskID               *uuid.UUID `json:"task_id,omitempty"`
	UserID               *uuid.UUID `json:"user_id,omitempty"`
	Date                 *time.Time `json:"date,omitempty"`
	Week                 *int       `json:"week,omitempty"`
	Month                *int       `json:"month,omitempty"`
	Year                 *int       `json:"year,omitempty"`
}

// TimeEntryAggregationParams represents the parameters for aggregating time entries
type TimeEntryAggregationParams struct {
	TaskID      *uuid.UUID `query:"task_id"`
	UserID      *uuid.UUID `query:"user_id"`
	ProjectID   *uuid.UUID `query:"project_id"`
	StartAfter  *time.Time `query:"start_after"`
	StartBefore *time.Time `query:"start_before"`
	GroupBy     string     `query:"group_by" default:"day"` // day, week, month, year, task, user
}
