package timeentry

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

// Handlers provides HTTP handlers for time entry management
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

// Create handles the creation of a new time entry
func (h *Handlers) Create(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Parse request body
	var req models.TimeEntryRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Create time entry
	response, err := h.service.Create(c.Request().Context(), req, userID)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Task not found")
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		if errors.Is(err, repository.ErrRunningTimeEntry) {
			return echo.NewHTTPError(http.StatusConflict, "You already have a running timer for this task")
		}
		if errors.Is(err, errors.New("end time cannot be before start time")) {
			return echo.NewHTTPError(http.StatusBadRequest, "End time cannot be before start time")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create time entry")
	}

	return c.JSON(http.StatusCreated, response)
}

// GetByID handles retrieving a time entry by ID
func (h *Handlers) GetByID(c echo.Context) error {
	// Get time entry ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid time entry ID")
	}

	// Get time entry
	response, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrTimeEntryNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Time entry not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve time entry")
	}

	return c.JSON(http.StatusOK, response)
}

// List handles retrieving time entries based on filter parameters
func (h *Handlers) List(c echo.Context) error {
	// Parse query parameters
	params := models.TimeEntryListParams{}
	
	// Parse task_id parameter
	taskIDParam := c.QueryParam("task_id")
	if taskIDParam != "" {
		taskID, err := uuid.Parse(taskIDParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid task_id parameter")
		}
		params.TaskID = &taskID
	}
	
	// Parse user_id parameter
	userIDParam := c.QueryParam("user_id")
	if userIDParam != "" {
		userID, err := uuid.Parse(userIDParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid user_id parameter")
		}
		params.UserID = &userID
	}
	
	// Parse start_after parameter
	startAfterParam := c.QueryParam("start_after")
	if startAfterParam != "" {
		startAfter, err := time.Parse(time.RFC3339, startAfterParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid start_after parameter")
		}
		params.StartAfter = &startAfter
	}
	
	// Parse start_before parameter
	startBeforeParam := c.QueryParam("start_before")
	if startBeforeParam != "" {
		startBefore, err := time.Parse(time.RFC3339, startBeforeParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid start_before parameter")
		}
		params.StartBefore = &startBefore
	}
	
	// Parse is_running parameter
	isRunningParam := c.QueryParam("is_running")
	if isRunningParam != "" {
		isRunning, err := strconv.ParseBool(isRunningParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid is_running parameter")
		}
		params.IsRunning = &isRunning
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
	
	// Get time entries
	timeEntries, total, err := h.service.List(c.Request().Context(), params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve time entries")
	}
	
	// Build response
	response := map[string]interface{}{
		"time_entries": timeEntries,
		"pagination": map[string]interface{}{
			"total":       total,
			"page":        params.Page,
			"page_size":   params.PageSize,
			"total_pages": (total + params.PageSize - 1) / params.PageSize,
		},
	}
	
	return c.JSON(http.StatusOK, response)
}

// Update handles updating a time entry
func (h *Handlers) Update(c echo.Context) error {
	// Get time entry ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid time entry ID")
	}
	
	// Parse request body
	var req models.TimeEntryRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	// Update time entry
	response, err := h.service.Update(c.Request().Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrTimeEntryNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Time entry not found")
		}
		if errors.Is(err, repository.ErrTaskNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Task not found")
		}
		if errors.Is(err, errors.New("end time cannot be before start time")) {
			return echo.NewHTTPError(http.StatusBadRequest, "End time cannot be before start time")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update time entry")
	}
	
	return c.JSON(http.StatusOK, response)
}

// Delete handles deleting a time entry
func (h *Handlers) Delete(c echo.Context) error {
	// Get time entry ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid time entry ID")
	}
	
	// Delete time entry
	err = h.service.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrTimeEntryNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Time entry not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete time entry")
	}
	
	return c.NoContent(http.StatusNoContent)
}

// StartTimer handles starting a timer for a task
func (h *Handlers) StartTimer(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}
	
	// Parse request body
	var req struct {
		TaskID      uuid.UUID `json:"task_id" validate:"required"`
		Description string    `json:"description" validate:"max=1000"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	// Start timer
	response, err := h.service.StartTimer(c.Request().Context(), req.TaskID, req.Description, userID)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Task not found")
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		if errors.Is(err, repository.ErrRunningTimeEntry) {
			return echo.NewHTTPError(http.StatusConflict, "You already have a running timer for this task")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to start timer")
	}
	
	return c.JSON(http.StatusOK, response)
}

// StopTimer handles stopping a running timer for a task
func (h *Handlers) StopTimer(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
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
	
	// Stop timer
	response, err := h.service.StopTimer(c.Request().Context(), req.TaskID, userID)
	if err != nil {
		if errors.Is(err, errors.New("no running timer found for this task")) {
			return echo.NewHTTPError(http.StatusNotFound, "No running timer found for this task")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to stop timer")
	}
	
	return c.JSON(http.StatusOK, response)
}

// GetRunningTimers handles retrieving all running timers for a user
func (h *Handlers) GetRunningTimers(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}
	
	// Get running timers
	responses, err := h.service.GetRunningTimers(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve running timers")
	}
	
	return c.JSON(http.StatusOK, responses)
}

// Aggregate handles aggregating time entries based on parameters
func (h *Handlers) Aggregate(c echo.Context) error {
	// Parse query parameters
	params := models.TimeEntryAggregationParams{}
	
	// Parse task_id parameter
	taskIDParam := c.QueryParam("task_id")
	if taskIDParam != "" {
		taskID, err := uuid.Parse(taskIDParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid task_id parameter")
		}
		params.TaskID = &taskID
	}
	
	// Parse user_id parameter
	userIDParam := c.QueryParam("user_id")
	if userIDParam != "" {
		userID, err := uuid.Parse(userIDParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid user_id parameter")
		}
		params.UserID = &userID
	}
	
	// Parse project_id parameter
	projectIDParam := c.QueryParam("project_id")
	if projectIDParam != "" {
		projectID, err := uuid.Parse(projectIDParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid project_id parameter")
		}
		params.ProjectID = &projectID
	}
	
	// Parse start_after parameter
	startAfterParam := c.QueryParam("start_after")
	if startAfterParam != "" {
		startAfter, err := time.Parse(time.RFC3339, startAfterParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid start_after parameter")
		}
		params.StartAfter = &startAfter
	}
	
	// Parse start_before parameter
	startBeforeParam := c.QueryParam("start_before")
	if startBeforeParam != "" {
		startBefore, err := time.Parse(time.RFC3339, startBeforeParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid start_before parameter")
		}
		params.StartBefore = &startBefore
	}
	
	// Parse group_by parameter
	groupByParam := c.QueryParam("group_by")
	if groupByParam != "" {
		validGroupBy := map[string]bool{
			"day":   true,
			"week":  true,
			"month": true,
			"year":  true,
			"task":  true,
			"user":  true,
		}
		
		if !validGroupBy[groupByParam] {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid group_by parameter")
		}
		
		params.GroupBy = groupByParam
	} else {
		params.GroupBy = "day"
	}
	
	// Get aggregated time entries
	aggregations, err := h.service.Aggregate(c.Request().Context(), params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to aggregate time entries")
	}
	
	return c.JSON(http.StatusOK, aggregations)
}

// RegisterRoutes registers the time entry routes
func (h *Handlers) RegisterRoutes(g *echo.Group, rbacMiddleware *middleware.RBACMiddleware) {
	timeEntryGroup := g.Group("/time-entries")
	
	// Routes that require time_entry:read permission
	timeEntryGroup.GET("", h.List, rbacMiddleware.RequirePermission(middleware.PermissionTimeEntryRead))
	timeEntryGroup.GET("/:id", h.GetByID, rbacMiddleware.RequirePermission(middleware.PermissionTimeEntryRead))
	timeEntryGroup.GET("/running", h.GetRunningTimers, rbacMiddleware.RequirePermission(middleware.PermissionTimeEntryRead))
	timeEntryGroup.GET("/aggregate", h.Aggregate, rbacMiddleware.RequirePermission(middleware.PermissionTimeEntryRead))
	
	// Routes that require time_entry:write permission
	timeEntryGroup.POST("", h.Create, rbacMiddleware.RequirePermission(middleware.PermissionTimeEntryWrite))
	timeEntryGroup.PUT("/:id", h.Update, rbacMiddleware.RequirePermission(middleware.PermissionTimeEntryWrite))
	timeEntryGroup.POST("/start", h.StartTimer, rbacMiddleware.RequirePermission(middleware.PermissionTimeEntryWrite))
	timeEntryGroup.POST("/stop", h.StopTimer, rbacMiddleware.RequirePermission(middleware.PermissionTimeEntryWrite))
	
	// Routes that require time_entry:delete permission
	timeEntryGroup.DELETE("/:id", h.Delete, rbacMiddleware.RequirePermission(middleware.PermissionTimeEntryDelete))
}
