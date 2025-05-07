package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// FileStorage defines the interface for file storage operations
type FileStorage interface {
	// SaveFile saves a file to storage
	SaveFile(file *multipart.FileHeader, taskID, userID uuid.UUID) (string, error)
	
	// GetFile retrieves a file from storage
	GetFile(storagePath string) (io.ReadCloser, error)
	
	// DeleteFile deletes a file from storage
	DeleteFile(storagePath string) error
}

// LocalFileStorage implements FileStorage using the local file system
type LocalFileStorage struct {
	basePath string
}

// NewLocalFileStorage creates a new LocalFileStorage
func NewLocalFileStorage(basePath string) FileStorage {
	// Create the base directory if it doesn't exist
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		err = os.MkdirAll(basePath, 0755)
		if err != nil {
			panic(fmt.Sprintf("failed to create storage directory: %v", err))
		}
	}
	
	return &LocalFileStorage{
		basePath: basePath,
	}
}

// SaveFile saves a file to storage
func (s *LocalFileStorage) SaveFile(file *multipart.FileHeader, taskID, userID uuid.UUID) (string, error) {
	// Create a unique filename
	fileExt := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s_%s_%s%s",
		taskID.String(),
		userID.String(),
		time.Now().Format("20060102150405"),
		fileExt,
	)
	
	// Create the directory structure
	dirPath := filepath.Join(s.basePath, taskID.String())
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
	}
	
	// Full path to the file
	filePath := filepath.Join(dirPath, fileName)
	
	// Open the source file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open source file: %w", err)
	}
	defer src.Close()
	
	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()
	
	// Copy the file content
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}
	
	// Return the relative storage path
	storagePath := filepath.Join(taskID.String(), fileName)
	return storagePath, nil
}

// GetFile retrieves a file from storage
func (s *LocalFileStorage) GetFile(storagePath string) (io.ReadCloser, error) {
	// Clean the path to prevent directory traversal attacks
	cleanPath := filepath.Clean(storagePath)
	if strings.Contains(cleanPath, "..") {
		return nil, fmt.Errorf("invalid file path")
	}
	
	// Full path to the file
	filePath := filepath.Join(s.basePath, cleanPath)
	
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %w", err)
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	
	return file, nil
}

// DeleteFile deletes a file from storage
func (s *LocalFileStorage) DeleteFile(storagePath string) error {
	// Clean the path to prevent directory traversal attacks
	cleanPath := filepath.Clean(storagePath)
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("invalid file path")
	}
	
	// Full path to the file
	filePath := filepath.Join(s.basePath, cleanPath)
	
	// Delete the file
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %w", err)
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}
	
	return nil
}
