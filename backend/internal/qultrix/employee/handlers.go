package employee

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/Jerinji2016/halooid/backend/pkg/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handlers provides HTTP handlers for employee management
type Handlers struct {
	service  Service
	validate *validator.Validate
}

// NewHandlers creates a new Handlers
func NewHandlers(service Service) *Handlers {
	return &Handlers{
		service:  service,
		validate: validator.New(),
	}
}

// Create handles the creation of a new employee
func (h *Handlers) Create(c echo.Context) error {
	// Get organization ID from path parameter
	orgIDParam := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	// Parse request body
	var req models.EmployeeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Create employee
	response, err := h.service.Create(c.Request().Context(), orgID, req)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		if errors.Is(err, repository.ErrEmployeeIDExists) {
			return echo.NewHTTPError(http.StatusConflict, "Employee ID already exists")
		}
		if errors.Is(err, repository.ErrUserAlreadyEmployee) {
			return echo.NewHTTPError(http.StatusConflict, "User is already an employee in this organization")
		}
		if errors.Is(err, repository.ErrManagerNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Manager not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create employee")
	}

	return c.JSON(http.StatusCreated, response)
}

// GetByID handles retrieving an employee by ID
func (h *Handlers) GetByID(c echo.Context) error {
	// Get employee ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid employee ID")
	}

	// Get employee
	response, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrEmployeeNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Employee not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve employee")
	}

	return c.JSON(http.StatusOK, response)
}

// GetByEmployeeID handles retrieving an employee by employee ID
func (h *Handlers) GetByEmployeeID(c echo.Context) error {
	// Get organization ID from path parameter
	orgIDParam := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	// Get employee ID from path parameter
	employeeID := c.Param("employee_id")
	if employeeID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Employee ID is required")
	}

	// Get employee
	response, err := h.service.GetByEmployeeID(c.Request().Context(), orgID, employeeID)
	if err != nil {
		if errors.Is(err, repository.ErrEmployeeNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Employee not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve employee")
	}

	return c.JSON(http.StatusOK, response)
}

// List handles retrieving employees based on filter parameters
func (h *Handlers) List(c echo.Context) error {
	// Get organization ID from path parameter
	orgIDParam := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	// Parse query parameters
	params := models.EmployeeListParams{
		OrganizationID: orgID,
		Department:     c.QueryParam("department"),
		Position:       c.QueryParam("position"),
		SortBy:         c.QueryParam("sort_by"),
		SortOrder:      c.QueryParam("sort_order"),
	}

	// Parse is_active parameter
	isActiveParam := c.QueryParam("is_active")
	if isActiveParam != "" {
		isActive, err := strconv.ParseBool(isActiveParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid is_active parameter")
		}
		params.IsActive = &isActive
	}

	// Parse page parameter
	pageParam := c.QueryParam("page")
	if pageParam != "" {
		page, err := strconv.Atoi(pageParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid page parameter")
		}
		params.Page = page
	} else {
		params.Page = 1
	}

	// Parse page_size parameter
	pageSizeParam := c.QueryParam("page_size")
	if pageSizeParam != "" {
		pageSize, err := strconv.Atoi(pageSizeParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid page_size parameter")
		}
		params.PageSize = pageSize
	} else {
		params.PageSize = 10
	}

	// Get employees
	employees, total, err := h.service.List(c.Request().Context(), params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve employees")
	}

	// Build response
	response := map[string]interface{}{
		"employees": employees,
		"pagination": map[string]interface{}{
			"total":      total,
			"page":       params.Page,
			"page_size":  params.PageSize,
			"total_pages": (total + params.PageSize - 1) / params.PageSize,
		},
	}

	return c.JSON(http.StatusOK, response)
}

// Update handles updating an employee
func (h *Handlers) Update(c echo.Context) error {
	// Get employee ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid employee ID")
	}

	// Parse request body
	var req models.EmployeeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Update employee
	response, err := h.service.Update(c.Request().Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrEmployeeNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Employee not found")
		}
		if errors.Is(err, repository.ErrEmployeeIDExists) {
			return echo.NewHTTPError(http.StatusConflict, "Employee ID already exists")
		}
		if errors.Is(err, repository.ErrManagerNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Manager not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update employee")
	}

	return c.JSON(http.StatusOK, response)
}

// Delete handles marking an employee as inactive
func (h *Handlers) Delete(c echo.Context) error {
	// Get employee ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid employee ID")
	}

	// Delete employee
	err = h.service.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrEmployeeNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Employee not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete employee")
	}

	return c.NoContent(http.StatusNoContent)
}

// RegisterRoutes registers the employee routes
func (h *Handlers) RegisterRoutes(g *echo.Group, rbacMiddleware *middleware.RBACMiddleware) {
	employeeGroup := g.Group("/employees")

	// Routes that require employee:read permission
	employeeGroup.GET("", h.List, rbacMiddleware.RequirePermission(middleware.PermissionEmployeeRead))
	employeeGroup.GET("/:id", h.GetByID, rbacMiddleware.RequirePermission(middleware.PermissionEmployeeRead))
	employeeGroup.GET("/by-employee-id/:employee_id", h.GetByEmployeeID, rbacMiddleware.RequirePermission(middleware.PermissionEmployeeRead))

	// Routes that require employee:write permission
	employeeGroup.POST("", h.Create, rbacMiddleware.RequirePermission(middleware.PermissionEmployeeWrite))
	employeeGroup.PUT("/:id", h.Update, rbacMiddleware.RequirePermission(middleware.PermissionEmployeeWrite))

	// Routes that require employee:delete permission
	employeeGroup.DELETE("/:id", h.Delete, rbacMiddleware.RequirePermission(middleware.PermissionEmployeeDelete))
}
