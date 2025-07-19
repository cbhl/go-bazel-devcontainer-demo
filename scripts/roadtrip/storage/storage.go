package storage

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// StorageClient interface for both GCS and MinIO
type StorageClient interface {
	UploadFile(ctx context.Context, localPath, remotePath string) error
	Close() error
}

// MockStorageClient implements StorageClient for testing
type MockStorageClient struct {
	uploadCount int
	bucketName  string
	projectID   string
}

// NewMockStorageClient creates a new mock storage client for testing
func NewMockStorageClient(projectID, bucketName string) *MockStorageClient {
	return &MockStorageClient{
		bucketName: bucketName,
		projectID:  projectID,
	}
}

// UploadFile simulates uploading a file
func (m *MockStorageClient) UploadFile(ctx context.Context, localPath, remotePath string) error {
	// Verify file exists
	if _, err := os.Stat(localPath); err != nil {
		return fmt.Errorf("file does not exist: %s", localPath)
	}

	m.uploadCount++
	slog.Info("Mock upload completed", "local", localPath, "remote", remotePath, "count", m.uploadCount)
	return nil
}

// Close closes the mock storage client
func (m *MockStorageClient) Close() error {
	return nil
}

// UploadManager handles batch uploads with progress tracking
type UploadManager struct {
	client StorageClient
}

// NewUploadManager creates a new upload manager
func NewUploadManager(client StorageClient) *UploadManager {
	return &UploadManager{
		client: client,
	}
}

// UploadFiles uploads multiple files with progress tracking
func (um *UploadManager) UploadFiles(ctx context.Context, localPaths []string, remotePrefix string) error {
	totalFiles := len(localPaths)
	slog.Info("Starting batch upload", "total_files", totalFiles, "remote_prefix", remotePrefix)

	for i, localPath := range localPaths {
		// Create remote path
		fileName := filepath.Base(localPath)
		remotePath := filepath.Join(remotePrefix, fileName)
		// Normalize path separators for cloud storage
		remotePath = strings.ReplaceAll(remotePath, "\\", "/")

		slog.Info("Uploading file", "progress", fmt.Sprintf("%d/%d", i+1, totalFiles), "file", fileName)

		if err := um.client.UploadFile(ctx, localPath, remotePath); err != nil {
			slog.Error("Failed to upload file", "file", localPath, "error", err)
			return fmt.Errorf("failed to upload %s: %w", localPath, err)
		}
	}

	slog.Info("Batch upload completed successfully", "total_files", totalFiles)
	return nil
}

// Close closes the upload manager and underlying client
func (um *UploadManager) Close() error {
	return um.client.Close()
}

// TODO: Implement real GCS and MinIO clients when dependencies are available
// For now, we use the mock client for testing and development