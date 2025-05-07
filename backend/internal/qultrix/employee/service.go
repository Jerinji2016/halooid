package employee

import (
	"context"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/google/uuid"
)

// Service provides employee management functionality
type Service interface {
	// Create creates a new employee
	Create(ctx context.Context, organizationID uuid.UUID, req models.EmployeeRequest) (*models.EmployeeResponse, error)
	
	// GetByID retrieves an employee by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.EmployeeResponse, error)
	
	// GetByEmployeeID retrieves an employee by employee ID
	GetByEmployeeID(ctx context.Context, organizationID uuid.UUID, employeeID string) (*models.EmployeeResponse, error)
	
	// GetByUserID retrieves an employee by user ID
	GetByUserID(ctx context.Context, organizationID, userID uuid.UUID) (*models.EmployeeResponse, error)
	
	// List retrieves employees based on filter parameters
	List(ctx context.Context, params models.EmployeeListParams) ([]models.EmployeeResponse, int, error)
	
	// Update updates an employee
	Update(ctx context.Context, id uuid.UUID, req models.EmployeeRequest) (*models.EmployeeResponse, error)
	
	// Delete marks an employee as inactive
	Delete(ctx context.Context, id uuid.UUID) error
}

// serviceImpl implements the Service interface
type serviceImpl struct {
	employeeRepo repository.EmployeeRepository
	userRepo     repository.UserRepository
}

// NewService creates a new employee service
func NewService(employeeRepo repository.EmployeeRepository, userRepo repository.UserRepository) Service {
	return &serviceImpl{
		employeeRepo: employeeRepo,
		userRepo:     userRepo,
	}
}

// Create creates a new employee
func (s *serviceImpl) Create(ctx context.Context, organizationID uuid.UUID, req models.EmployeeRequest) (*models.EmployeeResponse, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	employee := models.NewEmployee(req, organizationID)
	
	err = s.employeeRepo.Create(ctx, employee)
	if err != nil {
		return nil, err
	}
	
	// Get the complete employee record with user and manager information
	createdEmployee, err := s.employeeRepo.GetByID(ctx, employee.ID)
	if err != nil {
		return nil, err
	}
	
	response := createdEmployee.ToResponse()
	return &response, nil
}

// GetByID retrieves an employee by ID
func (s *serviceImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.EmployeeResponse, error) {
	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	response := employee.ToResponse()
	return &response, nil
}

// GetByEmployeeID retrieves an employee by employee ID
func (s *serviceImpl) GetByEmployeeID(ctx context.Context, organizationID uuid.UUID, employeeID string) (*models.EmployeeResponse, error) {
	employee, err := s.employeeRepo.GetByEmployeeID(ctx, organizationID, employeeID)
	if err != nil {
		return nil, err
	}
	
	response := employee.ToResponse()
	return &response, nil
}

// GetByUserID retrieves an employee by user ID
func (s *serviceImpl) GetByUserID(ctx context.Context, organizationID, userID uuid.UUID) (*models.EmployeeResponse, error) {
	employee, err := s.employeeRepo.GetByUserID(ctx, organizationID, userID)
	if err != nil {
		return nil, err
	}
	
	response := employee.ToResponse()
	return &response, nil
}

// List retrieves employees based on filter parameters
func (s *serviceImpl) List(ctx context.Context, params models.EmployeeListParams) ([]models.EmployeeResponse, int, error) {
	employees, total, err := s.employeeRepo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	
	responses := make([]models.EmployeeResponse, 0, len(employees))
	for _, emp := range employees {
		responses = append(responses, emp.ToResponse())
	}
	
	return responses, total, nil
}

// Update updates an employee
func (s *serviceImpl) Update(ctx context.Context, id uuid.UUID, req models.EmployeeRequest) (*models.EmployeeResponse, error) {
	// Check if employee exists
	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Update employee fields
	employee.EmployeeID = req.EmployeeID
	employee.Department = req.Department
	employee.Position = req.Position
	employee.HireDate = req.HireDate
	employee.ManagerID = req.ManagerID
	employee.Salary = req.Salary
	
	err = s.employeeRepo.Update(ctx, employee)
	if err != nil {
		return nil, err
	}
	
	// Get the updated employee record
	updatedEmployee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	response := updatedEmployee.ToResponse()
	return &response, nil
}

// Delete marks an employee as inactive
func (s *serviceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if employee exists
	_, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	return s.employeeRepo.Delete(ctx, id)
}
