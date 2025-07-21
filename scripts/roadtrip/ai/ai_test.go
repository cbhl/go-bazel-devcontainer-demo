package ai

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestMockAIClient_AnalyzeVideo(t *testing.T) {
	client := NewMockAIClient()
	defer client.Close()

	// Test default response
	response, err := client.AnalyzeVideo(context.Background(), "test_video.mp4")
	if err != nil {
		t.Errorf("AnalyzeVideo failed: %v", err)
	}

	if response.VideoPath != "test_video.mp4" {
		t.Errorf("Expected video path 'test_video.mp4', got '%s'", response.VideoPath)
	}

	if response.HasMusic {
		t.Error("Expected HasMusic to be false for default response")
	}
}

func TestMockAIClient_AddMockResponse(t *testing.T) {
	client := NewMockAIClient()
	defer client.Close()

	// Add a custom mock response
	expectedResponse := &VideoAnalysisResponse{
		Description: "Test video with music",
		HasMusic:    true,
		Song: Song{
			Title:  "Test Song",
			Artist: "Test Artist",
		},
		VideoPath: "test_video.mp4",
	}

	client.AddMockResponse("test_video.mp4", expectedResponse)

	// Test that we get the custom response
	response, err := client.AnalyzeVideo(context.Background(), "test_video.mp4")
	if err != nil {
		t.Errorf("AnalyzeVideo failed: %v", err)
	}

	if response.Description != expectedResponse.Description {
		t.Errorf("Expected description '%s', got '%s'", expectedResponse.Description, response.Description)
	}

	if !response.HasMusic {
		t.Error("Expected HasMusic to be true")
	}

	if response.Song.Title != expectedResponse.Song.Title {
		t.Errorf("Expected song title '%s', got '%s'", expectedResponse.Song.Title, response.Song.Title)
	}
}

func TestAnalysisManager_AnalyzeVideos(t *testing.T) {
	client := NewMockAIClient()
	defer client.Close()

	manager := NewAnalysisManager(client)
	defer manager.Close()

	videoPaths := []string{
		"video1.mp4",
		"video2.mp4",
		"video3.mp4",
	}

	responses, err := manager.AnalyzeVideos(context.Background(), videoPaths)
	if err != nil {
		t.Errorf("AnalyzeVideos failed: %v", err)
	}

	if len(responses) != len(videoPaths) {
		t.Errorf("Expected %d responses, got %d", len(videoPaths), len(responses))
	}

	for i, response := range responses {
		if response.VideoPath != videoPaths[i] {
			t.Errorf("Expected video path '%s', got '%s'", videoPaths[i], response.VideoPath)
		}
	}
}

func TestPromptManager_NewPromptManager(t *testing.T) {
	// Test with actual template file
	templatePath := "../prompts/video_analysis.tmpl"
	
	manager, err := NewPromptManager(templatePath)
	if err != nil {
		t.Fatalf("Failed to create prompt manager: %v", err)
	}

	// Test rendering with sample data
	data := struct {
		VideoPath string
	}{
		VideoPath: "test_video.mp4",
	}

	rendered, err := manager.RenderPrompt(data)
	if err != nil {
		t.Fatalf("Failed to render prompt: %v", err)
	}

	if rendered == "" {
		t.Error("Rendered prompt should not be empty")
	}

	if !strings.Contains(rendered, "test_video.mp4") {
		t.Error("Rendered prompt should contain the video path")
	}

	if !strings.Contains(rendered, "JSON format") {
		t.Error("Rendered prompt should contain JSON format instructions")
	}
}

func TestSong_JSON(t *testing.T) {
	song := Song{
		Title:  "Test Song",
		Artist: "Test Artist",
	}

	// Test JSON marshaling
	data, err := json.Marshal(song)
	if err != nil {
		t.Errorf("Failed to marshal song: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Song
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal song: %v", err)
	}

	if unmarshaled.Title != song.Title {
		t.Errorf("Expected title '%s', got '%s'", song.Title, unmarshaled.Title)
	}

	if unmarshaled.Artist != song.Artist {
		t.Errorf("Expected artist '%s', got '%s'", song.Artist, unmarshaled.Artist)
	}
}

func TestGeminiAIClient_NewGeminiAIClient(t *testing.T) {
	// This test requires GEMINI_API_KEY to be set
	if testing.Short() {
		t.Skip("skipping Gemini client test in short mode")
	}

	// Test client creation (will fail if no API key, which is expected)
	client, err := NewGeminiAIClient(context.Background())
	if err != nil {
		// Expected if no API key is set
		if strings.Contains(err.Error(), "GEMINI_API_KEY") {
			t.Skip("GEMINI_API_KEY not set, skipping test")
		}
		t.Fatalf("Unexpected error creating Gemini client: %v", err)
	}
	defer client.Close()

	// If we get here, the client was created successfully
	// We could add more tests here if we had a test video file
}