package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Common errors for comment repository
var (
	ErrCommentNotFound = errors.New("comment not found")
)

// CommentRepository defines the interface for comment data access
type CommentRepository interface {
	// Create creates a new comment
	Create(ctx context.Context, comment *models.Comment) error
	
	// GetByID retrieves a comment by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.Comment, error)
	
	// List retrieves comments based on filter parameters
	List(ctx context.Context, params models.CommentListParams) ([]models.Comment, int, error)
	
	// Update updates a comment
	Update(ctx context.Context, comment *models.Comment) error
	
	// Delete deletes a comment
	Delete(ctx context.Context, id uuid.UUID) error
	
	// AddMention adds a mention to a comment
	AddMention(ctx context.Context, commentID, userID uuid.UUID) error
	
	// GetMentions retrieves all mentions for a comment
	GetMentions(ctx context.Context, commentID uuid.UUID) ([]uuid.UUID, error)
}

// PostgresCommentRepository implements CommentRepository using PostgreSQL
type PostgresCommentRepository struct {
	db *sqlx.DB
}

// NewPostgresCommentRepository creates a new PostgresCommentRepository
func NewPostgresCommentRepository(db *sqlx.DB) CommentRepository {
	return &PostgresCommentRepository{db: db}
}

// Create creates a new comment
func (r *PostgresCommentRepository) Create(ctx context.Context, comment *models.Comment) error {
	query := `
		INSERT INTO taskodex.task_comments (
			id, task_id, user_id, content, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	
	_, err := r.db.ExecContext(
		ctx,
		query,
		comment.ID,
		comment.TaskID,
		comment.UserID,
		comment.Content,
		comment.CreatedAt,
		comment.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to insert comment: %w", err)
	}
	
	// Add mentions if any
	if len(comment.Mentions) > 0 {
		for _, userID := range comment.Mentions {
			err = r.AddMention(ctx, comment.ID, userID)
			if err != nil {
				return fmt.Errorf("failed to add mention: %w", err)
			}
		}
	}
	
	return nil
}

// GetByID retrieves a comment by ID
func (r *PostgresCommentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Comment, error) {
	query := `
		SELECT c.id, c.task_id, c.user_id, c.content, c.created_at, c.updated_at
		FROM taskodex.task_comments c
		WHERE c.id = $1
	`
	
	var comment models.Comment
	err := r.db.GetContext(ctx, &comment, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCommentNotFound
		}
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	
	// Get user
	userRepo := NewPostgresUserRepository(r.db)
	user, err := userRepo.GetByID(ctx, comment.UserID)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user != nil {
		comment.User = user
	}
	
	// Get task
	taskRepo := NewPostgresTaskRepository(r.db)
	task, err := taskRepo.GetByID(ctx, comment.TaskID)
	if err != nil && !errors.Is(err, ErrTaskNotFound) {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if task != nil {
		comment.Task = task
	}
	
	// Get mentions
	mentions, err := r.GetMentions(ctx, comment.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get mentions: %w", err)
	}
	comment.Mentions = mentions
	
	return &comment, nil
}

// List retrieves comments based on filter parameters
func (r *PostgresCommentRepository) List(ctx context.Context, params models.CommentListParams) ([]models.Comment, int, error) {
	// Build the query
	baseQuery := `
		FROM taskodex.task_comments c
		WHERE 1=1
	`
	
	// Add filters
	filters := []string{}
	args := []interface{}{}
	argIndex := 1
	
	if params.TaskID != nil {
		filters = append(filters, fmt.Sprintf("c.task_id = $%d", argIndex))
		args = append(args, *params.TaskID)
		argIndex++
	}
	
	if params.UserID != nil {
		filters = append(filters, fmt.Sprintf("c.user_id = $%d", argIndex))
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
		return nil, 0, fmt.Errorf("failed to count comments: %w", err)
	}
	
	// Add sorting and pagination
	validSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
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
		SELECT c.id, c.task_id, c.user_id, c.content, c.created_at, c.updated_at
		%s
		ORDER BY c.%s %s
		LIMIT %d OFFSET %d
	`, baseQuery, sortBy, sortOrder, params.PageSize, offset)
	
	// Execute the query
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()
	
	// Process the results
	comments := []models.Comment{}
	for rows.Next() {
		var comment models.Comment
		err := rows.StructScan(&comment)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan comment: %w", err)
		}
		
		// Get user
		userRepo := NewPostgresUserRepository(r.db)
		user, err := userRepo.GetByID(ctx, comment.UserID)
		if err == nil {
			comment.User = user
		}
		
		// Get mentions
		mentions, err := r.GetMentions(ctx, comment.ID)
		if err == nil {
			comment.Mentions = mentions
		}
		
		comments = append(comments, comment)
	}
	
	return comments, total, nil
}

// Update updates a comment
func (r *PostgresCommentRepository) Update(ctx context.Context, comment *models.Comment) error {
	query := `
		UPDATE taskodex.task_comments
		SET content = $1, updated_at = $2
		WHERE id = $3
	`
	
	comment.UpdatedAt = time.Now()
	
	_, err := r.db.ExecContext(
		ctx,
		query,
		comment.Content,
		comment.UpdatedAt,
		comment.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}
	
	// Update mentions
	// First, delete all existing mentions
	deleteMentionsQuery := `
		DELETE FROM taskodex.comment_mentions
		WHERE comment_id = $1
	`
	
	_, err = r.db.ExecContext(ctx, deleteMentionsQuery, comment.ID)
	if err != nil {
		return fmt.Errorf("failed to delete existing mentions: %w", err)
	}
	
	// Then, add new mentions
	if len(comment.Mentions) > 0 {
		for _, userID := range comment.Mentions {
			err = r.AddMention(ctx, comment.ID, userID)
			if err != nil {
				return fmt.Errorf("failed to add mention: %w", err)
			}
		}
	}
	
	return nil
}

// Delete deletes a comment
func (r *PostgresCommentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// First, delete all mentions
	deleteMentionsQuery := `
		DELETE FROM taskodex.comment_mentions
		WHERE comment_id = $1
	`
	
	_, err := r.db.ExecContext(ctx, deleteMentionsQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete mentions: %w", err)
	}
	
	// Then, delete the comment
	query := `
		DELETE FROM taskodex.task_comments
		WHERE id = $1
	`
	
	_, err = r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	
	return nil
}

// AddMention adds a mention to a comment
func (r *PostgresCommentRepository) AddMention(ctx context.Context, commentID, userID uuid.UUID) error {
	query := `
		INSERT INTO taskodex.comment_mentions (
			comment_id, user_id, created_at
		)
		VALUES ($1, $2, $3)
	`
	
	_, err := r.db.ExecContext(
		ctx,
		query,
		commentID,
		userID,
		time.Now(),
	)
	
	if err != nil {
		return fmt.Errorf("failed to insert mention: %w", err)
	}
	
	return nil
}

// GetMentions retrieves all mentions for a comment
func (r *PostgresCommentRepository) GetMentions(ctx context.Context, commentID uuid.UUID) ([]uuid.UUID, error) {
	query := `
		SELECT user_id
		FROM taskodex.comment_mentions
		WHERE comment_id = $1
	`
	
	var mentions []uuid.UUID
	err := r.db.SelectContext(ctx, &mentions, query, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get mentions: %w", err)
	}
	
	return mentions, nil
}
