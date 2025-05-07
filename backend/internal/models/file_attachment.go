package models

import (
	"time"

	"github.com/google/uuid"
)

// FileAttachment represents a file attached to a task
type FileAttachment struct {
	ID          uuid.UUID `json:"id" db:"id"`
	TaskID      uuid.UUID `json:"task_id" db:"task_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	FileName    string    `json:"file_name" db:"file_name"`
	FileSize    int64     `json:"file_size" db:"file_size"`
	ContentType string    `json:"content_type" db:"content_type"`
	StoragePath string    `json:"storage_path" db:"storage_path"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	
	// Related entities
	User        *User     `json:"user,omitempty" db:"-"`
	Task        *Task     `json:"task,omitempty" db:"-"`
}

// FileAttachmentResponse represents the file attachment data returned to clients
type FileAttachmentResponse struct {
	ID          uuid.UUID      `json:"id"`
	TaskID      uuid.UUID      `json:"task_id"`
	UserID      uuid.UUID      `json:"user_id"`
	FileName    string         `json:"file_name"`
	FileSize    int64          `json:"file_size"`
	ContentType string         `json:"content_type"`
	DownloadURL string         `json:"download_url"`
	CreatedAt   time.Time      `json:"created_at"`
	
	// Related entities
	User        *UserResponse  `json:"user,omitempty"`
	Task        *TaskResponse  `json:"task,omitempty"`
}

// ToResponse converts a FileAttachment to a FileAttachmentResponse
func (f *FileAttachment) ToResponse(baseURL string) FileAttachmentResponse {
	response := FileAttachmentResponse{
		ID:          f.ID,
		TaskID:      f.TaskID,
		UserID:      f.UserID,
		FileName:    f.FileName,
		FileSize:    f.FileSize,
		ContentType: f.ContentType,
		DownloadURL: baseURL + "/api/v1/taskodex/files/" + f.ID.String(),
		CreatedAt:   f.CreatedAt,
	}
	
	if f.User != nil {
		userResponse := f.User.ToResponse()
		response.User = &userResponse
	}
	
	if f.Task != nil {
		taskResponse := f.Task.ToResponse()
		response.Task = &taskResponse
	}
	
	return response
}

// NewFileAttachment creates a new FileAttachment
func NewFileAttachment(taskID, userID uuid.UUID, fileName, contentType string, fileSize int64, storagePath string) *FileAttachment {
	return &FileAttachment{
		ID:          uuid.New(),
		TaskID:      taskID,
		UserID:      userID,
		FileName:    fileName,
		FileSize:    fileSize,
		ContentType: contentType,
		StoragePath: storagePath,
		CreatedAt:   time.Now(),
	}
}

// FileAttachmentListParams represents the parameters for listing file attachments
type FileAttachmentListParams struct {
	TaskID    *uuid.UUID `query:"task_id"`
	UserID    *uuid.UUID `query:"user_id"`
	SortBy    string     `query:"sort_by" default:"created_at"`
	SortOrder string     `query:"sort_order" default:"desc"`
	Page      int        `query:"page" default:"1"`
	PageSize  int        `query:"page_size" default:"20"`
}
