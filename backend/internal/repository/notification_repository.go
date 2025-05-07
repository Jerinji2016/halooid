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

// Common errors for notification repository
var (
	ErrNotificationNotFound = errors.New("notification not found")
)

// NotificationRepository defines the interface for notification data access
type NotificationRepository interface {
	// Create creates a new notification
	Create(ctx context.Context, notification *models.Notification) error
	
	// GetByID retrieves a notification by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.Notification, error)
	
	// List retrieves notifications based on filter parameters
	List(ctx context.Context, params models.NotificationListParams) ([]models.Notification, int, error)
	
	// MarkAsRead marks a notification as read
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	
	// MarkAllAsRead marks all notifications for a user as read
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	
	// Delete deletes a notification
	Delete(ctx context.Context, id uuid.UUID) error
	
	// DeleteAllForUser deletes all notifications for a user
	DeleteAllForUser(ctx context.Context, userID uuid.UUID) error
}

// PostgresNotificationRepository implements NotificationRepository using PostgreSQL
type PostgresNotificationRepository struct {
	db *sqlx.DB
}

// NewPostgresNotificationRepository creates a new PostgresNotificationRepository
func NewPostgresNotificationRepository(db *sqlx.DB) NotificationRepository {
	return &PostgresNotificationRepository{db: db}
}

// Create creates a new notification
func (r *PostgresNotificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	query := `
		INSERT INTO notifications (
			id, user_id, type, title, message, resource_type, 
			resource_id, is_read, created_at, read_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	
	_, err := r.db.ExecContext(
		ctx,
		query,
		notification.ID,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.ResourceType,
		notification.ResourceID,
		notification.IsRead,
		notification.CreatedAt,
		notification.ReadAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to insert notification: %w", err)
	}
	
	return nil
}

// GetByID retrieves a notification by ID
func (r *PostgresNotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	query := `
		SELECT n.id, n.user_id, n.type, n.title, n.message, n.resource_type, 
			n.resource_id, n.is_read, n.created_at, n.read_at
		FROM notifications n
		WHERE n.id = $1
	`
	
	var notification models.Notification
	err := r.db.GetContext(ctx, &notification, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotificationNotFound
		}
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}
	
	// Get user
	userRepo := NewPostgresUserRepository(r.db)
	user, err := userRepo.GetByID(ctx, notification.UserID)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user != nil {
		notification.User = user
	}
	
	return &notification, nil
}

// List retrieves notifications based on filter parameters
func (r *PostgresNotificationRepository) List(ctx context.Context, params models.NotificationListParams) ([]models.Notification, int, error) {
	// Build the query
	baseQuery := `
		FROM notifications n
		WHERE n.user_id = $1
	`
	
	// Add filters
	filters := []string{}
	args := []interface{}{params.UserID}
	argIndex := 2
	
	if params.Type != nil {
		filters = append(filters, fmt.Sprintf("n.type = $%d", argIndex))
		args = append(args, *params.Type)
		argIndex++
	}
	
	if params.IsRead != nil {
		filters = append(filters, fmt.Sprintf("n.is_read = $%d", argIndex))
		args = append(args, *params.IsRead)
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
		return nil, 0, fmt.Errorf("failed to count notifications: %w", err)
	}
	
	// Add sorting and pagination
	validSortFields := map[string]bool{
		"created_at": true,
		"type":       true,
		"is_read":    true,
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
		SELECT n.id, n.user_id, n.type, n.title, n.message, n.resource_type, 
			n.resource_id, n.is_read, n.created_at, n.read_at
		%s
		ORDER BY n.%s %s
		LIMIT %d OFFSET %d
	`, baseQuery, sortBy, sortOrder, params.PageSize, offset)
	
	// Execute the query
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query notifications: %w", err)
	}
	defer rows.Close()
	
	// Process the results
	notifications := []models.Notification{}
	for rows.Next() {
		var notification models.Notification
		err := rows.StructScan(&notification)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}
	
	return notifications, total, nil
}

// MarkAsRead marks a notification as read
func (r *PostgresNotificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE notifications
		SET is_read = true, read_at = $1
		WHERE id = $2
	`
	
	now := time.Now()
	
	_, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	
	return nil
}

// MarkAllAsRead marks all notifications for a user as read
func (r *PostgresNotificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE notifications
		SET is_read = true, read_at = $1
		WHERE user_id = $2 AND is_read = false
	`
	
	now := time.Now()
	
	_, err := r.db.ExecContext(ctx, query, now, userID)
	if err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}
	
	return nil
}

// Delete deletes a notification
func (r *PostgresNotificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM notifications
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	
	return nil
}

// DeleteAllForUser deletes all notifications for a user
func (r *PostgresNotificationRepository) DeleteAllForUser(ctx context.Context, userID uuid.UUID) error {
	query := `
		DELETE FROM notifications
		WHERE user_id = $1
	`
	
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete all notifications for user: %w", err)
	}
	
	return nil
}
