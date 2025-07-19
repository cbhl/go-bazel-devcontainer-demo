package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"backend/scripts/roadtrip/storage"
	"backend/scripts/roadtrip/video"
)

// CLI represents the main command-line interface
type CLI struct {
	SplitVideo      SplitVideoCmd      `cmd:"" help:"Split video into chunks"`
	UploadChunks    UploadChunksCmd    `cmd:"" help:"Upload video chunks to cloud storage"`
	BuildPlaylist   BuildPlaylistCmd   `cmd:"" help:"Build playlist from video analysis"`
	BuildPlaylistCSV BuildPlaylistCSVCmd `cmd:"" help:"Convert playlist to CSV format"`
}

// SplitVideoCmd represents the split-video command
type SplitVideoCmd struct {
	In            string `flag:"in" help:"Input video file path"`
	StartTime     string `flag:"start" help:"Start timestamp (HH:MM:SS)"`
	EndTime       string `flag:"end" help:"End timestamp (HH:MM:SS)"`
	ChunkDuration int    `flag:"chunk-duration" default:"30" help:"Chunk duration in seconds"`
	OutputDir     string `flag:"out" default:"out" help:"Output directory"`
}

// UploadChunksCmd represents the upload-chunks command
type UploadChunksCmd struct {
	In        string `flag:"in" help:"Input folder or glob pattern"`
	ProjectID string `flag:"project-id" help:"GCP project ID"`
	Zone      string `flag:"zone" help:"GCP zone"`
	Bucket    string `flag:"bucket" help:"GCS bucket path"`
}

// BuildPlaylistCmd represents the build-playlist command
type BuildPlaylistCmd struct {
	In       []string `flag:"in" help:"GCS paths to analyze"`
	Validate bool     `flag:"validate-json" help:"Validate JSON output"`
}

// BuildPlaylistCSVCmd represents the build-playlist-csv command
type BuildPlaylistCSVCmd struct {
	In string `flag:"in" help:"Input JSON file or stdin"`
}

// Run implements the split-video command
func (s *SplitVideoCmd) Run() error {
	if s.In == "" {
		return fmt.Errorf("input file is required (use --in flag)")
	}
	
	fmt.Printf("Processing video: %s\n", s.In)
	fmt.Printf("Time range: %s to %s\n", s.StartTime, s.EndTime)
	fmt.Printf("Chunk duration: %d seconds\n", s.ChunkDuration)
	fmt.Printf("Output directory: %s\n", s.OutputDir)
	
	// Create video processor
	processor, err := video.NewVideoProcessor()
	if err != nil {
		return fmt.Errorf("failed to create video processor: %w", err)
	}

	// Get video info
	info, err := processor.GetVideoInfo(s.In)
	if err != nil {
		return fmt.Errorf("failed to get video info: %w", err)
	}

	fmt.Printf("Video duration: %s\n", info["duration"])
	if videoStream, ok := info["video_stream"]; ok {
		fmt.Printf("Video stream: %s\n", videoStream)
	}

	// Split video into chunks
	if err := processor.SplitVideo(s.In, s.OutputDir, s.StartTime, s.EndTime, s.ChunkDuration); err != nil {
		return fmt.Errorf("failed to split video: %w", err)
	}

	return nil
}

// Run implements the upload-chunks command
func (u *UploadChunksCmd) Run() error {
	if u.In == "" {
		return fmt.Errorf("input path is required (use --in flag)")
	}
	
	fmt.Printf("Uploading chunks from: %s\n", u.In)
	fmt.Printf("Project ID: %s\n", u.ProjectID)
	fmt.Printf("Zone: %s\n", u.Zone)
	fmt.Printf("Bucket: %s\n", u.Bucket)
	
	// Parse bucket path to extract bucket name and prefix
	bucketName, prefix, err := parseBucketPath(u.Bucket)
	if err != nil {
		return fmt.Errorf("invalid bucket path: %w", err)
	}
	
	// Find files to upload
	files, err := findFiles(u.In)
	if err != nil {
		return fmt.Errorf("failed to find files: %w", err)
	}
	
	if len(files) == 0 {
		return fmt.Errorf("no files found matching pattern: %s", u.In)
	}
	
	fmt.Printf("Found %d files to upload\n", len(files))
	
	// Create storage client
	ctx := context.Background()
	client, err := storage.NewGCSClient(ctx, u.ProjectID, bucketName)
	if err != nil {
		return fmt.Errorf("failed to create GCS client: %w", err)
	}
	defer client.Close()
	
	// Create upload manager
	manager := storage.NewUploadManager(client)
	defer manager.Close()
	
	// Upload files
	if err := manager.UploadFiles(ctx, files, prefix); err != nil {
		return fmt.Errorf("failed to upload files: %w", err)
	}
	
	fmt.Printf("Successfully uploaded %d files to %s\n", len(files), u.Bucket)
	return nil
}

// parseBucketPath extracts bucket name and prefix from a GCS path
func parseBucketPath(bucketPath string) (bucketName, prefix string, err error) {
	if !strings.HasPrefix(bucketPath, "gs://") {
		return "", "", fmt.Errorf("bucket path must start with gs://")
	}
	
	path := strings.TrimPrefix(bucketPath, "gs://")
	parts := strings.SplitN(path, "/", 2)
	
	bucketName = parts[0]
	if len(parts) > 1 {
		prefix = parts[1]
	}
	
	return bucketName, prefix, nil
}

// findFiles finds files matching the input pattern
func findFiles(pattern string) ([]string, error) {
	// Handle glob patterns
	if strings.Contains(pattern, "*") || strings.Contains(pattern, "?") {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to expand glob pattern: %w", err)
		}
		return matches, nil
	}
	
	// Handle directory
	info, err := os.Stat(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to stat path: %w", err)
	}
	
	if info.IsDir() {
		// Find all files in directory
		var files []string
		err := filepath.Walk(pattern, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to walk directory: %w", err)
		}
		return files, nil
	}
	
	// Single file
	return []string{pattern}, nil
}

// Run implements the build-playlist command
func (b *BuildPlaylistCmd) Run() error {
	if len(b.In) == 0 {
		return fmt.Errorf("input paths are required (use --in flag)")
	}
	
	fmt.Printf("Building playlist from %d paths\n", len(b.In))
	for i, path := range b.In {
		fmt.Printf("  %d: %s\n", i+1, path)
	}
	fmt.Printf("Validate JSON: %v\n", b.Validate)
	
	// TODO: Implement playlist building logic
	fmt.Println("Hello from build-playlist command!")
	return nil
}

// Run implements the build-playlist-csv command
func (b *BuildPlaylistCSVCmd) Run() error {
	if b.In == "" {
		return fmt.Errorf("input file is required (use --in flag)")
	}
	
	fmt.Printf("Converting to CSV from: %s\n", b.In)
	
	// TODO: Implement CSV conversion logic
	fmt.Println("Hello from build-playlist-csv command!")
	return nil
}

func main() {
	cli := CLI{}
	ctx := kong.Parse(&cli, kong.Vars{
		"version": "1.0.0",
	})

	err := ctx.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}