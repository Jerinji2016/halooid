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

// Common errors
var (
	ErrRoleNotFound        = errors.New("role not found")
	ErrPermissionNotFound  = errors.New("permission not found")
	ErrRoleNameExists      = errors.New("role name already exists")
	ErrPermissionNameExists = errors.New("permission name already exists")
)

// RoleRepository defines the interface for role data access
type RoleRepository interface {
	// CreateRole creates a new role
	CreateRole(ctx context.Context, role *models.Role) error
	
	// GetRoleByID retrieves a role by ID
	GetRoleByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
	
	// GetRoleByName retrieves a role by name
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
	
	// ListRoles retrieves all roles
	ListRoles(ctx context.Context) ([]models.Role, error)
	
	// UpdateRole updates a role
	UpdateRole(ctx context.Context, role *models.Role) error
	
	// DeleteRole deletes a role
	DeleteRole(ctx context.Context, id uuid.UUID) error
	
	// CreatePermission creates a new permission
	CreatePermission(ctx context.Context, permission *models.Permission) error
	
	// GetPermissionByID retrieves a permission by ID
	GetPermissionByID(ctx context.Context, id uuid.UUID) (*models.Permission, error)
	
	// GetPermissionByName retrieves a permission by name
	GetPermissionByName(ctx context.Context, name string) (*models.Permission, error)
	
	// ListPermissions retrieves all permissions
	ListPermissions(ctx context.Context) ([]models.Permission, error)
	
	// UpdatePermission updates a permission
	UpdatePermission(ctx context.Context, permission *models.Permission) error
	
	// DeletePermission deletes a permission
	DeletePermission(ctx context.Context, id uuid.UUID) error
	
	// AssignPermissionToRole assigns a permission to a role
	AssignPermissionToRole(ctx context.Context, roleID, permissionID uuid.UUID) error
	
	// RemovePermissionFromRole removes a permission from a role
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID uuid.UUID) error
	
	// GetRolePermissions retrieves all permissions for a role
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]models.Permission, error)
	
	// AssignRoleToUser assigns a role to a user
	AssignRoleToUser(ctx context.Context, userRole *models.UserRole) error
	
	// RemoveRoleFromUser removes a role from a user
	RemoveRoleFromUser(ctx context.Context, userID, roleID, organizationID uuid.UUID) error
	
	// GetUserRoles retrieves all roles for a user
	GetUserRoles(ctx context.Context, userID, organizationID uuid.UUID) ([]models.Role, error)
	
	// HasPermission checks if a user has a specific permission
	HasPermission(ctx context.Context, userID, organizationID uuid.UUID, permissionName string) (bool, error)
}

// PostgresRoleRepository implements RoleRepository using PostgreSQL
type PostgresRoleRepository struct {
	db *sqlx.DB
}

// NewPostgresRoleRepository creates a new PostgresRoleRepository
func NewPostgresRoleRepository(db *sqlx.DB) RoleRepository {
	return &PostgresRoleRepository{db: db}
}

// CreateRole creates a new role
func (r *PostgresRoleRepository) CreateRole(ctx context.Context, role *models.Role) error {
	// Check if role name already exists
	existingRole, err := r.GetRoleByName(ctx, role.Name)
	if err != nil && !errors.Is(err, ErrRoleNotFound) {
		return err
	}
	if existingRole != nil {
		return ErrRoleNameExists
	}

	query := `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = r.db.ExecContext(
		ctx,
		query,
		role.ID,
		role.Name,
		role.Description,
		role.CreatedAt,
		role.UpdatedAt,
	)

	if err != nil {
		return ErrDatabaseError
	}

	return nil
}

// GetRoleByID retrieves a role by ID
func (r *PostgresRoleRepository) GetRoleByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE id = $1
	`

	var role models.Role
	err := r.db.GetContext(ctx, &role, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRoleNotFound
		}
		return nil, ErrDatabaseError
	}

	// Get permissions for the role
	permissions, err := r.GetRolePermissions(ctx, id)
	if err != nil {
		return nil, err
	}
	role.Permissions = permissions

	return &role, nil
}

// GetRoleByName retrieves a role by name
func (r *PostgresRoleRepository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE name = $1
	`

	var role models.Role
	err := r.db.GetContext(ctx, &role, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRoleNotFound
		}
		return nil, ErrDatabaseError
	}

	// Get permissions for the role
	permissions, err := r.GetRolePermissions(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	role.Permissions = permissions

	return &role, nil
}

// ListRoles retrieves all roles
func (r *PostgresRoleRepository) ListRoles(ctx context.Context) ([]models.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		ORDER BY name
	`

	var roles []models.Role
	err := r.db.SelectContext(ctx, &roles, query)
	if err != nil {
		return nil, ErrDatabaseError
	}

	// Get permissions for each role
	for i := range roles {
		permissions, err := r.GetRolePermissions(ctx, roles[i].ID)
		if err != nil {
			return nil, err
		}
		roles[i].Permissions = permissions
	}

	return roles, nil
}

// UpdateRole updates a role
func (r *PostgresRoleRepository) UpdateRole(ctx context.Context, role *models.Role) error {
	// Check if role exists
	_, err := r.GetRoleByID(ctx, role.ID)
	if err != nil {
		return err
	}

	// Check if role name already exists (for a different role)
	existingRole, err := r.GetRoleByName(ctx, role.Name)
	if err != nil && !errors.Is(err, ErrRoleNotFound) {
		return err
	}
	if existingRole != nil && existingRole.ID != role.ID {
		return ErrRoleNameExists
	}

	query := `
		UPDATE roles
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4
	`

	role.UpdatedAt = time.Now()

	_, err = r.db.ExecContext(
		ctx,
		query,
		role.Name,
		role.Description,
		role.UpdatedAt,
		role.ID,
	)

	if err != nil {
		return ErrDatabaseError
	}

	return nil
}

// DeleteRole deletes a role
func (r *PostgresRoleRepository) DeleteRole(ctx context.Context, id uuid.UUID) error {
	// Check if role exists
	_, err := r.GetRoleByID(ctx, id)
	if err != nil {
		return err
	}

	// Start a transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return ErrDatabaseError
	}
	defer tx.Rollback()

	// Delete role permissions
	_, err = tx.ExecContext(ctx, "DELETE FROM role_permissions WHERE role_id = $1", id)
	if err != nil {
		return ErrDatabaseError
	}

	// Delete user roles
	_, err = tx.ExecContext(ctx, "DELETE FROM user_roles WHERE role_id = $1", id)
	if err != nil {
		return ErrDatabaseError
	}

	// Delete role
	_, err = tx.ExecContext(ctx, "DELETE FROM roles WHERE id = $1", id)
	if err != nil {
		return ErrDatabaseError
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return ErrDatabaseError
	}

	return nil
}

// CreatePermission creates a new permission
func (r *PostgresRoleRepository) CreatePermission(ctx context.Context, permission *models.Permission) error {
	// Check if permission name already exists
	existingPermission, err := r.GetPermissionByName(ctx, permission.Name)
	if err != nil && !errors.Is(err, ErrPermissionNotFound) {
		return err
	}
	if existingPermission != nil {
		return ErrPermissionNameExists
	}

	query := `
		INSERT INTO permissions (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = r.db.ExecContext(
		ctx,
		query,
		permission.ID,
		permission.Name,
		permission.Description,
		permission.CreatedAt,
		permission.UpdatedAt,
	)

	if err != nil {
		return ErrDatabaseError
	}

	return nil
}

// GetPermissionByID retrieves a permission by ID
func (r *PostgresRoleRepository) GetPermissionByID(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM permissions
		WHERE id = $1
	`

	var permission models.Permission
	err := r.db.GetContext(ctx, &permission, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPermissionNotFound
		}
		return nil, ErrDatabaseError
	}

	return &permission, nil
}

// GetPermissionByName retrieves a permission by name
func (r *PostgresRoleRepository) GetPermissionByName(ctx context.Context, name string) (*models.Permission, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM permissions
		WHERE name = $1
	`

	var permission models.Permission
	err := r.db.GetContext(ctx, &permission, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPermissionNotFound
		}
		return nil, ErrDatabaseError
	}

	return &permission, nil
}

// ListPermissions retrieves all permissions
func (r *PostgresRoleRepository) ListPermissions(ctx context.Context) ([]models.Permission, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM permissions
		ORDER BY name
	`

	var permissions []models.Permission
	err := r.db.SelectContext(ctx, &permissions, query)
	if err != nil {
		return nil, ErrDatabaseError
	}

	return permissions, nil
}

// UpdatePermission updates a permission
func (r *PostgresRoleRepository) UpdatePermission(ctx context.Context, permission *models.Permission) error {
	// Check if permission exists
	_, err := r.GetPermissionByID(ctx, permission.ID)
	if err != nil {
		return err
	}

	// Check if permission name already exists (for a different permission)
	existingPermission, err := r.GetPermissionByName(ctx, permission.Name)
	if err != nil && !errors.Is(err, ErrPermissionNotFound) {
		return err
	}
	if existingPermission != nil && existingPermission.ID != permission.ID {
		return ErrPermissionNameExists
	}

	query := `
		UPDATE permissions
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4
	`

	permission.UpdatedAt = time.Now()

	_, err = r.db.ExecContext(
		ctx,
		query,
		permission.Name,
		permission.Description,
		permission.UpdatedAt,
		permission.ID,
	)

	if err != nil {
		return ErrDatabaseError
	}

	return nil
}

// DeletePermission deletes a permission
func (r *PostgresRoleRepository) DeletePermission(ctx context.Context, id uuid.UUID) error {
	// Check if permission exists
	_, err := r.GetPermissionByID(ctx, id)
	if err != nil {
		return err
	}

	// Start a transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return ErrDatabaseError
	}
	defer tx.Rollback()

	// Delete role permissions
	_, err = tx.ExecContext(ctx, "DELETE FROM role_permissions WHERE permission_id = $1", id)
	if err != nil {
		return ErrDatabaseError
	}

	// Delete permission
	_, err = tx.ExecContext(ctx, "DELETE FROM permissions WHERE id = $1", id)
	if err != nil {
		return ErrDatabaseError
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return ErrDatabaseError
	}

	return nil
}

// AssignPermissionToRole assigns a permission to a role
func (r *PostgresRoleRepository) AssignPermissionToRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	// Check if role exists
	_, err := r.GetRoleByID(ctx, roleID)
	if err != nil {
		return err
	}

	// Check if permission exists
	_, err = r.GetPermissionByID(ctx, permissionID)
	if err != nil {
		return err
	}

	// Check if permission is already assigned to role
	query := `
		SELECT COUNT(*)
		FROM role_permissions
		WHERE role_id = $1 AND permission_id = $2
	`

	var count int
	err = r.db.GetContext(ctx, &count, query, roleID, permissionID)
	if err != nil {
		return ErrDatabaseError
	}

	if count > 0 {
		// Permission is already assigned to role, nothing to do
		return nil
	}

	// Assign permission to role
	query = `
		INSERT INTO role_permissions (role_id, permission_id)
		VALUES ($1, $2)
	`

	_, err = r.db.ExecContext(ctx, query, roleID, permissionID)
	if err != nil {
		return ErrDatabaseError
	}

	return nil
}

// RemovePermissionFromRole removes a permission from a role
func (r *PostgresRoleRepository) RemovePermissionFromRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	// Check if role exists
	_, err := r.GetRoleByID(ctx, roleID)
	if err != nil {
		return err
	}

	// Check if permission exists
	_, err = r.GetPermissionByID(ctx, permissionID)
	if err != nil {
		return err
	}

	// Remove permission from role
	query := `
		DELETE FROM role_permissions
		WHERE role_id = $1 AND permission_id = $2
	`

	_, err = r.db.ExecContext(ctx, query, roleID, permissionID)
	if err != nil {
		return ErrDatabaseError
	}

	return nil
}

// GetRolePermissions retrieves all permissions for a role
func (r *PostgresRoleRepository) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]models.Permission, error) {
	query := `
		SELECT p.id, p.name, p.description, p.created_at, p.updated_at
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.name
	`

	var permissions []models.Permission
	err := r.db.SelectContext(ctx, &permissions, query, roleID)
	if err != nil {
		return nil, ErrDatabaseError
	}

	return permissions, nil
}

// AssignRoleToUser assigns a role to a user
func (r *PostgresRoleRepository) AssignRoleToUser(ctx context.Context, userRole *models.UserRole) error {
	// Check if user exists
	userRepo := NewPostgresUserRepository(r.db)
	_, err := userRepo.GetByID(ctx, userRole.UserID)
	if err != nil {
		return err
	}

	// Check if role exists
	_, err = r.GetRoleByID(ctx, userRole.RoleID)
	if err != nil {
		return err
	}

	// Check if user already has this role in this organization
	query := `
		SELECT COUNT(*)
		FROM user_roles
		WHERE user_id = $1 AND role_id = $2 AND organization_id = $3
	`

	var count int
	err = r.db.GetContext(ctx, &count, query, userRole.UserID, userRole.RoleID, userRole.OrganizationID)
	if err != nil {
		return ErrDatabaseError
	}

	if count > 0 {
		// User already has this role in this organization, nothing to do
		return nil
	}

	// Assign role to user
	query = `
		INSERT INTO user_roles (user_id, role_id, organization_id)
		VALUES ($1, $2, $3)
	`

	_, err = r.db.ExecContext(ctx, query, userRole.UserID, userRole.RoleID, userRole.OrganizationID)
	if err != nil {
		return ErrDatabaseError
	}

	return nil
}

// RemoveRoleFromUser removes a role from a user
func (r *PostgresRoleRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID, organizationID uuid.UUID) error {
	// Check if user exists
	userRepo := NewPostgresUserRepository(r.db)
	_, err := userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Check if role exists
	_, err = r.GetRoleByID(ctx, roleID)
	if err != nil {
		return err
	}

	// Remove role from user
	query := `
		DELETE FROM user_roles
		WHERE user_id = $1 AND role_id = $2 AND organization_id = $3
	`

	_, err = r.db.ExecContext(ctx, query, userID, roleID, organizationID)
	if err != nil {
		return ErrDatabaseError
	}

	return nil
}

// GetUserRoles retrieves all roles for a user
func (r *PostgresRoleRepository) GetUserRoles(ctx context.Context, userID, organizationID uuid.UUID) ([]models.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1 AND ur.organization_id = $2
		ORDER BY r.name
	`

	var roles []models.Role
	err := r.db.SelectContext(ctx, &roles, query, userID, organizationID)
	if err != nil {
		return nil, ErrDatabaseError
	}

	// Get permissions for each role
	for i := range roles {
		permissions, err := r.GetRolePermissions(ctx, roles[i].ID)
		if err != nil {
			return nil, err
		}
		roles[i].Permissions = permissions
	}

	return roles, nil
}

// HasPermission checks if a user has a specific permission
func (r *PostgresRoleRepository) HasPermission(ctx context.Context, userID, organizationID uuid.UUID, permissionName string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM user_roles ur
		JOIN role_permissions rp ON ur.role_id = rp.role_id
		JOIN permissions p ON rp.permission_id = p.id
		WHERE ur.user_id = $1 AND ur.organization_id = $2 AND p.name = $3
	`

	var count int
	err := r.db.GetContext(ctx, &count, query, userID, organizationID, permissionName)
	if err != nil {
		return false, ErrDatabaseError
	}

	return count > 0, nil
}
