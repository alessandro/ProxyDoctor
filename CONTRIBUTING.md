# Contributing to ProxyDoctor

Thank you for your interest in contributing to ProxyDoctor! This document provides guidelines and information for contributors.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Code Style](#code-style)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Issue Guidelines](#issue-guidelines)
- [Writing New Checks](#writing-new-checks)
- [Community](#community)

## Getting Started

1. **Fork** the repository on GitHub
2. **Clone** your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/ProxyDoctor.git
   cd ProxyDoctor
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/francomano/ProxyDoctor.git
   ```
4. **Run setup**:
   ```bash
   ./setup.sh
   ```

## Development Setup

### Prerequisites

- **Go** >= 1.25
- **Git**

### Quick Start

```bash
# Install dependencies and run tests
./setup.sh

# Or manually:
go mod download
go test -v ./...
```

### Building

```bash
# Build CLI
go build -o bin/proxydoctor ./cmd/cli

# Build server
go build -o bin/proxydoctor-server ./cmd/server

# Or use convenience scripts
./run.sh cli --help
./run.sh server
```

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests for specific package
go test -v ./core/engine/...

# Run with coverage
go test -cover ./...
```

## Code Style

### Go Conventions

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` and `goimports` to format code
- Keep functions focused and small
- Write meaningful variable and function names
- Add comments for exported functions and complex logic

### Naming Conventions

```go
// Package names: lowercase, single word
package engine

// Exported functions: PascalCase
func (o *DiagnosisOrchestrator) RunDiagnosis() (*DiagnosisReport, error)

// Unexported functions: camelCase
func parseProxyURL(raw string) (*ProxyConfig, error)

// Interfaces: adjective or -er suffix
type Checker interface { ... }
type NetworkAdapter interface { ... }

// Structs: noun
type CheckResult struct { ... }
type DiagnosisReport struct { ... }
```

### Error Handling

```go
// Always handle errors explicitly
result, err := someOperation()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Use error wrapping for context
return fmt.Errorf("public_ip check failed: %w", err)
```

### File Organization

```
core/
├── check/          # Interface definitions
├── engine/         # Orchestration logic
├── adapters/       # Network implementations
├── checks/         # Individual check implementations
│   └── public_ip/
│       ├── check.go
│       └── resolver.go
└── utils/          # Shared helpers
```

## Making Changes

### Branch Naming

```bash
# Feature branches
git checkout -b feature/add-dns-leak-check

# Bug fixes
git checkout -b fix/timeout-parsing

# Documentation
git checkout -b docs/update-readme
```

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add DNS leak detection check

- Implement dns_leak checker with mock DNS server
- Add unit tests for leak detection logic
- Register check in CLI and server

Closes #5
```

```
fix: resolve timeout not being passed to orchestrator

The --timeout flag was parsed but never used.

Fixes #3
```

```
docs: add CONTRIBUTING.md guide

Add comprehensive guide for new contributors.
```

### Keep PRs Focused

- One feature or fix per PR
- Keep PRs small and reviewable
- Include tests for new functionality
- Update documentation as needed

## Testing

### Writing Tests

```go
package public_ip

import (
    "testing"
    "net/http/httptest"
)

func TestPublicIPCheck(t *testing.T) {
    // Arrange
    mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`{"ip": "203.0.113.42"}`))
    }))
    defer mockServer.Close()

    check := &PublicIPCheck{}
    ctx := ExecutionContext{
        URL: mockServer.URL,
    }

    // Act
    result := check.Execute(ctx)

    // Assert
    if result.Status != StatusPassed {
        t.Errorf("expected status passed, got %v", result.Status)
    }
}
```

### Test Requirements

- All new checks must include unit tests
- Mock external services (no real API calls in tests)
- Aim for >80% coverage on new code
- Tests should be deterministic (no flaky tests)

### Running Specific Tests

```bash
# Run all tests in a package
go test -v ./core/checks/public_ip/...

# Run a specific test
go test -v -run TestPublicIP ./core/checks/public_ip/...

# Run tests with race detector
go test -race ./...
```

## Pull Request Process

### Before Submitting

1. **Update your fork**:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run all tests**:
   ```bash
   go test -v ./...
   ```

3. **Build binaries**:
   ```bash
   go build -o bin/proxydoctor ./cmd/cli
   go build -o bin/proxydoctor-server ./cmd/server
   ```

4. **Update documentation** if needed

### PR Template

```markdown
## Description

Brief description of changes.

## Type of Change

- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing

- [ ] Unit tests added/updated
- [ ] All tests pass locally
- [ ] Manual testing performed

## Checklist

- [ ] Code follows project style
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] No new warnings
```

### Review Process

1. Maintainers will review your PR
2. Address any feedback or requested changes
3. Once approved, your PR will be merged

## Issue Guidelines

### Reporting Bugs

```markdown
## Bug Description

Clear description of the bug.

## Steps to Reproduce

1. Run command...
2. With flags...
3. See error...

## Expected Behavior

What should happen.

## Actual Behavior

What actually happens.

## Environment

- OS: [e.g., Ubuntu 22.04]
- Go version: [e.g., 1.25]
- ProxyDoctor version: [e.g., v0.2.0]
```

### Suggesting Features

```markdown
## Feature Description

Clear description of the feature.

## Use Case

Why this feature is needed.

## Proposed Implementation

How you think it could be implemented (optional).

## Acceptance Criteria

- [ ] Criterion 1
- [ ] Criterion 2
```

## Writing New Checks

### Check Interface

```go
type Checker interface {
    // Metadata
    ID() string
    Name() string
    Description() string
    Category() CheckCategory

    // Dependencies
    DependsOn() []string

    // Execution
    Execute(ctx ExecutionContext) CheckResult
}
```

### Example Check Structure

```
core/checks/your_check/
├── check.go          # Main check implementation
├── helper.go         # Optional helper functions
└── check_test.go     # Unit tests
```

### Check Registration

Register your check in:
1. `cmd/cli/commands/diagnose.go`
2. `cmd/cli/commands/list.go`
3. `cmd/server/main.go`

See existing checks in `core/checks/public_ip/` for reference.

## Community

- **Issues**: Use GitHub Issues for bugs and feature requests
- **Discussions**: Use GitHub Discussions for questions and ideas
- **PRs**: Welcome! See guidelines above

## Code of Conduct

Please read our [Code of Conduct](CODE_OF_CONDUCT.md) before contributing.

## License

By contributing, you agree that your contributions will be licensed under the GPL-3.0 License.

---

Thank you for contributing to ProxyDoctor! 🎉
