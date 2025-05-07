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

// Common errors for employee repository
var (
	ErrEmployeeNotFound      = errors.New("employee not found")
	ErrEmployeeIDExists      = errors.New("employee ID already exists")
	ErrUserAlreadyEmployee   = errors.New("user is already an employee in this organization")
	ErrManagerNotFound       = errors.New("manager not found")
)

// EmployeeRepository defines the interface for employee data access
type EmployeeRepository interface {
	// Create creates a new employee
	Create(ctx context.Context, employee *models.Employee) error
	
	// GetByID retrieves an employee by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.Employee, error)
	
	// GetByEmployeeID retrieves an employee by employee ID
	GetByEmployeeID(ctx context.Context, organizationID uuid.UUID, employeeID string) (*models.Employee, error)
	
	// GetByUserID retrieves an employee by user ID
	GetByUserID(ctx context.Context, organizationID, userID uuid.UUID) (*models.Employee, error)
	
	// List retrieves employees based on filter parameters
	List(ctx context.Context, params models.EmployeeListParams) ([]models.Employee, int, error)
	
	// Update updates an employee
	Update(ctx context.Context, employee *models.Employee) error
	
	// Delete marks an employee as inactive
	Delete(ctx context.Context, id uuid.UUID) error
	
	// GetDB returns the database connection
	GetDB() *sqlx.DB
}

// PostgresEmployeeRepository implements EmployeeRepository using PostgreSQL
type PostgresEmployeeRepository struct {
	db *sqlx.DB
}

// NewPostgresEmployeeRepository creates a new PostgresEmployeeRepository
func NewPostgresEmployeeRepository(db *sqlx.DB) EmployeeRepository {
	return &PostgresEmployeeRepository{db: db}
}

// Create creates a new employee
func (r *PostgresEmployeeRepository) Create(ctx context.Context, employee *models.Employee) error {
	// Check if employee ID already exists in the organization
	existingEmployee, err := r.GetByEmployeeID(ctx, employee.OrganizationID, employee.EmployeeID)
	if err != nil && !errors.Is(err, ErrEmployeeNotFound) {
		return err
	}
	if existingEmployee != nil {
		return ErrEmployeeIDExists
	}

	// Check if user is already an employee in this organization
	existingEmployee, err = r.GetByUserID(ctx, employee.OrganizationID, employee.UserID)
	if err != nil && !errors.Is(err, ErrEmployeeNotFound) {
		return err
	}
	if existingEmployee != nil {
		return ErrUserAlreadyEmployee
	}

	// Check if manager exists
	if employee.ManagerID != nil {
		_, err := r.GetByID(ctx, *employee.ManagerID)
		if err != nil {
			if errors.Is(err, ErrEmployeeNotFound) {
				return ErrManagerNotFound
			}
			return err
		}
	}

	query := `
		INSERT INTO qultrix.employees (
			id, organization_id, user_id, employee_id, department, position, 
			hire_date, manager_id, salary, is_active, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err = r.db.ExecContext(
		ctx,
		query,
		employee.ID,
		employee.OrganizationID,
		employee.UserID,
		employee.EmployeeID,
		employee.Department,
		employee.Position,
		employee.HireDate,
		employee.ManagerID,
		employee.Salary,
		employee.IsActive,
		employee.CreatedAt,
		employee.UpdatedAt,
	)

	if err != nil {
		return ErrDatabaseError
	}

	return nil
}

// GetByID retrieves an employee by ID
func (r *PostgresEmployeeRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
	query := `
		SELECT e.id, e.organization_id, e.user_id, e.employee_id, e.department, e.position, 
			e.hire_date, e.manager_id, e.salary, e.is_active, e.created_at, e.updated_at
		FROM qultrix.employees e
		WHERE e.id = $1
	`

	var employee models.Employee
	err := r.db.GetContext(ctx, &employee, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEmployeeNotFound
		}
		return nil, ErrDatabaseError
	}

	// Get user information
	userRepo := NewPostgresUserRepository(r.db)
	user, err := userRepo.GetByID(ctx, employee.UserID)
	if err != nil {
		return nil, err
	}
	employee.User = user

	// Get manager information if available
	if employee.ManagerID != nil {
		manager, err := r.GetByID(ctx, *employee.ManagerID)
		if err != nil && !errors.Is(err, ErrEmployeeNotFound) {
			return nil, err
		}
		if manager != nil {
			employee.Manager = manager
		}
	}

	return &employee, nil
}

// GetByEmployeeID retrieves an employee by employee ID
func (r *PostgresEmployeeRepository) GetByEmployeeID(ctx context.Context, organizationID uuid.UUID, employeeID string) (*models.Employee, error) {
	query := `
		SELECT e.id, e.organization_id, e.user_id, e.employee_id, e.department, e.position, 
			e.hire_date, e.manager_id, e.salary, e.is_active, e.created_at, e.updated_at
		FROM qultrix.employees e
		WHERE e.organization_id = $1 AND e.employee_id = $2
	`

	var employee models.Employee
	err := r.db.GetContext(ctx, &employee, query, organizationID, employeeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEmployeeNotFound
		}
		return nil, ErrDatabaseError
	}

	// Get user information
	userRepo := NewPostgresUserRepository(r.db)
	user, err := userRepo.GetByID(ctx, employee.UserID)
	if err != nil {
		return nil, err
	}
	employee.User = user

	// Get manager information if available
	if employee.ManagerID != nil {
		manager, err := r.GetByID(ctx, *employee.ManagerID)
		if err != nil && !errors.Is(err, ErrEmployeeNotFound) {
			return nil, err
		}
		if manager != nil {
			employee.Manager = manager
		}
	}

	return &employee, nil
}

// GetByUserID retrieves an employee by user ID
func (r *PostgresEmployeeRepository) GetByUserID(ctx context.Context, organizationID, userID uuid.UUID) (*models.Employee, error) {
	query := `
		SELECT e.id, e.organization_id, e.user_id, e.employee_id, e.department, e.position, 
			e.hire_date, e.manager_id, e.salary, e.is_active, e.created_at, e.updated_at
		FROM qultrix.employees e
		WHERE e.organization_id = $1 AND e.user_id = $2
	`

	var employee models.Employee
	err := r.db.GetContext(ctx, &employee, query, organizationID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEmployeeNotFound
		}
		return nil, ErrDatabaseError
	}

	// Get user information
	userRepo := NewPostgresUserRepository(r.db)
	user, err := userRepo.GetByID(ctx, employee.UserID)
	if err != nil {
		return nil, err
	}
	employee.User = user

	// Get manager information if available
	if employee.ManagerID != nil {
		manager, err := r.GetByID(ctx, *employee.ManagerID)
		if err != nil && !errors.Is(err, ErrEmployeeNotFound) {
			return nil, err
		}
		if manager != nil {
			employee.Manager = manager
		}
	}

	return &employee, nil
}

// List retrieves employees based on filter parameters
func (r *PostgresEmployeeRepository) List(ctx context.Context, params models.EmployeeListParams) ([]models.Employee, int, error) {
	// Build the query
	baseQuery := `
		FROM qultrix.employees e
		WHERE e.organization_id = :organization_id
	`

	// Add filters
	filters := []string{}
	args := map[string]interface{}{
		"organization_id": params.OrganizationID,
	}

	if params.Department != "" {
		filters = append(filters, "e.department = :department")
		args["department"] = params.Department
	}

	if params.Position != "" {
		filters = append(filters, "e.position = :position")
		args["position"] = params.Position
	}

	if params.IsActive != nil {
		filters = append(filters, "e.is_active = :is_active")
		args["is_active"] = *params.IsActive
	}

	if len(filters) > 0 {
		baseQuery += " AND " + strings.Join(filters, " AND ")
	}

	// Count total records
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	countStmt, err := r.db.PrepareNamedContext(ctx, countQuery)
	if err != nil {
		return nil, 0, ErrDatabaseError
	}
	defer countStmt.Close()

	err = countStmt.GetContext(ctx, &total, args)
	if err != nil {
		return nil, 0, ErrDatabaseError
	}

	// Add sorting and pagination
	validSortFields := map[string]bool{
		"employee_id": true,
		"department":  true,
		"position":    true,
		"hire_date":   true,
		"salary":      true,
	}

	sortBy := "employee_id"
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
		params.PageSize = 10
	}

	offset := (params.Page - 1) * params.PageSize

	// Build the final query
	query := fmt.Sprintf(`
		SELECT e.id, e.organization_id, e.user_id, e.employee_id, e.department, e.position, 
			e.hire_date, e.manager_id, e.salary, e.is_active, e.created_at, e.updated_at
		%s
		ORDER BY e.%s %s
		LIMIT %d OFFSET %d
	`, baseQuery, sortBy, sortOrder, params.PageSize, offset)

	// Execute the query
	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, 0, ErrDatabaseError
	}
	defer rows.Close()

	// Process the results
	employees := []models.Employee{}
	for rows.Next() {
		var employee models.Employee
		err := rows.StructScan(&employee)
		if err != nil {
			return nil, 0, ErrDatabaseError
		}
		employees = append(employees, employee)
	}

	// Get user information for each employee
	userRepo := NewPostgresUserRepository(r.db)
	for i := range employees {
		user, err := userRepo.GetByID(ctx, employees[i].UserID)
		if err != nil {
			return nil, 0, err
		}
		employees[i].User = user
	}

	return employees, total, nil
}

// Update updates an employee
func (r *PostgresEmployeeRepository) Update(ctx context.Context, employee *models.Employee) error {
	// Check if employee exists
	existingEmployee, err := r.GetByID(ctx, employee.ID)
	if err != nil {
		return err
	}

	// Check if employee ID already exists (for a different employee)
	if employee.EmployeeID != existingEmployee.EmployeeID {
		existingByEmployeeID, err := r.GetByEmployeeID(ctx, employee.OrganizationID, employee.EmployeeID)
		if err != nil && !errors.Is(err, ErrEmployeeNotFound) {
			return err
		}
		if existingByEmployeeID != nil && existingByEmployeeID.ID != employee.ID {
			return ErrEmployeeIDExists
		}
	}

	// Check if manager exists
	if employee.ManagerID != nil {
		// Prevent circular manager relationship
		if *employee.ManagerID == employee.ID {
			return ErrManagerNotFound
		}

		_, err := r.GetByID(ctx, *employee.ManagerID)
		if err != nil {
			if errors.Is(err, ErrEmployeeNotFound) {
				return ErrManagerNotFound
			}
			return err
		}
	}

	query := `
		UPDATE qultrix.employees
		SET employee_id = $1, department = $2, position = $3, hire_date = $4,
			manager_id = $5, salary = $6, is_active = $7, updated_at = $8
		WHERE id = $9
	`

	employee.UpdatedAt = time.Now()

	_, err = r.db.ExecContext(
		ctx,
		query,
		employee.EmployeeID,
		employee.Department,
		employee.Position,
		employee.HireDate,
		employee.ManagerID,
		employee.Salary,
		employee.IsActive,
		employee.UpdatedAt,
		employee.ID,
	)

	if err != nil {
		return ErrDatabaseError
	}

	return nil
}

// Delete marks an employee as inactive
func (r *PostgresEmployeeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE qultrix.employees
		SET is_active = false, updated_at = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return ErrDatabaseError
	}

	return nil
}

// GetDB returns the database connection
func (r *PostgresEmployeeRepository) GetDB() *sqlx.DB {
	return r.db
}
