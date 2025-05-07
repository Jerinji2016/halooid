package timeentry_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/Jerinji2016/halooid/backend/internal/taskodex/timeentry"
	"github.com/Jerinji2016/halooid/backend/internal/test"
	"github.com/Jerinji2016/halooid/backend/pkg/middleware"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestTimeEntryAPI(t *testing.T) {
	// Setup test environment
	tdb, prefix := test.SetupTestEnvironment(t)
	defer test.TeardownTestEnvironment(t, tdb, prefix)

	// Create test user, organization, project, and task
	testUser := tdb.CreateTestUser(t, prefix)
	testOrg := tdb.CreateTestOrganization(t, prefix, testUser.ID)
	testProject := tdb.CreateTestProject(t, prefix, testOrg.ID, testUser.ID)
	testTask := tdb.CreateTestTask(t, prefix, testProject.ID, testUser.ID)

	// Create repositories
	timeEntryRepo := repository.NewPostgresTimeEntryRepository(tdb.DB)
	taskRepo := repository.NewPostgresTaskRepository(tdb.DB)
	userRepo := repository.NewPostgresUserRepository(tdb.DB)

	// Create service and handlers
	timeEntryService := timeentry.NewService(timeEntryRepo, taskRepo, userRepo)
	timeEntryHandlers := timeentry.NewHandlers(timeEntryService)

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
	timeEntryHandlers.RegisterRoutes(g, mockRBACMiddleware)

	t.Run("CreateTimeEntry", func(t *testing.T) {
		// Create request body
		startTime := time.Now().Add(-1 * time.Hour)
		endTime := time.Now()
		duration := int(endTime.Sub(startTime).Minutes())
		
		reqBody := models.TimeEntryRequest{
			TaskID:          testTask.ID,
			StartTime:       startTime,
			EndTime:         &endTime,
			DurationMinutes: &duration,
			Description:     "Test time entry for integration tests",
		}
		reqJSON, _ := json.Marshal(reqBody)

		// Create request
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/time-entries")
		c.SetParamNames("org_id")
		c.SetParamValues(testOrg.ID.String())

		// Handle request
		err := timeEntryHandlers.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Parse response
		var response models.TimeEntryResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testTask.ID, response.TaskID)
		assert.Equal(t, testUser.ID, response.UserID)
		assert.Equal(t, reqBody.Description, response.Description)
		assert.NotNil(t, response.DurationMinutes)
		assert.Equal(t, *reqBody.DurationMinutes, *response.DurationMinutes)
	})

	// Create a test time entry for the remaining tests
	testTimeEntry := tdb.CreateTestTimeEntry(t, prefix, testTask.ID, testUser.ID)

	t.Run("GetTimeEntryByID", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/time-entries/:id")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testTimeEntry.ID.String())

		// Handle request
		err := timeEntryHandlers.GetByID(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response models.TimeEntryResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testTimeEntry.ID, response.ID)
		assert.Equal(t, testTask.ID, response.TaskID)
		assert.Equal(t, testUser.ID, response.UserID)
		assert.Equal(t, testTimeEntry.Description, response.Description)
	})

	t.Run("ListTimeEntries", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodGet, "/?page=1&page_size=10", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/time-entries")
		c.SetParamNames("org_id")
		c.SetParamValues(testOrg.ID.String())

		// Handle request
		err := timeEntryHandlers.List(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check that time entries array exists and contains at least one entry
		timeEntries, ok := response["time_entries"].([]interface{})
		assert.True(t, ok)
		assert.GreaterOrEqual(t, len(timeEntries), 1)

		// Check pagination
		pagination, ok := response["pagination"].(map[string]interface{})
		assert.True(t, ok)
		assert.GreaterOrEqual(t, int(pagination["total"].(float64)), 1)
	})

	t.Run("UpdateTimeEntry", func(t *testing.T) {
		// Create request body
		startTime := time.Now().Add(-2 * time.Hour)
		endTime := time.Now().Add(-1 * time.Hour)
		duration := int(endTime.Sub(startTime).Minutes())
		
		reqBody := models.TimeEntryRequest{
			TaskID:          testTask.ID,
			StartTime:       startTime,
			EndTime:         &endTime,
			DurationMinutes: &duration,
			Description:     "Updated time entry description",
		}
		reqJSON, _ := json.Marshal(reqBody)

		// Create request
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/time-entries/:id")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testTimeEntry.ID.String())

		// Handle request
		err := timeEntryHandlers.Update(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response models.TimeEntryResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testTimeEntry.ID, response.ID)
		assert.Equal(t, testTask.ID, response.TaskID)
		assert.Equal(t, testUser.ID, response.UserID)
		assert.Equal(t, reqBody.Description, response.Description)
		assert.NotNil(t, response.DurationMinutes)
		assert.Equal(t, *reqBody.DurationMinutes, *response.DurationMinutes)
	})

	t.Run("StartTimer", func(t *testing.T) {
		// Create a new task for this test to avoid conflicts with running timers
		newTask := tdb.CreateTestTask(t, fmt.Sprintf("%s_timer", prefix), testProject.ID, testUser.ID)
		
		// Create request body
		reqBody := struct {
			TaskID      string `json:"task_id"`
			Description string `json:"description"`
		}{
			TaskID:      newTask.ID.String(),
			Description: "Timer started for integration test",
		}
		reqJSON, _ := json.Marshal(reqBody)

		// Create request
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/time-entries/start")
		c.SetParamNames("org_id")
		c.SetParamValues(testOrg.ID.String())

		// Handle request
		err := timeEntryHandlers.StartTimer(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response models.TimeEntryResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, newTask.ID, response.TaskID)
		assert.Equal(t, testUser.ID, response.UserID)
		assert.Equal(t, reqBody.Description, response.Description)
		assert.Nil(t, response.EndTime)
		assert.Nil(t, response.DurationMinutes)

		// Store the timer ID for the next test
		timerID := response.ID

		// Test StopTimer
		t.Run("StopTimer", func(t *testing.T) {
			// Create request body
			reqBody := struct {
				TaskID string `json:"task_id"`
			}{
				TaskID: newTask.ID.String(),
			}
			reqJSON, _ := json.Marshal(reqBody)

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			// Set up context
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/organizations/:org_id/taskodex/time-entries/stop")
			c.SetParamNames("org_id")
			c.SetParamValues(testOrg.ID.String())

			// Handle request
			err := timeEntryHandlers.StopTimer(c)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)

			// Parse response
			var response models.TimeEntryResponse
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, timerID, response.ID)
			assert.Equal(t, newTask.ID, response.TaskID)
			assert.Equal(t, testUser.ID, response.UserID)
			assert.NotNil(t, response.EndTime)
			assert.NotNil(t, response.DurationMinutes)
		})
	})

	t.Run("GetRunningTimers", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/time-entries/running")
		c.SetParamNames("org_id")
		c.SetParamValues(testOrg.ID.String())

		// Handle request
		err := timeEntryHandlers.GetRunningTimers(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response - should be an array
		var response []models.TimeEntryResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
	})

	t.Run("Aggregate", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodGet, "/?group_by=day", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/time-entries/aggregate")
		c.SetParamNames("org_id")
		c.SetParamValues(testOrg.ID.String())

		// Handle request
		err := timeEntryHandlers.Aggregate(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response - should be an array of aggregations
		var response []models.TimeEntryAggregation
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(response), 1)
	})

	t.Run("DeleteTimeEntry", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/time-entries/:id")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testTimeEntry.ID.String())

		// Handle request
		err := timeEntryHandlers.Delete(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)

		// Verify time entry is deleted
		_, err = timeEntryRepo.GetByID(c.Request().Context(), testTimeEntry.ID)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrTimeEntryNotFound, err)
	})
}
