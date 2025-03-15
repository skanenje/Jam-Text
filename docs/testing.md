# Testing Documentation

## Overview
This document outlines the testing practices and patterns used in the Jam-Text project. The test suite covers unit tests, integration tests, and various testing techniques to ensure code quality and reliability.

## Test Structure

### Core Test Packages
- `cmd/main_test.go` - CLI executable tests
- `internal/simhash/similarity_test.go` - SimHash comparison tests
- `internal/cli/cli_test.go` - Command-line interface tests
- `internal/chunk/chunk_test.go` - Text chunking tests
- `internal/chunk/processor_test.go` - Chunk processing tests
- `internal/index/index_test.go` - Index management tests

## Testing Techniques

### Output Capture
Used to test CLI output and logging:
```go
func captureOutput(f func() error) (string, error) {
    oldStdout := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w
    // ... capture and return output
}
```

### Temporary Files
Creating and managing test files:
```go
tmpDir := t.TempDir()
tmpFile := filepath.Join(tmpDir, "test.txt")
// ... use temporary files
// cleanup handled automatically by t.TempDir()
```

### Table-Driven Tests
Used extensively for testing multiple scenarios:
```go
tests := []struct {
    name     string
    input    string
    want     string
    wantErr  bool
}{
    // ... test cases
}
```

## Key Test Functions

### SimHash Package Tests
- `TestNewDocumentSimilarity` - Validates similarity detector initialization
- `TestCompareDocuments` - Tests document comparison with various similarity levels
- `TestCompareFiles` - Tests file-based comparison
- `TestCompareFilesWithLongContent` - Tests handling of large documents

### CLI Tests
- `TestRunValidation` - Validates CLI argument parsing
- `TestGetContentHash` - Tests hash generation command
- `TestRunFuzzyCommand` - Tests fuzzy search functionality

### Chunk Processing Tests
- `TestReadChunk` - Tests chunk reading with different sizes and contexts
- `TestChunkWithMetadata` - Validates metadata handling
- `TestReadChunkWithDifferentFormats` - Tests various file format support
- `TestProcessorWithInvalidContent` - Tests edge cases and invalid inputs

### Index Tests
- `TestNew` - Validates index creation
- `TestSave` - Tests index persistence
- `TestLoad` - Tests index loading
- `TestLookup` - Tests exact hash lookup
- `TestFuzzyLookup` - Tests similarity-based lookup

## Running Tests

### Basic Test Execution
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run verbose tests
go test -v ./...
```

### Package-Specific Tests
```bash
# Test specific packages
go test ./internal/cli/...
go test ./internal/simhash/...
```

### Coverage Reports
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Test Categories

### Unit Tests
- Individual function testing
- Edge case validation
- Error handling verification

### Integration Tests
- Cross-package functionality
- File system operations
- CLI command execution

### Performance Tests
- Large document processing
- Concurrent operations
- Resource usage monitoring

## Best Practices

### Test Setup
1. Use `t.TempDir()` for temporary files
2. Initialize test data in `TestMain` when needed
3. Use helper functions for common setup tasks

### Test Organization
1. Group related test cases
2. Use clear, descriptive test names
3. Separate test utilities into helper functions

### Assertions
1. Use specific error messages
2. Check both positive and negative cases
3. Validate error types when testing error conditions

### Cleanup
1. Use `defer` for cleanup operations
2. Clean up temporary resources
3. Restore modified environment variables

## Common Test Patterns

### Error Testing
```go
if err := function(); (err != nil) != tt.wantErr {
    t.Errorf("function() error = %v, wantErr %v", err, tt.wantErr)
}
```

### Timeout Pattern
```go
select {
case <-done:
case <-time.After(500 * time.Millisecond):
    t.Fatal("Test timed out")
}
```

### Resource Cleanup
```go
defer func() {
    if err := cleanup(); err != nil {
        t.Errorf("cleanup failed: %v", err)
    }
}()
```

## Continuous Integration
- Tests run on every pull request
- Coverage reports generated automatically
- Performance benchmarks tracked over time