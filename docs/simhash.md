# SimHash Package

## Overview
The SimHash package implements text fingerprinting using random hyperplane projections for efficient similarity detection. It provides both frequency-based and n-gram based vectorization strategies with LSH (Locality-Sensitive Hashing) support.

## Core Components

### SimHash Type
- 64-bit fingerprint representation for efficient storage and comparison
- Built-in Hamming distance calculation for similarity metrics
- Thread-safe operations for concurrent processing

### Vectorization Strategies

#### 1. Frequency-based
- Word-level tokenization with MD5-based dimension mapping
- Suitable for longer documents with rich vocabulary
- Configurable vector dimensions (default: 128)
- Vector normalization for consistent results

#### 2. N-gram based
- Configurable n-gram size (default: 3)
- Normalized vectors for consistent comparison
- Optimized for short text handling
- Better preservation of word order information

### LSH Implementation
- Permutation table generation for fast similarity search
- Band signature computation with configurable sizes
- Thread-safe random number generation
- Optimized for high-dimensional binary vectors

## Usage Examples

```go
// Initialize similarity detector
detector := NewDocumentSimilarity()

// Compare two documents
similarity, details := detector.CompareDocuments(doc1, doc2)

// Generate hyperplanes for custom implementation
hyperplanes := GenerateHyperplanes(128, 64)

// Calculate SimHash with custom vectorizer
vectorizer := NewNGramVectorizer(128, 3)
hash := CalculateWithVectorizer(text, hyperplanes, vectorizer)
```

## Integration Points
- Works with `chunk` package for text segmentation
- Provides similarity metrics for `index` package
- Supports CLI operations through standardized interfaces

## Performance Considerations
- Use NGramVectorizer for texts shorter than 100 words
- FrequencyVectorizer performs better on longer documents
- Consider LSH for large-scale similarity searches
- Parallel hyperplane generation for improved performance

## Thread Safety
- All public methods are thread-safe
- Safe for concurrent use in worker pools
- Shared hyperplanes are immutable after creation

## Best Practices
- Initialize vectorizers once and reuse
- Use appropriate vector dimensions (128-256 recommended)
- Enable LSH for datasets larger than 10,000 documents
- Consider document length when choosing vectorization strategy
