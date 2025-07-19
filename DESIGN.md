# Roadtrip Playlist Tool Design Document

## Architecture Overview

The roadtrip playlist tool is a command-line application built in Go that processes vlog videos to extract music information and generate playlists. The tool follows a pipeline architecture:

```
Video Input → Split into Chunks → Upload to Cloud → AI Analysis → CSV Export
```

## Core Components

### 1. CLI Framework (Kong)
- **Purpose**: Command-line interface with subcommands
- **Structure**: 
  - `roadtrip` (main binary)
    - `split-video` (video processing)
    - `upload-chunks` (cloud storage)
    - `build-playlist` (AI analysis)
    - `build-playlist-csv` (data export)

### 2. Video Processing Module
- **Technology**: ffmpeg via os/exec
- **Functionality**: 
  - Split videos into time-based chunks
  - Use copy codec for fast processing
  - Progress tracking and output visibility
- **Input**: Video file, time range, chunk duration
- **Output**: Multiple video chunks in specified directory

### 3. Cloud Storage Module
- **Technology**: Google Cloud Storage client library
- **Functionality**:
  - Upload video chunks to GCS
  - Support for glob patterns
  - Progress tracking
- **Testing**: MinIO compatibility for local testing

### 4. AI Analysis Module
- **Technology**: Gemini 2.5 Flash API
- **Functionality**:
  - Video content analysis
  - Music detection and identification
  - Web search integration for song details
  - Environment variable configuration (GEMINI_API_KEY)
- **Output**: Verbatim Gemini response with optional JSON validation

### 5. Data Export Module
- **Technology**: Go CSV library
- **Functionality**:
  - Relaxed JSON parsing for non-standard Gemini outputs
  - JSON to CSV conversion with error handling
  - Standardized playlist format
  - Unit testable with mock data including malformed JSON

## Data Structures

### Video Analysis Response
```json
{
  "description": "string",
  "has_music": "boolean",
  "transcript": "string",
  "song_title": "string",
  "song_artist": "string",
  "web_search_title": "string",
  "web_search_artist": "string",
  "youtube_url": "string",
  "spotify_url": "string",
  "video_path": "string"
}
```

### CSV Output Format
```csv
description,has_music,transcript,song_title,song_artist,web_search_title,web_search_artist,youtube_url,spotify_url,video_path
```

## File Structure
```
scripts/roadtrip/
├── main.go                 # Main CLI entry point
├── BUILD                   # Bazel build rules
├── env/                    # Environment configuration
│   ├── env.go             # Environment variable handling
│   └── config.go          # Configuration structures
├── commands/               # Command implementations
│   ├── split_video.go
│   ├── upload_chunks.go
│   ├── build_playlist.go
│   └── build_playlist_csv.go
├── internal/               # Internal packages
│   ├── video/             # Video processing
│   ├── storage/           # Cloud storage
│   ├── ai/                # AI analysis
│   └── export/            # Data export
├── prompts/               # AI prompt templates
│   └── video_analysis.tmpl
└── testdata/              # Test data and examples
    ├── sample_video.mp4
    ├── sample_response.json
    ├── malformed_response.json
    └── expected_output.csv
```

## Dependencies

### Go Dependencies
- `github.com/alecthomas/kong` - CLI framework
- `cloud.google.com/go/storage` - GCS client
- `encoding/csv` - CSV processing
- `encoding/json` - JSON processing
- `os/exec` - ffmpeg execution
- `path/filepath` - File operations
- `github.com/vbauerster/mpb` - Progress bars
- `google.golang.org/api/generativeai/v1` - Gemini API client

### External Dependencies
- ffmpeg (via devcontainer feature)
- Google Cloud SDK
- Gemini 2.5 Flash API access
- MinIO (for local testing)

## Error Handling Strategy

1. **Graceful Degradation**: Continue processing other files if one fails
2. **Detailed Logging**: Log all operations with appropriate levels
3. **Validation**: Validate inputs before processing
4. **Recovery**: Implement retry logic for network operations

## Testing Strategy

1. **Unit Tests**: Test individual functions with mock data
2. **Integration Tests**: Test command workflows with test files
3. **MinIO Testing**: Use MinIO for cloud storage testing
4. **Mock AI Responses**: Use predefined responses for AI testing

## Security Considerations

1. **API Keys**: Use environment variables for sensitive data
2. **File Permissions**: Ensure proper file access controls
3. **Input Validation**: Validate all user inputs
4. **Error Messages**: Avoid exposing sensitive information in errors

## Performance Considerations

1. **Parallel Processing**: Process multiple chunks concurrently where possible
2. **Memory Management**: Stream large files instead of loading entirely
3. **Progress Tracking**: Provide real-time feedback for long operations
4. **Caching**: Cache AI responses to avoid duplicate API calls

## Configuration

The tool will support configuration via:
1. Environment variables (GEMINI_API_KEY, GCP project settings)
2. Command-line flags (highest priority)
3. Default GCS bucket: `cbhl-roadtrip-202507` in project `gen-lang-client-0629405113` in region `northamerica-northeast2`

## JSON Handling Strategy

### build-playlist Command
- Outputs verbatim Gemini response to stdout
- Optional `--validate-json` flag to check if response is valid JSON
- Continues processing even if JSON is malformed
- Provides clear error messages for API failures

### build-playlist-csv Command
- Implements relaxed JSON parsing to handle non-standard Gemini outputs
- Attempts to extract JSON from prose text using regex patterns
- Falls back to manual parsing if standard JSON parsing fails
- Provides warnings for malformed data but continues processing
- Supports filtering and data cleaning options

### Alternative Design: Filter Command
Consider adding a separate `filter` command that:
- Takes raw Gemini output and extracts/validates JSON
- Provides data cleaning and transformation
- Outputs standardized JSON for downstream processing
- Can be chained: `build-playlist | filter | build-playlist-csv`

## Future Enhancements

1. **Batch Processing**: Process multiple videos at once
2. **Playlist Formats**: Support multiple playlist formats (M3U, etc.)
3. **Music Recognition**: Integrate with music recognition APIs
4. **GUI Interface**: Optional web-based interface
5. **Cloud Functions**: Serverless processing for large videos
6. **Data Filtering**: Advanced filtering and transformation commands