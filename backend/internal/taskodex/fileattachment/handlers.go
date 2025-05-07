package fileattachment

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/Jerinji2016/halooid/backend/internal/storage"
	"github.com/Jerinji2016/halooid/backend/pkg/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handlers provides HTTP handlers for file attachment management
type Handlers struct {
	service     Service
	fileStorage storage.FileStorage
	validate    *validator.Validate
}

// NewHandlers creates a new Handlers
func NewHandlers(service Service, fileStorage storage.FileStorage) *Handlers {
	return &Handlers{
		service:     service,
		fileStorage: fileStorage,
		validate:    validator.New(),
	}
}

// Upload handles uploading a file and creating a file attachment
func (h *Handlers) Upload(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Get task ID from form
	taskIDParam := c.FormValue("task_id")
	taskID, err := uuid.Parse(taskIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file")
	}

	// Check file size (limit to 10MB)
	if file.Size > 10*1024*1024 {
		return echo.NewHTTPError(http.StatusBadRequest, "File size exceeds the limit (10MB)")
	}

	// Upload file
	response, err := h.service.Upload(c.Request().Context(), taskID, userID, file)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Task not found")
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to upload file")
	}

	return c.JSON(http.StatusCreated, response)
}

// GetByID handles retrieving a file attachment by ID
func (h *Handlers) GetByID(c echo.Context) error {
	// Get file attachment ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file attachment ID")
	}

	// Get file attachment
	response, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrFileAttachmentNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "File attachment not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve file attachment")
	}

	return c.JSON(http.StatusOK, response)
}

// List handles retrieving file attachments based on filter parameters
func (h *Handlers) List(c echo.Context) error {
	// Parse query parameters
	params := models.FileAttachmentListParams{}
	
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
	
	// Get file attachments
	fileAttachments, total, err := h.service.List(c.Request().Context(), params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve file attachments")
	}
	
	// Build response
	response := map[string]interface{}{
		"file_attachments": fileAttachments,
		"pagination": map[string]interface{}{
			"total":       total,
			"page":        params.Page,
			"page_size":   params.PageSize,
			"total_pages": (total + params.PageSize - 1) / params.PageSize,
		},
	}
	
	return c.JSON(http.StatusOK, response)
}

// Download handles downloading a file attachment
func (h *Handlers) Download(c echo.Context) error {
	// Get file attachment ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file attachment ID")
	}
	
	// Get file attachment
	fileAttachment, err := h.service.Download(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrFileAttachmentNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "File attachment not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve file attachment")
	}
	
	// Get file from storage
	file, err := h.fileStorage.GetFile(fileAttachment.StoragePath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve file")
	}
	defer file.Close()
	
	// Set response headers
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename="+fileAttachment.FileName)
	c.Response().Header().Set(echo.HeaderContentType, fileAttachment.ContentType)
	
	// Stream the file to the response
	_, err = io.Copy(c.Response().Writer, file)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to stream file")
	}
	
	return nil
}

// Delete handles deleting a file attachment
func (h *Handlers) Delete(c echo.Context) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}
	
	// Get file attachment ID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file attachment ID")
	}
	
	// Check if the file attachment belongs to the user
	fileAttachment, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrFileAttachmentNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "File attachment not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve file attachment")
	}
	
	if fileAttachment.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "You don't have permission to delete this file attachment")
	}
	
	// Delete file attachment
	err = h.service.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrFileAttachmentNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "File attachment not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete file attachment")
	}
	
	return c.NoContent(http.StatusNoContent)
}

// RegisterRoutes registers the file attachment routes
func (h *Handlers) RegisterRoutes(g *echo.Group, rbacMiddleware *middleware.RBACMiddleware) {
	fileGroup := g.Group("/files")
	
	// Routes that require file:read permission
	fileGroup.GET("", h.List, rbacMiddleware.RequirePermission(middleware.PermissionFileRead))
	fileGroup.GET("/:id", h.GetByID, rbacMiddleware.RequirePermission(middleware.PermissionFileRead))
	fileGroup.GET("/:id/download", h.Download, rbacMiddleware.RequirePermission(middleware.PermissionFileRead))
	
	// Routes that require file:write permission
	fileGroup.POST("", h.Upload, rbacMiddleware.RequirePermission(middleware.PermissionFileWrite))
	
	// Routes that require file:delete permission
	fileGroup.DELETE("/:id", h.Delete, rbacMiddleware.RequirePermission(middleware.PermissionFileDelete))
}
