# SimHash Package

Text fingerprinting using random hyperplane projections.

## Features
- 64-bit fingerprint representation
- Frequency-based and n-gram vectorization
- LSH support for fast similarity search
- Thread-safe operations

## Usage
```go
// Initialize detector
detector := NewDocumentSimilarity()

// Compare documents
similarity, details := detector.CompareDocuments(doc1, doc2)

// Custom implementation
hyperplanes := GenerateHyperplanes(128, 64)
vectorizer := NewNGramVectorizer(128, 3)
hash := CalculateWithVectorizer(text, hyperplanes, vectorizer)
```

## Best Practices
- Use NGramVectorizer for texts < 100 words
- Use FrequencyVectorizer for longer documents
- Configure LSH bands based on dataset size