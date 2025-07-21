package env

import (
	"os"
	"testing"
)

func TestDefaultGCSConfig(t *testing.T) {
	config := DefaultGCSConfig()
	if config.ProjectID != "gen-lang-client-0629405113" {
		t.Errorf("Expected ProjectID to be 'gen-lang-client-0629405113', got '%s'", config.ProjectID)
	}
	if config.Zone != "northamerica-northeast2" {
		t.Errorf("Expected Zone to be 'northamerica-northeast2', got '%s'", config.Zone)
	}
	if config.Bucket != "cbhl-roadtrip-202507" {
		t.Errorf("Expected Bucket to be 'cbhl-roadtrip-202507', got '%s'", config.Bucket)
	}
}

func TestDefaultVideoConfig(t *testing.T) {
	config := DefaultVideoConfig()
	if config.DefaultChunkDuration != 30 {
		t.Errorf("Expected DefaultChunkDuration to be 30, got %d", config.DefaultChunkDuration)
	}
	if config.DefaultOutputDir != "out" {
		t.Errorf("Expected DefaultOutputDir to be 'out', got '%s'", config.DefaultOutputDir)
	}
}

func TestConfigValidation(t *testing.T) {
	// Test with missing API key
	config := &Config{}
	err := config.Validate()
	if err == nil {
		t.Error("Expected error when API key is missing")
	}

	// Test with valid API key
	config.GeminiAPIKey = "test-key"
	err = config.Validate()
	if err != nil {
		t.Errorf("Expected no error with valid API key, got: %v", err)
	}
}

func TestLoadWithDefaults(t *testing.T) {
	// Save original environment
	originalAPIKey := os.Getenv("GEMINI_API_KEY")
	originalProjectID := os.Getenv("GCP_PROJECT_ID")
	originalZone := os.Getenv("GCP_ZONE")
	originalBucket := os.Getenv("GCS_BUCKET")

	// Clean up environment
	os.Unsetenv("GEMINI_API_KEY")
	os.Unsetenv("GCP_PROJECT_ID")
	os.Unsetenv("GCP_ZONE")
	os.Unsetenv("GCS_BUCKET")

	// Test that Load returns error without API key
	_, err := Load()
	if err == nil {
		t.Error("Expected error when GEMINI_API_KEY is not set")
	}

	// Set API key and test defaults
	os.Setenv("GEMINI_API_KEY", "test-key")
	config, err := Load()
	if err != nil {
		t.Errorf("Expected no error with API key set, got: %v", err)
	}

	if config.GeminiAPIKey != "test-key" {
		t.Errorf("Expected API key to be 'test-key', got '%s'", config.GeminiAPIKey)
	}
	if config.GCPProjectID != "gen-lang-client-0629405113" {
		t.Errorf("Expected default project ID, got '%s'", config.GCPProjectID)
	}
	if config.GCPZone != "northamerica-northeast2" {
		t.Errorf("Expected default zone, got '%s'", config.GCPZone)
	}
	if config.GCSBucket != "cbhl-roadtrip-202507" {
		t.Errorf("Expected default bucket, got '%s'", config.GCSBucket)
	}

	// Restore original environment
	if originalAPIKey != "" {
		os.Setenv("GEMINI_API_KEY", originalAPIKey)
	}
	if originalProjectID != "" {
		os.Setenv("GCP_PROJECT_ID", originalProjectID)
	}
	if originalZone != "" {
		os.Setenv("GCP_ZONE", originalZone)
	}
	if originalBucket != "" {
		os.Setenv("GCS_BUCKET", originalBucket)
	}
}