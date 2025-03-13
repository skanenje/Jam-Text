# SimHash Package

## Overview
The SimHash package implements text fingerprinting using the SimHash algorithm with LSH support.

## Core Components

### SimHash Type
- 64-bit fingerprint representation
- Hamming distance calculation
- Similarity comparison

### Vectorization
1. Frequency-based
   - Word-level tokenization
   - MD5 dimension mapping
   - Vector normalization

2. N-gram based
   - Configurable n-gram size
   - Normalized vectors
   - Short text handling

### LSH Implementation
- Permutation table generation
- Band signature computation
- Configurable band sizes

## Usage Examples

```go
// Generate hyperplanes
hyperplanes := simhash.GenerateHyperplanes(128, 64)

// Create vectorizer
vectorizer := simhash.NewFrequencyVectorizer(128)

// Calculate SimHash
hash := simhash.CalculateWithVectorizer(text, hyperplanes, vectorizer)
```

## Performance Considerations
- Parallel hyperplane generation
- Thread-safe random number generation
- Optimized vector operations