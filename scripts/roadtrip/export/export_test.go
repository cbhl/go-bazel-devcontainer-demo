package export

import (
	"bytes"
	"strings"
	"testing"
)

func TestCSVExporter_WriteHeader(t *testing.T) {
	var buf bytes.Buffer
	exporter := NewCSVExporter(&buf)

	err := exporter.WriteHeader()
	if err != nil {
		t.Errorf("WriteHeader failed: %v", err)
	}

	exporter.Flush()
	output := buf.String()
	expected := "description,has_music,transcript,song_title,song_artist,web_search_song_title,web_search_song_artist,youtube_url,spotify_url,video_path\n"
	if output != expected {
		t.Errorf("Expected header '%s', got '%s'", expected, output)
	}
}

func TestCSVExporter_WriteRecord(t *testing.T) {
	var buf bytes.Buffer
	exporter := NewCSVExporter(&buf)

	// Write header first
	exporter.WriteHeader()

	// Test record
	record := map[string]interface{}{
		"description": "Test video",
		"has_music":   true,
		"transcript":  "Hello world",
		"song_title":  "Test Song",
		"song_artist": "Test Artist",
		"video_path":  "gs://bucket/video.mp4",
	}

	err := exporter.WriteRecord(record)
	if err != nil {
		t.Errorf("WriteRecord failed: %v", err)
	}

	exporter.Flush()
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}

	// Check data line
	dataLine := lines[1]
	expected := "Test video,true,Hello world,Test Song,Test Artist,,,,,gs://bucket/video.mp4"
	if dataLine != expected {
		t.Errorf("Expected data line '%s', got '%s'", expected, dataLine)
	}
}

func TestRelaxedJSONParser_ParseJSON_Valid(t *testing.T) {
	parser := NewRelaxedJSONParser()

	validJSON := `{
		"description": "Test video",
		"has_music": true,
		"transcript": "Hello world",
		"song": {"title": "Test Song", "artist": "Test Artist"},
		"video_path": "gs://bucket/video.mp4"
	}`

	result, err := parser.ParseJSON(validJSON)
	if err != nil {
		t.Errorf("ParseJSON failed: %v", err)
	}

	if result["description"] != "Test video" {
		t.Errorf("Expected description 'Test video', got '%v'", result["description"])
	}

	if result["has_music"] != true {
		t.Errorf("Expected has_music true, got %v", result["has_music"])
	}
}

func TestRelaxedJSONParser_ParseJSON_Malformed(t *testing.T) {
	parser := NewRelaxedJSONParser()

	malformedJSON := `This is some prose text with JSON embedded: {"description": "Test video", "has_music": true} and more text`

	result, err := parser.ParseJSON(malformedJSON)
	if err != nil {
		t.Errorf("ParseJSON failed: %v", err)
	}

	if result["description"] != "Test video" {
		t.Errorf("Expected description 'Test video', got '%v'", result["description"])
	}

	if result["has_music"] != true {
		t.Errorf("Expected has_music true, got %v", result["has_music"])
	}
}

func TestRelaxedJSONParser_ParseJSON_ManualFallback(t *testing.T) {
	parser := NewRelaxedJSONParser()

	// Very malformed input that requires manual parsing
	malformedInput := `Here's what I found: The video has a description of "Amazing roadtrip video" and contains music: true. The song title is "Roadtrip Anthem" by "Travel Band". Video path is "gs://bucket/roadtrip.mp4".`

	result, err := parser.ParseJSON(malformedInput)
	if err != nil {
		t.Errorf("ParseJSON failed: %v", err)
	}

	if result["description"] != "Amazing roadtrip video" {
		t.Errorf("Expected description 'Amazing roadtrip video', got '%v'", result["description"])
	}

	if result["song_title"] != "Roadtrip Anthem" {
		t.Errorf("Expected song_title 'Roadtrip Anthem', got '%v'", result["song_title"])
	}
}

func TestExportManager_ExportFromString(t *testing.T) {
	var buf bytes.Buffer
	manager := NewExportManager(&buf)

	input := `{"description": "Video 1", "has_music": true, "video_path": "gs://bucket/video1.mp4"}
{"description": "Video 2", "has_music": false, "video_path": "gs://bucket/video2.mp4"}`

	err := manager.ExportFromString(input)
	if err != nil {
		t.Errorf("ExportFromString failed: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	
	if len(lines) != 3 { // header + 2 data lines
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}

	// Check that CSV is properly formatted
	if !strings.Contains(lines[0], "description,has_music") {
		t.Error("CSV header is missing expected fields")
	}
}

func TestExportManager_ExportFromString_WithMalformedJSON(t *testing.T) {
	var buf bytes.Buffer
	manager := NewExportManager(&buf)

	input := `{"description": "Video 1", "has_music": true, "video_path": "gs://bucket/video1.mp4"}
This is malformed JSON that should be skipped
{"description": "Video 3", "has_music": false, "video_path": "gs://bucket/video3.mp4"}`

	err := manager.ExportFromString(input)
	if err != nil {
		t.Errorf("ExportFromString failed: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	
	// Should have header + 2 valid data lines (malformed line skipped)
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}
}

func TestGetString(t *testing.T) {
	record := map[string]interface{}{
		"test_string": "hello",
		"test_int":    123,
	}

	if result := getString(record, "test_string"); result != "hello" {
		t.Errorf("Expected 'hello', got '%s'", result)
	}

	if result := getString(record, "test_int"); result != "" {
		t.Errorf("Expected empty string for non-string value, got '%s'", result)
	}

	if result := getString(record, "missing"); result != "" {
		t.Errorf("Expected empty string for missing key, got '%s'", result)
	}
}

func TestGetBoolString(t *testing.T) {
	record := map[string]interface{}{
		"test_true":  true,
		"test_false": false,
		"test_int":   123,
	}

	if result := getBoolString(record, "test_true"); result != "true" {
		t.Errorf("Expected 'true', got '%s'", result)
	}

	if result := getBoolString(record, "test_false"); result != "false" {
		t.Errorf("Expected 'false', got '%s'", result)
	}

	if result := getBoolString(record, "test_int"); result != "false" {
		t.Errorf("Expected 'false' for non-bool value, got '%s'", result)
	}

	if result := getBoolString(record, "missing"); result != "false" {
		t.Errorf("Expected 'false' for missing key, got '%s'", result)
	}
}