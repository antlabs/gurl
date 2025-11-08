# CI/CD Documentation

This document describes the CI/CD setup for the gurl project.

## Overview

The project uses GitHub Actions for continuous integration and deployment. The workflows are organized into several files:

- **ci.yml** - Main CI pipeline for testing and building
- **pr.yml** - Pull request validation
- **release.yml** - Release automation
- **codeql.yml** - Security analysis
- **benchmark.yml** - Performance benchmarking

## Workflows

### CI Workflow (`ci.yml`)

Runs on every push to main/master/develop branches and on pull requests.

**Jobs:**
- **Test**: Runs tests on Ubuntu, macOS, and Windows with Go 1.24.x
  - Downloads dependencies
  - Verifies dependencies
  - Runs tests with race detector
  - Generates coverage report
  - Uploads coverage to Codecov (Ubuntu only)

- **Lint**: Runs golangci-lint on Ubuntu
  - Checks code quality and style
  - Enforces best practices

- **Build**: Builds binaries for all platforms
  - Builds for current platform
  - Cross-compiles for Linux, macOS, Windows (amd64 and arm64)
  - Uploads artifacts for 7 days

### Pull Request Workflow (`pr.yml`)

Runs on pull request events (opened, synchronized, reopened).

**Jobs:**
- **Validate**: Validates PR changes
  - Checks code formatting
  - Verifies go.mod and go.sum are tidy
  - Runs tests with race detector
  - Runs benchmarks

- **Size Check**: Monitors binary size
  - Builds binary
  - Reports size
  - Warns if larger than 50MB

- **Security**: Scans for security issues
  - Runs Gosec security scanner
  - Uploads results to GitHub Security

### Release Workflow (`release.yml`)

Triggers on version tags (v*).

**Jobs:**
- **Release**: Creates GitHub release
  - Builds binaries for all platforms
  - Generates checksums
  - Creates release archives (tar.gz for Unix, zip for Windows)
  - Creates GitHub release with artifacts
  - Auto-generates release notes

- **Docker**: Builds and pushes Docker image
  - Builds multi-arch image (amd64, arm64)
  - Pushes to Docker Hub
  - Tags with version and latest

### CodeQL Workflow (`codeql.yml`)

Runs on push, pull requests, and weekly schedule.

**Jobs:**
- **Analyze**: Performs security analysis
  - Initializes CodeQL
  - Builds project
  - Analyzes code for vulnerabilities

### Benchmark Workflow (`benchmark.yml`)

Runs on push to main/master and pull requests.

**Jobs:**
- **Benchmark**: Runs performance benchmarks
  - Executes benchmarks
  - Stores results
  - Alerts on performance regressions (>150%)

## Configuration Files

### `.golangci.yml`

Configures golangci-lint with enabled linters:
- errcheck, gosimple, govet, ineffassign, staticcheck, unused
- gofmt, goimports, misspell, unconvert, unparam
- gocritic, gosec, revive

### `.github/dependabot.yml`

Configures Dependabot for automatic dependency updates:
- Go modules: weekly updates
- GitHub Actions: weekly updates

### `Dockerfile`

Multi-stage Docker build:
- Build stage: Compiles binary with Go 1.24
- Final stage: Minimal Alpine image with ca-certificates
- Runs as non-root user

### `.dockerignore`

Excludes unnecessary files from Docker build context.

## Secrets Required

For full functionality, configure these secrets in GitHub repository settings:

### Required for Docker builds:
- `DOCKER_USERNAME` - Docker Hub username
- `DOCKER_PASSWORD` - Docker Hub password or access token

### Optional:
- `CODECOV_TOKEN` - Codecov token for coverage reports (optional, works without)

## Creating a Release

To create a new release:

1. **Tag the version:**
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. **Automated process:**
   - Release workflow triggers automatically
   - Builds binaries for all platforms
   - Creates checksums
   - Creates GitHub release with artifacts
   - Builds and pushes Docker image

3. **Manual verification:**
   - Check GitHub releases page
   - Verify all artifacts are present
   - Test Docker image: `docker pull <username>/gurl:v1.0.0`

## Local Testing

### Test CI locally with act:

```bash
# Install act
brew install act  # macOS
# or
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash

# Run CI workflow
act -j test

# Run specific job
act -j lint
```

### Build Docker image locally:

```bash
docker build -t gurl:local .
docker run --rm gurl:local --version
```

### Run linter locally:

```bash
make lint
# or
golangci-lint run
```

## Monitoring

### Build Status

Check workflow status:
- GitHub Actions tab in repository
- Badge in README (add badges for visibility)

### Coverage

View coverage reports:
- Codecov dashboard (if configured)
- Coverage artifacts in CI runs

### Security

Review security alerts:
- GitHub Security tab
- CodeQL analysis results
- Dependabot alerts

## Troubleshooting

### Build Failures

1. **Go version mismatch:**
   - Update go-version in workflows
   - Update go.mod

2. **Dependency issues:**
   - Run `go mod tidy` locally
   - Commit go.mod and go.sum changes

3. **Test failures:**
   - Check test logs in Actions
   - Run tests locally: `make test`

### Release Issues

1. **Missing artifacts:**
   - Check build-all target in Makefile
   - Verify GOOS/GOARCH combinations

2. **Docker build fails:**
   - Verify Dockerfile syntax
   - Check Docker Hub credentials

### Lint Failures

1. **Fix formatting:**
   ```bash
   make fmt
   ```

2. **Fix specific issues:**
   ```bash
   golangci-lint run --fix
   ```

## Best Practices

1. **Before pushing:**
   - Run tests locally: `make test`
   - Run linter: `make lint`
   - Format code: `make fmt`

2. **Pull requests:**
   - Keep changes focused
   - Ensure all checks pass
   - Update tests for new features

3. **Releases:**
   - Use semantic versioning
   - Update CHANGELOG
   - Test release candidates

4. **Security:**
   - Review Dependabot PRs promptly
   - Address security alerts quickly
   - Keep dependencies updated

## Future Enhancements

Potential improvements:
- Add integration tests
- Add performance regression tests
- Add deployment to package managers (Homebrew, apt, etc.)
- Add automated changelog generation
- Add release candidate workflow
- Add nightly builds
