package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Common errors for file attachment repository
var (
	ErrFileAttachmentNotFound = errors.New("file attachment not found")
)

// FileAttachmentRepository defines the interface for file attachment data access
type FileAttachmentRepository interface {
	// Create creates a new file attachment
	Create(ctx context.Context, fileAttachment *models.FileAttachment) error
	
	// GetByID retrieves a file attachment by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.FileAttachment, error)
	
	// List retrieves file attachments based on filter parameters
	List(ctx context.Context, params models.FileAttachmentListParams) ([]models.FileAttachment, int, error)
	
	// Delete deletes a file attachment
	Delete(ctx context.Context, id uuid.UUID) error
}

// PostgresFileAttachmentRepository implements FileAttachmentRepository using PostgreSQL
type PostgresFileAttachmentRepository struct {
	db *sqlx.DB
}

// NewPostgresFileAttachmentRepository creates a new PostgresFileAttachmentRepository
func NewPostgresFileAttachmentRepository(db *sqlx.DB) FileAttachmentRepository {
	return &PostgresFileAttachmentRepository{db: db}
}

// Create creates a new file attachment
func (r *PostgresFileAttachmentRepository) Create(ctx context.Context, fileAttachment *models.FileAttachment) error {
	query := `
		INSERT INTO taskodex.task_file_attachments (
			id, task_id, user_id, file_name, file_size, content_type, storage_path, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err := r.db.ExecContext(
		ctx,
		query,
		fileAttachment.ID,
		fileAttachment.TaskID,
		fileAttachment.UserID,
		fileAttachment.FileName,
		fileAttachment.FileSize,
		fileAttachment.ContentType,
		fileAttachment.StoragePath,
		fileAttachment.CreatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to insert file attachment: %w", err)
	}
	
	return nil
}

// GetByID retrieves a file attachment by ID
func (r *PostgresFileAttachmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.FileAttachment, error) {
	query := `
		SELECT f.id, f.task_id, f.user_id, f.file_name, f.file_size, f.content_type, f.storage_path, f.created_at
		FROM taskodex.task_file_attachments f
		WHERE f.id = $1
	`
	
	var fileAttachment models.FileAttachment
	err := r.db.GetContext(ctx, &fileAttachment, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFileAttachmentNotFound
		}
		return nil, fmt.Errorf("failed to get file attachment: %w", err)
	}
	
	// Get user
	userRepo := NewPostgresUserRepository(r.db)
	user, err := userRepo.GetByID(ctx, fileAttachment.UserID)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user != nil {
		fileAttachment.User = user
	}
	
	// Get task
	taskRepo := NewPostgresTaskRepository(r.db)
	task, err := taskRepo.GetByID(ctx, fileAttachment.TaskID)
	if err != nil && !errors.Is(err, ErrTaskNotFound) {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if task != nil {
		fileAttachment.Task = task
	}
	
	return &fileAttachment, nil
}

// List retrieves file attachments based on filter parameters
func (r *PostgresFileAttachmentRepository) List(ctx context.Context, params models.FileAttachmentListParams) ([]models.FileAttachment, int, error) {
	// Build the query
	baseQuery := `
		FROM taskodex.task_file_attachments f
		WHERE 1=1
	`
	
	// Add filters
	filters := []string{}
	args := []interface{}{}
	argIndex := 1
	
	if params.TaskID != nil {
		filters = append(filters, fmt.Sprintf("f.task_id = $%d", argIndex))
		args = append(args, *params.TaskID)
		argIndex++
	}
	
	if params.UserID != nil {
		filters = append(filters, fmt.Sprintf("f.user_id = $%d", argIndex))
		args = append(args, *params.UserID)
		argIndex++
	}
	
	if len(filters) > 0 {
		baseQuery += " AND " + strings.Join(filters, " AND ")
	}
	
	// Count total records
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count file attachments: %w", err)
	}
	
	// Add sorting and pagination
	validSortFields := map[string]bool{
		"created_at": true,
		"file_name":  true,
		"file_size":  true,
	}
	
	sortBy := "created_at"
	if validSortFields[params.SortBy] {
		sortBy = params.SortBy
	}
	
	sortOrder := "DESC"
	if strings.ToLower(params.SortOrder) == "asc" {
		sortOrder = "ASC"
	}
	
	// Ensure page and page size are valid
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}
	
	offset := (params.Page - 1) * params.PageSize
	
	// Build the final query
	query := fmt.Sprintf(`
		SELECT f.id, f.task_id, f.user_id, f.file_name, f.file_size, f.content_type, f.storage_path, f.created_at
		%s
		ORDER BY f.%s %s
		LIMIT %d OFFSET %d
	`, baseQuery, sortBy, sortOrder, params.PageSize, offset)
	
	// Execute the query
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query file attachments: %w", err)
	}
	defer rows.Close()
	
	// Process the results
	fileAttachments := []models.FileAttachment{}
	for rows.Next() {
		var fileAttachment models.FileAttachment
		err := rows.StructScan(&fileAttachment)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan file attachment: %w", err)
		}
		
		// Get user
		userRepo := NewPostgresUserRepository(r.db)
		user, err := userRepo.GetByID(ctx, fileAttachment.UserID)
		if err == nil {
			fileAttachment.User = user
		}
		
		fileAttachments = append(fileAttachments, fileAttachment)
	}
	
	return fileAttachments, total, nil
}

// Delete deletes a file attachment
func (r *PostgresFileAttachmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM taskodex.task_file_attachments
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete file attachment: %w", err)
	}
	
	return nil
}
