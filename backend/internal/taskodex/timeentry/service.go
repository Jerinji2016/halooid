package timeentry

import (
	"context"
	"errors"
	"time"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/google/uuid"
)

// Service provides time entry management functionality
type Service interface {
	// Create creates a new time entry
	Create(ctx context.Context, req models.TimeEntryRequest, userID uuid.UUID) (*models.TimeEntryResponse, error)
	
	// GetByID retrieves a time entry by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.TimeEntryResponse, error)
	
	// List retrieves time entries based on filter parameters
	List(ctx context.Context, params models.TimeEntryListParams) ([]models.TimeEntryResponse, int, error)
	
	// Update updates a time entry
	Update(ctx context.Context, id uuid.UUID, req models.TimeEntryRequest) (*models.TimeEntryResponse, error)
	
	// Delete deletes a time entry
	Delete(ctx context.Context, id uuid.UUID) error
	
	// StartTimer starts a timer for a task
	StartTimer(ctx context.Context, taskID uuid.UUID, description string, userID uuid.UUID) (*models.TimeEntryResponse, error)
	
	// StopTimer stops a running timer for a task
	StopTimer(ctx context.Context, taskID uuid.UUID, userID uuid.UUID) (*models.TimeEntryResponse, error)
	
	// GetRunningTimers retrieves all running timers for a user
	GetRunningTimers(ctx context.Context, userID uuid.UUID) ([]models.TimeEntryResponse, error)
	
	// Aggregate aggregates time entries based on parameters
	Aggregate(ctx context.Context, params models.TimeEntryAggregationParams) ([]models.TimeEntryAggregation, error)
}

// serviceImpl implements the Service interface
type serviceImpl struct {
	timeEntryRepo repository.TimeEntryRepository
	taskRepo      repository.TaskRepository
	userRepo      repository.UserRepository
}

// NewService creates a new time entry service
func NewService(timeEntryRepo repository.TimeEntryRepository, taskRepo repository.TaskRepository, userRepo repository.UserRepository) Service {
	return &serviceImpl{
		timeEntryRepo: timeEntryRepo,
		taskRepo:      taskRepo,
		userRepo:      userRepo,
	}
}

// Create creates a new time entry
func (s *serviceImpl) Create(ctx context.Context, req models.TimeEntryRequest, userID uuid.UUID) (*models.TimeEntryResponse, error) {
	// Validate task
	_, err := s.taskRepo.GetByID(ctx, req.TaskID)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return nil, repository.ErrTaskNotFound
		}
		return nil, err
	}
	
	// Validate user
	_, err = s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}
	
	// Validate time entry
	if req.EndTime != nil && req.EndTime.Before(req.StartTime) {
		return nil, errors.New("end time cannot be before start time")
	}
	
	// Create time entry
	timeEntry := models.NewTimeEntry(req, userID)
	
	err = s.timeEntryRepo.Create(ctx, timeEntry)
	if err != nil {
		if errors.Is(err, repository.ErrRunningTimeEntry) {
			return nil, repository.ErrRunningTimeEntry
		}
		return nil, err
	}
	
	// Get the complete time entry record with related entities
	createdTimeEntry, err := s.timeEntryRepo.GetByID(ctx, timeEntry.ID)
	if err != nil {
		return nil, err
	}
	
	response := createdTimeEntry.ToResponse()
	return &response, nil
}

// GetByID retrieves a time entry by ID
func (s *serviceImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.TimeEntryResponse, error) {
	timeEntry, err := s.timeEntryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	response := timeEntry.ToResponse()
	return &response, nil
}

// List retrieves time entries based on filter parameters
func (s *serviceImpl) List(ctx context.Context, params models.TimeEntryListParams) ([]models.TimeEntryResponse, int, error) {
	timeEntries, total, err := s.timeEntryRepo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	
	responses := make([]models.TimeEntryResponse, 0, len(timeEntries))
	for _, timeEntry := range timeEntries {
		responses = append(responses, timeEntry.ToResponse())
	}
	
	return responses, total, nil
}

// Update updates a time entry
func (s *serviceImpl) Update(ctx context.Context, id uuid.UUID, req models.TimeEntryRequest) (*models.TimeEntryResponse, error) {
	// Check if time entry exists
	timeEntry, err := s.timeEntryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Validate task
	_, err = s.taskRepo.GetByID(ctx, req.TaskID)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return nil, repository.ErrTaskNotFound
		}
		return nil, err
	}
	
	// Validate time entry
	if req.EndTime != nil && req.EndTime.Before(req.StartTime) {
		return nil, errors.New("end time cannot be before start time")
	}
	
	// Update time entry fields
	timeEntry.TaskID = req.TaskID
	timeEntry.StartTime = req.StartTime
	timeEntry.EndTime = req.EndTime
	timeEntry.Description = req.Description
	
	// Calculate duration if end time is provided
	if req.EndTime != nil && req.DurationMinutes == nil {
		duration := int(req.EndTime.Sub(req.StartTime).Minutes())
		timeEntry.DurationMinutes = &duration
	} else {
		timeEntry.DurationMinutes = req.DurationMinutes
	}
	
	err = s.timeEntryRepo.Update(ctx, timeEntry)
	if err != nil {
		return nil, err
	}
	
	// Get the updated time entry record
	updatedTimeEntry, err := s.timeEntryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	response := updatedTimeEntry.ToResponse()
	return &response, nil
}

// Delete deletes a time entry
func (s *serviceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if time entry exists
	_, err := s.timeEntryRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	return s.timeEntryRepo.Delete(ctx, id)
}

// StartTimer starts a timer for a task
func (s *serviceImpl) StartTimer(ctx context.Context, taskID uuid.UUID, description string, userID uuid.UUID) (*models.TimeEntryResponse, error) {
	// Validate task
	_, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return nil, repository.ErrTaskNotFound
		}
		return nil, err
	}
	
	// Validate user
	_, err = s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}
	
	// Check if user already has a running timer for this task
	existingTimer, err := s.timeEntryRepo.GetRunningTimeEntry(ctx, userID, taskID)
	if err != nil && !errors.Is(err, repository.ErrTimeEntryNotFound) {
		return nil, err
	}
	if existingTimer != nil {
		return nil, repository.ErrRunningTimeEntry
	}
	
	// Create time entry
	now := time.Now()
	timeEntry := &models.TimeEntry{
		ID:          uuid.New(),
		TaskID:      taskID,
		UserID:      userID,
		StartTime:   now,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	err = s.timeEntryRepo.Create(ctx, timeEntry)
	if err != nil {
		return nil, err
	}
	
	// Get the complete time entry record with related entities
	createdTimeEntry, err := s.timeEntryRepo.GetByID(ctx, timeEntry.ID)
	if err != nil {
		return nil, err
	}
	
	response := createdTimeEntry.ToResponse()
	return &response, nil
}

// StopTimer stops a running timer for a task
func (s *serviceImpl) StopTimer(ctx context.Context, taskID uuid.UUID, userID uuid.UUID) (*models.TimeEntryResponse, error) {
	// Get running timer
	timeEntry, err := s.timeEntryRepo.GetRunningTimeEntry(ctx, userID, taskID)
	if err != nil {
		if errors.Is(err, repository.ErrTimeEntryNotFound) {
			return nil, errors.New("no running timer found for this task")
		}
		return nil, err
	}
	
	// Stop timer
	timeEntry.Stop()
	
	err = s.timeEntryRepo.Update(ctx, timeEntry)
	if err != nil {
		return nil, err
	}
	
	// Get the updated time entry record
	updatedTimeEntry, err := s.timeEntryRepo.GetByID(ctx, timeEntry.ID)
	if err != nil {
		return nil, err
	}
	
	response := updatedTimeEntry.ToResponse()
	return &response, nil
}

// GetRunningTimers retrieves all running timers for a user
func (s *serviceImpl) GetRunningTimers(ctx context.Context, userID uuid.UUID) ([]models.TimeEntryResponse, error) {
	// Validate user
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}
	
	timeEntries, err := s.timeEntryRepo.GetRunningTimeEntries(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	responses := make([]models.TimeEntryResponse, 0, len(timeEntries))
	for _, timeEntry := range timeEntries {
		// Get task
		task, err := s.taskRepo.GetByID(ctx, timeEntry.TaskID)
		if err != nil && !errors.Is(err, repository.ErrTaskNotFound) {
			return nil, err
		}
		if task != nil {
			timeEntry.Task = task
		}
		
		responses = append(responses, timeEntry.ToResponse())
	}
	
	return responses, nil
}

// Aggregate aggregates time entries based on parameters
func (s *serviceImpl) Aggregate(ctx context.Context, params models.TimeEntryAggregationParams) ([]models.TimeEntryAggregation, error) {
	return s.timeEntryRepo.Aggregate(ctx, params)
}
