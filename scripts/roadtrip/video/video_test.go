package video

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseTime(t *testing.T) {
	tests := []struct {
		name     string
		timeStr  string
		expected int
		hasError bool
	}{
		{"empty string", "", 0, false},
		{"valid time", "01:30:45", 5445, false},
		{"zero time", "00:00:00", 0, false},
		{"invalid format", "1:30", 0, true},
		{"invalid hours", "25:30:45", 0, true},
		{"invalid minutes", "01:60:45", 0, true},
		{"invalid seconds", "01:30:60", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseTime(tt.timeStr)
			if tt.hasError && err == nil {
				t.Errorf("parseTime() expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("parseTime() unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("parseTime() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestNewVideoProcessor(t *testing.T) {
	processor, err := NewVideoProcessor()
	if err != nil {
		// This test might fail if ffmpeg is not installed, which is expected in CI
		t.Logf("VideoProcessor creation failed (expected if ffmpeg not available): %v", err)
		return
	}

	if processor == nil {
		t.Error("NewVideoProcessor() returned nil")
	}

	if processor.ffmpegPath == "" {
		t.Error("ffmpegPath is empty")
	}
}

func TestVideoProcessor_SplitVideo_Validation(t *testing.T) {
	processor, err := NewVideoProcessor()
	if err != nil {
		t.Skipf("Skipping test - ffmpeg not available: %v", err)
	}

	// Test with non-existent input file
	err = processor.SplitVideo("nonexistent.mp4", "out", "00:00:00", "00:01:00", 30)
	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	// Test with invalid time range
	err = processor.SplitVideo("test.mp4", "out", "00:01:00", "00:00:00", 30)
	if err == nil {
		t.Error("Expected error for invalid time range")
	}
}

func TestVideoProcessor_GetVideoInfo(t *testing.T) {
	processor, err := NewVideoProcessor()
	if err != nil {
		t.Skipf("Skipping test - ffmpeg not available: %v", err)
	}

	// Test with non-existent file
	_, err = processor.GetVideoInfo("nonexistent.mp4")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestCreateOutputDirectory(t *testing.T) {
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "test", "nested", "dir")

	// Test directory creation
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Errorf("Failed to create directory: %v", err)
	}

	// Verify directory exists
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Error("Directory was not created")
	}
}