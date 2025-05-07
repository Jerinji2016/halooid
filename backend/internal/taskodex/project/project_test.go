package project_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/Jerinji2016/halooid/backend/internal/taskodex/project"
	"github.com/Jerinji2016/halooid/backend/internal/test"
	"github.com/Jerinji2016/halooid/backend/pkg/middleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestProjectAPI(t *testing.T) {
	// Setup test environment
	tdb, prefix := test.SetupTestEnvironment(t)
	defer test.TeardownTestEnvironment(t, tdb, prefix)

	// Create test user and organization
	testUser := tdb.CreateTestUser(t, prefix)
	testOrg := tdb.CreateTestOrganization(t, prefix, testUser.ID)

	// Create repositories
	projectRepo := repository.NewPostgresProjectRepository(tdb.DB)
	userRepo := repository.NewPostgresUserRepository(tdb.DB)

	// Create service and handlers
	projectService := project.NewService(projectRepo, userRepo)
	projectHandlers := project.NewHandlers(projectService)

	// Setup Echo
	e := echo.New()
	g := e.Group("/api/v1/organizations/:org_id/taskodex")

	// Create a mock RBAC middleware that always allows access
	mockRBACMiddleware := &middleware.RBACMiddleware{
		RequirePermissionFunc: func(permission string) echo.MiddlewareFunc {
			return func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					// Set user ID in context
					c.Set("user_id", testUser.ID.String())
					return next(c)
				}
			}
		},
	}

	// Register routes
	projectHandlers.RegisterRoutes(g, mockRBACMiddleware)

	t.Run("CreateProject", func(t *testing.T) {
		// Create request body
		reqBody := models.ProjectRequest{
			Name:        fmt.Sprintf("%s Test Project", prefix),
			Description: "Test project for integration tests",
			Status:      models.ProjectStatusPlanning,
		}
		reqJSON, _ := json.Marshal(reqBody)

		// Create request
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/projects")
		c.SetParamNames("org_id")
		c.SetParamValues(testOrg.ID.String())

		// Handle request
		err := projectHandlers.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Parse response
		var response models.ProjectResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, reqBody.Name, response.Name)
		assert.Equal(t, reqBody.Description, response.Description)
		assert.Equal(t, reqBody.Status, response.Status)
		assert.Equal(t, testOrg.ID, response.OrganizationID)
		assert.Equal(t, testUser.ID, response.CreatedBy)
	})

	// Create a test project for the remaining tests
	testProject := tdb.CreateTestProject(t, prefix, testOrg.ID, testUser.ID)

	t.Run("GetProjectByID", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/projects/:id")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testProject.ID.String())

		// Handle request
		err := projectHandlers.GetByID(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response models.ProjectResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testProject.ID, response.ID)
		assert.Equal(t, testProject.Name, response.Name)
		assert.Equal(t, testProject.Description, response.Description)
		assert.Equal(t, testProject.Status, response.Status)
		assert.Equal(t, testOrg.ID, response.OrganizationID)
		assert.Equal(t, testUser.ID, response.CreatedBy)
	})

	t.Run("ListProjects", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodGet, "/?page=1&page_size=10", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/projects")
		c.SetParamNames("org_id")
		c.SetParamValues(testOrg.ID.String())

		// Handle request
		err := projectHandlers.List(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check that projects array exists and contains at least one project
		projects, ok := response["projects"].([]interface{})
		assert.True(t, ok)
		assert.GreaterOrEqual(t, len(projects), 1)

		// Check pagination
		pagination, ok := response["pagination"].(map[string]interface{})
		assert.True(t, ok)
		assert.GreaterOrEqual(t, int(pagination["total"].(float64)), 1)
	})

	t.Run("UpdateProject", func(t *testing.T) {
		// Create request body
		reqBody := models.ProjectRequest{
			Name:        fmt.Sprintf("%s Updated Project", prefix),
			Description: "Updated project description",
			Status:      models.ProjectStatusActive,
		}
		reqJSON, _ := json.Marshal(reqBody)

		// Create request
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/projects/:id")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testProject.ID.String())

		// Handle request
		err := projectHandlers.Update(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response models.ProjectResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testProject.ID, response.ID)
		assert.Equal(t, reqBody.Name, response.Name)
		assert.Equal(t, reqBody.Description, response.Description)
		assert.Equal(t, reqBody.Status, response.Status)
		assert.Equal(t, testOrg.ID, response.OrganizationID)
		assert.Equal(t, testUser.ID, response.CreatedBy)
	})

	// Create a test task for the project
	testTask := tdb.CreateTestTask(t, prefix, testProject.ID, testUser.ID)

	t.Run("GetProjectTasks", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodGet, "/?page=1&page_size=10", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/projects/:id/tasks")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testProject.ID.String())

		// Handle request
		err := projectHandlers.GetTasks(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check that tasks array exists and contains at least one task
		tasks, ok := response["tasks"].([]interface{})
		assert.True(t, ok)
		assert.GreaterOrEqual(t, len(tasks), 1)

		// Check pagination
		pagination, ok := response["pagination"].(map[string]interface{})
		assert.True(t, ok)
		assert.GreaterOrEqual(t, int(pagination["total"].(float64)), 1)
	})

	t.Run("RemoveTask", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/projects/:id/tasks/:task_id")
		c.SetParamNames("org_id", "id", "task_id")
		c.SetParamValues(testOrg.ID.String(), testProject.ID.String(), testTask.ID.String())

		// Handle request
		err := projectHandlers.RemoveTask(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("DeleteProject", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/projects/:id")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testProject.ID.String())

		// Handle request
		err := projectHandlers.Delete(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)

		// Verify project is deleted
		_, err = projectRepo.GetByID(c.Request().Context(), testProject.ID)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrProjectNotFound, err)
	})
}
