package env

// GCSConfig holds Google Cloud Storage configuration
type GCSConfig struct {
	ProjectID string
	Zone      string
	Bucket    string
}

// GeminiConfig holds Gemini API configuration
type GeminiConfig struct {
	APIKey string
}

// VideoConfig holds video processing configuration
type VideoConfig struct {
	DefaultChunkDuration int
	DefaultOutputDir     string
}

// DefaultVideoConfig returns default video processing settings
func DefaultVideoConfig() *VideoConfig {
	return &VideoConfig{
		DefaultChunkDuration: 30,
		DefaultOutputDir:     "out",
	}
}

// DefaultGCSConfig returns default GCS settings
func DefaultGCSConfig() *GCSConfig {
	return &GCSConfig{
		ProjectID: "gen-lang-client-0629405113",
		Zone:      "northamerica-northeast2",
		Bucket:    "cbhl-roadtrip-202507",
	}
}