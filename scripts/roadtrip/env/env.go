package env

import (
	"fmt"
	"os"
)

// Config holds environment-based configuration
type Config struct {
	GeminiAPIKey string
	GCPProjectID string
	GCPZone      string
	GCSBucket    string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		GCPProjectID: os.Getenv("GCP_PROJECT_ID"),
		GCPZone:      os.Getenv("GCP_ZONE"),
		GCSBucket:    os.Getenv("GCS_BUCKET"),
	}

	// Set defaults if not provided
	if config.GCPProjectID == "" {
		config.GCPProjectID = "gen-lang-client-0629405113"
	}
	if config.GCPZone == "" {
		config.GCPZone = "northamerica-northeast2"
	}
	if config.GCSBucket == "" {
		config.GCSBucket = "cbhl-roadtrip-202507"
	}

	// Validate required environment variables
	if config.GeminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
	}

	return config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.GeminiAPIKey == "" {
		return fmt.Errorf("GEMINI_API_KEY is required")
	}
	return nil
}