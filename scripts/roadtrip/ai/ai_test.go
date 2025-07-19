package ai

import (
	"context"
	"encoding/json"
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
	// This test requires the template file to exist
	// For now, we'll skip it in CI
	if testing.Short() {
		t.Skip("skipping prompt manager test in short mode")
	}

	// TODO: Add test with actual template file
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