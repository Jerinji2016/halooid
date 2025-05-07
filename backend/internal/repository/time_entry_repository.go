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

// Common errors for time entry repository
var (
	ErrTimeEntryNotFound = errors.New("time entry not found")
	ErrRunningTimeEntry  = errors.New("user already has a running time entry for this task")
)

// TimeEntryRepository defines the interface for time entry data access
type TimeEntryRepository interface {
	// Create creates a new time entry
	Create(ctx context.Context, timeEntry *models.TimeEntry) error
	
	// GetByID retrieves a time entry by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.TimeEntry, error)
	
	// List retrieves time entries based on filter parameters
	List(ctx context.Context, params models.TimeEntryListParams) ([]models.TimeEntry, int, error)
	
	// Update updates a time entry
	Update(ctx context.Context, timeEntry *models.TimeEntry) error
	
	// Delete deletes a time entry
	Delete(ctx context.Context, id uuid.UUID) error
	
	// GetRunningTimeEntry retrieves a running time entry for a user and task
	GetRunningTimeEntry(ctx context.Context, userID, taskID uuid.UUID) (*models.TimeEntry, error)
	
	// GetRunningTimeEntries retrieves all running time entries for a user
	GetRunningTimeEntries(ctx context.Context, userID uuid.UUID) ([]models.TimeEntry, error)
	
	// Aggregate aggregates time entries based on parameters
	Aggregate(ctx context.Context, params models.TimeEntryAggregationParams) ([]models.TimeEntryAggregation, error)
}

// PostgresTimeEntryRepository implements TimeEntryRepository using PostgreSQL
type PostgresTimeEntryRepository struct {
	db *sqlx.DB
}

// NewPostgresTimeEntryRepository creates a new PostgresTimeEntryRepository
func NewPostgresTimeEntryRepository(db *sqlx.DB) TimeEntryRepository {
	return &PostgresTimeEntryRepository{db: db}
}

// Create creates a new time entry
func (r *PostgresTimeEntryRepository) Create(ctx context.Context, timeEntry *models.TimeEntry) error {
	// Check if user already has a running time entry for this task
	if timeEntry.EndTime == nil {
		existingEntry, err := r.GetRunningTimeEntry(ctx, timeEntry.UserID, timeEntry.TaskID)
		if err != nil && !errors.Is(err, ErrTimeEntryNotFound) {
			return err
		}
		if existingEntry != nil {
			return ErrRunningTimeEntry
		}
	}
	
	query := `
		INSERT INTO taskodex.task_time_entries (
			id, task_id, user_id, start_time, end_time, 
			duration_minutes, description, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	_, err := r.db.ExecContext(
		ctx,
		query,
		timeEntry.ID,
		timeEntry.TaskID,
		timeEntry.UserID,
		timeEntry.StartTime,
		timeEntry.EndTime,
		timeEntry.DurationMinutes,
		timeEntry.Description,
		timeEntry.CreatedAt,
		timeEntry.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to insert time entry: %w", err)
	}
	
	return nil
}

// GetByID retrieves a time entry by ID
func (r *PostgresTimeEntryRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.TimeEntry, error) {
	query := `
		SELECT te.id, te.task_id, te.user_id, te.start_time, te.end_time, 
			te.duration_minutes, te.description, te.created_at, te.updated_at
		FROM taskodex.task_time_entries te
		WHERE te.id = $1
	`
	
	var timeEntry models.TimeEntry
	err := r.db.GetContext(ctx, &timeEntry, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTimeEntryNotFound
		}
		return nil, fmt.Errorf("failed to get time entry: %w", err)
	}
	
	// Get task
	taskRepo := NewPostgresTaskRepository(r.db)
	task, err := taskRepo.GetByID(ctx, timeEntry.TaskID)
	if err != nil && !errors.Is(err, ErrTaskNotFound) {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if task != nil {
		timeEntry.Task = task
	}
	
	// Get user
	userRepo := NewPostgresUserRepository(r.db)
	user, err := userRepo.GetByID(ctx, timeEntry.UserID)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user != nil {
		timeEntry.User = user
	}
	
	return &timeEntry, nil
}

// List retrieves time entries based on filter parameters
func (r *PostgresTimeEntryRepository) List(ctx context.Context, params models.TimeEntryListParams) ([]models.TimeEntry, int, error) {
	// Build the query
	baseQuery := `
		FROM taskodex.task_time_entries te
		WHERE 1=1
	`
	
	// Add filters
	filters := []string{}
	args := []interface{}{}
	argIndex := 1
	
	if params.TaskID != nil {
		filters = append(filters, fmt.Sprintf("te.task_id = $%d", argIndex))
		args = append(args, *params.TaskID)
		argIndex++
	}
	
	if params.UserID != nil {
		filters = append(filters, fmt.Sprintf("te.user_id = $%d", argIndex))
		args = append(args, *params.UserID)
		argIndex++
	}
	
	if params.StartAfter != nil {
		filters = append(filters, fmt.Sprintf("te.start_time >= $%d", argIndex))
		args = append(args, *params.StartAfter)
		argIndex++
	}
	
	if params.StartBefore != nil {
		filters = append(filters, fmt.Sprintf("te.start_time <= $%d", argIndex))
		args = append(args, *params.StartBefore)
		argIndex++
	}
	
	if params.IsRunning != nil {
		if *params.IsRunning {
			filters = append(filters, "te.end_time IS NULL")
		} else {
			filters = append(filters, "te.end_time IS NOT NULL")
		}
	}
	
	if len(filters) > 0 {
		baseQuery += " AND " + strings.Join(filters, " AND ")
	}
	
	// Count total records
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count time entries: %w", err)
	}
	
	// Add sorting and pagination
	validSortFields := map[string]bool{
		"start_time":  true,
		"end_time":    true,
		"created_at":  true,
		"updated_at":  true,
	}
	
	sortBy := "start_time"
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
		SELECT te.id, te.task_id, te.user_id, te.start_time, te.end_time, 
			te.duration_minutes, te.description, te.created_at, te.updated_at
		%s
		ORDER BY te.%s %s
		LIMIT %d OFFSET %d
	`, baseQuery, sortBy, sortOrder, params.PageSize, offset)
	
	// Execute the query
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query time entries: %w", err)
	}
	defer rows.Close()
	
	// Process the results
	timeEntries := []models.TimeEntry{}
	for rows.Next() {
		var timeEntry models.TimeEntry
		err := rows.StructScan(&timeEntry)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan time entry: %w", err)
		}
		timeEntries = append(timeEntries, timeEntry)
	}
	
	return timeEntries, total, nil
}

// Update updates a time entry
func (r *PostgresTimeEntryRepository) Update(ctx context.Context, timeEntry *models.TimeEntry) error {
	query := `
		UPDATE taskodex.task_time_entries
		SET start_time = $1, end_time = $2, duration_minutes = $3, 
			description = $4, updated_at = $5
		WHERE id = $6
	`
	
	timeEntry.UpdatedAt = time.Now()
	
	_, err := r.db.ExecContext(
		ctx,
		query,
		timeEntry.StartTime,
		timeEntry.EndTime,
		timeEntry.DurationMinutes,
		timeEntry.Description,
		timeEntry.UpdatedAt,
		timeEntry.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update time entry: %w", err)
	}
	
	return nil
}

// Delete deletes a time entry
func (r *PostgresTimeEntryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM taskodex.task_time_entries
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete time entry: %w", err)
	}
	
	return nil
}

// GetRunningTimeEntry retrieves a running time entry for a user and task
func (r *PostgresTimeEntryRepository) GetRunningTimeEntry(ctx context.Context, userID, taskID uuid.UUID) (*models.TimeEntry, error) {
	query := `
		SELECT te.id, te.task_id, te.user_id, te.start_time, te.end_time, 
			te.duration_minutes, te.description, te.created_at, te.updated_at
		FROM taskodex.task_time_entries te
		WHERE te.user_id = $1 AND te.task_id = $2 AND te.end_time IS NULL
		ORDER BY te.start_time DESC
		LIMIT 1
	`
	
	var timeEntry models.TimeEntry
	err := r.db.GetContext(ctx, &timeEntry, query, userID, taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTimeEntryNotFound
		}
		return nil, fmt.Errorf("failed to get running time entry: %w", err)
	}
	
	return &timeEntry, nil
}

// GetRunningTimeEntries retrieves all running time entries for a user
func (r *PostgresTimeEntryRepository) GetRunningTimeEntries(ctx context.Context, userID uuid.UUID) ([]models.TimeEntry, error) {
	query := `
		SELECT te.id, te.task_id, te.user_id, te.start_time, te.end_time, 
			te.duration_minutes, te.description, te.created_at, te.updated_at
		FROM taskodex.task_time_entries te
		WHERE te.user_id = $1 AND te.end_time IS NULL
		ORDER BY te.start_time DESC
	`
	
	timeEntries := []models.TimeEntry{}
	err := r.db.SelectContext(ctx, &timeEntries, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get running time entries: %w", err)
	}
	
	return timeEntries, nil
}

// Aggregate aggregates time entries based on parameters
func (r *PostgresTimeEntryRepository) Aggregate(ctx context.Context, params models.TimeEntryAggregationParams) ([]models.TimeEntryAggregation, error) {
	// Build the base query
	baseQuery := `
		FROM taskodex.task_time_entries te
		LEFT JOIN taskodex.tasks t ON te.task_id = t.id
		WHERE 1=1
	`
	
	// Add filters
	filters := []string{}
	args := []interface{}{}
	argIndex := 1
	
	if params.TaskID != nil {
		filters = append(filters, fmt.Sprintf("te.task_id = $%d", argIndex))
		args = append(args, *params.TaskID)
		argIndex++
	}
	
	if params.UserID != nil {
		filters = append(filters, fmt.Sprintf("te.user_id = $%d", argIndex))
		args = append(args, *params.UserID)
		argIndex++
	}
	
	if params.ProjectID != nil {
		filters = append(filters, fmt.Sprintf("t.project_id = $%d", argIndex))
		args = append(args, *params.ProjectID)
		argIndex++
	}
	
	if params.StartAfter != nil {
		filters = append(filters, fmt.Sprintf("te.start_time >= $%d", argIndex))
		args = append(args, *params.StartAfter)
		argIndex++
	}
	
	if params.StartBefore != nil {
		filters = append(filters, fmt.Sprintf("te.start_time <= $%d", argIndex))
		args = append(args, *params.StartBefore)
		argIndex++
	}
	
	if len(filters) > 0 {
		baseQuery += " AND " + strings.Join(filters, " AND ")
	}
	
	// Build the group by clause based on the groupBy parameter
	var groupBy, selectFields string
	
	switch params.GroupBy {
	case "day":
		groupBy = "DATE(te.start_time)"
		selectFields = "DATE(te.start_time) AS date"
	case "week":
		groupBy = "EXTRACT(YEAR FROM te.start_time), EXTRACT(WEEK FROM te.start_time)"
		selectFields = "EXTRACT(YEAR FROM te.start_time) AS year, EXTRACT(WEEK FROM te.start_time) AS week"
	case "month":
		groupBy = "EXTRACT(YEAR FROM te.start_time), EXTRACT(MONTH FROM te.start_time)"
		selectFields = "EXTRACT(YEAR FROM te.start_time) AS year, EXTRACT(MONTH FROM te.start_time) AS month"
	case "year":
		groupBy = "EXTRACT(YEAR FROM te.start_time)"
		selectFields = "EXTRACT(YEAR FROM te.start_time) AS year"
	case "task":
		groupBy = "te.task_id"
		selectFields = "te.task_id"
	case "user":
		groupBy = "te.user_id"
		selectFields = "te.user_id"
	default:
		groupBy = "DATE(te.start_time)"
		selectFields = "DATE(te.start_time) AS date"
	}
	
	// Build the final query
	query := fmt.Sprintf(`
		SELECT %s, 
			SUM(COALESCE(te.duration_minutes, 
				CASE 
					WHEN te.end_time IS NOT NULL THEN EXTRACT(EPOCH FROM (te.end_time - te.start_time))/60 
					ELSE EXTRACT(EPOCH FROM (NOW() - te.start_time))/60 
				END
			))::INTEGER AS total_duration_minutes
		%s
		GROUP BY %s
		ORDER BY %s
	`, selectFields, baseQuery, groupBy, groupBy)
	
	// Execute the query
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate time entries: %w", err)
	}
	defer rows.Close()
	
	// Process the results
	aggregations := []models.TimeEntryAggregation{}
	
	for rows.Next() {
		var agg models.TimeEntryAggregation
		
		switch params.GroupBy {
		case "day":
			var date time.Time
			err = rows.Scan(&date, &agg.TotalDurationMinutes)
			agg.Date = &date
		case "week":
			var year, week int
			err = rows.Scan(&year, &week, &agg.TotalDurationMinutes)
			agg.Year = &year
			agg.Week = &week
		case "month":
			var year, month int
			err = rows.Scan(&year, &month, &agg.TotalDurationMinutes)
			agg.Year = &year
			agg.Month = &month
		case "year":
			var year int
			err = rows.Scan(&year, &agg.TotalDurationMinutes)
			agg.Year = &year
		case "task":
			var taskID uuid.UUID
			err = rows.Scan(&taskID, &agg.TotalDurationMinutes)
			agg.TaskID = &taskID
		case "user":
			var userID uuid.UUID
			err = rows.Scan(&userID, &agg.TotalDurationMinutes)
			agg.UserID = &userID
		}
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan aggregation: %w", err)
		}
		
		aggregations = append(aggregations, agg)
	}
	
	return aggregations, nil
}
