package test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Jerinji2016/halooid/backend/internal/config"
	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

// TestDB represents a test database connection
type TestDB struct {
	DB *sqlx.DB
}

// NewTestDB creates a new test database connection
func NewTestDB(t *testing.T) *TestDB {
	// Load test configuration
	cfg, err := config.LoadConfig("../../config/test.yaml")
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	// Connect to the database
	db, err := sqlx.Connect("postgres", cfg.Database.URL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	return &TestDB{DB: db}
}

// Close closes the database connection
func (tdb *TestDB) Close() {
	if tdb.DB != nil {
		tdb.DB.Close()
	}
}

// CleanupTestData cleans up test data from the database
func (tdb *TestDB) CleanupTestData(t *testing.T, prefix string) {
	// List of tables to clean up in reverse order of dependencies
	tables := []string{
		"taskodex.task_time_entries",
		"taskodex.task_file_attachments",
		"taskodex.comment_mentions",
		"taskodex.task_comments",
		"taskodex.task_tags",
		"taskodex.tasks",
		"taskodex.projects",
		"notifications",
		"users",
		"organizations",
	}

	for _, table := range tables {
		_, err := tdb.DB.Exec(fmt.Sprintf("DELETE FROM %s WHERE id::text LIKE '%s%%'", table, prefix))
		if err != nil {
			t.Logf("Warning: Failed to clean up table %s: %v", table, err)
		}
	}
}

// CreateTestUser creates a test user
func (tdb *TestDB) CreateTestUser(t *testing.T, prefix string) *models.User {
	userRepo := repository.NewPostgresUserRepository(tdb.DB)

	userID := uuid.New()
	user := &models.User{
		ID:        userID,
		Email:     fmt.Sprintf("%s_test_user@example.com", prefix),
		FirstName: "Test",
		LastName:  "User",
		Username:  fmt.Sprintf("%s_testuser", prefix),
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := userRepo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return user
}

// CreateTestOrganization creates a test organization
func (tdb *TestDB) CreateTestOrganization(t *testing.T, prefix string, createdBy uuid.UUID) *models.Organization {
	orgRepo := repository.NewPostgresOrganizationRepository(tdb.DB)

	orgID := uuid.New()
	org := &models.Organization{
		ID:          orgID,
		Name:        fmt.Sprintf("%s Test Organization", prefix),
		Description: "Test organization for integration tests",
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := orgRepo.Create(context.Background(), org)
	if err != nil {
		t.Fatalf("Failed to create test organization: %v", err)
	}

	return org
}

// CreateTestProject creates a test project
func (tdb *TestDB) CreateTestProject(t *testing.T, prefix string, orgID, createdBy uuid.UUID) *models.Project {
	projectRepo := repository.NewPostgresProjectRepository(tdb.DB)

	projectID := uuid.New()
	project := &models.Project{
		ID:             projectID,
		OrganizationID: orgID,
		Name:           fmt.Sprintf("%s Test Project", prefix),
		Description:    "Test project for integration tests",
		Status:         models.ProjectStatusPlanning,
		CreatedBy:      createdBy,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := projectRepo.Create(context.Background(), project)
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}

	return project
}

// CreateTestTask creates a test task
func (tdb *TestDB) CreateTestTask(t *testing.T, prefix string, projectID, createdBy uuid.UUID) *models.Task {
	taskRepo := repository.NewPostgresTaskRepository(tdb.DB)

	taskID := uuid.New()
	task := &models.Task{
		ID:          taskID,
		ProjectID:   &projectID,
		Title:       fmt.Sprintf("%s Test Task", prefix),
		Description: "Test task for integration tests",
		Status:      models.TaskStatusTodo,
		Priority:    models.TaskPriorityMedium,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := taskRepo.Create(context.Background(), task)
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	return task
}

// CreateTestComment creates a test comment
func (tdb *TestDB) CreateTestComment(t *testing.T, prefix string, taskID, userID uuid.UUID) *models.Comment {
	commentRepo := repository.NewPostgresCommentRepository(tdb.DB)

	commentID := uuid.New()
	comment := &models.Comment{
		ID:        commentID,
		TaskID:    taskID,
		UserID:    userID,
		Content:   fmt.Sprintf("%s Test Comment", prefix),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := commentRepo.Create(context.Background(), comment)
	if err != nil {
		t.Fatalf("Failed to create test comment: %v", err)
	}

	return comment
}

// CreateTestTimeEntry creates a test time entry
func (tdb *TestDB) CreateTestTimeEntry(t *testing.T, prefix string, taskID, userID uuid.UUID) *models.TimeEntry {
	timeEntryRepo := repository.NewPostgresTimeEntryRepository(tdb.DB)

	timeEntryID := uuid.New()
	now := time.Now()
	duration := 60 // 60 minutes

	timeEntry := &models.TimeEntry{
		ID:              timeEntryID,
		TaskID:          taskID,
		UserID:          userID,
		StartTime:       now.Add(-time.Duration(duration) * time.Minute),
		EndTime:         &now,
		DurationMinutes: &duration,
		Description:     fmt.Sprintf("%s Test Time Entry", prefix),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err := timeEntryRepo.Create(context.Background(), timeEntry)
	if err != nil {
		t.Fatalf("Failed to create test time entry: %v", err)
	}

	return timeEntry
}

// MockContext creates a mock Echo context for testing
func MockContext(method, path string) (echo.Context, *echo.Echo) {
	e := echo.New()
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, e
}

// SetupTestEnvironment sets up the test environment
func SetupTestEnvironment(t *testing.T) (*TestDB, string) {
	// Generate a unique prefix for test data
	prefix := fmt.Sprintf("test_%s_", uuid.New().String()[:8])

	// Create a test database connection
	tdb := NewTestDB(t)

	return tdb, prefix
}

// TeardownTestEnvironment tears down the test environment
func TeardownTestEnvironment(t *testing.T, tdb *TestDB, prefix string) {
	// Clean up test data
	tdb.CleanupTestData(t, prefix)

	// Close the database connection
	tdb.Close()
}
