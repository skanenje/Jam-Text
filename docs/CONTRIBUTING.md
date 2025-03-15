# Contributing to Jam-Text

Thank you for your interest in contributing to Jam-Text! This document provides guidelines and instructions for contributing.

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/your-username/Jam-Text.git
   cd Jam-Text
   ```
3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/skanenje/Jam-Text.git
   ```
4. Create a new branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Setup

1. Install Go 1.24.1 or higher
2. Install required dependencies:
   ```bash
   # For PDF support
   sudo apt-get install poppler-utils

   # For DOCX support
   sudo apt-get install pandoc
   ```
3. Build the project:
   ```bash
   make
   ```

## Code Standards

### Go Code
- Follow the [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Ensure all tests pass: `go test ./...`
- Maintain test coverage above 80%
- Document all exported functions and types

### Commit Messages
- Use clear, descriptive commit messages
- Follow the format:
  ```
  type: brief description

  Detailed description of changes and reasoning.
  ```
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

## Testing

1. Add tests for new features
2. Run the test suite:
   ```bash
   # Run all tests
   go test ./...

   # Run with coverage
   go test -cover ./...

   # Generate coverage report
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out
   ```

## Pull Request Process

1. Update documentation for new features
2. Add or update tests as needed
3. Ensure all tests pass
4. Update the README.md if needed
5. Create a Pull Request with:
   - Clear description of changes
   - Link to related issues
   - Screenshots for UI changes
   - List of tested scenarios

## Documentation

- Update relevant documentation in `docs/`
- Follow existing documentation style
- Include code examples where appropriate
- Update package documentation for new features

## Community

- Be respectful and inclusive
- Help others in discussions
- Report bugs and issues
- Suggest improvements

## License

By contributing, you agree that your contributions will be licensed under the specified [![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE).

## Questions?

Feel free to:
- Open an issue for questions
- Contact the maintainers

Thank you for contributing to Jam-Text!