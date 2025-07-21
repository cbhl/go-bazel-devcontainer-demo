package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

// CSVExporter handles conversion from JSON to CSV format
type CSVExporter struct {
	writer *csv.Writer
}

// NewCSVExporter creates a new CSV exporter
func NewCSVExporter(w io.Writer) *CSVExporter {
	return &CSVExporter{
		writer: csv.NewWriter(w),
	}
}

// WriteHeader writes the CSV header
func (e *CSVExporter) WriteHeader() error {
	header := []string{
		"description",
		"has_music",
		"transcript",
		"song_title",
		"song_artist",
		"web_search_song_title",
		"web_search_song_artist",
		"youtube_url",
		"spotify_url",
		"video_path",
	}
	return e.writer.Write(header)
}

// WriteRecord writes a single record to CSV
func (e *CSVExporter) WriteRecord(record map[string]interface{}) error {
	row := []string{
		getString(record, "description"),
		getBoolString(record, "has_music"),
		getString(record, "transcript"),
		getString(record, "song_title"),
		getString(record, "song_artist"),
		getString(record, "web_search_song_title"),
		getString(record, "web_search_song_artist"),
		getString(record, "youtube_url"),
		getString(record, "spotify_url"),
		getString(record, "video_path"),
	}
	return e.writer.Write(row)
}

// Flush flushes the CSV writer
func (e *CSVExporter) Flush() {
	e.writer.Flush()
}

// getString safely extracts a string value from a map
func getString(record map[string]interface{}, key string) string {
	if val, ok := record[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// getBoolString safely extracts a boolean value and converts it to string
func getBoolString(record map[string]interface{}, key string) string {
	if val, ok := record[key]; ok {
		if b, ok := val.(bool); ok {
			if b {
				return "true"
			}
			return "false"
		}
	}
	return "false"
}

// RelaxedJSONParser handles parsing of non-standard JSON outputs
type RelaxedJSONParser struct {
	// Regex patterns for extracting JSON from prose text
	jsonPatterns []*regexp.Regexp
}

// NewRelaxedJSONParser creates a new relaxed JSON parser
func NewRelaxedJSONParser() *RelaxedJSONParser {
	return &RelaxedJSONParser{
		jsonPatterns: []*regexp.Regexp{
			regexp.MustCompile(`\{[^{}]*"description"[^{}]*\}`),
			regexp.MustCompile(`\{.*?\}`),
		},
	}
}

// ParseJSON attempts to parse JSON with fallback mechanisms
func (p *RelaxedJSONParser) ParseJSON(input string) (map[string]interface{}, error) {
	// Try standard JSON parsing first
	var result map[string]interface{}
	err := json.Unmarshal([]byte(input), &result)
	if err == nil {
		return result, nil
	}

	slog.Warn("Standard JSON parsing failed, attempting relaxed parsing", "error", err)

	// Try to extract JSON from prose text
	for _, pattern := range p.jsonPatterns {
		matches := pattern.FindAllString(input, -1)
		for _, match := range matches {
			err := json.Unmarshal([]byte(match), &result)
			if err == nil {
				slog.Info("Successfully parsed JSON using relaxed parsing", "pattern", pattern.String())
				return result, nil
			}
		}
	}

	// Manual parsing fallback
	return p.parseManually(input)
}

// parseManually attempts to extract key-value pairs manually
func (p *RelaxedJSONParser) parseManually(input string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	
	// Extract common patterns
	patterns := map[string]*regexp.Regexp{
		"description": regexp.MustCompile(`description\s+of\s+"([^"]*)"`),
		"has_music":   regexp.MustCompile(`music:\s*(true|false)`),
		"transcript":  regexp.MustCompile(`"transcript"\s*:\s*"([^"]*)"`),
		"song_title":  regexp.MustCompile(`title\s+is\s+"([^"]*)"`),
		"song_artist": regexp.MustCompile(`by\s+"([^"]*)"`),
		"video_path":  regexp.MustCompile(`path\s+is\s+"([^"]*)"`),
	}

	for key, pattern := range patterns {
		matches := pattern.FindStringSubmatch(input)
		if len(matches) > 1 {
			if key == "has_music" {
				result[key] = matches[1] == "true"
			} else {
				result[key] = matches[1]
			}
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("failed to extract any data from input")
	}

	slog.Info("Successfully parsed data using manual parsing", "extracted_fields", len(result))
	return result, nil
}

// ExportManager handles the complete export process
type ExportManager struct {
	parser   *RelaxedJSONParser
	exporter *CSVExporter
}

// NewExportManager creates a new export manager
func NewExportManager(output io.Writer) *ExportManager {
	return &ExportManager{
		parser:   NewRelaxedJSONParser(),
		exporter: NewCSVExporter(output),
	}
}

// ExportFromFile exports data from a JSON file to CSV
func (em *ExportManager) ExportFromFile(inputFile string) error {
	// Read input file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	return em.ExportFromString(string(data))
}

// ExportFromString exports data from a JSON string to CSV
func (em *ExportManager) ExportFromString(input string) error {
	// Write CSV header
	if err := em.exporter.WriteHeader(); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Split input into lines (one JSON object per line)
	lines := strings.Split(strings.TrimSpace(input), "\n")
	
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		slog.Info("Processing line", "line_number", i+1)

		// Parse JSON
		record, err := em.parser.ParseJSON(line)
		if err != nil {
			slog.Warn("Failed to parse line, skipping", "line_number", i+1, "error", err)
			continue
		}

		// Write to CSV
		if err := em.exporter.WriteRecord(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	em.exporter.Flush()
	slog.Info("Export completed successfully")
	return nil
}