package organization

import (
	"context"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/google/uuid"
)

// Service provides organization management functionality
type Service interface {
	// Create creates a new organization
	Create(ctx context.Context, req models.OrganizationRequest) (*models.OrganizationResponse, error)
	
	// GetByID retrieves an organization by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationResponse, error)
	
	// GetByName retrieves an organization by name
	GetByName(ctx context.Context, name string) (*models.OrganizationResponse, error)
	
	// List retrieves all organizations
	List(ctx context.Context) ([]models.OrganizationResponse, error)
	
	// Update updates an organization
	Update(ctx context.Context, id uuid.UUID, req models.OrganizationRequest) (*models.OrganizationResponse, error)
	
	// Delete marks an organization as inactive
	Delete(ctx context.Context, id uuid.UUID) error
	
	// AddUser adds a user to an organization
	AddUser(ctx context.Context, req models.AddUserToOrganizationRequest) error
	
	// RemoveUser removes a user from an organization
	RemoveUser(ctx context.Context, organizationID, userID uuid.UUID) error
	
	// GetUsers retrieves all users in an organization
	GetUsers(ctx context.Context, organizationID uuid.UUID) ([]models.UserResponse, error)
	
	// GetOrganizations retrieves all organizations a user belongs to
	GetOrganizations(ctx context.Context, userID uuid.UUID) ([]models.OrganizationResponse, error)
}

// serviceImpl implements the Service interface
type serviceImpl struct {
	orgRepo repository.OrganizationRepository
}

// NewService creates a new organization service
func NewService(orgRepo repository.OrganizationRepository) Service {
	return &serviceImpl{
		orgRepo: orgRepo,
	}
}

// Create creates a new organization
func (s *serviceImpl) Create(ctx context.Context, req models.OrganizationRequest) (*models.OrganizationResponse, error) {
	organization := models.NewOrganization(req)
	
	err := s.orgRepo.Create(ctx, organization)
	if err != nil {
		return nil, err
	}
	
	response := organization.ToResponse()
	return &response, nil
}

// GetByID retrieves an organization by ID
func (s *serviceImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationResponse, error) {
	organization, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	response := organization.ToResponse()
	return &response, nil
}

// GetByName retrieves an organization by name
func (s *serviceImpl) GetByName(ctx context.Context, name string) (*models.OrganizationResponse, error) {
	organization, err := s.orgRepo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	
	response := organization.ToResponse()
	return &response, nil
}

// List retrieves all organizations
func (s *serviceImpl) List(ctx context.Context) ([]models.OrganizationResponse, error) {
	organizations, err := s.orgRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	
	responses := make([]models.OrganizationResponse, 0, len(organizations))
	for _, org := range organizations {
		responses = append(responses, org.ToResponse())
	}
	
	return responses, nil
}

// Update updates an organization
func (s *serviceImpl) Update(ctx context.Context, id uuid.UUID, req models.OrganizationRequest) (*models.OrganizationResponse, error) {
	organization, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	organization.Name = req.Name
	organization.Description = req.Description
	
	err = s.orgRepo.Update(ctx, organization)
	if err != nil {
		return nil, err
	}
	
	response := organization.ToResponse()
	return &response, nil
}

// Delete marks an organization as inactive
func (s *serviceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return s.orgRepo.Delete(ctx, id)
}

// AddUser adds a user to an organization
func (s *serviceImpl) AddUser(ctx context.Context, req models.AddUserToOrganizationRequest) error {
	organizationUser := &models.OrganizationUser{
		OrganizationID: req.OrganizationID,
		UserID:         req.UserID,
	}
	
	return s.orgRepo.AddUser(ctx, organizationUser)
}

// RemoveUser removes a user from an organization
func (s *serviceImpl) RemoveUser(ctx context.Context, organizationID, userID uuid.UUID) error {
	return s.orgRepo.RemoveUser(ctx, organizationID, userID)
}

// GetUsers retrieves all users in an organization
func (s *serviceImpl) GetUsers(ctx context.Context, organizationID uuid.UUID) ([]models.UserResponse, error) {
	users, err := s.orgRepo.GetUsers(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	
	responses := make([]models.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}
	
	return responses, nil
}

// GetOrganizations retrieves all organizations a user belongs to
func (s *serviceImpl) GetOrganizations(ctx context.Context, userID uuid.UUID) ([]models.OrganizationResponse, error) {
	organizations, err := s.orgRepo.GetOrganizations(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	responses := make([]models.OrganizationResponse, 0, len(organizations))
	for _, org := range organizations {
		responses = append(responses, org.ToResponse())
	}
	
	return responses, nil
}
