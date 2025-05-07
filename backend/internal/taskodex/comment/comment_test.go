package comment_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/notification"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/Jerinji2016/halooid/backend/internal/taskodex/comment"
	"github.com/Jerinji2016/halooid/backend/internal/test"
	"github.com/Jerinji2016/halooid/backend/pkg/middleware"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCommentAPI(t *testing.T) {
	// Setup test environment
	tdb, prefix := test.SetupTestEnvironment(t)
	defer test.TeardownTestEnvironment(t, tdb, prefix)

	// Create test user, organization, project, and task
	testUser := tdb.CreateTestUser(t, prefix)
	testOrg := tdb.CreateTestOrganization(t, prefix, testUser.ID)
	testProject := tdb.CreateTestProject(t, prefix, testOrg.ID, testUser.ID)
	testTask := tdb.CreateTestTask(t, prefix, testProject.ID, testUser.ID)

	// Create repositories
	commentRepo := repository.NewPostgresCommentRepository(tdb.DB)
	taskRepo := repository.NewPostgresTaskRepository(tdb.DB)
	userRepo := repository.NewPostgresUserRepository(tdb.DB)
	notificationRepo := repository.NewPostgresNotificationRepository(tdb.DB)

	// Create services
	notificationService := notification.NewService(notificationRepo, userRepo)
	commentService := comment.NewService(commentRepo, taskRepo, userRepo, notificationService)
	commentHandlers := comment.NewHandlers(commentService)

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
	commentHandlers.RegisterRoutes(g, mockRBACMiddleware)

	t.Run("CreateComment", func(t *testing.T) {
		// Create request body
		reqBody := models.CommentRequest{
			TaskID:  testTask.ID,
			Content: fmt.Sprintf("This is a test comment for task %s", testTask.ID),
		}
		reqJSON, _ := json.Marshal(reqBody)

		// Create request
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/comments")
		c.SetParamNames("org_id")
		c.SetParamValues(testOrg.ID.String())

		// Handle request
		err := commentHandlers.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Parse response
		var response models.CommentResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testTask.ID, response.TaskID)
		assert.Equal(t, testUser.ID, response.UserID)
		assert.Equal(t, reqBody.Content, response.Content)
	})

	// Create a test comment for the remaining tests
	testComment := tdb.CreateTestComment(t, prefix, testTask.ID, testUser.ID)

	t.Run("GetCommentByID", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/comments/:id")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testComment.ID.String())

		// Handle request
		err := commentHandlers.GetByID(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response models.CommentResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testComment.ID, response.ID)
		assert.Equal(t, testTask.ID, response.TaskID)
		assert.Equal(t, testUser.ID, response.UserID)
		assert.Equal(t, testComment.Content, response.Content)
	})

	t.Run("ListComments", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?task_id=%s&page=1&page_size=10", testTask.ID), nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/comments")
		c.SetParamNames("org_id")
		c.SetParamValues(testOrg.ID.String())

		// Handle request
		err := commentHandlers.List(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check that comments array exists and contains at least one comment
		comments, ok := response["comments"].([]interface{})
		assert.True(t, ok)
		assert.GreaterOrEqual(t, len(comments), 1)

		// Check pagination
		pagination, ok := response["pagination"].(map[string]interface{})
		assert.True(t, ok)
		assert.GreaterOrEqual(t, int(pagination["total"].(float64)), 1)
	})

	t.Run("UpdateComment", func(t *testing.T) {
		// Create request body
		reqBody := models.CommentRequest{
			TaskID:  testTask.ID,
			Content: fmt.Sprintf("This is an updated comment for task %s", testTask.ID),
		}
		reqJSON, _ := json.Marshal(reqBody)

		// Create request
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/comments/:id")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testComment.ID.String())

		// Handle request
		err := commentHandlers.Update(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response models.CommentResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testComment.ID, response.ID)
		assert.Equal(t, testTask.ID, response.TaskID)
		assert.Equal(t, testUser.ID, response.UserID)
		assert.Equal(t, reqBody.Content, response.Content)
	})

	t.Run("DeleteComment", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/comments/:id")
		c.SetParamNames("org_id", "id")
		c.SetParamValues(testOrg.ID.String(), testComment.ID.String())

		// Handle request
		err := commentHandlers.Delete(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)

		// Verify comment is deleted
		_, err = commentRepo.GetByID(c.Request().Context(), testComment.ID)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrCommentNotFound, err)
	})
}
