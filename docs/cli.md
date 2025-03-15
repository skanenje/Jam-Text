# CLI Package

Jamtext is a powerful text analysis tool designed for content creators, researchers, and developers who need to:
- Detect similar or duplicate content across large text collections
- Identify potential plagiarism in academic or professional writing
- Monitor content for inappropriate or unwanted material
- Compare and analyze text similarities in documents
- Create searchable indexes for large text repositories

## Real-World Applications

### Content Creation & Publishing
- Detect duplicate content across your blog posts or articles
- Ensure originality in user-generated content
- Find similar content for internal linking and references

### Academic & Research
- Check student submissions for potential plagiarism
- Analyze text corpus for similar patterns
- Index and search through research papers

### Content Moderation
- Automatically flag inappropriate content
- Monitor user submissions in real-time
- Maintain content quality standards

## Commands
- `index` - Create searchable index from text documents
- `lookup` - Perform exact SimHash lookup for matching content
- `fuzzy` - Find similar content using fuzzy SimHash matching
- `hash` - Generate document fingerprint for comparison
- `compare` - Compare two documents for similarity
- `moderate` - Screen content against moderation rules

## Usage
```bash
./textindex -c <command> [options]

Options:
  -i string      Input file path
  -o string      Output file path
  -s int         Chunk size (default: 4096)
  -overlap int   Overlap size (default: 256)
  -threshold int Similarity threshold (default: 3)
```

## Examples

### Content Indexing
```bash
# Create searchable index from a book
# Index with custom overlap for better matching
./textindex -c index -i content.txt -o content.idx -s 2048 -overlap 512
```

### Similarity Detection
```bash
# Generate hash for comparison
HASH=$(textindex -c hash -i article.txt)

# Find similar content
./textindex -c fuzzy -i database.idx -h $HASH -threshold 5

# Direct document comparison
./textindex -c compare -i original.txt -i2 submission.txt -o report.txt
```

### Content Moderation
```bash
# Check content against moderation rules
./textindex -c moderate -i submission.txt -wordlist forbidden.txt -level strict

# Lenient moderation with context
./textindex  -c moderate -i post.txt -wordlist rules.txt -level lenient -context 100
```

## Performance Tips
- Use larger chunk sizes (8192+) for better performance on large documents
- Reduce overlap for faster indexing at the cost of accuracy
- Adjust threshold based on your similarity requirements
- Use LSH bands for faster similarity search in large datasets

## Integration Examples
```bash
# Git pre-commit hook
#!/bin/bash
HASH=$(textindex -c hash -i "$1")
./textindex -c fuzzy -i repository.idx -h $HASH -threshold 3

# CI/CD Pipeline
./textindex -c moderate -i release-notes.txt -wordlist company-rules.txt || exit 1

# Batch Processing
for file in content/*.txt; do
./textindex -c index -i "$file" -o "indexes/$(basename "$file").idx"
done
```

For detailed implementation examples and API documentation, see the [SimHash Package](simhash.md) and [Chunking System](chunking.md) documentation.
