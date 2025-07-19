# Go Bazel DevContainer Demo

Trying out GitPod Flex + Devcontainers + Go + Bazel.

Also installs the gcloud cli in the Dockerfile.

Last updated: Jul 17, 2025.

## Roadtrip Playlist Tool

This repository includes a comprehensive **Roadtrip Playlist Tool** built with Go and Bazel. The tool helps create playlists from vlog videos using AI analysis, video processing, and cloud storage integration.

### 🎯 Features

- **Video Processing**: Split large video files into manageable chunks using ffmpeg
- **Cloud Storage**: Upload video chunks to Google Cloud Storage
- **AI Analysis**: Analyze video content using Gemini 2.5 Flash to identify music
- **Data Export**: Convert analysis results to CSV format with relaxed JSON parsing
- **Comprehensive Testing**: Unit tests with mock clients for development and testing

### 📚 Documentation

The playlist tool documentation is located in `scripts/roadtrip/`:

- **[README.md](scripts/roadtrip/README.md)** - Main tool documentation with usage examples
- **[PLAN.md](scripts/roadtrip/PLAN.md)** - Development plan and project phases
- **[DESIGN.md](scripts/roadtrip/DESIGN.md)** - System design and architecture
- **[CLI_STANDARDS.md](scripts/roadtrip/CLI_STANDARDS.md)** - CLI development standards
- **[BAZEL.md](scripts/roadtrip/BAZEL.md)** - Bazel build system documentation

### 🚀 Quick Start

```bash
# Build the tool
bazel build //scripts/roadtrip:roadtrip

# Run with help
bazel run //scripts/roadtrip:roadtrip -- --help

# Split a video
bazel run //scripts/roadtrip:roadtrip -- split-video --in video.mp4 --start 00:00:00 --end 01:30:00

# Analyze videos for music (requires GEMINI_API_KEY)
bazel run //scripts/roadtrip:roadtrip -- build-playlist --in gs://bucket/video1.mp4 --validate
```

### 🔧 Environment Setup

Go here and follow the instructions:

- https://app.gitpod.io/#https://github.com/cbhl/go-bazel-devcontainer-demo

If the devcontainer or Dockerfile fails to build, you get both pieces!

## Running

```
bazel run //backend:backend_server
```
