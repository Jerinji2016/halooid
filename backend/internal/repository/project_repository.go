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

// Common errors for project repository
var (
	ErrProjectNotFound = errors.New("project not found")
	ErrProjectNameExists = errors.New("project name already exists in this organization")
)

// ProjectRepository defines the interface for project data access
type ProjectRepository interface {
	// Create creates a new project
	Create(ctx context.Context, project *models.Project) error
	
	// GetByID retrieves a project by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error)
	
	// GetByName retrieves a project by name within an organization
	GetByName(ctx context.Context, organizationID uuid.UUID, name string) (*models.Project, error)
	
	// List retrieves projects based on filter parameters
	List(ctx context.Context, params models.ProjectListParams) ([]models.Project, int, error)
	
	// Update updates a project
	Update(ctx context.Context, project *models.Project) error
	
	// Delete deletes a project
	Delete(ctx context.Context, id uuid.UUID) error
}

// PostgresProjectRepository implements ProjectRepository using PostgreSQL
type PostgresProjectRepository struct {
	db *sqlx.DB
}

// NewPostgresProjectRepository creates a new PostgresProjectRepository
func NewPostgresProjectRepository(db *sqlx.DB) ProjectRepository {
	return &PostgresProjectRepository{db: db}
}

// Create creates a new project
func (r *PostgresProjectRepository) Create(ctx context.Context, project *models.Project) error {
	// Check if project name already exists in this organization
	existingProject, err := r.GetByName(ctx, project.OrganizationID, project.Name)
	if err != nil && !errors.Is(err, ErrProjectNotFound) {
		return fmt.Errorf("failed to check if project name exists: %w", err)
	}
	
	if existingProject != nil {
		return ErrProjectNameExists
	}
	
	query := `
		INSERT INTO taskodex.projects (
			id, organization_id, name, description, status, 
			start_date, end_date, created_by, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	
	_, err = r.db.ExecContext(
		ctx,
		query,
		project.ID,
		project.OrganizationID,
		project.Name,
		project.Description,
		project.Status,
		project.StartDate,
		project.EndDate,
		project.CreatedBy,
		project.CreatedAt,
		project.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to insert project: %w", err)
	}
	
	return nil
}

// GetByID retrieves a project by ID
func (r *PostgresProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	query := `
		SELECT p.id, p.organization_id, p.name, p.description, p.status, 
			p.start_date, p.end_date, p.created_by, p.created_at, p.updated_at
		FROM taskodex.projects p
		WHERE p.id = $1
	`
	
	var project models.Project
	err := r.db.GetContext(ctx, &project, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	
	// Get creator
	userRepo := NewPostgresUserRepository(r.db)
	creator, err := userRepo.GetByID(ctx, project.CreatedBy)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, fmt.Errorf("failed to get creator: %w", err)
	}
	if creator != nil {
		project.Creator = creator
	}
	
	return &project, nil
}

// GetByName retrieves a project by name within an organization
func (r *PostgresProjectRepository) GetByName(ctx context.Context, organizationID uuid.UUID, name string) (*models.Project, error) {
	query := `
		SELECT p.id, p.organization_id, p.name, p.description, p.status, 
			p.start_date, p.end_date, p.created_by, p.created_at, p.updated_at
		FROM taskodex.projects p
		WHERE p.organization_id = $1 AND p.name = $2
	`
	
	var project models.Project
	err := r.db.GetContext(ctx, &project, query, organizationID, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	
	// Get creator
	userRepo := NewPostgresUserRepository(r.db)
	creator, err := userRepo.GetByID(ctx, project.CreatedBy)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, fmt.Errorf("failed to get creator: %w", err)
	}
	if creator != nil {
		project.Creator = creator
	}
	
	return &project, nil
}

// List retrieves projects based on filter parameters
func (r *PostgresProjectRepository) List(ctx context.Context, params models.ProjectListParams) ([]models.Project, int, error) {
	// Build the query
	baseQuery := `
		FROM taskodex.projects p
		WHERE p.organization_id = $1
	`
	
	// Add filters
	filters := []string{}
	args := []interface{}{params.OrganizationID}
	argIndex := 2
	
	if params.Status != nil {
		filters = append(filters, fmt.Sprintf("p.status = $%d", argIndex))
		args = append(args, *params.Status)
		argIndex++
	}
	
	if params.CreatedBy != nil {
		filters = append(filters, fmt.Sprintf("p.created_by = $%d", argIndex))
		args = append(args, *params.CreatedBy)
		argIndex++
	}
	
	if params.SearchTerm != nil && *params.SearchTerm != "" {
		filters = append(filters, fmt.Sprintf("(p.name ILIKE $%d OR p.description ILIKE $%d)", argIndex, argIndex))
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
		return nil, 0, fmt.Errorf("failed to count projects: %w", err)
	}
	
	// Add sorting and pagination
	validSortFields := map[string]bool{
		"name":       true,
		"status":     true,
		"start_date": true,
		"end_date":   true,
		"created_at": true,
		"updated_at": true,
	}
	
	sortBy := "name"
	if validSortFields[params.SortBy] {
		sortBy = params.SortBy
	}
	
	sortOrder := "ASC"
	if strings.ToLower(params.SortOrder) == "desc" {
		sortOrder = "DESC"
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
		SELECT p.id, p.organization_id, p.name, p.description, p.status, 
			p.start_date, p.end_date, p.created_by, p.created_at, p.updated_at
		%s
		ORDER BY p.%s %s
		LIMIT %d OFFSET %d
	`, baseQuery, sortBy, sortOrder, params.PageSize, offset)
	
	// Execute the query
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()
	
	// Process the results
	projects := []models.Project{}
	for rows.Next() {
		var project models.Project
		err := rows.StructScan(&project)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, project)
	}
	
	return projects, total, nil
}

// Update updates a project
func (r *PostgresProjectRepository) Update(ctx context.Context, project *models.Project) error {
	// Check if project exists
	existingProject, err := r.GetByID(ctx, project.ID)
	if err != nil {
		return err
	}
	
	// Check if project name already exists (for a different project)
	if project.Name != existingProject.Name {
		nameExists, err := r.GetByName(ctx, project.OrganizationID, project.Name)
		if err != nil && !errors.Is(err, ErrProjectNotFound) {
			return fmt.Errorf("failed to check if project name exists: %w", err)
		}
		
		if nameExists != nil && nameExists.ID != project.ID {
			return ErrProjectNameExists
		}
	}
	
	query := `
		UPDATE taskodex.projects
		SET name = $1, description = $2, status = $3, start_date = $4, 
			end_date = $5, updated_at = $6
		WHERE id = $7
	`
	
	project.UpdatedAt = time.Now()
	
	_, err = r.db.ExecContext(
		ctx,
		query,
		project.Name,
		project.Description,
		project.Status,
		project.StartDate,
		project.EndDate,
		project.UpdatedAt,
		project.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}
	
	return nil
}

// Delete deletes a project
func (r *PostgresProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if project exists
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	// Check if project has tasks
	var taskCount int
	err = r.db.GetContext(ctx, &taskCount, "SELECT COUNT(*) FROM taskodex.tasks WHERE project_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to check if project has tasks: %w", err)
	}
	
	if taskCount > 0 {
		return fmt.Errorf("cannot delete project with tasks")
	}
	
	// Delete project
	_, err = r.db.ExecContext(ctx, "DELETE FROM taskodex.projects WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	
	return nil
}
