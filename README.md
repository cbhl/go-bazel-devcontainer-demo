# Roadtrip Playlist Tool

A command-line tool in Go to help create roadtrip playlists from vlog videos using AI analysis, ffmpeg, and cloud storage.

## Features

- **Video Processing**: Split large video files into manageable chunks using ffmpeg
- **Cloud Storage**: Upload video chunks to Google Cloud Storage or MinIO
- **AI Analysis**: Analyze video content using Gemini 2.5 Flash to identify music
- **Data Export**: Convert analysis results to CSV format with relaxed JSON parsing
- **Comprehensive Testing**: Unit tests with mock clients for development and testing

## Installation

### Prerequisites

- Go 1.23 or later
- ffmpeg
- Google Cloud SDK (for GCS uploads)
- Gemini 2.5 Flash API access (for AI analysis)

### Building

```bash
# Clone the repository
git clone <repository-url>
cd roadtrip-playlist-tool

# Build the tool
go build -o roadtrip scripts/roadtrip/

# Or use Bazel
bazel build //scripts/roadtrip:roadtrip
```

## Usage

The tool provides four main commands:

### 1. Split Video (`split-video`)

Split a video file into chunks for processing.

```bash
./roadtrip split-video \
  --in input_video.mp4 \
  --start 00:00:00 \
  --end 01:30:00 \
  --chunk-duration 30 \
  --out output_directory
```

**Flags:**
- `--in`: Input video file path (required)
- `--start`: Start timestamp in HH:MM:SS format
- `--end`: End timestamp in HH:MM:SS format
- `--chunk-duration`: Duration of each chunk in seconds (default: 30)
- `--out`: Output directory (default: "out")

### 2. Upload Chunks (`upload-chunks`)

Upload video chunks to cloud storage.

```bash
./roadtrip upload-chunks \
  --in "out/*.mp4" \
  --project-id my-gcp-project \
  --bucket gs://my-bucket/chunks
```

**Flags:**
- `--in`: Input folder, file, or glob pattern (required)
- `--project-id`: GCP project ID (required)
- `--zone`: GCP zone (optional)
- `--bucket`: GCS bucket path in format `gs://bucket-name/prefix` (required)

### 3. Build Playlist (`build-playlist`)

Analyze videos using AI to identify music and build playlists.

```bash
./roadtrip build-playlist \
  --in gs://bucket/video1.mp4 \
  --in gs://bucket/video2.mp4 \
  --validate
```

**Flags:**
- `--in`: GCS paths to analyze (can be specified multiple times)
- `--validate`: Validate and output JSON format

### 4. Build Playlist CSV (`build-playlist-csv`)

Convert analysis results to CSV format.

```bash
# From file
./roadtrip build-playlist-csv --in analysis_results.json

# From stdin
cat analysis_results.json | ./roadtrip build-playlist-csv --in -
```

**Flags:**
- `--in`: Input JSON file or "-" for stdin (required)

## Configuration

### Environment Variables

- `GEMINI_API_KEY`: API key for Gemini 2.5 Flash (for AI analysis)
- `GOOGLE_APPLICATION_CREDENTIALS`: Path to GCP service account key (for GCS uploads)

### Authentication

#### Google Cloud Storage

1. Install Google Cloud SDK
2. Authenticate: `gcloud auth application-default login`
3. Or set service account key: `export GOOGLE_APPLICATION_CREDENTIALS=/path/to/key.json`

#### Gemini API

1. Get API key from Google AI Studio
2. Set environment variable: `export GEMINI_API_KEY=your-api-key`

## Development

### Project Structure

```
scripts/roadtrip/
├── main.go              # CLI entry point
├── video/               # Video processing with ffmpeg
├── storage/             # Cloud storage (GCS/MinIO)
├── ai/                  # AI analysis (Gemini 2.5 Flash)
├── export/              # Data export (JSON to CSV)
├── prompts/             # AI prompt templates
└── testdata/            # Test data files
```

### Testing

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./scripts/roadtrip/video/...
go test ./scripts/roadtrip/storage/...
go test ./scripts/roadtrip/ai/...
go test ./scripts/roadtrip/export/...

# Run with verbose output
go test -v ./...
```

### Mock Clients

The tool includes mock clients for development and testing:

- `MockStorageClient`: Simulates cloud storage uploads
- `MockAIClient`: Simulates AI analysis responses
- No external dependencies required for testing

## Examples

### Complete Workflow

```bash
# 1. Split a 2-hour video into 30-second chunks
./roadtrip split-video \
  --in roadtrip_vlog.mp4 \
  --start 00:00:00 \
  --end 02:00:00 \
  --chunk-duration 30 \
  --out chunks/

# 2. Upload chunks to GCS
./roadtrip upload-chunks \
  --in "chunks/*.mp4" \
  --project-id my-project \
  --bucket gs://my-bucket/roadtrip-chunks/

# 3. Analyze videos for music
./roadtrip build-playlist \
  --in gs://my-bucket/roadtrip-chunks/chunk_001.mp4 \
  --in gs://my-bucket/roadtrip-chunks/chunk_002.mp4 \
  --validate > analysis.json

# 4. Convert to CSV
./roadtrip build-playlist-csv --in analysis.json > playlist.csv
```

### Advanced Usage

```bash
# Process multiple videos with glob patterns
./roadtrip upload-chunks \
  --in "videos/*.mp4" \
  --project-id my-project \
  --bucket gs://my-bucket/videos/

# Analyze with JSON validation
./roadtrip build-playlist \
  --in gs://bucket/video1.mp4 \
  --in gs://bucket/video2.mp4 \
  --validate | jq '.[] | select(.has_music)'

# Handle malformed JSON with relaxed parsing
echo 'Prose text with JSON: {"description": "Video", "has_music": true}' | \
  ./roadtrip build-playlist-csv --in -
```

## Troubleshooting

### Common Issues

1. **ffmpeg not found**: Install ffmpeg and ensure it's in your PATH
2. **GCS authentication failed**: Check your GCP credentials and permissions
3. **Gemini API errors**: Verify your API key and quota limits
4. **JSON parsing errors**: The tool includes relaxed parsing for malformed JSON

### Debug Mode

Enable debug logging by setting the log level:

```bash
export LOG_LEVEL=debug
./roadtrip <command> [flags]
```

### Error Handling

The tool includes comprehensive error handling:

- Graceful degradation when external services are unavailable
- Detailed error messages with context
- Fallback mechanisms for malformed data
- Progress tracking for long-running operations

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

[Add your license information here]

## Support

For issues and questions:

1. Check the troubleshooting section
2. Review the test examples
3. Open an issue with detailed error information
4. Include relevant logs and configuration details
