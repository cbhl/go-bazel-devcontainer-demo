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

## Environment Setup

### Installing Bazelisk and Bazel
```bash
# Download bazelisk
curl -L https://github.com/bazelbuild/bazelisk/releases/latest/download/bazelisk-linux-amd64 -o ~/bazelisk

# Make executable
chmod +x ~/bazelisk

# Add to PATH (add to ~/.bashrc for persistence)
export PATH=$PATH:~/

# Verify installation
bazelisk version
```

### Initial Configuration
1. **Create root BUILD file** (`BUILD`):
   ```python
   load("@gazelle//:def.bzl", "gazelle")
   
   # gazelle:prefix backend
   gazelle(name = "gazelle")
   ```

2. **Configure MODULE.bazel**:
   ```python
   module(
       name = "backend",
       version = "0.1.0",
   )
   
   bazel_dep(name = "rules_go", version = "0.47.0")
   bazel_dep(name = "gazelle", version = "0.36.0")
   
   go_deps = use_extension("@gazelle//:extensions.bzl", "go_deps")
   go_deps.from_file(go_mod = "//scripts/roadtrip:go.mod")
   use_repo(go_deps, "com_github_alecthomas_kong")
   
   go_register_toolchains_ext = use_extension("@rules_go//go:extensions.bzl", "go_register_toolchains")
   go_rules_dependencies_ext = use_extension("@rules_go//go:extensions.bzl", "go_rules_dependencies")
   gazelle_ext = use_extension("@gazelle//:extensions.bzl", "gazelle")
   
   use_repo(go_register_toolchains_ext, "go_register_toolchains")
   use_repo(go_rules_dependencies_ext, "go_rules_dependencies")
   use_repo(gazelle_ext, "gazelle_extension")
   ```

3. **Create go.mod file** in your Go package directory
4. **Run gazelle to generate BUILD files**:
   ```bash
   bazelisk run //:gazelle -- update
   ```

## Common Commands

### Building Binaries
```bash
# Build a specific binary
bazelisk build //scripts/roadtrip:roadtrip

# Build all binaries
bazelisk build //...

# Build with verbose output
bazelisk build --verbose_failures //scripts/roadtrip:roadtrip
```

### Running Binaries
```bash
# Run a binary
bazelisk run //scripts/roadtrip:roadtrip

# Run with arguments
bazelisk run //scripts/roadtrip:roadtrip -- --help

# Run the built binary directly
bazel-bin/scripts/roadtrip/roadtrip_/roadtrip --help
```

### Testing
```bash
# Run all tests
bazelisk test //...

# Run specific tests
bazelisk test //scripts/roadtrip:roadtrip_test

# Run tests with coverage
bazelisk coverage //scripts/roadtrip/...

# Run tests for a specific package and its subpackages
bazelisk test //scripts/roadtrip/...
```

### Updating Dependencies
```bash
# Update BUILD files after adding dependencies
bazelisk run //:gazelle -- update

# Apply fixes to BUILD files
bazelisk run //:gazelle -- fix

# Update external dependencies
bazelisk run //:gazelle -- update-repos -from_file=go.mod
```

## BUILD File Structure

### Example Binary BUILD Rule (Gazelle-Generated)
```python
load("@rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "roadtrip_lib",
    srcs = ["main.go"],
    importpath = "backend/scripts/roadtrip",
    visibility = ["//scripts/roadtrip:__subpackages__"],
    deps = ["@com_github_alecthomas_kong//:kong"],
)

go_binary(
    name = "roadtrip",
    embed = [":roadtrip_lib"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "roadtrip_test",
    srcs = ["main_test.go"],
    data = glob(["testdata/**"]),
    embed = [":roadtrip_lib"],
)
```

### Example Library BUILD Rule (Gazelle-Generated)
```python
load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "env",
    srcs = [
        "config.go",
        "env.go",
    ],
    importpath = "backend/scripts/roadtrip/env",
    visibility = ["//scripts/roadtrip:__subpackages__"],
)

go_test(
    name = "env_test",
    srcs = ["env_test.go"],
    embed = [":env"],
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

1. **Missing BUILD file**: Run `bazelisk run //:gazelle -- update`
2. **Dependency not found**: Check `go.mod` and run gazelle update
3. **Build failures**: Use `--verbose_failures` for detailed error messages
4. **Test failures**: Check test dependencies and visibility settings
5. **Glob pattern errors**: Ensure directories exist or use `allow_empty = True`
6. **Gazelle not found**: Check MODULE.bazel configuration and gazelle extension setup

### Debugging Commands
```bash
# Show dependency graph
bazelisk query --noimplicit_deps "deps(//scripts/roadtrip:roadtrip)" --output=graph

# Show build configuration
bazelisk config

# Clean and rebuild
bazelisk clean --expunge
bazelisk build //...

# Check what targets are available
bazelisk query //scripts/roadtrip/...

# Show detailed build information
bazelisk build --verbose_failures //scripts/roadtrip:roadtrip
```

## Integration with CI/CD

### Pre-commit Checks
```bash
# Format BUILD files
bazelisk run //:buildifier

# Update dependencies
bazelisk run //:gazelle -- update

# Apply fixes
bazelisk run //:gazelle -- fix

# Run tests
bazelisk test //...

# Build all targets
bazelisk build //...
```

### CI Pipeline
1. Install Bazelisk and dependencies
2. Run `bazelisk test //...`
3. Run `bazelisk build //...`
4. Generate and upload test coverage reports

## Roadtrip Tool Specific Setup

### Project Structure
```
scripts/roadtrip/
├── main.go                 # Main CLI entry point
├── main_test.go           # CLI tests
├── go.mod                 # Go dependencies
├── BUILD                  # Gazelle-generated build rules
├── env/                   # Environment configuration
│   ├── env.go
│   ├── config.go
│   ├── env_test.go
│   └── BUILD
└── testdata/              # Test data directory
    └── .gitkeep
```

### Key Configuration Files

**MODULE.bazel** (root):
```python
module(
    name = "backend",
    version = "0.1.0",
)

bazel_dep(name = "rules_go", version = "0.47.0")
bazel_dep(name = "gazelle", version = "0.36.0")

go_deps = use_extension("@gazelle//:extensions.bzl", "go_deps")
go_deps.from_file(go_mod = "//scripts/roadtrip:go.mod")
use_repo(go_deps, "com_github_alecthomas_kong")

go_register_toolchains_ext = use_extension("@rules_go//go:extensions.bzl", "go_register_toolchains")
go_rules_dependencies_ext = use_extension("@rules_go//go:extensions.bzl", "go_rules_dependencies")
gazelle_ext = use_extension("@gazelle//:extensions.bzl", "gazelle")

use_repo(go_register_toolchains_ext, "go_register_toolchains")
use_repo(go_rules_dependencies_ext, "go_rules_dependencies")
use_repo(gazelle_ext, "gazelle_extension")
```

**BUILD** (root):
```python
load("@gazelle//:def.bzl", "gazelle")

# gazelle:prefix backend
gazelle(name = "gazelle")
```

### Lessons Learned

1. **Use bazelisk**: Always use bazelisk instead of bazel directly for version management
2. **Gazelle is essential**: Let gazelle generate BUILD files rather than writing them manually
3. **Module system**: Use the new Bazel module system (MODULE.bazel) instead of WORKSPACE
4. **Test data directories**: Create testdata directories with .gitkeep files to satisfy glob patterns
5. **Dependency management**: Use go_deps extension for Go dependencies in MODULE.bazel
6. **Library separation**: Separate library and binary targets for better testability

### Common Workflow

```bash
# 1. Add new Go dependencies to go.mod
cd scripts/roadtrip
go get github.com/new/dependency

# 2. Update BUILD files
cd /workspace
bazelisk run //:gazelle -- update

# 3. Build and test
bazelisk build //scripts/roadtrip:roadtrip
bazelisk test //scripts/roadtrip/...

# 4. Run the binary
bazelisk run //scripts/roadtrip:roadtrip -- --help
```

## Resources

- [Bazelisk Documentation](https://github.com/bazelbuild/bazelisk)
- [Bazel Go Rules Documentation](https://github.com/bazelbuild/rules_go)
- [Gazelle Documentation](https://github.com/bazelbuild/bazel-gazelle)
- [Bazel Best Practices](https://bazel.build/configure/best-practices)
- [Bazel Module System](https://bazel.build/concepts/external/overview)