# Bazel Build System Documentation

## Overview
This repository uses Bazel as the primary build system for all binaries, libraries, and scripts. Bazel provides reproducible builds, dependency management, and efficient caching.

## Key Principles

### 1. All Binaries Need BUILD Rules
- Every Go binary must have a corresponding `BUILD` file
- Libraries should be organized in `BUILD` files with appropriate visibility
- Scripts and tools should have `BUILD` rules for consistency

### 2. Dependency Management
- Use `gazelle` to automatically generate and update `BUILD` rules
- Run `bazel run //:gazelle -- update` after adding new dependencies
- Keep `go.mod` and `go.sum` files in sync with Bazel dependencies

### 3. Testing Requirements
- Use `build_test` rules to verify binaries build correctly
- Write unit tests for critical business logic and error handling
- Use interfaces and mocks for external dependencies (network, file system, etc.)
- Consider test coverage vs. complexity - prioritize clean, testable interfaces over 100% coverage
- Run tests after each significant change: `bazel test //...`

## Common Commands

### Building Binaries
```bash
# Build a specific binary
bazel build //scripts/roadtrip:roadtrip

# Build all binaries
bazel build //...

# Build with verbose output
bazel build --verbose_failures //scripts/roadtrip:roadtrip
```

### Running Binaries
```bash
# Run a binary
bazel run //scripts/roadtrip:roadtrip

# Run with arguments
bazel run //scripts/roadtrip:roadtrip -- --help
```

### Testing
```bash
# Run all tests
bazel test //...

# Run specific tests
bazel test //scripts/roadtrip:roadtrip_test

# Run tests with coverage
bazel coverage //scripts/roadtrip/...
```

### Updating Dependencies
```bash
# Update BUILD files after adding dependencies
bazel run //:gazelle -- update

# Update external dependencies
bazel run //:gazelle -- update-repos -from_file=go.mod
```

## BUILD File Structure

### Example Binary BUILD Rule
```python
load("@rules_go//go:def.bzl", "go_binary", "go_library")

go_binary(
    name = "roadtrip",
    srcs = ["main.go"],
    deps = [
        "//internal/commands",
        "//internal/video",
        "//internal/storage",
        "//internal/ai",
        "//internal/export",
    ],
    visibility = ["//visibility:public"],
)

# Build test to verify the binary builds
build_test(
    name = "roadtrip_build_test",
    targets = [":roadtrip"],
)
```

### Example Library BUILD Rule
```python
load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "video",
    srcs = ["video.go"],
    importpath = "backend/scripts/roadtrip/internal/video",
    deps = [
        "//internal/storage",
    ],
    visibility = ["//internal:__subpackages__"],
)

go_test(
    name = "video_test",
    srcs = ["video_test.go"],
    embed = [":video"],
    deps = [
        "//internal/storage",
    ],
)
```

## Best Practices

### 1. File Organization
- Keep `BUILD` files close to the code they build
- Use consistent naming conventions
- Group related targets in the same `BUILD` file

### 2. Dependencies
- Minimize dependency depth
- Use `visibility` to control access
- Prefer internal packages over external dependencies

### 3. Testing
- Write tests for critical business logic and error handling
- Use build tests to verify compilation
- Test both success and failure cases
- Use interfaces and mocks for external dependencies
- Prioritize clean, testable interfaces over 100% coverage

### 4. Performance
- Use `deps` instead of `data` for code dependencies
- Avoid unnecessary dependencies
- Use `select` for platform-specific code

## Troubleshooting

### Common Issues

1. **Missing BUILD file**: Run `bazel run //:gazelle -- update`
2. **Dependency not found**: Check `go.mod` and run gazelle update
3. **Build failures**: Use `--verbose_failures` for detailed error messages
4. **Test failures**: Check test dependencies and visibility settings

### Debugging Commands
```bash
# Show dependency graph
bazel query --noimplicit_deps "deps(//scripts/roadtrip:roadtrip)" --output=graph

# Show build configuration
bazel config

# Clean and rebuild
bazel clean --expunge
bazel build //...
```

## Integration with CI/CD

### Pre-commit Checks
```bash
# Format BUILD files
bazel run //:buildifier

# Update dependencies
bazel run //:gazelle -- update

# Run tests
bazel test //...

# Build all targets
bazel build //...
```

### CI Pipeline
1. Install Bazel and dependencies
2. Run `bazel test //...`
3. Run `bazel build //...`
4. Generate and upload test coverage reports

## Resources

- [Bazel Go Rules Documentation](https://github.com/bazelbuild/rules_go)
- [Gazelle Documentation](https://github.com/bazelbuild/bazel-gazelle)
- [Bazel Best Practices](https://bazel.build/configure/best-practices)