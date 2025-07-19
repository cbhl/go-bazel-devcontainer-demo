package ai

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"text/template"
)

// VideoAnalysisRequest represents a request for video analysis
type VideoAnalysisRequest struct {
	VideoPath string
}

// VideoAnalysisResponse represents the response from video analysis
type VideoAnalysisResponse struct {
	Description    string `json:"description"`
	HasMusic       bool   `json:"has_music"`
	Transcript     string `json:"transcript"`
	Song           Song   `json:"song"`
	WebSearchSong  Song   `json:"web_search_song"`
	URLs           URLs   `json:"urls"`
	VideoPath      string `json:"video_path"`
}

// Song represents song information
type Song struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

// URLs represents links to music platforms
type URLs struct {
	YouTube string `json:"youtube"`
	Spotify string `json:"spotify"`
}

// AIClient interface for AI analysis
type AIClient interface {
	AnalyzeVideo(ctx context.Context, videoPath string) (*VideoAnalysisResponse, error)
	Close() error
}

// MockAIClient implements AIClient for testing
type MockAIClient struct {
	responses map[string]*VideoAnalysisResponse
}

// NewMockAIClient creates a new mock AI client
func NewMockAIClient() *MockAIClient {
	return &MockAIClient{
		responses: make(map[string]*VideoAnalysisResponse),
	}
}

// AddMockResponse adds a mock response for a video path
func (m *MockAIClient) AddMockResponse(videoPath string, response *VideoAnalysisResponse) {
	m.responses[videoPath] = response
}

// AnalyzeVideo performs mock video analysis
func (m *MockAIClient) AnalyzeVideo(ctx context.Context, videoPath string) (*VideoAnalysisResponse, error) {
	// Check if we have a mock response
	if response, exists := m.responses[videoPath]; exists {
		slog.Info("Using mock response for video", "path", videoPath)
		return response, nil
	}

	// Return a default mock response
	response := &VideoAnalysisResponse{
		Description: "Mock video analysis",
		HasMusic:    false,
		Transcript:  "No transcript available",
		Song: Song{
			Title:  "",
			Artist: "",
		},
		WebSearchSong: Song{
			Title:  "",
			Artist: "",
		},
		URLs: URLs{
			YouTube: "",
			Spotify: "",
		},
		VideoPath: videoPath,
	}

	slog.Info("Generated mock response for video", "path", videoPath)
	return response, nil
}

// Close closes the mock AI client
func (m *MockAIClient) Close() error {
	return nil
}

// AnalysisManager handles batch video analysis
type AnalysisManager struct {
	client AIClient
}

// NewAnalysisManager creates a new analysis manager
func NewAnalysisManager(client AIClient) *AnalysisManager {
	return &AnalysisManager{
		client: client,
	}
}

// AnalyzeVideos performs analysis on multiple videos
func (am *AnalysisManager) AnalyzeVideos(ctx context.Context, videoPaths []string) ([]*VideoAnalysisResponse, error) {
	var responses []*VideoAnalysisResponse

	for i, videoPath := range videoPaths {
		slog.Info("Analyzing video", "progress", fmt.Sprintf("%d/%d", i+1, len(videoPaths)), "path", videoPath)

		response, err := am.client.AnalyzeVideo(ctx, videoPath)
		if err != nil {
			slog.Error("Failed to analyze video", "path", videoPath, "error", err)
			return nil, fmt.Errorf("failed to analyze %s: %w", videoPath, err)
		}

		responses = append(responses, response)
	}

	slog.Info("Video analysis completed", "total_videos", len(videoPaths))
	return responses, nil
}

// Close closes the analysis manager and underlying client
func (am *AnalysisManager) Close() error {
	return am.client.Close()
}

// PromptManager handles prompt template loading and rendering
type PromptManager struct {
	template *template.Template
}

// NewPromptManager creates a new prompt manager
func NewPromptManager(templatePath string) (*PromptManager, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &PromptManager{
		template: tmpl,
	}, nil
}

// RenderPrompt renders a prompt with the given data
func (pm *PromptManager) RenderPrompt(data interface{}) (string, error) {
	var buf bytes.Buffer
	err := pm.template.ExecuteTemplate(&buf, "video_analysis.tmpl", data)
	if err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}

// TODO: Implement real Gemini 2.5 Flash client when API access is available
// For now, we use the mock client for testing and development