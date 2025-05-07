package task

import (
	"context"
	"errors"

	"github.com/Jerinji2016/halooid/backend/internal/models"
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
}

// serviceImpl implements the Service interface
type serviceImpl struct {
	taskRepo    repository.TaskRepository
	projectRepo repository.ProjectRepository
	userRepo    repository.UserRepository
}

// NewService creates a new task service
func NewService(taskRepo repository.TaskRepository, projectRepo repository.ProjectRepository, userRepo repository.UserRepository) Service {
	return &serviceImpl{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
		userRepo:    userRepo,
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
