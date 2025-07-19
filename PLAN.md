# Roadtrip Playlist Tool Development Plan

## Overview
Build a command-line tool in Go to help create roadtrip playlists from vlog videos using Gemini 2.5 Flash, ffmpeg, and the Kong library.

## Development Steps

### Phase 1: Environment Setup
- [ ] Add ffmpeg devcontainer feature: `ghcr.io/devcontainers-extra/features/ffmpeg-apt-get:1`
- [ ] Create Bazel documentation for the repository
- [ ] Document Go/Kong CLI standards

### Phase 2: Basic CLI Foundation
- [ ] Create `scripts/roadtrip/` directory structure
- [ ] Implement basic Kong CLI program in `scripts/roadtrip/main.go` with "hello world"
- [ ] Add BUILD rule for the binary
- [ ] Add build_test rule for verification
- [ ] Verify binary builds and runs correctly

**Pause for human review here**

### Phase 3: Video Processing
- [ ] Add `split-video` command with flags:
  - Input video filename
  - Start timestamp
  - End timestamp  
  - Chunk duration
  - Output folder (default: `out/`)
- [ ] Implement ffmpeg integration with copy codec
- [ ] Add progress bar and command output visibility
- [ ] Test with sample video files

### Phase 4: Cloud Storage Integration
- [ ] Add `upload-chunks` command with flags:
  - Input folder/glob pattern
  - GCP project-id
  - GCP zone
  - GCS bucket path
- [ ] Implement GCS upload functionality
- [ ] Add unit tests with MinIO compatibility
- [ ] Test upload functionality

### Phase 5: AI Analysis
- [ ] Create `scripts/roadtrip/prompts/` directory
- [ ] Design Gemini 2.5 Flash prompt template for video analysis
- [ ] Add `build-playlist` command for GCS path processing
- [ ] Implement JSON response parsing

### Phase 6: Data Export
- [ ] Add `build-playlist-csv` command
- [ ] Implement JSON to CSV conversion
- [ ] Add comprehensive unit tests with example data
- [ ] Test CSV output format

### Phase 7: Documentation
- [ ] Create comprehensive README
- [ ] Add usage examples
- [ ] Document all commands and flags
- [ ] Add troubleshooting guide

## Success Criteria
- [ ] All commands work end-to-end
- [ ] Unit tests pass
- [ ] Build system works correctly
- [ ] Documentation is complete
- [ ] Tool can process real vlog videos and generate playlists

## Dependencies
- Go 1.x
- Bazel build system
- ffmpeg
- Kong CLI library
- Google Cloud SDK
- Gemini 2.5 Flash API access