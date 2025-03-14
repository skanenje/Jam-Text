# Jam-Text Core Features

Jam-Text provides three powerful text analysis capabilities: similarity detection, duplication detection, and content moderation. This guide explains each feature and its usage.

## 1. Similarity Detection

Uses SimHash fingerprinting with random hyperplane projections to detect similar text content, even with minor variations.

### Key Capabilities
- 64-bit fingerprint representation
- Configurable similarity thresholds
- Support for both frequency-based and n-gram vectorization
- Detailed similarity reports

### Usage Examples

#### Direct Document Comparison
```bash
# Compare two documents with detailed report
jamtext -c compare -i doc1.txt -i2 doc2.txt -o report.txt
```

#### Fuzzy Search
```bash
# Generate hash for reference document
HASH=$(jamtext -c hash -i reference.txt)

# Search for similar content with threshold
jamtext -c fuzzy -i corpus.idx -h $HASH -threshold 5
```

#### Best Practices
- Use n-gram vectorization for texts under 100 words
- Use frequency vectorization for longer documents
- Adjust threshold based on desired precision (3-5 recommended)

## 2. Duplication Detection

Identifies duplicate or near-duplicate content across large text collections using LSH (Locality-Sensitive Hashing).

### Key Capabilities
- Efficient large-scale similarity search
- Configurable chunk sizes for granular detection
- Overlapping text analysis
- Context preservation

### Usage Examples

#### Index Creation for Duplication Detection
```bash
# Create searchable index with optimal settings
jamtext -c index -i source.txt -o index.dat \
    -s 4096 \
    -overlap 256 \
    -boundary=true \
    -boundary-chars=".!?\n" \
    -preserve-nl=true
```

#### Duplication Search
```bash
# Search for duplicates with context
jamtext -c lookup -i index.dat -h $HASH \
    -context-before 100 \
    -context-after 100
```

#### Best Practices
- Use smaller chunk sizes (2048-4096) for precise detection
- Enable boundary splitting for natural text segments
- Adjust overlap size based on content structure

## 3. Content Moderation

Screens text content against customizable word lists with context-aware detection.

### Key Capabilities
- Configurable moderation levels (strict/lenient)
- Context preservation around matches
- Detailed occurrence reporting
- Custom wordlist support

### Usage Examples

#### Basic Moderation
```bash
# Check content against wordlist
jamtext -c moderate -i content.txt \
    -wordlist words.txt \
    -level strict \
    -context 50
```

#### Advanced Usage
```bash
# Verbose mode with detailed reporting
jamtext -c moderate -i content.txt \
    -wordlist words.txt \
    -level strict \
    -context 50 \
    -v
```

#### Best Practices
- Use "strict" mode for exact word matching
- Use "lenient" mode for partial matching
- Maintain separate wordlists for different content types
- Adjust context size based on content structure

## Integration Examples

### Combined Workflow
```bash
# 1. Create index for similarity search
jamtext -c index -i corpus.txt -o corpus.idx

# 2. Check for inappropriate content
jamtext -c moderate -i new_doc.txt -wordlist words.txt

# 3. If moderation passes, search for duplicates
HASH=$(jamtext -c hash -i new_doc.txt)
jamtext -c fuzzy -i corpus.idx -h $HASH -threshold 3
```

## Performance Considerations

### Similarity Detection
- Use appropriate vectorization strategy based on document length
- Enable LSH for datasets larger than 10,000 documents
- Initialize vectorizers once and reuse

### Duplication Detection
- Balance chunk size with detection granularity
- Use appropriate overlap for content type
- Consider sharding for very large datasets

### Content Moderation
- Optimize wordlist size and organization
- Use appropriate moderation level for content type
- Consider context size impact on performance

## Additional Resources
- [CLI Documentation](../cli.md) for detailed command options
- [SimHash Documentation](../simhash.md) for similarity detection details
- [Performance Tuning](performance.md) for optimization guidelines