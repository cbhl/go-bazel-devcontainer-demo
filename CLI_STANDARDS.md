# CLI Development Standards

## Overview
All command-line tools in this repository should be built in Go using the Kong library for consistent, maintainable, and user-friendly interfaces.

## Kong Library

### Why Kong?
- **Type Safety**: Compile-time validation of command-line arguments
- **Automatic Help**: Built-in help generation with proper formatting
- **Subcommands**: Clean support for complex command hierarchies
- **Validation**: Built-in validation and custom validation support
- **Testing**: Easy to test with structured command definitions

## Command Structure

### Basic CLI Structure
```go
package main

import (
    "fmt"
    "os"
    
    "github.com/alecthomas/kong"
)

type CLI struct {
    Version bool `help:"Show version information"`
    
    SplitVideo SplitVideoCmd `cmd:"" help:"Split video into chunks"`
    UploadChunks UploadChunksCmd `cmd:"" help:"Upload video chunks to cloud storage"`
    BuildPlaylist BuildPlaylistCmd `cmd:"" help:"Build playlist from video analysis"`
    BuildPlaylistCSV BuildPlaylistCSVCmd `cmd:"" help:"Convert playlist to CSV format"`
}

type SplitVideoCmd struct {
    InputFile    string `arg:"" help:"Input video file path"`
    StartTime    string `flag:"start" help:"Start timestamp (HH:MM:SS)"`
    EndTime      string `flag:"end" help:"End timestamp (HH:MM:SS)"`
    ChunkDuration int    `flag:"duration" default:"30" help:"Chunk duration in seconds"`
    OutputDir    string `flag:"output" default:"out" help:"Output directory"`
    Verbose      bool   `flag:"verbose" help:"Enable verbose output"`
}

func (s *SplitVideoCmd) Run() error {
    // Implementation here
    return nil
}

func main() {
    cli := CLI{}
    ctx := kong.Parse(&cli)
    
    if cli.Version {
        fmt.Println("roadtrip v1.0.0")
        return
    }
    
    err := ctx.Run()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

## Command Design Patterns

### 1. Flag Naming Conventions
- Use kebab-case for flag names: `--input-file`, `--chunk-duration`
- Use descriptive names that clearly indicate purpose
- Provide sensible defaults where appropriate
- Always include help text for every flag

### 2. Argument Validation
```go
type SplitVideoCmd struct {
    InputFile string `arg:"" validate:"file"`
    StartTime string `flag:"start" validate:"time"`
    EndTime   string `flag:"end" validate:"time"`
}

// Custom validation
func (s *SplitVideoCmd) Validate() error {
    if s.StartTime >= s.EndTime {
        return fmt.Errorf("start time must be before end time")
    }
    return nil
}
```

### 3. Error Handling
```go
func (s *SplitVideoCmd) Run() error {
    // Validate inputs
    if err := s.validate(); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    // Process with proper error context
    if err := s.processVideo(); err != nil {
        return fmt.Errorf("failed to process video: %w", err)
    }
    
    return nil
}
```

### 4. Progress and Output
```go
func (s *SplitVideoCmd) Run() error {
    if s.Verbose {
        fmt.Printf("Processing video: %s\n", s.InputFile)
        fmt.Printf("Time range: %s to %s\n", s.StartTime, s.EndTime)
    }
    
    // Show progress
    progress := progressbar.Default(100)
    defer progress.Finish()
    
    // Update progress during processing
    progress.Add(10)
    
    return nil
}
```

## Testing Patterns

### Unit Testing Commands
```go
func TestSplitVideoCmd_Validate(t *testing.T) {
    tests := []struct {
        name    string
        cmd     SplitVideoCmd
        wantErr bool
    }{
        {
            name: "valid times",
            cmd: SplitVideoCmd{
                StartTime: "00:00:00",
                EndTime:   "00:01:00",
            },
            wantErr: false,
        },
        {
            name: "invalid times",
            cmd: SplitVideoCmd{
                StartTime: "00:01:00",
                EndTime:   "00:00:00",
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.cmd.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("SplitVideoCmd.Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Integration Testing
```go
func TestSplitVideoCmd_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    // Create temporary test file
    tmpDir := t.TempDir()
    testFile := filepath.Join(tmpDir, "test.mp4")
    
    // Run command
    cmd := SplitVideoCmd{
        InputFile:    testFile,
        StartTime:    "00:00:00",
        EndTime:      "00:00:30",
        ChunkDuration: 10,
        OutputDir:    tmpDir,
    }
    
    err := cmd.Run()
    if err != nil {
        t.Errorf("SplitVideoCmd.Run() error = %v", err)
    }
    
    // Verify output files exist
    // ...
}
```

## Best Practices

### 1. Command Organization
- Group related commands in separate files
- Use consistent naming patterns
- Keep commands focused on single responsibilities

### 2. User Experience
- Provide clear, actionable error messages
- Include examples in help text
- Support both short and long flag names where appropriate
- Use consistent output formatting

### 3. Configuration
- Support environment variables for sensitive data
- Use configuration files for complex settings
- Provide sensible defaults
- Allow override via command-line flags

### 4. Logging and Output
- Use structured logging for debugging
- Provide progress indicators for long operations
- Support quiet mode for scripting
- Use consistent output formats (JSON, CSV, etc.)

## Example Help Output
```
Usage: roadtrip <command> [flags]

A tool for building roadtrip playlists from vlog videos.

Commands:
  split-video     Split video into chunks for processing
  upload-chunks   Upload video chunks to cloud storage
  build-playlist  Build playlist from video analysis
  build-playlist-csv Convert playlist to CSV format

Flags:
  --version    Show version information
  --help, -h   Show this help message

Examples:
  roadtrip split-video input.mp4 --start 00:00:00 --end 00:05:00
  roadtrip upload-chunks out/*.mp4 --bucket gs://my-bucket
  roadtrip build-playlist gs://my-bucket/chunk-*.mp4
```

## Dependencies

### Required Dependencies
```go
go get github.com/alecthomas/kong
go get github.com/schollz/progressbar/v3  // For progress bars
go get github.com/sirupsen/logrus         // For structured logging
```

### Optional Dependencies
```go
go get github.com/spf13/viper             // For configuration management
go get github.com/urfave/cli/v2           // Alternative CLI framework
```

## Migration from Other CLI Libraries

### From Cobra
- Replace `cobra.Command` with Kong struct tags
- Convert `RunE` functions to `Run` methods
- Replace `PersistentFlags` with embedded structs
- Update help text to use Kong's format

### From urfave/cli
- Replace `cli.App` with Kong struct
- Convert `Action` functions to `Run` methods
- Replace `cli.StringFlag` with struct tags
- Update command registration to use struct embedding