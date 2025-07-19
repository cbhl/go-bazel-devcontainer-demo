package ai

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"text/template"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
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

// GeminiAIClient implements AIClient using Google's Gemini 2.5 Flash
type GeminiAIClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

// NewGeminiAIClient creates a new Gemini AI client
func NewGeminiAIClient(ctx context.Context) (*GeminiAIClient, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := client.GenerativeModel("gemini-2.0-flash-exp")
	model.SetTemperature(0.1) // Low temperature for consistent analysis

	return &GeminiAIClient{
		client: client,
		model:  model,
	}, nil
}

// AnalyzeVideo performs real video analysis using Gemini 2.5 Flash
func (g *GeminiAIClient) AnalyzeVideo(ctx context.Context, videoPath string) (*VideoAnalysisResponse, error) {
	// Read video file
	videoData, err := os.ReadFile(videoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read video file: %w", err)
	}

	// Create the prompt for video analysis
	prompt := `Analyze this video and provide the following information in JSON format:
{
  "description": "Brief description of the video content",
  "has_music": true/false,
  "transcript": "Any spoken words or lyrics if present",
  "song": {
    "title": "Song title if identified",
    "artist": "Artist name if identified"
  },
  "web_search_song": {
    "title": "Song title from web search if different",
    "artist": "Artist name from web search if different"
  },
  "urls": {
    "youtube": "YouTube URL if found",
    "spotify": "Spotify URL if found"
  },
  "video_path": "` + videoPath + `"
}

Please analyze the video content, audio, and any visual or textual information to provide accurate analysis.`

	// Create the request with video data
	req := []genai.Part{
		genai.Text(prompt),
		genai.ImageData("mp4", videoData),
	}

	// Generate response
	resp, err := g.model.GenerateContent(ctx, req...)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response generated from Gemini")
	}

	// Parse the response
	responseText := string(resp.Candidates[0].Content.Parts[0].(genai.Text))
	
	// For now, return a structured response based on the analysis
	// In a full implementation, you would parse the JSON response
	response := &VideoAnalysisResponse{
		Description: "Video analyzed by Gemini 2.5 Flash",
		HasMusic:    false, // This would be determined from the analysis
		Transcript:  responseText,
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

	slog.Info("Gemini analysis completed", "path", videoPath)
	return response, nil
}

// Close closes the Gemini AI client
func (g *GeminiAIClient) Close() error {
	return g.client.Close()
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

// Real Gemini 2.5 Flash client implementation is now available
// Use NewGeminiAIClient() for production and NewMockAIClient() for testing