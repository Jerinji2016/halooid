package models

import (
	"time"

	"github.com/google/uuid"
)

// Comment represents a comment on a task
type Comment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	TaskID    uuid.UUID `json:"task_id" db:"task_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	
	// Related entities
	User      *User     `json:"user,omitempty" db:"-"`
	Task      *Task     `json:"task,omitempty" db:"-"`
	Mentions  []uuid.UUID `json:"mentions,omitempty" db:"-"`
}

// CommentRequest represents the data needed to create or update a comment
type CommentRequest struct {
	TaskID  uuid.UUID `json:"task_id" validate:"required,uuid4"`
	Content string    `json:"content" validate:"required,min=1,max=5000"`
}

// CommentResponse represents the comment data returned to clients
type CommentResponse struct {
	ID        uuid.UUID      `json:"id"`
	TaskID    uuid.UUID      `json:"task_id"`
	UserID    uuid.UUID      `json:"user_id"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	
	// Related entities
	User      *UserResponse  `json:"user,omitempty"`
	Task      *TaskResponse  `json:"task,omitempty"`
	Mentions  []uuid.UUID    `json:"mentions,omitempty"`
}

// ToResponse converts a Comment to a CommentResponse
func (c *Comment) ToResponse() CommentResponse {
	response := CommentResponse{
		ID:        c.ID,
		TaskID:    c.TaskID,
		UserID:    c.UserID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Mentions:  c.Mentions,
	}
	
	if c.User != nil {
		userResponse := c.User.ToResponse()
		response.User = &userResponse
	}
	
	if c.Task != nil {
		taskResponse := c.Task.ToResponse()
		response.Task = &taskResponse
	}
	
	return response
}

// NewComment creates a new Comment from a CommentRequest
func NewComment(req CommentRequest, userID uuid.UUID) *Comment {
	now := time.Now()
	
	return &Comment{
		ID:        uuid.New(),
		TaskID:    req.TaskID,
		UserID:    userID,
		Content:   req.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ExtractMentions extracts @mentions from the comment content
// Returns a list of user IDs that were mentioned
func (c *Comment) ExtractMentions(content string, userIDMap map[string]uuid.UUID) []uuid.UUID {
	// This is a simplified implementation
	// In a real-world scenario, you would use a more sophisticated approach
	// such as regular expressions to extract mentions
	
	// For now, we'll just return an empty list
	return []uuid.UUID{}
}

// CommentListParams represents the parameters for listing comments
type CommentListParams struct {
	TaskID    *uuid.UUID `query:"task_id"`
	UserID    *uuid.UUID `query:"user_id"`
	SortBy    string     `query:"sort_by" default:"created_at"`
	SortOrder string     `query:"sort_order" default:"desc"`
	Page      int        `query:"page" default:"1"`
	PageSize  int        `query:"page_size" default:"20"`
}
