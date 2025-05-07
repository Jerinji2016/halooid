package notification

import (
	"context"
	"fmt"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/google/uuid"
)

// Service provides notification management functionality
type Service interface {
	// Create creates a new notification
	Create(ctx context.Context, notification *models.Notification) error
	
	// GetByID retrieves a notification by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.NotificationResponse, error)
	
	// List retrieves notifications based on filter parameters
	List(ctx context.Context, params models.NotificationListParams) ([]models.NotificationResponse, int, error)
	
	// MarkAsRead marks a notification as read
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	
	// MarkAllAsRead marks all notifications for a user as read
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	
	// Delete deletes a notification
	Delete(ctx context.Context, id uuid.UUID) error
	
	// DeleteAllForUser deletes all notifications for a user
	DeleteAllForUser(ctx context.Context, userID uuid.UUID) error
	
	// NotifyTaskAssigned notifies a user that a task has been assigned to them
	NotifyTaskAssigned(ctx context.Context, task *models.Task) error
	
	// NotifyTaskStatusUpdate notifies relevant users about a task status update
	NotifyTaskStatusUpdate(ctx context.Context, task *models.Task, oldStatus models.TaskStatus) error
	
	// NotifyTaskDueSoon notifies assigned users about tasks that are due soon
	NotifyTaskDueSoon(ctx context.Context, task *models.Task, daysUntilDue int) error
	
	// NotifyTaskOverdue notifies assigned users about overdue tasks
	NotifyTaskOverdue(ctx context.Context, task *models.Task) error
}

// serviceImpl implements the Service interface
type serviceImpl struct {
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
}

// NewService creates a new notification service
func NewService(notificationRepo repository.NotificationRepository, userRepo repository.UserRepository) Service {
	return &serviceImpl{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
	}
}

// Create creates a new notification
func (s *serviceImpl) Create(ctx context.Context, notification *models.Notification) error {
	return s.notificationRepo.Create(ctx, notification)
}

// GetByID retrieves a notification by ID
func (s *serviceImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.NotificationResponse, error) {
	notification, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	response := notification.ToResponse()
	return &response, nil
}

// List retrieves notifications based on filter parameters
func (s *serviceImpl) List(ctx context.Context, params models.NotificationListParams) ([]models.NotificationResponse, int, error) {
	notifications, total, err := s.notificationRepo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	
	responses := make([]models.NotificationResponse, 0, len(notifications))
	for _, notification := range notifications {
		responses = append(responses, notification.ToResponse())
	}
	
	return responses, total, nil
}

// MarkAsRead marks a notification as read
func (s *serviceImpl) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	return s.notificationRepo.MarkAsRead(ctx, id)
}

// MarkAllAsRead marks all notifications for a user as read
func (s *serviceImpl) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return s.notificationRepo.MarkAllAsRead(ctx, userID)
}

// Delete deletes a notification
func (s *serviceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return s.notificationRepo.Delete(ctx, id)
}

// DeleteAllForUser deletes all notifications for a user
func (s *serviceImpl) DeleteAllForUser(ctx context.Context, userID uuid.UUID) error {
	return s.notificationRepo.DeleteAllForUser(ctx, userID)
}

// NotifyTaskAssigned notifies a user that a task has been assigned to them
func (s *serviceImpl) NotifyTaskAssigned(ctx context.Context, task *models.Task) error {
	if task.AssignedTo == nil {
		return nil
	}
	
	// Get assignee
	assignee, err := s.userRepo.GetByID(ctx, *task.AssignedTo)
	if err != nil {
		return err
	}
	
	// Get creator
	creator, err := s.userRepo.GetByID(ctx, task.CreatedBy)
	if err != nil {
		return err
	}
	
	// Create notification
	title := "Task Assigned"
	message := fmt.Sprintf("%s %s assigned you a task: %s", creator.FirstName, creator.LastName, task.Title)
	
	notification := models.NewNotification(
		*task.AssignedTo,
		models.NotificationTypeTaskAssigned,
		title,
		message,
		"task",
		task.ID,
	)
	
	return s.notificationRepo.Create(ctx, notification)
}

// NotifyTaskStatusUpdate notifies relevant users about a task status update
func (s *serviceImpl) NotifyTaskStatusUpdate(ctx context.Context, task *models.Task, oldStatus models.TaskStatus) error {
	// Determine who should be notified
	var userIDsToNotify []uuid.UUID
	
	// Always notify the creator if they're not the one who updated the task
	userIDsToNotify = append(userIDsToNotify, task.CreatedBy)
	
	// Notify the assignee if there is one and they're not the creator
	if task.AssignedTo != nil && *task.AssignedTo != task.CreatedBy {
		userIDsToNotify = append(userIDsToNotify, *task.AssignedTo)
	}
	
	// Create notifications
	title := "Task Status Updated"
	message := fmt.Sprintf("Task '%s' status changed from %s to %s", task.Title, oldStatus, task.Status)
	
	for _, userID := range userIDsToNotify {
		notification := models.NewNotification(
			userID,
			models.NotificationTypeTaskStatusUpdate,
			title,
			message,
			"task",
			task.ID,
		)
		
		err := s.notificationRepo.Create(ctx, notification)
		if err != nil {
			return err
		}
	}
	
	return nil
}

// NotifyTaskDueSoon notifies assigned users about tasks that are due soon
func (s *serviceImpl) NotifyTaskDueSoon(ctx context.Context, task *models.Task, daysUntilDue int) error {
	if task.AssignedTo == nil || task.DueDate == nil {
		return nil
	}
	
	// Create notification
	title := "Task Due Soon"
	message := fmt.Sprintf("Task '%s' is due in %d days", task.Title, daysUntilDue)
	
	notification := models.NewNotification(
		*task.AssignedTo,
		models.NotificationTypeTaskDueSoon,
		title,
		message,
		"task",
		task.ID,
	)
	
	return s.notificationRepo.Create(ctx, notification)
}

// NotifyTaskOverdue notifies assigned users about overdue tasks
func (s *serviceImpl) NotifyTaskOverdue(ctx context.Context, task *models.Task) error {
	if task.AssignedTo == nil || task.DueDate == nil {
		return nil
	}
	
	// Create notification
	title := "Task Overdue"
	message := fmt.Sprintf("Task '%s' is overdue", task.Title)
	
	notification := models.NewNotification(
		*task.AssignedTo,
		models.NotificationTypeTaskOverdue,
		title,
		message,
		"task",
		task.ID,
	)
	
	return s.notificationRepo.Create(ctx, notification)
}
