# CLI Package

JamText is a powerful text analysis tool designed for content creators, researchers, and developers who need to:
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
jamtext -c <command> [options]

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
jamtext -c index -i book.txt -o book.idx -s 4096

# Index with custom overlap for better matching
jamtext -c index -i content.txt -o content.idx -s 2048 -overlap 512
```

### Similarity Detection
```bash
# Generate hash for comparison
HASH=$(jamtext -c hash -i article.txt)

# Find similar content
jamtext -c fuzzy -i database.idx -h $HASH -threshold 5

# Direct document comparison
jamtext -c compare -i original.txt -i2 submission.txt -o report.txt
```

### Content Moderation
```bash
# Check content against moderation rules
jamtext -c moderate -i submission.txt -wordlist forbidden.txt -level strict

# Lenient moderation with context
jamtext -c moderate -i post.txt -wordlist rules.txt -level lenient -context 100
```

## Content Moderation
The `moderate` command allows you to screen content against predefined word lists and rules. It supports two moderation modes: strict and lenient.

### Moderation Options
```bash
jamtext -c moderate [options]

Options:
  -i string        Input file to moderate
  -wordlist string Path to wordlist file
  -level string    Moderation level (strict|lenient) (default: strict)
  -context int     Context size in characters (default: 50)
```

### Wordlist Format
The wordlist file should contain one entry per line in the format:
```
word:severity
```
Where severity can be:
- high: Always flag
- medium: Flag in strict mode
- low: Only flag in specific contexts

Example wordlist:
```
offensive_word1:high
questionable_word2:medium
context_dependent_word3:low
```

### Moderation Levels

#### Strict Mode
- Matches whole words only
- Case-insensitive matching
- Flags all severity levels
- No context consideration
```bash
jamtext -c moderate -i content.txt -wordlist rules.txt -level strict
```

#### Lenient Mode
- Considers word context
- Ignores low-severity matches
- Provides surrounding context
- Allows partial matches
```bash
jamtext -c moderate -i content.txt -wordlist rules.txt -level lenient -context 100
```

### Integration Examples

#### CI/CD Pipeline
```bash
#!/bin/bash
# Fail build if moderation finds issues
jamtext -c moderate -i release-notes.txt -wordlist company-rules.txt || exit 1
```

#### Batch Processing
```bash
#!/bin/bash
# Process multiple files
for file in content/*.txt; do
    echo "Checking $file..."
    jamtext -c moderate -i "$file" -wordlist rules.txt -level strict
done
```

#### Real-time Moderation
```bash
#!/bin/bash
# Monitor new content
inotifywait -m /content/incoming -e create |
while read path action file; do
    jamtext -c moderate -i "$path/$file" -wordlist rules.txt -level strict
done
```

### Exit Codes
- 0: No issues found
- 1: Moderation issues detected
- 2: Processing error

### Performance Considerations
- Use strict mode for smaller documents
- Use lenient mode with larger context for longer content
- Batch process multiple files for better performance
- Consider using parallel processing for large datasets

## Performance Tips
- Use larger chunk sizes (8192+) for better performance on large documents
- Reduce overlap for faster indexing at the cost of accuracy
- Adjust threshold based on your similarity requirements
- Use LSH bands for faster similarity search in large datasets

## Integration Examples
```bash
# Git pre-commit hook
#!/bin/bash
HASH=$(jamtext -c hash -i "$1")
jamtext -c fuzzy -i repository.idx -h $HASH -threshold 3

# CI/CD Pipeline
jamtext -c moderate -i release-notes.txt -wordlist company-rules.txt || exit 1

# Batch Processing
for file in content/*.txt; do
    jamtext -c index -i "$file" -o "indexes/$(basename "$file").idx"
done
```

For detailed implementation examples and API documentation, see the [SimHash Package](simhash.md) and [Chunking System](chunking.md) documentation.
