package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestGCSClient_NewGCSClient(t *testing.T) {
	// This test requires GCP credentials, so we'll skip it in CI
	if testing.Short() {
		t.Skip("skipping GCS test in short mode")
	}

	ctx := context.Background()
	client, err := NewGCSClient(ctx, "test-project", "test-bucket")
	if err != nil {
		t.Skipf("GCS test skipped: %v", err)
	}
	defer client.Close()

	if client.projectID != "test-project" {
		t.Errorf("Expected project ID 'test-project', got '%s'", client.projectID)
	}
	if client.bucketName != "test-bucket" {
		t.Errorf("Expected bucket name 'test-bucket', got '%s'", client.bucketName)
	}
}

func TestMinIOClient_NewMinIOClient(t *testing.T) {
	// This test requires MinIO server, so we'll skip it in CI
	if testing.Short() {
		t.Skip("skipping MinIO test in short mode")
	}

	client, err := NewMinIOClient("localhost:9000", "minioadmin", "minioadmin", "test-bucket", false)
	if err != nil {
		t.Skipf("MinIO test skipped: %v", err)
	}
	defer client.Close()

	if client.bucketName != "test-bucket" {
		t.Errorf("Expected bucket name 'test-bucket', got '%s'", client.bucketName)
	}
}

func TestUploadManager_UploadFiles(t *testing.T) {
	// Create temporary test files
	tmpDir := t.TempDir()
	testFiles := []string{
		filepath.Join(tmpDir, "test1.txt"),
		filepath.Join(tmpDir, "test2.txt"),
	}

	// Create test files
	for _, file := range testFiles {
		if err := os.WriteFile(file, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Create a mock storage client for testing
	mockClient := &MockStorageClient{}
	manager := NewUploadManager(mockClient)
	defer manager.Close()

	ctx := context.Background()
	err := manager.UploadFiles(ctx, testFiles, "test-prefix")
	if err != nil {
		t.Errorf("UploadFiles failed: %v", err)
	}

	// Verify that all files were uploaded
	expectedUploads := len(testFiles)
	if mockClient.uploadCount != expectedUploads {
		t.Errorf("Expected %d uploads, got %d", expectedUploads, mockClient.uploadCount)
	}
}

// MockStorageClient for testing
type MockStorageClient struct {
	uploadCount int
}

func (m *MockStorageClient) UploadFile(ctx context.Context, localPath, remotePath string) error {
	m.uploadCount++
	return nil
}

func (m *MockStorageClient) Close() error {
	return nil
}