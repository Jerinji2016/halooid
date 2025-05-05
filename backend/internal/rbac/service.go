package rbac

import (
	"context"
	"errors"
	"time"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/google/uuid"
)

// Common errors
var (
	ErrRoleNotFound        = errors.New("role not found")
	ErrPermissionNotFound  = errors.New("permission not found")
	ErrRoleNameExists      = errors.New("role name already exists")
	ErrPermissionNameExists = errors.New("permission name already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidPermission   = errors.New("invalid permission")
	ErrPermissionDenied    = errors.New("permission denied")
)

// Service provides RBAC functionality
type Service interface {
	// CreateRole creates a new role
	CreateRole(ctx context.Context, req models.RoleRequest) (*models.RoleResponse, error)
	
	// GetRoleByID retrieves a role by ID
	GetRoleByID(ctx context.Context, id uuid.UUID) (*models.RoleResponse, error)
	
	// GetRoleByName retrieves a role by name
	GetRoleByName(ctx context.Context, name string) (*models.RoleResponse, error)
	
	// ListRoles retrieves all roles
	ListRoles(ctx context.Context) ([]models.RoleResponse, error)
	
	// UpdateRole updates a role
	UpdateRole(ctx context.Context, id uuid.UUID, req models.RoleRequest) (*models.RoleResponse, error)
	
	// DeleteRole deletes a role
	DeleteRole(ctx context.Context, id uuid.UUID) error
	
	// CreatePermission creates a new permission
	CreatePermission(ctx context.Context, req models.PermissionRequest) (*models.PermissionResponse, error)
	
	// GetPermissionByID retrieves a permission by ID
	GetPermissionByID(ctx context.Context, id uuid.UUID) (*models.PermissionResponse, error)
	
	// GetPermissionByName retrieves a permission by name
	GetPermissionByName(ctx context.Context, name string) (*models.PermissionResponse, error)
	
	// ListPermissions retrieves all permissions
	ListPermissions(ctx context.Context) ([]models.PermissionResponse, error)
	
	// UpdatePermission updates a permission
	UpdatePermission(ctx context.Context, id uuid.UUID, req models.PermissionRequest) (*models.PermissionResponse, error)
	
	// DeletePermission deletes a permission
	DeletePermission(ctx context.Context, id uuid.UUID) error
	
	// AssignRoleToUser assigns a role to a user
	AssignRoleToUser(ctx context.Context, req models.AssignRoleRequest) error
	
	// RemoveRoleFromUser removes a role from a user
	RemoveRoleFromUser(ctx context.Context, userID, roleID, organizationID uuid.UUID) error
	
	// GetUserRoles retrieves all roles for a user
	GetUserRoles(ctx context.Context, userID, organizationID uuid.UUID) ([]models.RoleResponse, error)
	
	// HasPermission checks if a user has a specific permission
	HasPermission(ctx context.Context, userID, organizationID uuid.UUID, permissionName string) (bool, error)
}

// serviceImpl implements the Service interface
type serviceImpl struct {
	roleRepo repository.RoleRepository
	userRepo repository.UserRepository
}

// NewService creates a new RBAC service
func NewService(roleRepo repository.RoleRepository, userRepo repository.UserRepository) Service {
	return &serviceImpl{
		roleRepo: roleRepo,
		userRepo: userRepo,
	}
}

// CreateRole creates a new role
func (s *serviceImpl) CreateRole(ctx context.Context, req models.RoleRequest) (*models.RoleResponse, error) {
	// Create role
	now := time.Now()
	role := &models.Role{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Create role
	err := s.roleRepo.CreateRole(ctx, role)
	if err != nil {
		if errors.Is(err, repository.ErrRoleNameExists) {
			return nil, ErrRoleNameExists
		}
		return nil, err
	}

	// Assign permissions to role
	for _, permissionName := range req.Permissions {
		permission, err := s.roleRepo.GetPermissionByName(ctx, permissionName)
		if err != nil {
			if errors.Is(err, repository.ErrPermissionNotFound) {
				return nil, ErrPermissionNotFound
			}
			return nil, err
		}

		err = s.roleRepo.AssignPermissionToRole(ctx, role.ID, permission.ID)
		if err != nil {
			return nil, err
		}
	}

	// Get role with permissions
	role, err = s.roleRepo.GetRoleByID(ctx, role.ID)
	if err != nil {
		return nil, err
	}

	// Convert to response
	response := role.ToResponse()
	return &response, nil
}

// GetRoleByID retrieves a role by ID
func (s *serviceImpl) GetRoleByID(ctx context.Context, id uuid.UUID) (*models.RoleResponse, error) {
	role, err := s.roleRepo.GetRoleByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrRoleNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	response := role.ToResponse()
	return &response, nil
}

// GetRoleByName retrieves a role by name
func (s *serviceImpl) GetRoleByName(ctx context.Context, name string) (*models.RoleResponse, error) {
	role, err := s.roleRepo.GetRoleByName(ctx, name)
	if err != nil {
		if errors.Is(err, repository.ErrRoleNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	response := role.ToResponse()
	return &response, nil
}

// ListRoles retrieves all roles
func (s *serviceImpl) ListRoles(ctx context.Context) ([]models.RoleResponse, error) {
	roles, err := s.roleRepo.ListRoles(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]models.RoleResponse, 0, len(roles))
	for _, role := range roles {
		responses = append(responses, role.ToResponse())
	}

	return responses, nil
}

// UpdateRole updates a role
func (s *serviceImpl) UpdateRole(ctx context.Context, id uuid.UUID, req models.RoleRequest) (*models.RoleResponse, error) {
	// Get role
	role, err := s.roleRepo.GetRoleByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrRoleNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	// Update role
	role.Name = req.Name
	role.Description = req.Description
	role.UpdatedAt = time.Now()

	err = s.roleRepo.UpdateRole(ctx, role)
	if err != nil {
		if errors.Is(err, repository.ErrRoleNameExists) {
			return nil, ErrRoleNameExists
		}
		return nil, err
	}

	// Get current permissions
	currentPermissions, err := s.roleRepo.GetRolePermissions(ctx, id)
	if err != nil {
		return nil, err
	}

	// Create maps for easier comparison
	currentPermissionMap := make(map[string]uuid.UUID)
	for _, p := range currentPermissions {
		currentPermissionMap[p.Name] = p.ID
	}

	newPermissionMap := make(map[string]bool)
	for _, p := range req.Permissions {
		newPermissionMap[p] = true
	}

	// Remove permissions that are no longer needed
	for _, p := range currentPermissions {
		if !newPermissionMap[p.Name] {
			err = s.roleRepo.RemovePermissionFromRole(ctx, id, p.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	// Add new permissions
	for _, permissionName := range req.Permissions {
		if _, exists := currentPermissionMap[permissionName]; !exists {
			permission, err := s.roleRepo.GetPermissionByName(ctx, permissionName)
			if err != nil {
				if errors.Is(err, repository.ErrPermissionNotFound) {
					return nil, ErrPermissionNotFound
				}
				return nil, err
			}

			err = s.roleRepo.AssignPermissionToRole(ctx, id, permission.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	// Get updated role
	role, err = s.roleRepo.GetRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := role.ToResponse()
	return &response, nil
}

// DeleteRole deletes a role
func (s *serviceImpl) DeleteRole(ctx context.Context, id uuid.UUID) error {
	err := s.roleRepo.DeleteRole(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrRoleNotFound) {
			return ErrRoleNotFound
		}
		return err
	}

	return nil
}

// CreatePermission creates a new permission
func (s *serviceImpl) CreatePermission(ctx context.Context, req models.PermissionRequest) (*models.PermissionResponse, error) {
	// Create permission
	now := time.Now()
	permission := &models.Permission{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Create permission
	err := s.roleRepo.CreatePermission(ctx, permission)
	if err != nil {
		if errors.Is(err, repository.ErrPermissionNameExists) {
			return nil, ErrPermissionNameExists
		}
		return nil, err
	}

	// Convert to response
	response := permission.ToResponse()
	return &response, nil
}

// GetPermissionByID retrieves a permission by ID
func (s *serviceImpl) GetPermissionByID(ctx context.Context, id uuid.UUID) (*models.PermissionResponse, error) {
	permission, err := s.roleRepo.GetPermissionByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrPermissionNotFound) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}

	response := permission.ToResponse()
	return &response, nil
}

// GetPermissionByName retrieves a permission by name
func (s *serviceImpl) GetPermissionByName(ctx context.Context, name string) (*models.PermissionResponse, error) {
	permission, err := s.roleRepo.GetPermissionByName(ctx, name)
	if err != nil {
		if errors.Is(err, repository.ErrPermissionNotFound) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}

	response := permission.ToResponse()
	return &response, nil
}

// ListPermissions retrieves all permissions
func (s *serviceImpl) ListPermissions(ctx context.Context) ([]models.PermissionResponse, error) {
	permissions, err := s.roleRepo.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]models.PermissionResponse, 0, len(permissions))
	for _, permission := range permissions {
		responses = append(responses, permission.ToResponse())
	}

	return responses, nil
}

// UpdatePermission updates a permission
func (s *serviceImpl) UpdatePermission(ctx context.Context, id uuid.UUID, req models.PermissionRequest) (*models.PermissionResponse, error) {
	// Get permission
	permission, err := s.roleRepo.GetPermissionByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrPermissionNotFound) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}

	// Update permission
	permission.Name = req.Name
	permission.Description = req.Description
	permission.UpdatedAt = time.Now()

	err = s.roleRepo.UpdatePermission(ctx, permission)
	if err != nil {
		if errors.Is(err, repository.ErrPermissionNameExists) {
			return nil, ErrPermissionNameExists
		}
		return nil, err
	}

	// Get updated permission
	permission, err = s.roleRepo.GetPermissionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := permission.ToResponse()
	return &response, nil
}

// DeletePermission deletes a permission
func (s *serviceImpl) DeletePermission(ctx context.Context, id uuid.UUID) error {
	err := s.roleRepo.DeletePermission(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrPermissionNotFound) {
			return ErrPermissionNotFound
		}
		return err
	}

	return nil
}

// AssignRoleToUser assigns a role to a user
func (s *serviceImpl) AssignRoleToUser(ctx context.Context, req models.AssignRoleRequest) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Check if role exists
	_, err = s.roleRepo.GetRoleByID(ctx, req.RoleID)
	if err != nil {
		if errors.Is(err, repository.ErrRoleNotFound) {
			return ErrRoleNotFound
		}
		return err
	}

	// Assign role to user
	userRole := &models.UserRole{
		UserID:         req.UserID,
		RoleID:         req.RoleID,
		OrganizationID: req.OrganizationID,
	}

	err = s.roleRepo.AssignRoleToUser(ctx, userRole)
	if err != nil {
		return err
	}

	return nil
}

// RemoveRoleFromUser removes a role from a user
func (s *serviceImpl) RemoveRoleFromUser(ctx context.Context, userID, roleID, organizationID uuid.UUID) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Check if role exists
	_, err = s.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		if errors.Is(err, repository.ErrRoleNotFound) {
			return ErrRoleNotFound
		}
		return err
	}

	// Remove role from user
	err = s.roleRepo.RemoveRoleFromUser(ctx, userID, roleID, organizationID)
	if err != nil {
		return err
	}

	return nil
}

// GetUserRoles retrieves all roles for a user
func (s *serviceImpl) GetUserRoles(ctx context.Context, userID, organizationID uuid.UUID) ([]models.RoleResponse, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Get user roles
	roles, err := s.roleRepo.GetUserRoles(ctx, userID, organizationID)
	if err != nil {
		return nil, err
	}

	responses := make([]models.RoleResponse, 0, len(roles))
	for _, role := range roles {
		responses = append(responses, role.ToResponse())
	}

	return responses, nil
}

// HasPermission checks if a user has a specific permission
func (s *serviceImpl) HasPermission(ctx context.Context, userID, organizationID uuid.UUID, permissionName string) (bool, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return false, ErrUserNotFound
		}
		return false, err
	}

	// Check if permission exists
	_, err = s.roleRepo.GetPermissionByName(ctx, permissionName)
	if err != nil {
		if errors.Is(err, repository.ErrPermissionNotFound) {
			return false, ErrInvalidPermission
		}
		return false, err
	}

	// Check if user has permission
	hasPermission, err := s.roleRepo.HasPermission(ctx, userID, organizationID, permissionName)
	if err != nil {
		return false, err
	}

	return hasPermission, nil
}
