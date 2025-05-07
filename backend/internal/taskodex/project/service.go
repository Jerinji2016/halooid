package project

import (
	"context"
	"errors"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/google/uuid"
)

// Service provides project management functionality
type Service interface {
	// Create creates a new project
	Create(ctx context.Context, req models.ProjectRequest, createdBy uuid.UUID) (*models.ProjectResponse, error)
	
	// GetByID retrieves a project by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.ProjectResponse, error)
	
	// GetByName retrieves a project by name within an organization
	GetByName(ctx context.Context, organizationID uuid.UUID, name string) (*models.ProjectResponse, error)
	
	// List retrieves projects based on filter parameters
	List(ctx context.Context, params models.ProjectListParams) ([]models.ProjectResponse, int, error)
	
	// Update updates a project
	Update(ctx context.Context, id uuid.UUID, req models.ProjectRequest) (*models.ProjectResponse, error)
	
	// Delete deletes a project
	Delete(ctx context.Context, id uuid.UUID) error
	
	// GetTasks retrieves all tasks for a project
	GetTasks(ctx context.Context, projectID uuid.UUID, params models.TaskListParams) ([]models.TaskResponse, int, error)
	
	// AddTask adds a task to a project
	AddTask(ctx context.Context, projectID uuid.UUID, taskID uuid.UUID) error
	
	// RemoveTask removes a task from a project
	RemoveTask(ctx context.Context, projectID uuid.UUID, taskID uuid.UUID) error
}

// serviceImpl implements the Service interface
type serviceImpl struct {
	projectRepo repository.ProjectRepository
	taskRepo    repository.TaskRepository
	userRepo    repository.UserRepository
}

// NewService creates a new project service
func NewService(projectRepo repository.ProjectRepository, taskRepo repository.TaskRepository, userRepo repository.UserRepository) Service {
	return &serviceImpl{
		projectRepo: projectRepo,
		taskRepo:    taskRepo,
		userRepo:    userRepo,
	}
}

// Create creates a new project
func (s *serviceImpl) Create(ctx context.Context, req models.ProjectRequest, createdBy uuid.UUID) (*models.ProjectResponse, error) {
	// Create project
	project := models.NewProject(req, createdBy)
	
	err := s.projectRepo.Create(ctx, project)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNameExists) {
			return nil, repository.ErrProjectNameExists
		}
		return nil, err
	}
	
	// Get the complete project record with related entities
	createdProject, err := s.projectRepo.GetByID(ctx, project.ID)
	if err != nil {
		return nil, err
	}
	
	response := createdProject.ToResponse()
	return &response, nil
}

// GetByID retrieves a project by ID
func (s *serviceImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.ProjectResponse, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	response := project.ToResponse()
	return &response, nil
}

// GetByName retrieves a project by name within an organization
func (s *serviceImpl) GetByName(ctx context.Context, organizationID uuid.UUID, name string) (*models.ProjectResponse, error) {
	project, err := s.projectRepo.GetByName(ctx, organizationID, name)
	if err != nil {
		return nil, err
	}
	
	response := project.ToResponse()
	return &response, nil
}

// List retrieves projects based on filter parameters
func (s *serviceImpl) List(ctx context.Context, params models.ProjectListParams) ([]models.ProjectResponse, int, error) {
	projects, total, err := s.projectRepo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	
	responses := make([]models.ProjectResponse, 0, len(projects))
	for _, project := range projects {
		responses = append(responses, project.ToResponse())
	}
	
	return responses, total, nil
}

// Update updates a project
func (s *serviceImpl) Update(ctx context.Context, id uuid.UUID, req models.ProjectRequest) (*models.ProjectResponse, error) {
	// Check if project exists
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Update project fields
	project.OrganizationID = req.OrganizationID
	project.Name = req.Name
	project.Description = req.Description
	project.Status = req.Status
	project.StartDate = req.StartDate
	project.EndDate = req.EndDate
	
	err = s.projectRepo.Update(ctx, project)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNameExists) {
			return nil, repository.ErrProjectNameExists
		}
		return nil, err
	}
	
	// Get the updated project record
	updatedProject, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	response := updatedProject.ToResponse()
	return &response, nil
}

// Delete deletes a project
func (s *serviceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if project exists
	_, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	return s.projectRepo.Delete(ctx, id)
}

// GetTasks retrieves all tasks for a project
func (s *serviceImpl) GetTasks(ctx context.Context, projectID uuid.UUID, params models.TaskListParams) ([]models.TaskResponse, int, error) {
	// Check if project exists
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, 0, err
	}
	
	// Set project ID in params
	params.ProjectID = &projectID
	
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

// AddTask adds a task to a project
func (s *serviceImpl) AddTask(ctx context.Context, projectID uuid.UUID, taskID uuid.UUID) error {
	// Check if project exists
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}
	
	// Check if task exists
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}
	
	// Update task's project ID
	task.ProjectID = &projectID
	
	return s.taskRepo.Update(ctx, task)
}

// RemoveTask removes a task from a project
func (s *serviceImpl) RemoveTask(ctx context.Context, projectID uuid.UUID, taskID uuid.UUID) error {
	// Check if project exists
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}
	
	// Check if task exists
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}
	
	// Check if task belongs to the project
	if task.ProjectID == nil || *task.ProjectID != projectID {
		return errors.New("task does not belong to the project")
	}
	
	// Remove task's project ID
	task.ProjectID = nil
	
	return s.taskRepo.Update(ctx, task)
}
