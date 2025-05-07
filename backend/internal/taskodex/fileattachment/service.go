package fileattachment

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/notification"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/Jerinji2016/halooid/backend/internal/storage"
	"github.com/google/uuid"
)

// Service provides file attachment management functionality
type Service interface {
	// Upload uploads a file and creates a file attachment
	Upload(ctx context.Context, taskID, userID uuid.UUID, file *multipart.FileHeader) (*models.FileAttachmentResponse, error)
	
	// GetByID retrieves a file attachment by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.FileAttachmentResponse, error)
	
	// List retrieves file attachments based on filter parameters
	List(ctx context.Context, params models.FileAttachmentListParams) ([]models.FileAttachmentResponse, int, error)
	
	// Download downloads a file attachment
	Download(ctx context.Context, id uuid.UUID) (*models.FileAttachment, error)
	
	// Delete deletes a file attachment
	Delete(ctx context.Context, id uuid.UUID) error
}

// serviceImpl implements the Service interface
type serviceImpl struct {
	fileAttachmentRepo repository.FileAttachmentRepository
	taskRepo           repository.TaskRepository
	userRepo           repository.UserRepository
	fileStorage        storage.FileStorage
	notificationSvc    notification.Service
	baseURL            string
}

// NewService creates a new file attachment service
func NewService(
	fileAttachmentRepo repository.FileAttachmentRepository,
	taskRepo repository.TaskRepository,
	userRepo repository.UserRepository,
	fileStorage storage.FileStorage,
	notificationSvc notification.Service,
	baseURL string,
) Service {
	return &serviceImpl{
		fileAttachmentRepo: fileAttachmentRepo,
		taskRepo:           taskRepo,
		userRepo:           userRepo,
		fileStorage:        fileStorage,
		notificationSvc:    notificationSvc,
		baseURL:            baseURL,
	}
}

// Upload uploads a file and creates a file attachment
func (s *serviceImpl) Upload(ctx context.Context, taskID, userID uuid.UUID, file *multipart.FileHeader) (*models.FileAttachmentResponse, error) {
	// Validate task
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return nil, repository.ErrTaskNotFound
		}
		return nil, err
	}
	
	// Validate user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}
	
	// Save file to storage
	storagePath, err := s.fileStorage.SaveFile(file, taskID, userID)
	if err != nil {
		return nil, err
	}
	
	// Create file attachment record
	fileAttachment := models.NewFileAttachment(
		taskID,
		userID,
		file.Filename,
		file.Header.Get("Content-Type"),
		file.Size,
		storagePath,
	)
	
	err = s.fileAttachmentRepo.Create(ctx, fileAttachment)
	if err != nil {
		// If there's an error creating the record, try to clean up the file
		_ = s.fileStorage.DeleteFile(storagePath)
		return nil, err
	}
	
	// Get the complete file attachment record with related entities
	createdFileAttachment, err := s.fileAttachmentRepo.GetByID(ctx, fileAttachment.ID)
	if err != nil {
		return nil, err
	}
	
	// Send notification to task creator if they're not the uploader
	if task.CreatedBy != userID {
		notification := models.NewNotification(
			task.CreatedBy,
			models.NotificationTypeTaskComment,
			"New file uploaded to your task",
			user.FirstName+" "+user.LastName+" uploaded a file to your task: "+task.Title,
			"task",
			task.ID,
		)
		
		err = s.notificationSvc.Create(ctx, notification)
		if err != nil {
			// Log the error but don't fail the operation
			// TODO: Add proper logging
		}
	}
	
	// Send notification to task assignee if they're not the uploader and not the creator
	if task.AssignedTo != nil && *task.AssignedTo != userID && *task.AssignedTo != task.CreatedBy {
		notification := models.NewNotification(
			*task.AssignedTo,
			models.NotificationTypeTaskComment,
			"New file uploaded to your assigned task",
			user.FirstName+" "+user.LastName+" uploaded a file to a task assigned to you: "+task.Title,
			"task",
			task.ID,
		)
		
		err = s.notificationSvc.Create(ctx, notification)
		if err != nil {
			// Log the error but don't fail the operation
			// TODO: Add proper logging
		}
	}
	
	response := createdFileAttachment.ToResponse(s.baseURL)
	return &response, nil
}

// GetByID retrieves a file attachment by ID
func (s *serviceImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.FileAttachmentResponse, error) {
	fileAttachment, err := s.fileAttachmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	response := fileAttachment.ToResponse(s.baseURL)
	return &response, nil
}

// List retrieves file attachments based on filter parameters
func (s *serviceImpl) List(ctx context.Context, params models.FileAttachmentListParams) ([]models.FileAttachmentResponse, int, error) {
	fileAttachments, total, err := s.fileAttachmentRepo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	
	responses := make([]models.FileAttachmentResponse, 0, len(fileAttachments))
	for _, fileAttachment := range fileAttachments {
		responses = append(responses, fileAttachment.ToResponse(s.baseURL))
	}
	
	return responses, total, nil
}

// Download downloads a file attachment
func (s *serviceImpl) Download(ctx context.Context, id uuid.UUID) (*models.FileAttachment, error) {
	fileAttachment, err := s.fileAttachmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return fileAttachment, nil
}

// Delete deletes a file attachment
func (s *serviceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	// Get file attachment
	fileAttachment, err := s.fileAttachmentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	// Delete file from storage
	err = s.fileStorage.DeleteFile(fileAttachment.StoragePath)
	if err != nil {
		return err
	}
	
	// Delete file attachment record
	return s.fileAttachmentRepo.Delete(ctx, id)
}
