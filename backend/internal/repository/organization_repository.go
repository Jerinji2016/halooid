package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Common errors for organization repository
var (
	ErrOrganizationNotFound = errors.New("organization not found")
	ErrOrganizationNameExists = errors.New("organization name already exists")
)

// OrganizationRepository defines the interface for organization data access
type OrganizationRepository interface {
	// Create creates a new organization
	Create(ctx context.Context, organization *models.Organization) error
	
	// GetByID retrieves an organization by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error)
	
	// GetByName retrieves an organization by name
	GetByName(ctx context.Context, name string) (*models.Organization, error)
	
	// List retrieves all organizations
	List(ctx context.Context) ([]models.Organization, error)
	
	// Update updates an organization
	Update(ctx context.Context, organization *models.Organization) error
	
	// Delete marks an organization as inactive
	Delete(ctx context.Context, id uuid.UUID) error
	
	// AddUser adds a user to an organization
	AddUser(ctx context.Context, organizationUser *models.OrganizationUser) error
	
	// RemoveUser removes a user from an organization
	RemoveUser(ctx context.Context, organizationID, userID uuid.UUID) error
	
	// GetUsers retrieves all users in an organization
	GetUsers(ctx context.Context, organizationID uuid.UUID) ([]models.User, error)
	
	// GetOrganizations retrieves all organizations a user belongs to
	GetOrganizations(ctx context.Context, userID uuid.UUID) ([]models.Organization, error)
	
	// GetDB returns the database connection
	GetDB() *sqlx.DB
}

// PostgresOrganizationRepository implements OrganizationRepository using PostgreSQL
type PostgresOrganizationRepository struct {
	db *sqlx.DB
}

// NewPostgresOrganizationRepository creates a new PostgresOrganizationRepository
func NewPostgresOrganizationRepository(db *sqlx.DB) OrganizationRepository {
	return &PostgresOrganizationRepository{db: db}
}

// Create creates a new organization
func (r *PostgresOrganizationRepository) Create(ctx context.Context, organization *models.Organization) error {
	// Check if organization name already exists
	existingOrg, err := r.GetByName(ctx, organization.Name)
	if err != nil && !errors.Is(err, ErrOrganizationNotFound) {
		return err
	}
	if existingOrg != nil {
		return ErrOrganizationNameExists
	}

	query := `
		INSERT INTO organizations (id, name, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = r.db.ExecContext(
		ctx,
		query,
		organization.ID,
		organization.Name,
		organization.Description,
		organization.IsActive,
		organization.CreatedAt,
		organization.UpdatedAt,
	)

	if err != nil {
		return ErrDatabaseError
	}

	return nil
}

// GetByID retrieves an organization by ID
func (r *PostgresOrganizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM organizations
		WHERE id = $1
	`

	var organization models.Organization
	err := r.db.GetContext(ctx, &organization, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOrganizationNotFound
		}
		return nil, ErrDatabaseError
	}

	// Get users for the organization
	users, err := r.GetUsers(ctx, id)
	if err != nil {
		return nil, err
	}
	organization.Users = users

	return &organization, nil
}

// GetByName retrieves an organization by name
func (r *PostgresOrganizationRepository) GetByName(ctx context.Context, name string) (*models.Organization, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM organizations
		WHERE name = $1
	`

	var organization models.Organization
	err := r.db.GetContext(ctx, &organization, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOrganizationNotFound
		}
		return nil, ErrDatabaseError
	}

	// Get users for the organization
	users, err := r.GetUsers(ctx, organization.ID)
	if err != nil {
		return nil, err
	}
	organization.Users = users

	return &organization, nil
}

// List retrieves all organizations
func (r *PostgresOrganizationRepository) List(ctx context.Context) ([]models.Organization, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM organizations
		WHERE is_active = true
		ORDER BY name
	`

	var organizations []models.Organization
	err := r.db.SelectContext(ctx, &organizations, query)
	if err != nil {
		return nil, ErrDatabaseError
	}

	// Get users for each organization
	for i := range organizations {
		users, err := r.GetUsers(ctx, organizations[i].ID)
		if err != nil {
			return nil, err
		}
		organizations[i].Users = users
	}

	return organizations, nil
}

// Update updates an organization
func (r *PostgresOrganizationRepository) Update(ctx context.Context, organization *models.Organization) error {
	// Check if organization exists
	_, err := r.GetByID(ctx, organization.ID)
	if err != nil {
		return err
	}

	// Check if organization name already exists (for a different organization)
	existingOrg, err := r.GetByName(ctx, organization.Name)
	if err != nil && !errors.Is(err, ErrOrganizationNotFound) {
		return err
	}
	if existingOrg != nil && existingOrg.ID != organization.ID {
		return ErrOrganizationNameExists
	}

	query := `
		UPDATE organizations
		SET name = $1, description = $2, is_active = $3, updated_at = $4
		WHERE id = $5
	`

	organization.UpdatedAt = time.Now()

	_, err = r.db.ExecContext(
		ctx,
		query,
		organization.Name,
		organization.Description,
		organization.IsActive,
		organization.UpdatedAt,
		organization.ID,
	)

	if err != nil {
		return ErrDatabaseError
	}

	return nil
}

// Delete marks an organization as inactive
func (r *PostgresOrganizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE organizations
		SET is_active = false, updated_at = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return ErrDatabaseError
	}

	return nil
}

// AddUser adds a user to an organization
func (r *PostgresOrganizationRepository) AddUser(ctx context.Context, organizationUser *models.OrganizationUser) error {
	// Check if organization exists
	_, err := r.GetByID(ctx, organizationUser.OrganizationID)
	if err != nil {
		return err
	}

	// Check if user exists
	userRepo := NewPostgresUserRepository(r.db)
	_, err = userRepo.GetByID(ctx, organizationUser.UserID)
	if err != nil {
		return err
	}

	// Check if user is already in the organization
	query := `
		SELECT COUNT(*)
		FROM organization_users
		WHERE organization_id = $1 AND user_id = $2
	`

	var count int
	err = r.db.GetContext(ctx, &count, query, organizationUser.OrganizationID, organizationUser.UserID)
	if err != nil {
		return ErrDatabaseError
	}

	if count > 0 {
		// User is already in the organization, nothing to do
		return nil
	}

	// Add user to organization
	query = `
		INSERT INTO organization_users (organization_id, user_id)
		VALUES ($1, $2)
	`

	_, err = r.db.ExecContext(ctx, query, organizationUser.OrganizationID, organizationUser.UserID)
	if err != nil {
		return ErrDatabaseError
	}

	return nil
}

// RemoveUser removes a user from an organization
func (r *PostgresOrganizationRepository) RemoveUser(ctx context.Context, organizationID, userID uuid.UUID) error {
	query := `
		DELETE FROM organization_users
		WHERE organization_id = $1 AND user_id = $2
	`

	_, err := r.db.ExecContext(ctx, query, organizationID, userID)
	if err != nil {
		return ErrDatabaseError
	}

	return nil
}

// GetUsers retrieves all users in an organization
func (r *PostgresOrganizationRepository) GetUsers(ctx context.Context, organizationID uuid.UUID) ([]models.User, error) {
	query := `
		SELECT u.id, u.email, u.password_hash, u.first_name, u.last_name, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN organization_users ou ON u.id = ou.user_id
		WHERE ou.organization_id = $1 AND u.is_active = true
		ORDER BY u.email
	`

	var users []models.User
	err := r.db.SelectContext(ctx, &users, query, organizationID)
	if err != nil {
		return nil, ErrDatabaseError
	}

	return users, nil
}

// GetOrganizations retrieves all organizations a user belongs to
func (r *PostgresOrganizationRepository) GetOrganizations(ctx context.Context, userID uuid.UUID) ([]models.Organization, error) {
	query := `
		SELECT o.id, o.name, o.description, o.is_active, o.created_at, o.updated_at
		FROM organizations o
		JOIN organization_users ou ON o.id = ou.organization_id
		WHERE ou.user_id = $1 AND o.is_active = true
		ORDER BY o.name
	`

	var organizations []models.Organization
	err := r.db.SelectContext(ctx, &organizations, query, userID)
	if err != nil {
		return nil, ErrDatabaseError
	}

	return organizations, nil
}

// GetDB returns the database connection
func (r *PostgresOrganizationRepository) GetDB() *sqlx.DB {
	return r.db
}
