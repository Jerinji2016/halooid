package employee

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEmployeeRepository is a mock implementation of the EmployeeRepository interface
type MockEmployeeRepository struct {
	mock.Mock
}

func (m *MockEmployeeRepository) Create(ctx context.Context, employee *models.Employee) error {
	args := m.Called(ctx, employee)
	return args.Error(0)
}

func (m *MockEmployeeRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) GetByEmployeeID(ctx context.Context, organizationID uuid.UUID, employeeID string) (*models.Employee, error) {
	args := m.Called(ctx, organizationID, employeeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) GetByUserID(ctx context.Context, organizationID, userID uuid.UUID) (*models.Employee, error) {
	args := m.Called(ctx, organizationID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) List(ctx context.Context, params models.EmployeeListParams) ([]models.Employee, int, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]models.Employee), args.Int(1), args.Error(2)
}

func (m *MockEmployeeRepository) Update(ctx context.Context, employee *models.Employee) error {
	args := m.Called(ctx, employee)
	return args.Error(0)
}

func (m *MockEmployeeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEmployeeRepository) GetDB() *sqlx.DB {
	args := m.Called()
	return args.Get(0).(*sqlx.DB)
}

// MockUserRepository is a mock implementation of the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) GetDB() *sqlx.DB {
	args := m.Called()
	return args.Get(0).(*sqlx.DB)
}

func TestCreate(t *testing.T) {
	// Setup
	ctx := context.Background()
	employeeRepo := new(MockEmployeeRepository)
	userRepo := new(MockUserRepository)
	service := NewService(employeeRepo, userRepo)

	orgID := uuid.New()
	userID := uuid.New()
	employeeID := "EMP001"
	hireDate := time.Now()

	user := &models.User{
		ID:        userID,
		Email:     "john.doe@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	req := models.EmployeeRequest{
		UserID:     userID,
		EmployeeID: employeeID,
		Department: "Engineering",
		Position:   "Software Engineer",
		HireDate:   hireDate,
		Salary:     75000,
	}

	employee := &models.Employee{
		ID:             uuid.New(),
		OrganizationID: orgID,
		UserID:         userID,
		EmployeeID:     employeeID,
		Department:     "Engineering",
		Position:       "Software Engineer",
		HireDate:       hireDate,
		Salary:         75000,
		IsActive:       true,
		User:           user,
	}

	// Expectations
	userRepo.On("GetByID", ctx, userID).Return(user, nil)
	employeeRepo.On("Create", ctx, mock.AnythingOfType("*models.Employee")).Return(nil)
	employeeRepo.On("GetByID", ctx, mock.AnythingOfType("uuid.UUID")).Return(employee, nil)

	// Execute
	response, err := service.Create(ctx, orgID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, employee.ID, response.ID)
	assert.Equal(t, employee.EmployeeID, response.EmployeeID)
	assert.Equal(t, employee.Department, response.Department)
	assert.Equal(t, employee.Position, response.Position)
	assert.Equal(t, employee.Salary, response.Salary)
	assert.Equal(t, user.ID, response.User.ID)
	assert.Equal(t, user.Email, response.User.Email)

	// Verify expectations
	userRepo.AssertExpectations(t)
	employeeRepo.AssertExpectations(t)
}

func TestGetByID(t *testing.T) {
	// Setup
	ctx := context.Background()
	employeeRepo := new(MockEmployeeRepository)
	userRepo := new(MockUserRepository)
	service := NewService(employeeRepo, userRepo)

	id := uuid.New()
	orgID := uuid.New()
	userID := uuid.New()
	employeeID := "EMP001"
	hireDate := time.Now()

	user := &models.User{
		ID:        userID,
		Email:     "john.doe@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	employee := &models.Employee{
		ID:             id,
		OrganizationID: orgID,
		UserID:         userID,
		EmployeeID:     employeeID,
		Department:     "Engineering",
		Position:       "Software Engineer",
		HireDate:       hireDate,
		Salary:         75000,
		IsActive:       true,
		User:           user,
	}

	// Expectations
	employeeRepo.On("GetByID", ctx, id).Return(employee, nil)

	// Execute
	response, err := service.GetByID(ctx, id)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, employee.ID, response.ID)
	assert.Equal(t, employee.EmployeeID, response.EmployeeID)
	assert.Equal(t, employee.Department, response.Department)
	assert.Equal(t, employee.Position, response.Position)
	assert.Equal(t, employee.Salary, response.Salary)
	assert.Equal(t, user.ID, response.User.ID)
	assert.Equal(t, user.Email, response.User.Email)

	// Verify expectations
	employeeRepo.AssertExpectations(t)
}

func TestGetByID_NotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	employeeRepo := new(MockEmployeeRepository)
	userRepo := new(MockUserRepository)
	service := NewService(employeeRepo, userRepo)

	id := uuid.New()

	// Expectations
	employeeRepo.On("GetByID", ctx, id).Return(nil, repository.ErrEmployeeNotFound)

	// Execute
	response, err := service.GetByID(ctx, id)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, repository.ErrEmployeeNotFound, err)

	// Verify expectations
	employeeRepo.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	// Setup
	ctx := context.Background()
	employeeRepo := new(MockEmployeeRepository)
	userRepo := new(MockUserRepository)
	service := NewService(employeeRepo, userRepo)

	id := uuid.New()
	orgID := uuid.New()
	userID := uuid.New()
	employeeID := "EMP001"
	hireDate := time.Now()

	user := &models.User{
		ID:        userID,
		Email:     "john.doe@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	employee := &models.Employee{
		ID:             id,
		OrganizationID: orgID,
		UserID:         userID,
		EmployeeID:     employeeID,
		Department:     "Engineering",
		Position:       "Software Engineer",
		HireDate:       hireDate,
		Salary:         75000,
		IsActive:       true,
		User:           user,
	}

	req := models.EmployeeRequest{
		UserID:     userID,
		EmployeeID: employeeID,
		Department: "Product",
		Position:   "Product Manager",
		HireDate:   hireDate,
		Salary:     85000,
	}

	updatedEmployee := &models.Employee{
		ID:             id,
		OrganizationID: orgID,
		UserID:         userID,
		EmployeeID:     employeeID,
		Department:     "Product",
		Position:       "Product Manager",
		HireDate:       hireDate,
		Salary:         85000,
		IsActive:       true,
		User:           user,
	}

	// Expectations
	employeeRepo.On("GetByID", ctx, id).Return(employee, nil)
	employeeRepo.On("Update", ctx, mock.AnythingOfType("*models.Employee")).Return(nil)
	employeeRepo.On("GetByID", ctx, id).Return(updatedEmployee, nil)

	// Execute
	response, err := service.Update(ctx, id, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, updatedEmployee.ID, response.ID)
	assert.Equal(t, updatedEmployee.EmployeeID, response.EmployeeID)
	assert.Equal(t, updatedEmployee.Department, response.Department)
	assert.Equal(t, updatedEmployee.Position, response.Position)
	assert.Equal(t, updatedEmployee.Salary, response.Salary)
	assert.Equal(t, user.ID, response.User.ID)
	assert.Equal(t, user.Email, response.User.Email)

	// Verify expectations
	employeeRepo.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	// Setup
	ctx := context.Background()
	employeeRepo := new(MockEmployeeRepository)
	userRepo := new(MockUserRepository)
	service := NewService(employeeRepo, userRepo)

	id := uuid.New()
	orgID := uuid.New()
	userID := uuid.New()
	employeeID := "EMP001"
	hireDate := time.Now()

	user := &models.User{
		ID:        userID,
		Email:     "john.doe@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	employee := &models.Employee{
		ID:             id,
		OrganizationID: orgID,
		UserID:         userID,
		EmployeeID:     employeeID,
		Department:     "Engineering",
		Position:       "Software Engineer",
		HireDate:       hireDate,
		Salary:         75000,
		IsActive:       true,
		User:           user,
	}

	// Expectations
	employeeRepo.On("GetByID", ctx, id).Return(employee, nil)
	employeeRepo.On("Delete", ctx, id).Return(nil)

	// Execute
	err := service.Delete(ctx, id)

	// Assert
	assert.NoError(t, err)

	// Verify expectations
	employeeRepo.AssertExpectations(t)
}
