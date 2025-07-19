package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/api/option"
)

// StorageClient interface for both GCS and MinIO
type StorageClient interface {
	UploadFile(ctx context.Context, localPath, remotePath string) error
	Close() error
}

// GCSClient implements StorageClient for Google Cloud Storage
type GCSClient struct {
	client     *storage.Client
	bucketName string
	projectID  string
}

// MinIOClient implements StorageClient for MinIO
type MinIOClient struct {
	client     *minio.Client
	bucketName string
}

// MockStorageClient implements StorageClient for testing
type MockStorageClient struct {
	UploadCount int
	bucketName  string
	projectID   string
}

// NewGCSClient creates a new GCS client
func NewGCSClient(ctx context.Context, projectID, bucketName string) (*GCSClient, error) {
	client, err := storage.NewClient(ctx, option.WithScopes(storage.ScopeReadWrite))
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	return &GCSClient{
		client:     client,
		bucketName: bucketName,
		projectID:  projectID,
	}, nil
}

// NewMinIOClient creates a new MinIO client for testing
func NewMinIOClient(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (*MinIOClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return &MinIOClient{
		client:     client,
		bucketName: bucketName,
	}, nil
}

// NewMockStorageClient creates a new mock storage client for testing
func NewMockStorageClient(projectID, bucketName string) *MockStorageClient {
	return &MockStorageClient{
		bucketName: bucketName,
		projectID:  projectID,
	}
}

// UploadFile uploads a file to GCS
func (g *GCSClient) UploadFile(ctx context.Context, localPath, remotePath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", localPath, err)
	}
	defer file.Close()

	bucket := g.client.Bucket(g.bucketName)
	obj := bucket.Object(remotePath)
	writer := obj.NewWriter(ctx)

	// Set metadata
	writer.ObjectAttrs.Metadata = map[string]string{
		"uploaded-by": "roadtrip-tool",
		"upload-time": time.Now().Format(time.RFC3339),
	}

	if _, err := io.Copy(writer, file); err != nil {
		writer.Close()
		return fmt.Errorf("failed to copy file to GCS: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close GCS writer: %w", err)
	}

	slog.Info("Uploaded file to GCS", "local", localPath, "remote", remotePath)
	return nil
}

// Close closes the GCS client
func (g *GCSClient) Close() error {
	return g.client.Close()
}

// UploadFile uploads a file to MinIO
func (m *MinIOClient) UploadFile(ctx context.Context, localPath, remotePath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", localPath, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	_, err = m.client.PutObject(ctx, m.bucketName, remotePath, file, fileInfo.Size(), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
		UserMetadata: map[string]string{
			"uploaded-by": "roadtrip-tool",
			"upload-time": time.Now().Format(time.RFC3339),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	slog.Info("Uploaded file to MinIO", "local", localPath, "remote", remotePath)
	return nil
}

// Close closes the MinIO client
func (m *MinIOClient) Close() error {
	return nil // MinIO client doesn't need explicit closing
}

// UploadFile simulates uploading a file
func (m *MockStorageClient) UploadFile(ctx context.Context, localPath, remotePath string) error {
	// Verify file exists
	if _, err := os.Stat(localPath); err != nil {
		return fmt.Errorf("file does not exist: %s", localPath)
	}

	m.UploadCount++
	slog.Info("Mock upload completed", "local", localPath, "remote", remotePath, "count", m.UploadCount)
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