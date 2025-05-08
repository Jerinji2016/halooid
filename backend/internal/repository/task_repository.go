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

// Common errors for task repository
var (
	ErrTaskNotFound = errors.New("task not found")
)

// TaskRepository defines the interface for task data access
type TaskRepository interface {
	// Create creates a new task
	Create(ctx context.Context, task *models.Task) error

	// GetByID retrieves a task by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error)

	// List retrieves tasks based on filter parameters
	List(ctx context.Context, params models.TaskListParams) ([]models.Task, int, error)

	// Update updates a task
	Update(ctx context.Context, task *models.Task) error

	// Delete deletes a task
	Delete(ctx context.Context, id uuid.UUID) error

	// AddTag adds a tag to a task
	AddTag(ctx context.Context, taskID uuid.UUID, tag string) error

	// RemoveTag removes a tag from a task
	RemoveTag(ctx context.Context, taskID uuid.UUID, tag string) error

	// GetTags retrieves all tags for a task
	GetTags(ctx context.Context, taskID uuid.UUID) ([]string, error)
}

// PostgresTaskRepository implements TaskRepository using PostgreSQL
type PostgresTaskRepository struct {
	db *sqlx.DB
}

// NewPostgresTaskRepository creates a new PostgresTaskRepository
func NewPostgresTaskRepository(db *sqlx.DB) TaskRepository {
	return &PostgresTaskRepository{db: db}
}

// Create creates a new task
func (r *PostgresTaskRepository) Create(ctx context.Context, task *models.Task) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert task
	query := `
		INSERT INTO taskodex.tasks (
			id, project_id, title, description, status, priority,
			due_date, created_by, assigned_to, estimated_hours,
			actual_hours, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err = tx.ExecContext(
		ctx,
		query,
		task.ID,
		task.ProjectID,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.DueDate,
		task.CreatedBy,
		task.AssignedTo,
		task.EstimatedHours,
		task.ActualHours,
		task.CreatedAt,
		task.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	// Insert tags if any
	if len(task.Tags) > 0 {
		for _, tag := range task.Tags {
			err = r.addTagTx(ctx, tx, task.ID, tag)
			if err != nil {
				return fmt.Errorf("failed to add tag: %w", err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetByID retrieves a task by ID
func (r *PostgresTaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	query := `
		SELECT t.id, t.project_id, t.title, t.description, t.status, t.priority,
			t.due_date, t.created_by, t.assigned_to, t.estimated_hours,
			t.actual_hours, t.created_at, t.updated_at
		FROM taskodex.tasks t
		WHERE t.id = $1
	`

	var task models.Task
	err := r.db.GetContext(ctx, &task, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	// Get tags
	tags, err := r.GetTags(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	task.Tags = tags

	// Get related entities
	if task.ProjectID != nil {
		projectRepo := NewPostgresProjectRepository(r.db)
		project, err := projectRepo.GetByID(ctx, *task.ProjectID)
		if err != nil && !errors.Is(err, ErrProjectNotFound) {
			return nil, fmt.Errorf("failed to get project: %w", err)
		}
		if project != nil {
			task.Project = project
		}
	}

	userRepo := NewPostgresUserRepository(r.db)

	creator, err := userRepo.GetByID(ctx, task.CreatedBy)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, fmt.Errorf("failed to get creator: %w", err)
	}
	if creator != nil {
		task.Creator = creator
	}

	if task.AssignedTo != nil {
		assignee, err := userRepo.GetByID(ctx, *task.AssignedTo)
		if err != nil && !errors.Is(err, ErrUserNotFound) {
			return nil, fmt.Errorf("failed to get assignee: %w", err)
		}
		if assignee != nil {
			task.Assignee = assignee
		}
	}

	return &task, nil
}

// List retrieves tasks based on filter parameters
func (r *PostgresTaskRepository) List(ctx context.Context, params models.TaskListParams) ([]models.Task, int, error) {
	// Build the query
	baseQuery := `
		FROM taskodex.tasks t
		WHERE 1=1
	`

	// Add filters
	filters := []string{}
	args := []interface{}{}
	argIndex := 1

	if params.ProjectID != nil {
		filters = append(filters, fmt.Sprintf("t.project_id = $%d", argIndex))
		args = append(args, *params.ProjectID)
		argIndex++
	}

	if params.Status != nil {
		filters = append(filters, fmt.Sprintf("t.status = $%d", argIndex))
		args = append(args, *params.Status)
		argIndex++
	}

	if params.Priority != nil {
		filters = append(filters, fmt.Sprintf("t.priority = $%d", argIndex))
		args = append(args, *params.Priority)
		argIndex++
	}

	if params.CreatedBy != nil {
		filters = append(filters, fmt.Sprintf("t.created_by = $%d", argIndex))
		args = append(args, *params.CreatedBy)
		argIndex++
	}

	if params.AssignedTo != nil {
		filters = append(filters, fmt.Sprintf("t.assigned_to = $%d", argIndex))
		args = append(args, *params.AssignedTo)
		argIndex++
	}

	if params.DueBefore != nil {
		filters = append(filters, fmt.Sprintf("t.due_date <= $%d", argIndex))
		args = append(args, *params.DueBefore)
		argIndex++
	}

	if params.DueAfter != nil {
		filters = append(filters, fmt.Sprintf("t.due_date >= $%d", argIndex))
		args = append(args, *params.DueAfter)
		argIndex++
	}

	if params.SearchTerm != nil && *params.SearchTerm != "" {
		filters = append(filters, fmt.Sprintf("(t.title ILIKE $%d OR t.description ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+*params.SearchTerm+"%")
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
		return nil, 0, fmt.Errorf("failed to count tasks: %w", err)
	}

	// Add sorting and pagination
	validSortFields := map[string]bool{
		"title":      true,
		"status":     true,
		"priority":   true,
		"due_date":   true,
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
		SELECT t.id, t.project_id, t.title, t.description, t.status, t.priority,
			t.due_date, t.created_by, t.assigned_to, t.estimated_hours,
			t.actual_hours, t.created_at, t.updated_at
		%s
		ORDER BY t.%s %s
		LIMIT %d OFFSET %d
	`, baseQuery, sortBy, sortOrder, params.PageSize, offset)

	// Execute the query
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	// Process the results
	tasks := []models.Task{}
	for rows.Next() {
		var task models.Task
		err := rows.StructScan(&task)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	// Get tags and related entities for each task
	for i := range tasks {
		tags, err := r.GetTags(ctx, tasks[i].ID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get tags: %w", err)
		}
		tasks[i].Tags = tags
	}

	return tasks, total, nil
}

// Update updates a task
func (r *PostgresTaskRepository) Update(ctx context.Context, task *models.Task) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Update task
	query := `
		UPDATE taskodex.tasks
		SET project_id = $1, title = $2, description = $3, status = $4, priority = $5,
			due_date = $6, assigned_to = $7, estimated_hours = $8, actual_hours = $9,
			updated_at = $10
		WHERE id = $11
	`

	task.UpdatedAt = time.Now()

	_, err = tx.ExecContext(
		ctx,
		query,
		task.ProjectID,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.DueDate,
		task.AssignedTo,
		task.EstimatedHours,
		task.ActualHours,
		task.UpdatedAt,
		task.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	// Get existing tags
	existingTags, err := r.GetTags(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing tags: %w", err)
	}

	// Add new tags and remove old ones
	existingTagMap := make(map[string]bool)
	for _, tag := range existingTags {
		existingTagMap[tag] = true
	}

	newTagMap := make(map[string]bool)
	for _, tag := range task.Tags {
		newTagMap[tag] = true
	}

	// Add new tags
	for _, tag := range task.Tags {
		if !existingTagMap[tag] {
			err = r.addTagTx(ctx, tx, task.ID, tag)
			if err != nil {
				return fmt.Errorf("failed to add tag: %w", err)
			}
		}
	}

	// Remove old tags
	for _, tag := range existingTags {
		if !newTagMap[tag] {
			err = r.removeTagTx(ctx, tx, task.ID, tag)
			if err != nil {
				return fmt.Errorf("failed to remove tag: %w", err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Delete deletes a task
func (r *PostgresTaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Delete tags
	_, err = tx.ExecContext(ctx, "DELETE FROM taskodex.task_tags WHERE task_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete task tags: %w", err)
	}

	// Delete task
	_, err = tx.ExecContext(ctx, "DELETE FROM taskodex.tasks WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// AddTag adds a tag to a task
func (r *PostgresTaskRepository) AddTag(ctx context.Context, taskID uuid.UUID, tag string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = r.addTagTx(ctx, tx, taskID, tag)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// addTagTx adds a tag to a task within a transaction
func (r *PostgresTaskRepository) addTagTx(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID, tag string) error {
	// Check if task exists
	var exists bool
	err := tx.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM taskodex.tasks WHERE id = $1)", taskID)
	if err != nil {
		return fmt.Errorf("failed to check if task exists: %w", err)
	}

	if !exists {
		return ErrTaskNotFound
	}

	// Insert tag
	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO taskodex.task_tags (task_id, tag) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		taskID,
		tag,
	)

	if err != nil {
		return fmt.Errorf("failed to insert tag: %w", err)
	}

	return nil
}

// RemoveTag removes a tag from a task
func (r *PostgresTaskRepository) RemoveTag(ctx context.Context, taskID uuid.UUID, tag string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = r.removeTagTx(ctx, tx, taskID, tag)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// removeTagTx removes a tag from a task within a transaction
func (r *PostgresTaskRepository) removeTagTx(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID, tag string) error {
	// Check if task exists
	var exists bool
	err := tx.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM taskodex.tasks WHERE id = $1)", taskID)
	if err != nil {
		return fmt.Errorf("failed to check if task exists: %w", err)
	}

	if !exists {
		return ErrTaskNotFound
	}

	// Delete tag
	_, err = tx.ExecContext(
		ctx,
		"DELETE FROM taskodex.task_tags WHERE task_id = $1 AND tag = $2",
		taskID,
		tag,
	)

	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

// GetTags retrieves all tags for a task
func (r *PostgresTaskRepository) GetTags(ctx context.Context, taskID uuid.UUID) ([]string, error) {
	query := "SELECT tag FROM taskodex.task_tags WHERE task_id = $1"

	var tags []string
	err := r.db.SelectContext(ctx, &tags, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	return tags, nil
}
