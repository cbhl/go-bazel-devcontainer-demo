package video

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// VideoProcessor handles video processing operations
type VideoProcessor struct {
	ffmpegPath string
}

// NewVideoProcessor creates a new video processor instance
func NewVideoProcessor() (*VideoProcessor, error) {
	// Check if ffmpeg is available
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, fmt.Errorf("ffmpeg not found in PATH: %w", err)
	}

	return &VideoProcessor{
		ffmpegPath: ffmpegPath,
	}, nil
}

// SplitVideo splits a video into chunks
func (vp *VideoProcessor) SplitVideo(inputFile, outputDir, startTime, endTime string, chunkDuration int) error {
	// Validate input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputFile)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Parse time range
	start, err := parseTime(startTime)
	if err != nil {
		return fmt.Errorf("invalid start time: %w", err)
	}

	end, err := parseTime(endTime)
	if err != nil {
		return fmt.Errorf("invalid end time: %w", err)
	}

	if start >= end {
		return fmt.Errorf("start time must be before end time")
	}

	// Calculate number of chunks
	duration := end - start
	numChunks := int(duration) / chunkDuration
	if int(duration)%chunkDuration != 0 {
		numChunks++
	}

	fmt.Printf("Splitting video into %d chunks of %d seconds each\n", numChunks, chunkDuration)
	fmt.Printf("Time range: %s to %s (duration: %d seconds)\n", startTime, endTime, duration)

	// Split video into chunks
	for i := 0; i < numChunks; i++ {
		chunkStart := start + (i * chunkDuration)
		chunkEnd := chunkStart + chunkDuration
		if chunkEnd > end {
			chunkEnd = end
		}

		outputFile := filepath.Join(outputDir, fmt.Sprintf("chunk_%03d.mp4", i+1))
		
		fmt.Printf("Processing chunk %d/%d: %s\n", i+1, numChunks, outputFile)
		
		if err := vp.extractChunk(inputFile, outputFile, chunkStart, chunkEnd); err != nil {
			return fmt.Errorf("failed to extract chunk %d: %w", i+1, err)
		}
	}

	fmt.Printf("Successfully created %d video chunks in %s\n", numChunks, outputDir)
	return nil
}

// extractChunk extracts a single chunk from the video
func (vp *VideoProcessor) extractChunk(inputFile, outputFile string, start, end int) error {
	// Build ffmpeg command with re-encoding for reliable seeking
	args := []string{
		"-i", inputFile,
		"-ss", fmt.Sprintf("%d", start),
		"-t", fmt.Sprintf("%d", end-start),
		"-c:v", "libx264", // Use H.264 codec for compatibility
		"-preset", "fast", // Use fast preset for reasonable speed
		"-crf", "23", // Use constant rate factor for good quality
		"-avoid_negative_ts", "make_zero",
		"-y", // Overwrite output file if it exists
		outputFile,
	}

	cmd := exec.Command(vp.ffmpegPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Debug output removed for cleaner interface
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg command failed: %w", err)
	}

	return nil
}

// parseTime converts a time string (HH:MM:SS) to seconds
func parseTime(timeStr string) (int, error) {
	if timeStr == "" {
		return 0, nil
	}

	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("time format must be HH:MM:SS")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hours: %s", parts[0])
	}
	if hours < 0 || hours > 23 {
		return 0, fmt.Errorf("hours must be between 0 and 23: %d", hours)
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes: %s", parts[1])
	}
	if minutes < 0 || minutes > 59 {
		return 0, fmt.Errorf("minutes must be between 0 and 59: %d", minutes)
	}

	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf("invalid seconds: %s", parts[2])
	}
	if seconds < 0 || seconds > 59 {
		return 0, fmt.Errorf("seconds must be between 0 and 59: %d", seconds)
	}

	return hours*3600 + minutes*60 + seconds, nil
}

// GetVideoInfo returns basic information about a video file
func (vp *VideoProcessor) GetVideoInfo(inputFile string) (map[string]string, error) {
	// Use ffprobe to get video information
	ffprobePath, err := exec.LookPath("ffprobe")
	if err != nil {
		return nil, fmt.Errorf("ffprobe not found in PATH: %w", err)
	}

	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("input file does not exist: %s", inputFile)
	}

	args := []string{
		"-v", "quiet",
		"-show_entries", "format=duration",
		"-of", "csv=p=0",
		inputFile,
	}

	// Debug output removed for cleaner interface
	
	cmd := exec.Command(ffprobePath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get video info: %w (output: %s)", err, string(output))
	}

	duration := strings.TrimSpace(string(output))
	
	// Convert duration to HH:MM:SS format
	if durationFloat, err := strconv.ParseFloat(duration, 64); err == nil {
		hours := int(durationFloat) / 3600
		minutes := int(durationFloat) % 3600 / 60
		seconds := int(durationFloat) % 60
		duration = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	info := make(map[string]string)
	info["duration"] = duration
	info["duration_seconds"] = strings.TrimSpace(string(output))

	return info, nil
}