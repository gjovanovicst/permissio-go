# Contributing to Permissio.io Go SDK

Thank you for your interest in contributing to the Permissio.io Go SDK! This document provides guidelines and steps for contributing.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct.

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in [Issues](https://github.com/permissio/permissio-go/issues)
2. If not, create a new issue with:
   - A clear, descriptive title
   - Steps to reproduce the bug
   - Expected vs actual behavior
   - Go version and SDK version
   - Any relevant code snippets or error messages

### Suggesting Features

1. Check existing issues for similar suggestions
2. Create a new issue with the `enhancement` label
3. Describe the feature and its use case

### Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature-name`
3. Make your changes
4. Write or update tests as needed
5. Ensure all tests pass: `go test ./...`
6. Run linting: `golangci-lint run`
7. Commit with clear messages following [Conventional Commits](https://www.conventionalcommits.org/)
8. Push to your fork and create a Pull Request

## Development Setup

### Prerequisites

- Go 1.21 or later
- golangci-lint (for linting)

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Code Style

- Follow standard Go conventions and idioms
- Use `gofmt` or `goimports` for formatting
- Document exported functions and types
- Keep functions focused and small
- Write meaningful test cases

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `test:` - Test changes
- `refactor:` - Code refactoring
- `chore:` - Maintenance tasks

Example: `feat: add tenant filtering to role assignment`

## Release Process

Releases are managed by the maintainers. Version bumps follow [Semantic Versioning](https://semver.org/).

## Questions?

Feel free to open an issue for any questions or reach out to the maintainers.

Thank you for contributing! 🎉
