package comment

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

// Handlers provides HTTP handlers for comment management
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

// Create handles the creation of a new comment
func (h *Handlers) Create(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Parse request body
	var req models.CommentRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Create comment
	response, err := h.service.Create(c.Request().Context(), req, userID)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Task not found")
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create comment")
	}

	return c.JSON(http.StatusCreated, response)
}

// GetByID handles retrieving a comment by ID
func (h *Handlers) GetByID(c echo.Context) error {
	// Get comment ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid comment ID")
	}

	// Get comment
	response, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrCommentNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Comment not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve comment")
	}

	return c.JSON(http.StatusOK, response)
}

// List handles retrieving comments based on filter parameters
func (h *Handlers) List(c echo.Context) error {
	// Parse query parameters
	params := models.CommentListParams{}
	
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
	
	// Get comments
	comments, total, err := h.service.List(c.Request().Context(), params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve comments")
	}
	
	// Build response
	response := map[string]interface{}{
		"comments": comments,
		"pagination": map[string]interface{}{
			"total":       total,
			"page":        params.Page,
			"page_size":   params.PageSize,
			"total_pages": (total + params.PageSize - 1) / params.PageSize,
		},
	}
	
	return c.JSON(http.StatusOK, response)
}

// Update handles updating a comment
func (h *Handlers) Update(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}
	
	// Get comment ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid comment ID")
	}
	
	// Check if the comment belongs to the user
	comment, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrCommentNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Comment not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve comment")
	}
	
	if comment.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "You don't have permission to update this comment")
	}
	
	// Parse request body
	var req models.CommentRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	// Update comment
	response, err := h.service.Update(c.Request().Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrCommentNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Comment not found")
		}
		if errors.Is(err, repository.ErrTaskNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Task not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update comment")
	}
	
	return c.JSON(http.StatusOK, response)
}

// Delete handles deleting a comment
func (h *Handlers) Delete(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}
	
	// Get comment ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid comment ID")
	}
	
	// Check if the comment belongs to the user
	comment, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrCommentNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Comment not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve comment")
	}
	
	if comment.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "You don't have permission to delete this comment")
	}
	
	// Delete comment
	err = h.service.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrCommentNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Comment not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete comment")
	}
	
	return c.NoContent(http.StatusNoContent)
}

// RegisterRoutes registers the comment routes
func (h *Handlers) RegisterRoutes(g *echo.Group, rbacMiddleware *middleware.RBACMiddleware) {
	commentGroup := g.Group("/comments")
	
	// Routes that require comment:read permission
	commentGroup.GET("", h.List, rbacMiddleware.RequirePermission(middleware.PermissionCommentRead))
	commentGroup.GET("/:id", h.GetByID, rbacMiddleware.RequirePermission(middleware.PermissionCommentRead))
	
	// Routes that require comment:write permission
	commentGroup.POST("", h.Create, rbacMiddleware.RequirePermission(middleware.PermissionCommentWrite))
	commentGroup.PUT("/:id", h.Update, rbacMiddleware.RequirePermission(middleware.PermissionCommentWrite))
	
	// Routes that require comment:delete permission
	commentGroup.DELETE("/:id", h.Delete, rbacMiddleware.RequirePermission(middleware.PermissionCommentDelete))
}
