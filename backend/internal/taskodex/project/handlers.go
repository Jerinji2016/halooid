package project

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/Jerinji2016/halooid/backend/pkg/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handlers provides HTTP handlers for project management
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

// Create handles the creation of a new project
func (h *Handlers) Create(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Get organization ID from path parameter
	orgIDParam := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	// Parse request body
	var req models.ProjectRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Set organization ID from path parameter
	req.OrganizationID = orgID

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Create project
	response, err := h.service.Create(c.Request().Context(), req, userID)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNameExists) {
			return echo.NewHTTPError(http.StatusConflict, "Project name already exists in this organization")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create project")
	}

	return c.JSON(http.StatusCreated, response)
}

// GetByID handles retrieving a project by ID
func (h *Handlers) GetByID(c echo.Context) error {
	// Get project ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid project ID")
	}

	// Get project
	response, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Project not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve project")
	}

	return c.JSON(http.StatusOK, response)
}

// List handles retrieving projects based on filter parameters
func (h *Handlers) List(c echo.Context) error {
	// Get organization ID from path parameter
	orgIDParam := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	// Parse query parameters
	params := models.ProjectListParams{
		OrganizationID: orgID,
	}
	
	// Parse status parameter
	statusParam := c.QueryParam("status")
	if statusParam != "" {
		status := models.ProjectStatus(statusParam)
		params.Status = &status
	}
	
	// Parse created_by parameter
	createdByParam := c.QueryParam("created_by")
	if createdByParam != "" {
		createdBy, err := uuid.Parse(createdByParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid created_by parameter")
		}
		params.CreatedBy = &createdBy
	}
	
	// Parse search parameter
	searchParam := c.QueryParam("search")
	if searchParam != "" {
		params.SearchTerm = &searchParam
	}
	
	// Parse sort_by parameter
	sortByParam := c.QueryParam("sort_by")
	if sortByParam != "" {
		params.SortBy = sortByParam
	}
	
	// Parse sort_order parameter
	sortOrderParam := c.QueryParam("sort_order")
	if sortOrderParam != "" {
		params.SortOrder = sortOrderParam
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
		params.PageSize = 20
	}
	
	// Get projects
	projects, total, err := h.service.List(c.Request().Context(), params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve projects")
	}
	
	// Build response
	response := map[string]interface{}{
		"projects": projects,
		"pagination": map[string]interface{}{
			"total":       total,
			"page":        params.Page,
			"page_size":   params.PageSize,
			"total_pages": (total + params.PageSize - 1) / params.PageSize,
		},
	}
	
	return c.JSON(http.StatusOK, response)
}

// Update handles updating a project
func (h *Handlers) Update(c echo.Context) error {
	// Get project ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid project ID")
	}
	
	// Get organization ID from path parameter
	orgIDParam := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}
	
	// Parse request body
	var req models.ProjectRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	// Set organization ID from path parameter
	req.OrganizationID = orgID
	
	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	// Update project
	response, err := h.service.Update(c.Request().Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Project not found")
		}
		if errors.Is(err, repository.ErrProjectNameExists) {
			return echo.NewHTTPError(http.StatusConflict, "Project name already exists in this organization")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update project")
	}
	
	return c.JSON(http.StatusOK, response)
}

// Delete handles deleting a project
func (h *Handlers) Delete(c echo.Context) error {
	// Get project ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid project ID")
	}
	
	// Delete project
	err = h.service.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Project not found")
		}
		if errors.Is(err, errors.New("cannot delete project with tasks")) {
			return echo.NewHTTPError(http.StatusBadRequest, "Cannot delete project with tasks")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete project")
	}
	
	return c.NoContent(http.StatusNoContent)
}

// GetTasks handles retrieving all tasks for a project
func (h *Handlers) GetTasks(c echo.Context) error {
	// Get project ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid project ID")
	}
	
	// Parse query parameters for task filtering
	params := models.TaskListParams{}
	
	// Parse status parameter
	statusParam := c.QueryParam("status")
	if statusParam != "" {
		status := models.TaskStatus(statusParam)
		params.Status = &status
	}
	
	// Parse priority parameter
	priorityParam := c.QueryParam("priority")
	if priorityParam != "" {
		priority := models.TaskPriority(priorityParam)
		params.Priority = &priority
	}
	
	// Parse assigned_to parameter
	assignedToParam := c.QueryParam("assigned_to")
	if assignedToParam != "" {
		assignedTo, err := uuid.Parse(assignedToParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid assigned_to parameter")
		}
		params.AssignedTo = &assignedTo
	}
	
	// Parse due_before parameter
	dueBeforeParam := c.QueryParam("due_before")
	if dueBeforeParam != "" {
		dueBefore, err := time.Parse(time.RFC3339, dueBeforeParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid due_before parameter")
		}
		params.DueBefore = &dueBefore
	}
	
	// Parse due_after parameter
	dueAfterParam := c.QueryParam("due_after")
	if dueAfterParam != "" {
		dueAfter, err := time.Parse(time.RFC3339, dueAfterParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid due_after parameter")
		}
		params.DueAfter = &dueAfter
	}
	
	// Parse search parameter
	searchParam := c.QueryParam("search")
	if searchParam != "" {
		params.SearchTerm = &searchParam
	}
	
	// Parse sort_by parameter
	sortByParam := c.QueryParam("sort_by")
	if sortByParam != "" {
		params.SortBy = sortByParam
	}
	
	// Parse sort_order parameter
	sortOrderParam := c.QueryParam("sort_order")
	if sortOrderParam != "" {
		params.SortOrder = sortOrderParam
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
		params.PageSize = 20
	}
	
	// Get tasks
	tasks, total, err := h.service.GetTasks(c.Request().Context(), id, params)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Project not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve tasks")
	}
	
	// Build response
	response := map[string]interface{}{
		"tasks": tasks,
		"pagination": map[string]interface{}{
			"total":       total,
			"page":        params.Page,
			"page_size":   params.PageSize,
			"total_pages": (total + params.PageSize - 1) / params.PageSize,
		},
	}
	
	return c.JSON(http.StatusOK, response)
}

// AddTask handles adding a task to a project
func (h *Handlers) AddTask(c echo.Context) error {
	// Get project ID from path parameter
	projectIDParam := c.Param("id")
	projectID, err := uuid.Parse(projectIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid project ID")
	}
	
	// Parse request body
	var req struct {
		TaskID uuid.UUID `json:"task_id" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	// Add task to project
	err = h.service.AddTask(c.Request().Context(), projectID, req.TaskID)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Project not found")
		}
		if errors.Is(err, repository.ErrTaskNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Task not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add task to project")
	}
	
	return c.NoContent(http.StatusNoContent)
}

// RemoveTask handles removing a task from a project
func (h *Handlers) RemoveTask(c echo.Context) error {
	// Get project ID from path parameter
	projectIDParam := c.Param("id")
	projectID, err := uuid.Parse(projectIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid project ID")
	}
	
	// Get task ID from path parameter
	taskIDParam := c.Param("task_id")
	taskID, err := uuid.Parse(taskIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}
	
	// Remove task from project
	err = h.service.RemoveTask(c.Request().Context(), projectID, taskID)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Project not found")
		}
		if errors.Is(err, repository.ErrTaskNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Task not found")
		}
		if errors.Is(err, errors.New("task does not belong to the project")) {
			return echo.NewHTTPError(http.StatusBadRequest, "Task does not belong to the project")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove task from project")
	}
	
	return c.NoContent(http.StatusNoContent)
}

// RegisterRoutes registers the project routes
func (h *Handlers) RegisterRoutes(g *echo.Group, rbacMiddleware *middleware.RBACMiddleware) {
	projectGroup := g.Group("/projects")
	
	// Routes that require project:read permission
	projectGroup.GET("", h.List, rbacMiddleware.RequirePermission(middleware.PermissionProjectRead))
	projectGroup.GET("/:id", h.GetByID, rbacMiddleware.RequirePermission(middleware.PermissionProjectRead))
	projectGroup.GET("/:id/tasks", h.GetTasks, rbacMiddleware.RequirePermission(middleware.PermissionProjectRead))
	
	// Routes that require project:write permission
	projectGroup.POST("", h.Create, rbacMiddleware.RequirePermission(middleware.PermissionProjectWrite))
	projectGroup.PUT("/:id", h.Update, rbacMiddleware.RequirePermission(middleware.PermissionProjectWrite))
	projectGroup.POST("/:id/tasks", h.AddTask, rbacMiddleware.RequirePermission(middleware.PermissionProjectWrite))
	projectGroup.DELETE("/:id/tasks/:task_id", h.RemoveTask, rbacMiddleware.RequirePermission(middleware.PermissionProjectWrite))
	
	// Routes that require project:delete permission
	projectGroup.DELETE("/:id", h.Delete, rbacMiddleware.RequirePermission(middleware.PermissionProjectDelete))
}
