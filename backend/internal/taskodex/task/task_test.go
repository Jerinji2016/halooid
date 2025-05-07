package task_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/notification"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/Jerinji2016/halooid/backend/internal/taskodex/task"
	"github.com/Jerinji2016/halooid/backend/internal/test"
	"github.com/Jerinji2016/halooid/backend/pkg/middleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestTaskAPI(t *testing.T) {
	// Setup test environment
	tdb, prefix := test.SetupTestEnvironment(t)
	defer test.TeardownTestEnvironment(t, tdb, prefix)

	// Create test user and organization
	testUser := tdb.CreateTestUser(t, prefix)
	testOrg := tdb.CreateTestOrganization(t, prefix, testUser.ID)
	testProject := tdb.CreateTestProject(t, prefix, testOrg.ID, testUser.ID)

	// Create repositories
	taskRepo := repository.NewPostgresTaskRepository(tdb.DB)
	projectRepo := repository.NewPostgresProjectRepository(tdb.DB)
	userRepo := repository.NewPostgresUserRepository(tdb.DB)
	notificationRepo := repository.NewPostgresNotificationRepository(tdb.DB)

	// Create services
	notificationService := notification.NewService(notificationRepo, userRepo)
	taskService := task.NewService(taskRepo, projectRepo, userRepo, notificationService)
	taskHandlers := task.NewHandlers(taskService)

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
	taskHandlers.RegisterRoutes(g, mockRBACMiddleware)

	t.Run("CreateTask", func(t *testing.T) {
		// Create request body
		dueDate := time.Now().Add(7 * 24 * time.Hour)
		reqBody := models.TaskRequest{
			Title:       fmt.Sprintf("%s Test Task", prefix),
			Description: "Test task for integration tests",
			Status:      models.TaskStatusTodo,
			Priority:    models.TaskPriorityMedium,
			ProjectID:   &testProject.ID,
			DueDate:     &dueDate,
		}
		reqJSON, _ := json.Marshal(reqBody)

		// Create request
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/tasks")
		c.SetParamNames("org_id")
		c.SetParamValues(testOrg.ID.String())

		// Handle request
		err := taskHandlers.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Parse response
		var response models.TaskResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, reqBody.Title, response.Title)
		assert.Equal(t, reqBody.Description, response.Description)
		assert.Equal(t, reqBody.Status, response.Status)
		assert.Equal(t, reqBody.Priority, response.Priority)
		assert.Equal(t, testProject.ID, *response.ProjectID)
		assert.Equal(t, testUser.ID, response.CreatedBy)
	})

	// Create a test task for the remaining tests
	testTask := tdb.CreateTestTask(t, prefix, testProject.ID, testUser.ID)

	t.Run("GetTaskByID", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/tasks/:id")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testTask.ID.String())

		// Handle request
		err := taskHandlers.GetByID(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response models.TaskResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testTask.ID, response.ID)
		assert.Equal(t, testTask.Title, response.Title)
		assert.Equal(t, testTask.Description, response.Description)
		assert.Equal(t, testTask.Status, response.Status)
		assert.Equal(t, testTask.Priority, response.Priority)
		assert.Equal(t, testProject.ID, *response.ProjectID)
		assert.Equal(t, testUser.ID, response.CreatedBy)
	})

	t.Run("ListTasks", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodGet, "/?page=1&page_size=10", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/tasks")
		c.SetParamNames("org_id")
		c.SetParamValues(testOrg.ID.String())

		// Handle request
		err := taskHandlers.List(c)
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

	t.Run("UpdateTask", func(t *testing.T) {
		// Create request body
		dueDate := time.Now().Add(14 * 24 * time.Hour)
		reqBody := models.TaskRequest{
			Title:       fmt.Sprintf("%s Updated Task", prefix),
			Description: "Updated task description",
			Status:      models.TaskStatusInProgress,
			Priority:    models.TaskPriorityHigh,
			ProjectID:   &testProject.ID,
			DueDate:     &dueDate,
		}
		reqJSON, _ := json.Marshal(reqBody)

		// Create request
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/tasks/:id")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testTask.ID.String())

		// Handle request
		err := taskHandlers.Update(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response models.TaskResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testTask.ID, response.ID)
		assert.Equal(t, reqBody.Title, response.Title)
		assert.Equal(t, reqBody.Description, response.Description)
		assert.Equal(t, reqBody.Status, response.Status)
		assert.Equal(t, reqBody.Priority, response.Priority)
		assert.Equal(t, testProject.ID, *response.ProjectID)
		assert.Equal(t, testUser.ID, response.CreatedBy)
	})

	t.Run("AddTag", func(t *testing.T) {
		// Create request body
		reqBody := struct {
			Tag string `json:"tag"`
		}{
			Tag: "test-tag",
		}
		reqJSON, _ := json.Marshal(reqBody)

		// Create request
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/tasks/:id/tags")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testTask.ID.String())

		// Handle request
		err := taskHandlers.AddTag(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)

		// Verify tag was added
		task, err := taskRepo.GetByID(c.Request().Context(), testTask.ID)
		assert.NoError(t, err)
		assert.Contains(t, task.Tags, "test-tag")
	})

	t.Run("RemoveTag", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/tasks/:id/tags/:tag")
		c.SetParamNames("org_id", "id", "tag")
		c.SetParamValues(testOrg.ID.String(), testTask.ID.String(), "test-tag")

		// Handle request
		err := taskHandlers.RemoveTag(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)

		// Verify tag was removed
		task, err := taskRepo.GetByID(c.Request().Context(), testTask.ID)
		assert.NoError(t, err)
		assert.NotContains(t, task.Tags, "test-tag")
	})

	t.Run("AssignTask", func(t *testing.T) {
		// Create request body
		reqBody := struct {
			UserID uuid.UUID `json:"user_id"`
		}{
			UserID: testUser.ID,
		}
		reqJSON, _ := json.Marshal(reqBody)

		// Create request
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/tasks/:id/assign")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testTask.ID.String())

		// Handle request
		err := taskHandlers.AssignTask(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response models.TaskResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testTask.ID, response.ID)
		assert.NotNil(t, response.AssignedTo)
		assert.Equal(t, testUser.ID, *response.AssignedTo)
	})

	t.Run("UpdateTaskStatus", func(t *testing.T) {
		// Create request body
		reqBody := struct {
			Status models.TaskStatus `json:"status"`
		}{
			Status: models.TaskStatusDone,
		}
		reqJSON, _ := json.Marshal(reqBody)

		// Create request
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/tasks/:id/status")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testTask.ID.String())

		// Handle request
		err := taskHandlers.UpdateTaskStatus(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response models.TaskResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testTask.ID, response.ID)
		assert.Equal(t, models.TaskStatusDone, response.Status)
	})

	t.Run("UnassignTask", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/tasks/:id/unassign")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testTask.ID.String())

		// Handle request
		err := taskHandlers.UnassignTask(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response models.TaskResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testTask.ID, response.ID)
		assert.Nil(t, response.AssignedTo)
	})

	t.Run("DeleteTask", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/tasks/:id")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testTask.ID.String())

		// Handle request
		err := taskHandlers.Delete(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)

		// Verify task is deleted
		_, err = taskRepo.GetByID(c.Request().Context(), testTask.ID)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrTaskNotFound, err)
	})
}
