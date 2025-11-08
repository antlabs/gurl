# Contributing to gurl

Thank you for your interest in contributing to gurl! This document provides guidelines and instructions for contributing.

## Getting Started

### Prerequisites

- Go 1.24 or later
- Git
- Make

### Setting Up Development Environment

1. **Fork and clone the repository:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/gurl.git
   cd gurl
   ```

2. **Install dependencies:**
   ```bash
   make deps
   ```

3. **Set up development tools:**
   ```bash
   make dev-setup
   ```

4. **Build the project:**
   ```bash
   make build
   ```

5. **Run tests:**
   ```bash
   make test
   ```

## Development Workflow

### Before Making Changes

1. Create a new branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make sure you're up to date:
   ```bash
   git pull origin main
   ```

### Making Changes

1. **Write code** following Go best practices
2. **Add tests** for new functionality
3. **Update documentation** if needed
4. **Format code:**
   ```bash
   make fmt
   ```
5. **Run linter:**
   ```bash
   make lint
   ```
6. **Run tests:**
   ```bash
   make test
   ```

### Committing Changes

1. **Write clear commit messages:**
   ```
   feat: add new feature
   fix: resolve bug in X
   docs: update README
   test: add tests for Y
   refactor: improve Z implementation
   ```

2. **Commit your changes:**
   ```bash
   git add .
   git commit -m "feat: your feature description"
   ```

### Submitting Pull Request

1. **Push to your fork:**
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create Pull Request** on GitHub

3. **Ensure CI passes:**
   - All tests pass
   - Linter checks pass
   - Code coverage is maintained

4. **Address review comments** if any

## Code Style

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Follow golangci-lint recommendations
- Write clear, self-documenting code
- Add comments for exported functions and types

### Example:

```go
// ProcessRequest handles HTTP request processing with rate limiting.
// It returns the response or an error if the request fails.
func ProcessRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
    // Implementation
}
```

## Testing

### Writing Tests

- Write unit tests for all new functions
- Use table-driven tests when appropriate
- Mock external dependencies
- Aim for high test coverage

### Example Test:

```go
func TestProcessRequest(t *testing.T) {
    tests := []struct {
        name    string
        input   *http.Request
        want    *http.Response
        wantErr bool
    }{
        {
            name:    "successful request",
            input:   createTestRequest(),
            want:    &http.Response{StatusCode: 200},
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ProcessRequest(context.Background(), tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ProcessRequest() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            // Add assertions
        })
    }
}
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test -v -run TestProcessRequest ./...

# Run benchmarks
make bench
```

## Documentation

### Code Documentation

- Document all exported functions, types, and constants
- Use godoc format
- Include examples when helpful

### README Updates

- Update README.md for new features
- Add examples for new functionality
- Update feature list if applicable

## Pull Request Guidelines

### PR Title

Use conventional commit format:
- `feat: add new feature`
- `fix: resolve bug`
- `docs: update documentation`
- `test: add tests`
- `refactor: improve code`
- `perf: performance improvement`
- `chore: maintenance tasks`

### PR Description

Include:
1. **What** - What changes were made
2. **Why** - Why these changes are needed
3. **How** - How the changes work
4. **Testing** - How the changes were tested
5. **Screenshots** - If UI changes (for terminal UI)

### Example PR Description:

```markdown
## Description
Adds rate limiting compensation to account for request processing time.

## Motivation
Current implementation doesn't account for request duration, leading to lower actual rates than configured.

## Changes
- Modified runConnection() to track request duration
- Added compensation logic to adjust sleep time
- Updated tests to verify compensation

## Testing
- Added unit tests for compensation logic
- Verified with benchmark tests
- Tested with various rate limits and request durations

## Related Issues
Fixes #123
```

## Issue Guidelines

### Reporting Bugs

Include:
1. **Description** - Clear description of the bug
2. **Steps to Reproduce** - How to reproduce the issue
3. **Expected Behavior** - What should happen
4. **Actual Behavior** - What actually happens
5. **Environment** - OS, Go version, gurl version
6. **Logs** - Relevant error messages or logs

### Feature Requests

Include:
1. **Description** - What feature you'd like
2. **Motivation** - Why this feature is needed
3. **Use Case** - How you would use it
4. **Alternatives** - Other solutions you've considered

## Code Review Process

1. **Automated checks** must pass:
   - CI tests
   - Linter
   - Code coverage

2. **Manual review** by maintainers:
   - Code quality
   - Design decisions
   - Documentation

3. **Address feedback** promptly

4. **Approval** required before merge

## Release Process

Releases are automated via GitHub Actions:

1. Maintainer creates version tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. CI automatically:
   - Builds binaries for all platforms
   - Creates GitHub release
   - Publishes Docker image

## Getting Help

- **Documentation**: Check [docs/CI_CD.md](../docs/CI_CD.md)
- **Issues**: Search existing issues
- **Discussions**: Start a discussion for questions
- **Contact**: Open an issue for help

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on the code, not the person
- Help others learn and grow

## License

By contributing, you agree that your contributions will be licensed under the project's license.

## Recognition

Contributors are recognized in:
- GitHub contributors page
- Release notes
- Project documentation

Thank you for contributing to gurl! 🚀
