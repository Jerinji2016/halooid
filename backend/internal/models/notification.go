package models

import (
	"time"

	"github.com/google/uuid"
)

// NotificationType represents the type of notification
type NotificationType string

// Notification types
const (
	NotificationTypeTaskAssigned     NotificationType = "task_assigned"
	NotificationTypeTaskStatusUpdate NotificationType = "task_status_update"
	NotificationTypeTaskDueSoon      NotificationType = "task_due_soon"
	NotificationTypeTaskOverdue      NotificationType = "task_overdue"
	NotificationTypeTaskComment      NotificationType = "task_comment"
	NotificationTypeTaskMention      NotificationType = "task_mention"
)

// Notification represents a notification in the system
type Notification struct {
	ID           uuid.UUID        `json:"id" db:"id"`
	UserID       uuid.UUID        `json:"user_id" db:"user_id"`
	Type         NotificationType `json:"type" db:"type"`
	Title        string           `json:"title" db:"title"`
	Message      string           `json:"message" db:"message"`
	ResourceType string           `json:"resource_type" db:"resource_type"`
	ResourceID   uuid.UUID        `json:"resource_id" db:"resource_id"`
	IsRead       bool             `json:"is_read" db:"is_read"`
	CreatedAt    time.Time        `json:"created_at" db:"created_at"`
	ReadAt       *time.Time       `json:"read_at,omitempty" db:"read_at"`
	
	// Related entities
	User         *User            `json:"user,omitempty" db:"-"`
}

// NotificationResponse represents the notification data returned to clients
type NotificationResponse struct {
	ID           uuid.UUID        `json:"id"`
	UserID       uuid.UUID        `json:"user_id"`
	Type         NotificationType `json:"type"`
	Title        string           `json:"title"`
	Message      string           `json:"message"`
	ResourceType string           `json:"resource_type"`
	ResourceID   uuid.UUID        `json:"resource_id"`
	IsRead       bool             `json:"is_read"`
	CreatedAt    time.Time        `json:"created_at"`
	ReadAt       *time.Time       `json:"read_at,omitempty"`
	
	// Related entities
	User         *UserResponse    `json:"user,omitempty"`
}

// ToResponse converts a Notification to a NotificationResponse
func (n *Notification) ToResponse() NotificationResponse {
	response := NotificationResponse{
		ID:           n.ID,
		UserID:       n.UserID,
		Type:         n.Type,
		Title:        n.Title,
		Message:      n.Message,
		ResourceType: n.ResourceType,
		ResourceID:   n.ResourceID,
		IsRead:       n.IsRead,
		CreatedAt:    n.CreatedAt,
		ReadAt:       n.ReadAt,
	}
	
	if n.User != nil {
		userResponse := n.User.ToResponse()
		response.User = &userResponse
	}
	
	return response
}

// NewNotification creates a new Notification
func NewNotification(userID uuid.UUID, notificationType NotificationType, title, message, resourceType string, resourceID uuid.UUID) *Notification {
	return &Notification{
		ID:           uuid.New(),
		UserID:       userID,
		Type:         notificationType,
		Title:        title,
		Message:      message,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		IsRead:       false,
		CreatedAt:    time.Now(),
	}
}

// MarkAsRead marks the notification as read
func (n *Notification) MarkAsRead() {
	n.IsRead = true
	now := time.Now()
	n.ReadAt = &now
}

// NotificationListParams represents the parameters for listing notifications
type NotificationListParams struct {
	UserID     uuid.UUID         `query:"user_id" validate:"required,uuid4"`
	Type       *NotificationType `query:"type"`
	IsRead     *bool             `query:"is_read"`
	SortBy     string            `query:"sort_by" default:"created_at"`
	SortOrder  string            `query:"sort_order" default:"desc"`
	Page       int               `query:"page" default:"1"`
	PageSize   int               `query:"page_size" default:"20"`
}
