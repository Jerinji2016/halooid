package organization

import (
	"errors"
	"net/http"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handlers provides HTTP handlers for organization management
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

// Create handles the creation of a new organization
func (h *Handlers) Create(c echo.Context) error {
	var req models.OrganizationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	response, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, repository.ErrOrganizationNameExists) {
			return echo.NewHTTPError(http.StatusConflict, "Organization name already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create organization")
	}

	return c.JSON(http.StatusCreated, response)
}

// GetByID handles retrieving an organization by ID
func (h *Handlers) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	response, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrOrganizationNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve organization")
	}

	return c.JSON(http.StatusOK, response)
}

// List handles retrieving all organizations
func (h *Handlers) List(c echo.Context) error {
	responses, err := h.service.List(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve organizations")
	}

	return c.JSON(http.StatusOK, responses)
}

// Update handles updating an organization
func (h *Handlers) Update(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	var req models.OrganizationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	response, err := h.service.Update(c.Request().Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrOrganizationNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
		}
		if errors.Is(err, repository.ErrOrganizationNameExists) {
			return echo.NewHTTPError(http.StatusConflict, "Organization name already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update organization")
	}

	return c.JSON(http.StatusOK, response)
}

// Delete handles marking an organization as inactive
func (h *Handlers) Delete(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	err = h.service.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrOrganizationNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete organization")
	}

	return c.NoContent(http.StatusNoContent)
}

// AddUser handles adding a user to an organization
func (h *Handlers) AddUser(c echo.Context) error {
	var req models.AddUserToOrganizationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := h.service.AddUser(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, repository.ErrOrganizationNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add user to organization")
	}

	return c.NoContent(http.StatusNoContent)
}

// RemoveUser handles removing a user from an organization
func (h *Handlers) RemoveUser(c echo.Context) error {
	orgIDParam := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	err = h.service.RemoveUser(c.Request().Context(), orgID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove user from organization")
	}

	return c.NoContent(http.StatusNoContent)
}

// GetUsers handles retrieving all users in an organization
func (h *Handlers) GetUsers(c echo.Context) error {
	orgIDParam := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	responses, err := h.service.GetUsers(c.Request().Context(), orgID)
	if err != nil {
		if errors.Is(err, repository.ErrOrganizationNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve users")
	}

	return c.JSON(http.StatusOK, responses)
}

// GetOrganizations handles retrieving all organizations a user belongs to
func (h *Handlers) GetOrganizations(c echo.Context) error {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	responses, err := h.service.GetOrganizations(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve organizations")
	}

	return c.JSON(http.StatusOK, responses)
}
