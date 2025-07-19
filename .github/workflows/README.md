# GitHub Actions Workflows

This directory contains GitHub Actions workflows for continuous integration.

## Workflows

### bazel-test.yml
**Full repository testing** - Runs on all pushes and pull requests.

**Triggers:**
- Push to main/master
- Pull request to main/master

**Steps:**
1. Setup Go 1.21
2. Install Bazelisk
3. Install dependencies (ffmpeg, MinIO, Google Cloud SDK)
4. Cache Bazel build artifacts
5. Verify Bazel setup and list available targets
6. Run all tests: `bazelisk test //...`
7. Build all targets: `bazelisk build //...`
8. Show test results
9. Upload test results as artifacts

**Use case:** Full CI pipeline for the entire repository.

### roadtrip-test.yml
**Focused roadtrip tool testing** - Runs only when roadtrip tool files change.

**Triggers:**
- Push to main/master (only if roadtrip files changed)
- Pull request to main/master (only if roadtrip files changed)

**Steps:**
1. Setup Go 1.21
2. Install Bazelisk
3. Cache Bazel build artifacts
4. Verify Bazel setup and list available targets
5. Run roadtrip tests: `bazelisk test //scripts/roadtrip/...`
6. Build roadtrip binary: `bazelisk build //scripts/roadtrip:roadtrip`
7. Test CLI functionality

**Use case:** Fast feedback for roadtrip tool development.

## Local Testing

To test the workflows locally:

```bash
# Test roadtrip tool specifically
bazelisk test //scripts/roadtrip/...

# Test everything
bazelisk test //...

# Build everything
bazelisk build //...
```

## Cache Strategy

Both workflows use Bazel caching to speed up builds:
- Cache key: `{os}-bazel-{hash-of-MODULE.bazel-files}`
- Cache paths: `~/.cache/bazel`, `~/.cache/bazelisk`
- Fallback: `{os}-bazel-` (partial cache hits)

## Dependencies

The full workflow installs:
- **ffmpeg**: For video processing
- **MinIO**: For local object storage testing
- **Google Cloud SDK**: For GCS operations

The focused workflow only installs:
- **Go**: For building the tool
- **Bazelisk**: For build system management