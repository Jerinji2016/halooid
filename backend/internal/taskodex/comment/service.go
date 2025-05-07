package comment

import (
	"context"
	"errors"
	"regexp"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/notification"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/google/uuid"
)

// Service provides comment management functionality
type Service interface {
	// Create creates a new comment
	Create(ctx context.Context, req models.CommentRequest, userID uuid.UUID) (*models.CommentResponse, error)
	
	// GetByID retrieves a comment by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.CommentResponse, error)
	
	// List retrieves comments based on filter parameters
	List(ctx context.Context, params models.CommentListParams) ([]models.CommentResponse, int, error)
	
	// Update updates a comment
	Update(ctx context.Context, id uuid.UUID, req models.CommentRequest) (*models.CommentResponse, error)
	
	// Delete deletes a comment
	Delete(ctx context.Context, id uuid.UUID) error
}

// serviceImpl implements the Service interface
type serviceImpl struct {
	commentRepo     repository.CommentRepository
	taskRepo        repository.TaskRepository
	userRepo        repository.UserRepository
	notificationSvc notification.Service
}

// NewService creates a new comment service
func NewService(
	commentRepo repository.CommentRepository,
	taskRepo repository.TaskRepository,
	userRepo repository.UserRepository,
	notificationSvc notification.Service,
) Service {
	return &serviceImpl{
		commentRepo:     commentRepo,
		taskRepo:        taskRepo,
		userRepo:        userRepo,
		notificationSvc: notificationSvc,
	}
}

// Create creates a new comment
func (s *serviceImpl) Create(ctx context.Context, req models.CommentRequest, userID uuid.UUID) (*models.CommentResponse, error) {
	// Validate task
	task, err := s.taskRepo.GetByID(ctx, req.TaskID)
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
	
	// Create comment
	comment := models.NewComment(req, userID)
	
	// Extract mentions
	mentions := s.extractMentions(ctx, req.Content)
	comment.Mentions = mentions
	
	err = s.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}
	
	// Get the complete comment record with related entities
	createdComment, err := s.commentRepo.GetByID(ctx, comment.ID)
	if err != nil {
		return nil, err
	}
	
	// Send notifications for mentions
	for _, mentionedUserID := range mentions {
		// Create a notification for the mentioned user
		notification := models.NewNotification(
			mentionedUserID,
			models.NotificationTypeTaskMention,
			"You were mentioned in a comment",
			user.FirstName+" "+user.LastName+" mentioned you in a comment on task: "+task.Title,
			"task",
			task.ID,
		)
		
		err = s.notificationSvc.Create(ctx, notification)
		if err != nil {
			// Log the error but don't fail the operation
			// TODO: Add proper logging
		}
	}
	
	// Send notification to task creator if they're not the commenter
	if task.CreatedBy != userID {
		notification := models.NewNotification(
			task.CreatedBy,
			models.NotificationTypeTaskComment,
			"New comment on your task",
			user.FirstName+" "+user.LastName+" commented on your task: "+task.Title,
			"task",
			task.ID,
		)
		
		err = s.notificationSvc.Create(ctx, notification)
		if err != nil {
			// Log the error but don't fail the operation
			// TODO: Add proper logging
		}
	}
	
	// Send notification to task assignee if they're not the commenter and not the creator
	if task.AssignedTo != nil && *task.AssignedTo != userID && *task.AssignedTo != task.CreatedBy {
		notification := models.NewNotification(
			*task.AssignedTo,
			models.NotificationTypeTaskComment,
			"New comment on your assigned task",
			user.FirstName+" "+user.LastName+" commented on a task assigned to you: "+task.Title,
			"task",
			task.ID,
		)
		
		err = s.notificationSvc.Create(ctx, notification)
		if err != nil {
			// Log the error but don't fail the operation
			// TODO: Add proper logging
		}
	}
	
	response := createdComment.ToResponse()
	return &response, nil
}

// GetByID retrieves a comment by ID
func (s *serviceImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.CommentResponse, error) {
	comment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	response := comment.ToResponse()
	return &response, nil
}

// List retrieves comments based on filter parameters
func (s *serviceImpl) List(ctx context.Context, params models.CommentListParams) ([]models.CommentResponse, int, error) {
	comments, total, err := s.commentRepo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	
	responses := make([]models.CommentResponse, 0, len(comments))
	for _, comment := range comments {
		responses = append(responses, comment.ToResponse())
	}
	
	return responses, total, nil
}

// Update updates a comment
func (s *serviceImpl) Update(ctx context.Context, id uuid.UUID, req models.CommentRequest) (*models.CommentResponse, error) {
	// Check if comment exists
	comment, err := s.commentRepo.GetByID(ctx, id)
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
	
	// Update comment fields
	comment.Content = req.Content
	
	// Extract mentions
	mentions := s.extractMentions(ctx, req.Content)
	comment.Mentions = mentions
	
	err = s.commentRepo.Update(ctx, comment)
	if err != nil {
		return nil, err
	}
	
	// Get the updated comment record
	updatedComment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	response := updatedComment.ToResponse()
	return &response, nil
}

// Delete deletes a comment
func (s *serviceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if comment exists
	_, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	return s.commentRepo.Delete(ctx, id)
}

// extractMentions extracts @mentions from the comment content
func (s *serviceImpl) extractMentions(ctx context.Context, content string) []uuid.UUID {
	// Use a regular expression to find @username mentions
	re := regexp.MustCompile(`@([a-zA-Z0-9_]+)`)
	matches := re.FindAllStringSubmatch(content, -1)
	
	// Map to store unique user IDs
	mentionedUserIDs := make(map[uuid.UUID]bool)
	
	for _, match := range matches {
		if len(match) > 1 {
			username := match[1]
			
			// Look up the user by username
			user, err := s.userRepo.GetByUsername(ctx, username)
			if err == nil && user != nil {
				mentionedUserIDs[user.ID] = true
			}
		}
	}
	
	// Convert map keys to slice
	result := make([]uuid.UUID, 0, len(mentionedUserIDs))
	for userID := range mentionedUserIDs {
		result = append(result, userID)
	}
	
	return result
}
