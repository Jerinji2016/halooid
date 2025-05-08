package notification_test

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
	"github.com/Jerinji2016/halooid/backend/internal/test"
	"github.com/Jerinji2016/halooid/backend/pkg/middleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNotificationAPI(t *testing.T) {
	// Setup test environment
	tdb, prefix := test.SetupTestEnvironment(t)
	defer test.TeardownTestEnvironment(t, tdb, prefix)

	// Create test user, organization, project, and task
	testUser := tdb.CreateTestUser(t, prefix)
	testOrg := tdb.CreateTestOrganization(t, prefix, testUser.ID)
	testProject := tdb.CreateTestProject(t, prefix, testOrg.ID, testUser.ID)
	testTask := tdb.CreateTestTask(t, prefix, testProject.ID, testUser.ID)

	// Create repositories
	notificationRepo := repository.NewPostgresNotificationRepository(tdb.DB)
	userRepo := repository.NewPostgresUserRepository(tdb.DB)

	// Create services
	notificationService := notification.NewService(notificationRepo, userRepo)
	notificationHandlers := notification.NewHandlers(notificationService)

	// Setup Echo
	e := echo.New()
	g := e.Group("/api/v1/organizations/:org_id")

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
	notificationHandlers.RegisterRoutes(g)

	// Create a test notification
	testNotification := &models.Notification{
		ID:           uuid.New(),
		UserID:       testUser.ID,
		Type:         models.NotificationTypeTaskAssigned,
		Title:        fmt.Sprintf("%s Test Notification", prefix),
		Message:      "This is a test notification for integration tests",
		ResourceType: "task",
		ResourceID:   testTask.ID,
		IsRead:       false,
		CreatedAt:    time.Now(),
	}
	err := notificationRepo.Create(nil, testNotification)
	assert.NoError(t, err)

	t.Run("GetNotificationByID", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/notifications/:id")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testNotification.ID.String())

		// Handle request
		err := notificationHandlers.GetByID(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response models.NotificationResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testNotification.ID, response.ID)
		assert.Equal(t, testNotification.UserID, response.UserID)
		assert.Equal(t, testNotification.Type, response.Type)
		assert.Equal(t, testNotification.Title, response.Title)
		assert.Equal(t, testNotification.Message, response.Message)
		assert.Equal(t, testNotification.ResourceType, response.ResourceType)
		assert.Equal(t, testNotification.ResourceID, response.ResourceID)
		assert.Equal(t, testNotification.IsRead, response.IsRead)
	})

	t.Run("ListNotifications", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodGet, "/?page=1&page_size=10", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/notifications")
		c.SetParamNames("org_id")
		c.SetParamValues(testOrg.ID.String())

		// Handle request
		err := notificationHandlers.List(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check that notifications array exists and contains at least one notification
		notifications, ok := response["notifications"].([]interface{})
		assert.True(t, ok)
		assert.GreaterOrEqual(t, len(notifications), 1)

		// Check pagination
		pagination, ok := response["pagination"].(map[string]interface{})
		assert.True(t, ok)
		assert.GreaterOrEqual(t, int(pagination["total"].(float64)), 1)
	})

	t.Run("MarkAsRead", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodPut, "/", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/notifications/:id/read")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testNotification.ID.String())

		// Handle request
		err := notificationHandlers.MarkAsRead(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)

		// Verify notification is marked as read
		notification, err := notificationRepo.GetByID(c.Request().Context(), testNotification.ID)
		assert.NoError(t, err)
		assert.True(t, notification.IsRead)
		assert.NotNil(t, notification.ReadAt)
	})

	t.Run("MarkAllAsRead", func(t *testing.T) {
		// Create another unread notification
		anotherNotification := &models.Notification{
			ID:           uuid.New(),
			UserID:       testUser.ID,
			Type:         models.NotificationTypeTaskStatusUpdate,
			Title:        fmt.Sprintf("%s Another Test Notification", prefix),
			Message:      "This is another test notification for integration tests",
			ResourceType: "task",
			ResourceID:   testTask.ID,
			IsRead:       false,
			CreatedAt:    time.Now(),
		}
		err := notificationRepo.Create(nil, anotherNotification)
		assert.NoError(t, err)

		// Create request
		req := httptest.NewRequest(http.MethodPut, "/", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/notifications/read-all")
		c.SetParamNames("org_id")
		c.SetParamValues(testOrg.ID.String())

		// Handle request
		err = notificationHandlers.MarkAllAsRead(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)

		// Verify all notifications are marked as read
		params := models.NotificationListParams{
			UserID: testUser.ID,
			Page:   1,
			PageSize: 100,
		}
		isRead := false
		params.IsRead = &isRead
		
		unreadNotifications, _, err := notificationRepo.List(c.Request().Context(), params)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(unreadNotifications))
	})

	t.Run("Delete", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/notifications/:id")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testNotification.ID.String())

		// Handle request
		err := notificationHandlers.Delete(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)

		// Verify notification is deleted
		_, err = notificationRepo.GetByID(c.Request().Context(), testNotification.ID)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotificationNotFound, err)
	})
}
