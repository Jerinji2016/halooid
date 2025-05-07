package task

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/notification"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/google/uuid"
)

// Service provides task management functionality
type Service interface {
	// Create creates a new task
	Create(ctx context.Context, req models.TaskRequest, createdBy uuid.UUID) (*models.TaskResponse, error)

	// GetByID retrieves a task by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.TaskResponse, error)

	// List retrieves tasks based on filter parameters
	List(ctx context.Context, params models.TaskListParams) ([]models.TaskResponse, int, error)

	// Update updates a task
	Update(ctx context.Context, id uuid.UUID, req models.TaskRequest) (*models.TaskResponse, error)

	// Delete deletes a task
	Delete(ctx context.Context, id uuid.UUID) error

	// AddTag adds a tag to a task
	AddTag(ctx context.Context, taskID uuid.UUID, tag string) error

	// RemoveTag removes a tag from a task
	RemoveTag(ctx context.Context, taskID uuid.UUID, tag string) error

	// AssignTask assigns a task to a user
	AssignTask(ctx context.Context, taskID uuid.UUID, userID uuid.UUID, assignedBy uuid.UUID) (*models.TaskResponse, error)

	// UnassignTask removes the assignment of a task
	UnassignTask(ctx context.Context, taskID uuid.UUID, unassignedBy uuid.UUID) (*models.TaskResponse, error)

	// UpdateTaskStatus updates the status of a task
	UpdateTaskStatus(ctx context.Context, taskID uuid.UUID, status models.TaskStatus, updatedBy uuid.UUID) (*models.TaskResponse, error)

	// GetTasksByAssignee retrieves tasks assigned to a user
	GetTasksByAssignee(ctx context.Context, userID uuid.UUID, params models.TaskListParams) ([]models.TaskResponse, int, error)

	// GetOverdueTasks retrieves overdue tasks
	GetOverdueTasks(ctx context.Context, params models.TaskListParams) ([]models.TaskResponse, int, error)

	// GetTasksDueSoon retrieves tasks due within a specified number of days
	GetTasksDueSoon(ctx context.Context, days int, params models.TaskListParams) ([]models.TaskResponse, int, error)
}

// serviceImpl implements the Service interface
type serviceImpl struct {
	taskRepo        repository.TaskRepository
	projectRepo     repository.ProjectRepository
	userRepo        repository.UserRepository
	notificationSvc notification.Service
}

// NewService creates a new task service
func NewService(
	taskRepo repository.TaskRepository,
	projectRepo repository.ProjectRepository,
	userRepo repository.UserRepository,
	notificationSvc notification.Service,
) Service {
	return &serviceImpl{
		taskRepo:        taskRepo,
		projectRepo:     projectRepo,
		userRepo:        userRepo,
		notificationSvc: notificationSvc,
	}
}

// Create creates a new task
func (s *serviceImpl) Create(ctx context.Context, req models.TaskRequest, createdBy uuid.UUID) (*models.TaskResponse, error) {
	// Validate project if provided
	if req.ProjectID != nil {
		_, err := s.projectRepo.GetByID(ctx, *req.ProjectID)
		if err != nil {
			if errors.Is(err, repository.ErrProjectNotFound) {
				return nil, repository.ErrProjectNotFound
			}
			return nil, err
		}
	}

	// Validate assignee if provided
	if req.AssignedTo != nil {
		_, err := s.userRepo.GetByID(ctx, *req.AssignedTo)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				return nil, repository.ErrUserNotFound
			}
			return nil, err
		}
	}

	// Create task
	task := models.NewTask(req, createdBy)

	err := s.taskRepo.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	// Get the complete task record with related entities
	createdTask, err := s.taskRepo.GetByID(ctx, task.ID)
	if err != nil {
		return nil, err
	}

	response := createdTask.ToResponse()
	return &response, nil
}

// GetByID retrieves a task by ID
func (s *serviceImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.TaskResponse, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := task.ToResponse()
	return &response, nil
}

// List retrieves tasks based on filter parameters
func (s *serviceImpl) List(ctx context.Context, params models.TaskListParams) ([]models.TaskResponse, int, error) {
	tasks, total, err := s.taskRepo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		responses = append(responses, task.ToResponse())
	}

	return responses, total, nil
}

// Update updates a task
func (s *serviceImpl) Update(ctx context.Context, id uuid.UUID, req models.TaskRequest) (*models.TaskResponse, error) {
	// Check if task exists
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate project if provided
	if req.ProjectID != nil {
		_, err := s.projectRepo.GetByID(ctx, *req.ProjectID)
		if err != nil {
			if errors.Is(err, repository.ErrProjectNotFound) {
				return nil, repository.ErrProjectNotFound
			}
			return nil, err
		}
	}

	// Validate assignee if provided
	if req.AssignedTo != nil {
		_, err := s.userRepo.GetByID(ctx, *req.AssignedTo)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				return nil, repository.ErrUserNotFound
			}
			return nil, err
		}
	}

	// Update task fields
	task.ProjectID = req.ProjectID
	task.Title = req.Title
	task.Description = req.Description
	task.Status = req.Status
	task.Priority = req.Priority
	task.DueDate = req.DueDate
	task.AssignedTo = req.AssignedTo
	task.EstimatedHours = req.EstimatedHours
	task.Tags = req.Tags

	err = s.taskRepo.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	// Get the updated task record
	updatedTask, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := updatedTask.ToResponse()
	return &response, nil
}

// Delete deletes a task
func (s *serviceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if task exists
	_, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.taskRepo.Delete(ctx, id)
}

// AddTag adds a tag to a task
func (s *serviceImpl) AddTag(ctx context.Context, taskID uuid.UUID, tag string) error {
	// Check if task exists
	_, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}

	return s.taskRepo.AddTag(ctx, taskID, tag)
}

// RemoveTag removes a tag from a task
func (s *serviceImpl) RemoveTag(ctx context.Context, taskID uuid.UUID, tag string) error {
	// Check if task exists
	_, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}

	return s.taskRepo.RemoveTag(ctx, taskID, tag)
}

// AssignTask assigns a task to a user
func (s *serviceImpl) AssignTask(ctx context.Context, taskID uuid.UUID, userID uuid.UUID, assignedBy uuid.UUID) (*models.TaskResponse, error) {
	// Check if task exists
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Check if user exists
	_, err = s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	// Update task's assignee
	task.AssignedTo = &userID
	task.UpdatedAt = time.Now()

	err = s.taskRepo.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	// Get the updated task record
	updatedTask, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Send notification to the assignee
	err = s.notificationSvc.NotifyTaskAssigned(ctx, updatedTask)
	if err != nil {
		// Log the error but don't fail the operation
		// TODO: Add proper logging
		fmt.Println("Failed to send task assignment notification:", err)
	}

	response := updatedTask.ToResponse()
	return &response, nil
}

// UnassignTask removes the assignment of a task
func (s *serviceImpl) UnassignTask(ctx context.Context, taskID uuid.UUID, unassignedBy uuid.UUID) (*models.TaskResponse, error) {
	// Check if task exists
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Check if task is assigned
	if task.AssignedTo == nil {
		return nil, errors.New("task is not assigned to anyone")
	}

	// Update task's assignee
	task.AssignedTo = nil
	task.UpdatedAt = time.Now()

	err = s.taskRepo.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	// Get the updated task record
	updatedTask, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	response := updatedTask.ToResponse()
	return &response, nil
}

// UpdateTaskStatus updates the status of a task
func (s *serviceImpl) UpdateTaskStatus(ctx context.Context, taskID uuid.UUID, status models.TaskStatus, updatedBy uuid.UUID) (*models.TaskResponse, error) {
	// Check if task exists
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Check if status is valid
	validStatuses := map[models.TaskStatus]bool{
		models.TaskStatusTodo:       true,
		models.TaskStatusInProgress: true,
		models.TaskStatusReview:     true,
		models.TaskStatusDone:       true,
		models.TaskStatusCancelled:  true,
	}

	if !validStatuses[status] {
		return nil, errors.New("invalid task status")
	}

	// Save old status for notification
	oldStatus := task.Status

	// Update task's status
	task.Status = status
	task.UpdatedAt = time.Now()

	err = s.taskRepo.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	// Get the updated task record
	updatedTask, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Send notification about status update
	err = s.notificationSvc.NotifyTaskStatusUpdate(ctx, updatedTask, oldStatus)
	if err != nil {
		// Log the error but don't fail the operation
		// TODO: Add proper logging
		fmt.Println("Failed to send task status update notification:", err)
	}

	response := updatedTask.ToResponse()
	return &response, nil
}

// GetTasksByAssignee retrieves tasks assigned to a user
func (s *serviceImpl) GetTasksByAssignee(ctx context.Context, userID uuid.UUID, params models.TaskListParams) ([]models.TaskResponse, int, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, 0, repository.ErrUserNotFound
		}
		return nil, 0, err
	}

	// Set assignee in params
	params.AssignedTo = &userID

	// Get tasks
	tasks, total, err := s.taskRepo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		responses = append(responses, task.ToResponse())
	}

	return responses, total, nil
}

// GetOverdueTasks retrieves overdue tasks
func (s *serviceImpl) GetOverdueTasks(ctx context.Context, params models.TaskListParams) ([]models.TaskResponse, int, error) {
	// Set due date filter to find overdue tasks
	now := time.Now()
	params.DueBefore = &now

	// Exclude completed and cancelled tasks
	todoStatus := models.TaskStatusTodo
	inProgressStatus := models.TaskStatusInProgress
	reviewStatus := models.TaskStatusReview

	// We need to get tasks with these statuses separately and combine them
	// since we can only set one status filter at a time

	// Get todo tasks
	params.Status = &todoStatus
	todoTasks, _, err := s.taskRepo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Get in progress tasks
	params.Status = &inProgressStatus
	inProgressTasks, _, err := s.taskRepo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Get review tasks
	params.Status = &reviewStatus
	reviewTasks, _, err := s.taskRepo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Combine tasks
	allTasks := append(todoTasks, inProgressTasks...)
	allTasks = append(allTasks, reviewTasks...)

	// Filter out tasks without due dates
	overdueTasks := make([]models.Task, 0)
	for _, task := range allTasks {
		if task.DueDate != nil {
			overdueTasks = append(overdueTasks, task)
		}
	}

	// Convert to responses
	responses := make([]models.TaskResponse, 0, len(overdueTasks))
	for _, task := range overdueTasks {
		responses = append(responses, task.ToResponse())
	}

	return responses, len(responses), nil
}

// GetTasksDueSoon retrieves tasks due within a specified number of days
func (s *serviceImpl) GetTasksDueSoon(ctx context.Context, days int, params models.TaskListParams) ([]models.TaskResponse, int, error) {
	// Set due date filter to find tasks due soon
	now := time.Now()
	dueAfter := now
	dueBefore := now.AddDate(0, 0, days)

	params.DueAfter = &dueAfter
	params.DueBefore = &dueBefore

	// Exclude completed and cancelled tasks
	todoStatus := models.TaskStatusTodo
	inProgressStatus := models.TaskStatusInProgress
	reviewStatus := models.TaskStatusReview

	// We need to get tasks with these statuses separately and combine them
	// since we can only set one status filter at a time

	// Get todo tasks
	params.Status = &todoStatus
	todoTasks, _, err := s.taskRepo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Get in progress tasks
	params.Status = &inProgressStatus
	inProgressTasks, _, err := s.taskRepo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Get review tasks
	params.Status = &reviewStatus
	reviewTasks, _, err := s.taskRepo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Combine tasks
	allTasks := append(todoTasks, inProgressTasks...)
	allTasks = append(allTasks, reviewTasks...)

	// Filter out tasks without due dates
	dueSoonTasks := make([]models.Task, 0)
	for _, task := range allTasks {
		if task.DueDate != nil {
			dueSoonTasks = append(dueSoonTasks, task)
		}
	}

	// Convert to responses
	responses := make([]models.TaskResponse, 0, len(dueSoonTasks))
	for _, task := range dueSoonTasks {
		responses = append(responses, task.ToResponse())
	}

	return responses, len(responses), nil
}
