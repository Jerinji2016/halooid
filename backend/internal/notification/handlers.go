package notification

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

// Handlers provides HTTP handlers for notification management
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

// GetByID handles retrieving a notification by ID
func (h *Handlers) GetByID(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Get notification ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid notification ID")
	}

	// Get notification
	response, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotificationNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Notification not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve notification")
	}

	// Check if the notification belongs to the user
	if response.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "You don't have permission to access this notification")
	}

	return c.JSON(http.StatusOK, response)
}

// List handles retrieving notifications based on filter parameters
func (h *Handlers) List(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Parse query parameters
	params := models.NotificationListParams{
		UserID: userID,
	}
	
	// Parse is_read parameter
	isReadParam := c.QueryParam("is_read")
	if isReadParam != "" {
		isRead, err := strconv.ParseBool(isReadParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid is_read parameter")
		}
		params.IsRead = &isRead
	}
	
	// Parse type parameter
	typeParam := c.QueryParam("type")
	if typeParam != "" {
		notificationType := models.NotificationType(typeParam)
		params.Type = &notificationType
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
	
	// Get notifications
	notifications, total, err := h.service.List(c.Request().Context(), params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve notifications")
	}
	
	// Build response
	response := map[string]interface{}{
		"notifications": notifications,
		"pagination": map[string]interface{}{
			"total":       total,
			"page":        params.Page,
			"page_size":   params.PageSize,
			"total_pages": (total + params.PageSize - 1) / params.PageSize,
		},
	}
	
	return c.JSON(http.StatusOK, response)
}

// MarkAsRead handles marking a notification as read
func (h *Handlers) MarkAsRead(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Get notification ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid notification ID")
	}

	// Check if the notification belongs to the user
	notification, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotificationNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Notification not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve notification")
	}

	if notification.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "You don't have permission to access this notification")
	}

	// Mark notification as read
	err = h.service.MarkAsRead(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to mark notification as read")
	}

	return c.NoContent(http.StatusNoContent)
}

// MarkAllAsRead handles marking all notifications for a user as read
func (h *Handlers) MarkAllAsRead(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Mark all notifications as read
	err = h.service.MarkAllAsRead(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to mark all notifications as read")
	}

	return c.NoContent(http.StatusNoContent)
}

// Delete handles deleting a notification
func (h *Handlers) Delete(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Get notification ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid notification ID")
	}

	// Check if the notification belongs to the user
	notification, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotificationNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Notification not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve notification")
	}

	if notification.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "You don't have permission to access this notification")
	}

	// Delete notification
	err = h.service.Delete(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete notification")
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteAll handles deleting all notifications for a user
func (h *Handlers) DeleteAll(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Delete all notifications
	err = h.service.DeleteAllForUser(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete all notifications")
	}

	return c.NoContent(http.StatusNoContent)
}

// RegisterRoutes registers the notification routes
func (h *Handlers) RegisterRoutes(g *echo.Group) {
	notificationGroup := g.Group("/notifications")
	
	notificationGroup.GET("", h.List)
	notificationGroup.GET("/:id", h.GetByID)
	notificationGroup.PUT("/:id/read", h.MarkAsRead)
	notificationGroup.PUT("/read-all", h.MarkAllAsRead)
	notificationGroup.DELETE("/:id", h.Delete)
	notificationGroup.DELETE("", h.DeleteAll)
}
