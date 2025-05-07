package fileattachment_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Jerinji2016/halooid/backend/internal/config"
	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/notification"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/Jerinji2016/halooid/backend/internal/storage"
	"github.com/Jerinji2016/halooid/backend/internal/taskodex/fileattachment"
	"github.com/Jerinji2016/halooid/backend/internal/test"
	"github.com/Jerinji2016/halooid/backend/pkg/middleware"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestFileAttachmentAPI(t *testing.T) {
	// Setup test environment
	tdb, prefix := test.SetupTestEnvironment(t)
	defer test.TeardownTestEnvironment(t, tdb, prefix)

	// Create test user, organization, project, and task
	testUser := tdb.CreateTestUser(t, prefix)
	testOrg := tdb.CreateTestOrganization(t, prefix, testUser.ID)
	testProject := tdb.CreateTestProject(t, prefix, testOrg.ID, testUser.ID)
	testTask := tdb.CreateTestTask(t, prefix, testProject.ID, testUser.ID)

	// Load config
	cfg, err := config.LoadConfig("../../../config/test.yaml")
	assert.NoError(t, err)

	// Create storage directory if it doesn't exist
	err = os.MkdirAll(cfg.Storage.BasePath, 0755)
	assert.NoError(t, err)

	// Create repositories
	fileAttachmentRepo := repository.NewPostgresFileAttachmentRepository(tdb.DB)
	taskRepo := repository.NewPostgresTaskRepository(tdb.DB)
	userRepo := repository.NewPostgresUserRepository(tdb.DB)
	notificationRepo := repository.NewPostgresNotificationRepository(tdb.DB)

	// Create services
	fileStorage := storage.NewLocalFileStorage(cfg.Storage.BasePath)
	notificationService := notification.NewService(notificationRepo, userRepo)
	fileAttachmentService := fileattachment.NewService(
		fileAttachmentRepo,
		taskRepo,
		userRepo,
		fileStorage,
		notificationService,
		cfg.API.BaseURL,
	)
	fileAttachmentHandlers := fileattachment.NewHandlers(fileAttachmentService, fileStorage)

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
	fileAttachmentHandlers.RegisterRoutes(g, mockRBACMiddleware)

	t.Run("UploadFile", func(t *testing.T) {
		// Create a temporary test file
		tempFile, err := os.CreateTemp("", "test-file-*.txt")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		// Write some content to the file
		content := []byte("This is a test file for file attachment API testing.")
		_, err = tempFile.Write(content)
		assert.NoError(t, err)
		tempFile.Close()

		// Create a new multipart request
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add task_id field
		err = writer.WriteField("task_id", testTask.ID.String())
		assert.NoError(t, err)

		// Add file field
		part, err := writer.CreateFormFile("file", filepath.Base(tempFile.Name()))
		assert.NoError(t, err)

		// Open the file again for reading
		file, err := os.Open(tempFile.Name())
		assert.NoError(t, err)
		defer file.Close()

		// Copy the file content to the form field
		_, err = io.Copy(part, file)
		assert.NoError(t, err)

		// Close the writer
		err = writer.Close()
		assert.NoError(t, err)

		// Create request
		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		rec := httptest.NewRecorder()

		// Set up context
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/organizations/:org_id/taskodex/files")
		c.SetParamNames("org_id")
		c.SetParamValues(testOrg.ID.String())

		// Handle request
		err = fileAttachmentHandlers.Upload(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Parse response
		var response models.FileAttachmentResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testTask.ID, response.TaskID)
		assert.Equal(t, testUser.ID, response.UserID)
		assert.Equal(t, filepath.Base(tempFile.Name()), response.FileName)
		assert.NotEmpty(t, response.DownloadURL)

		// Store the file attachment ID for the remaining tests
		fileAttachmentID := response.ID

		t.Run("GetFileAttachmentByID", func(t *testing.T) {
			// Create request
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			// Set up context
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/organizations/:org_id/taskodex/files/:id")
			c.SetParamNames("org_id", "id")
			c.SetParamValues(testOrg.ID.String(), fileAttachmentID.String())

			// Handle request
			err := fileAttachmentHandlers.GetByID(c)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)

			// Parse response
			var response models.FileAttachmentResponse
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, fileAttachmentID, response.ID)
			assert.Equal(t, testTask.ID, response.TaskID)
			assert.Equal(t, testUser.ID, response.UserID)
			assert.Equal(t, filepath.Base(tempFile.Name()), response.FileName)
			assert.NotEmpty(t, response.DownloadURL)
		})

		t.Run("ListFileAttachments", func(t *testing.T) {
			// Create request
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?task_id=%s&page=1&page_size=10", testTask.ID), nil)
			rec := httptest.NewRecorder()

			// Set up context
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/organizations/:org_id/taskodex/files")
			c.SetParamNames("org_id")
			c.SetParamValues(testOrg.ID.String())

			// Handle request
			err := fileAttachmentHandlers.List(c)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)

			// Parse response
			var response map[string]interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Check that file_attachments array exists and contains at least one file attachment
			fileAttachments, ok := response["file_attachments"].([]interface{})
			assert.True(t, ok)
			assert.GreaterOrEqual(t, len(fileAttachments), 1)

			// Check pagination
			pagination, ok := response["pagination"].(map[string]interface{})
			assert.True(t, ok)
			assert.GreaterOrEqual(t, int(pagination["total"].(float64)), 1)
		})

		t.Run("DeleteFileAttachment", func(t *testing.T) {
			// Create request
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			rec := httptest.NewRecorder()

			// Set up context
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/organizations/:org_id/taskodex/files/:id")
			c.SetParamNames("org_id", "id")
			c.SetParamValues(testOrg.ID.String(), fileAttachmentID.String())

			// Handle request
			err := fileAttachmentHandlers.Delete(c)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNoContent, rec.Code)

			// Verify file attachment is deleted
			_, err = fileAttachmentRepo.GetByID(c.Request().Context(), fileAttachmentID)
			assert.Error(t, err)
			assert.Equal(t, repository.ErrFileAttachmentNotFound, err)
		})
	})
}
