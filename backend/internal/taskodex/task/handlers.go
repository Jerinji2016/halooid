package task

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

// Handlers provides HTTP handlers for task management
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

// Create handles the creation of a new task
func (h *Handlers) Create(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Parse request body
	var req models.TaskRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Create task
	response, err := h.service.Create(c.Request().Context(), req, userID)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Project not found")
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create task")
	}

	return c.JSON(http.StatusCreated, response)
}

// GetByID handles retrieving a task by ID
func (h *Handlers) GetByID(c echo.Context) error {
	// Get task ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}

	// Get task
	response, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Task not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve task")
	}

	return c.JSON(http.StatusOK, response)
}

// List handles retrieving tasks based on filter parameters
func (h *Handlers) List(c echo.Context) error {
	// Parse query parameters
	params := models.TaskListParams{}

	// Parse project_id parameter
	projectIDParam := c.QueryParam("project_id")
	if projectIDParam != "" {
		projectID, err := uuid.Parse(projectIDParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid project_id parameter")
		}
		params.ProjectID = &projectID
	}

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

	// Parse created_by parameter
	createdByParam := c.QueryParam("created_by")
	if createdByParam != "" {
		createdBy, err := uuid.Parse(createdByParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid created_by parameter")
		}
		params.CreatedBy = &createdBy
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
	tasks, total, err := h.service.List(c.Request().Context(), params)
	if err != nil {
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

// Update handles updating a task
func (h *Handlers) Update(c echo.Context) error {
	// Get task ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}

	// Parse request body
	var req models.TaskRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Update task
	response, err := h.service.Update(c.Request().Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Task not found")
		}
		if errors.Is(err, repository.ErrProjectNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Project not found")
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update task")
	}

	return c.JSON(http.StatusOK, response)
}

// Delete handles deleting a task
func (h *Handlers) Delete(c echo.Context) error {
	// Get task ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}

	// Delete task
	err = h.service.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Task not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete task")
	}

	return c.NoContent(http.StatusNoContent)
}

// AddTag handles adding a tag to a task
func (h *Handlers) AddTag(c echo.Context) error {
	// Get task ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}

	// Parse request body
	var req struct {
		Tag string `json:"tag" validate:"required,max=50"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Add tag
	err = h.service.AddTag(c.Request().Context(), id, req.Tag)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Task not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add tag")
	}

	return c.NoContent(http.StatusNoContent)
}

// RemoveTag handles removing a tag from a task
func (h *Handlers) RemoveTag(c echo.Context) error {
	// Get task ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}

	// Get tag from path parameter
	tag := c.Param("tag")
	if tag == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Tag is required")
	}

	// Remove tag
	err = h.service.RemoveTag(c.Request().Context(), id, tag)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Task not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove tag")
	}

	return c.NoContent(http.StatusNoContent)
}

// AssignTask handles assigning a task to a user
func (h *Handlers) AssignTask(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Get task ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}

	// Parse request body
	var req struct {
		UserID uuid.UUID `json:"user_id" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Assign task
	response, err := h.service.AssignTask(c.Request().Context(), id, req.UserID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Task not found")
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to assign task")
	}

	return c.JSON(http.StatusOK, response)
}

// UnassignTask handles removing the assignment of a task
func (h *Handlers) UnassignTask(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Get task ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}

	// Unassign task
	response, err := h.service.UnassignTask(c.Request().Context(), id, userID)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Task not found")
		}
		if errors.Is(err, errors.New("task is not assigned to anyone")) {
			return echo.NewHTTPError(http.StatusBadRequest, "Task is not assigned to anyone")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to unassign task")
	}

	return c.JSON(http.StatusOK, response)
}

// UpdateTaskStatus handles updating the status of a task
func (h *Handlers) UpdateTaskStatus(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Get task ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}

	// Parse request body
	var req struct {
		Status models.TaskStatus `json:"status" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Update task status
	response, err := h.service.UpdateTaskStatus(c.Request().Context(), id, req.Status, userID)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Task not found")
		}
		if errors.Is(err, errors.New("invalid task status")) {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid task status")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update task status")
	}

	return c.JSON(http.StatusOK, response)
}

// GetTasksByAssignee handles retrieving tasks assigned to a user
func (h *Handlers) GetTasksByAssignee(c echo.Context) error {
	// Get user ID from path parameter
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	// Parse query parameters
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
	tasks, total, err := h.service.GetTasksByAssignee(c.Request().Context(), userID, params)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
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

// GetOverdueTasks handles retrieving overdue tasks
func (h *Handlers) GetOverdueTasks(c echo.Context) error {
	// Parse query parameters
	params := models.TaskListParams{}

	// Parse project_id parameter
	projectIDParam := c.QueryParam("project_id")
	if projectIDParam != "" {
		projectID, err := uuid.Parse(projectIDParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid project_id parameter")
		}
		params.ProjectID = &projectID
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

	// Get overdue tasks
	tasks, total, err := h.service.GetOverdueTasks(c.Request().Context(), params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve overdue tasks")
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

// GetTasksDueSoon handles retrieving tasks due within a specified number of days
func (h *Handlers) GetTasksDueSoon(c echo.Context) error {
	// Parse days parameter
	daysParam := c.Param("days")
	days, err := strconv.Atoi(daysParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid days parameter")
	}

	// Parse query parameters
	params := models.TaskListParams{}

	// Parse project_id parameter
	projectIDParam := c.QueryParam("project_id")
	if projectIDParam != "" {
		projectID, err := uuid.Parse(projectIDParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid project_id parameter")
		}
		params.ProjectID = &projectID
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

	// Get tasks due soon
	tasks, total, err := h.service.GetTasksDueSoon(c.Request().Context(), days, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve tasks due soon")
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

// RegisterRoutes registers the task routes
func (h *Handlers) RegisterRoutes(g *echo.Group, rbacMiddleware *middleware.RBACMiddleware) {
	taskGroup := g.Group("/tasks")

	// Routes that require task:read permission
	taskGroup.GET("", h.List, rbacMiddleware.RequirePermission(middleware.PermissionTaskRead))
	taskGroup.GET("/:id", h.GetByID, rbacMiddleware.RequirePermission(middleware.PermissionTaskRead))
	taskGroup.GET("/assignee/:user_id", h.GetTasksByAssignee, rbacMiddleware.RequirePermission(middleware.PermissionTaskRead))
	taskGroup.GET("/overdue", h.GetOverdueTasks, rbacMiddleware.RequirePermission(middleware.PermissionTaskRead))
	taskGroup.GET("/due-soon/:days", h.GetTasksDueSoon, rbacMiddleware.RequirePermission(middleware.PermissionTaskRead))

	// Routes that require task:write permission
	taskGroup.POST("", h.Create, rbacMiddleware.RequirePermission(middleware.PermissionTaskWrite))
	taskGroup.PUT("/:id", h.Update, rbacMiddleware.RequirePermission(middleware.PermissionTaskWrite))
	taskGroup.POST("/:id/tags", h.AddTag, rbacMiddleware.RequirePermission(middleware.PermissionTaskWrite))
	taskGroup.DELETE("/:id/tags/:tag", h.RemoveTag, rbacMiddleware.RequirePermission(middleware.PermissionTaskWrite))
	taskGroup.POST("/:id/assign", h.AssignTask, rbacMiddleware.RequirePermission(middleware.PermissionTaskWrite))
	taskGroup.POST("/:id/unassign", h.UnassignTask, rbacMiddleware.RequirePermission(middleware.PermissionTaskWrite))
	taskGroup.PUT("/:id/status", h.UpdateTaskStatus, rbacMiddleware.RequirePermission(middleware.PermissionTaskWrite))

	// Routes that require task:delete permission
	taskGroup.DELETE("/:id", h.Delete, rbacMiddleware.RequirePermission(middleware.PermissionTaskDelete))
}
