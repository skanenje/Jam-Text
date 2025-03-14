package cli

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"jamtext/internal/chunk"
	"jamtext/internal/index"
	"jamtext/internal/simhash"
)

func Run(args []string) error {
	fs := flag.NewFlagSet("textindex", flag.ExitOnError)

	verbose := fs.Bool("v", false, "Enable Verbose output")
	logFile := fs.String("log", "", "Log file path(default: stderr)")

	// Basic commands
	cmd := fs.String("c", "", "Command to run")
	input := fs.String("i", "", "Input file path")
	output := fs.String("o", "", "Output file path")
	size := fs.Int("s", 4096, "Chunk size in bytes")
	hashStr := fs.String("h", "", "SimHash value to lookup")
	secondInput := fs.String("i2", "", "Second input file for comparison")

	// Advanced commands
	overlapSize := fs.Int("overlap", 256, "Overlap size in bytes")
	splitBoundary := fs.Bool("boundary", true, "Split on text boundaries")
	boundaryChars := fs.String("boundary-chars", ".!?\n", "Text boundaries to split on")
	maxChunkSize := fs.Int("max-size", 6144, "Maximum chunk size in bytes")
	preserveNewlines := fs.Bool("preserve-nl", true, "Preserve newlines in chunks")
	indexDir := fs.String("index-dir", "", "Directory to store index shards")
	contextBefore := fs.Int("context-before", 100, "Number of bytes to include before chunk")
	contextAfter := fs.Int("context-after", 100, "Number of bytes to include after chunk")
	threshold := fs.Int("threshold", 3, "Threshold for fuzzy lookup")

	// Content moderation flags
	wordlistPath := fs.String("wordlist", "", "Path to wordlist file")
	modLevel := fs.String("level", "strict", "Moderation level (strict|lenient)")
	contextSize := fs.Int("context", 50, "Context size for matches")

	fs.Parse(args[1:])

	// Setup logger
	var logger *log.Logger
	if *logFile != "" {
		f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
		if err != nil {
			return fmt.Errorf("error opening log file: %v", err)
		}
		defer f.Close()
		logger = log.New(f, "", log.LstdFlags)
	} else if *verbose {
		logger = log.New(os.Stderr, "", log.LstdFlags)
	} else {
		logger = log.New(io.Discard, "", 0) // Discard logs unless verbose or log file specified
	}

	switch *cmd {
	case "index":
		if *input == "" || *output == "" {
			return fmt.Errorf("input and output file paths must be specified")
		}

		hyperplanes := simhash.GenerateHyperplanes(simhash.VectorDimensions, simhash.NumHyperplanes)

		opts := chunk.ChunkOptions{
			ChunkSize:        *size,
			OverlapSize:      *overlapSize,
			SplitOnBoundary:  *splitBoundary,
			BoundaryChars:    *boundaryChars,
			MaxChunkSize:     *maxChunkSize,
			PreserveNewlines: *preserveNewlines,
			Logger:           logger,
			Verbose:          *verbose,
		}

		start := time.Now()

		idx, err := chunk.ProcessFile(*input, opts, hyperplanes, *indexDir)
		if err != nil {
			return err
		}

		if err := index.Save(idx, *output); err != nil {
			return err
		}

		stats := idx.Stats()
		fmt.Printf("Indexed %d unique hashes with %d total positions in %v\n",
			stats["unique_hashes"],
			stats["total_positions"],
			time.Since(start))
		fmt.Printf("Created %d shards\n", stats["shards"])

		return nil

	case "lookup":
		if *input == "" || *hashStr == "" {
			return fmt.Errorf("input and hash must be specified")
		}

		idx, err := index.Load(*input)
		if err != nil {
			return err
		}

		var hash simhash.SimHash
		if _, err := fmt.Sscanf(*hashStr, "%x", &hash); err != nil {
			return fmt.Errorf("invalid hash: %w", err)
		}

		matches, err := idx.Lookup(hash)
		if err != nil {
			return err
		}

		if len(matches) == 0 {
			return fmt.Errorf("no matches found for hash %x", hash)
		}

		formatLookupOutput(idx.SourceFile, hash, map[simhash.SimHash][]int64{hash: matches}, idx.ChunkSize, *contextBefore, *contextAfter)

		defer idx.Close()
		return nil

	case "stats":
		if *input == "" {
			return fmt.Errorf("input file must be specified")
		}

		idx, err := index.Load(*input)
		if err != nil {
			return err
		}

		stats := idx.Stats()
		fmt.Println("Index Statistics:")
		fmt.Printf("Source file: %s\n", stats["source_file"])
		fmt.Printf("Chunk size: %d bytes\n", stats["chunk_size"])
		fmt.Printf("Created: %v\n", stats["created"])
		fmt.Printf("Shards: %d\n", stats["shards"])
		fmt.Printf("Unique hashes: %d\n", stats["unique_hashes"])
		fmt.Printf("Total positions: %d\n", stats["total_positions"])

		defer idx.Close()
		return nil

	case "fuzzy":
		if *input == "" || *hashStr == "" {
			return fmt.Errorf("input and hash must be specified")
		}

		idx, err := index.Load(*input)
		if err != nil {
			return err
		}

		var hash simhash.SimHash
		if _, err := fmt.Sscanf(*hashStr, "%x", &hash); err != nil {
			return fmt.Errorf("invalid hash: %w", err)
		}

		resultMap, exists := idx.FuzzyLookup(hash, *threshold)
		if !exists {
			fmt.Printf("No exact matches found. Trying with increased threshold...\n")
			// Try with a higher threshold
			resultMap, exists = idx.FuzzyLookup(hash, *threshold+2)
			if !exists {
				return fmt.Errorf("no similar hashes found within threshold %d", *threshold+2)
			}
		}

		fmt.Printf("Found %d similar hashes for SimHash %x\n", len(resultMap), hash)

		// Get first hash and its positions from the result map
		var firstPositions []int64
		for _, p := range resultMap {
			firstPositions = p
			break
		}

		originalChunk, _, _, err := chunk.ReadChunk(idx.SourceFile, firstPositions[0], idx.ChunkSize, 0, 0)
		if err != nil {
			return fmt.Errorf("failed to read original chunk: %w", err)
		}

		fmt.Printf("Searching for matches to text chunk:\n%s\n\n", originalChunk)

		for similarHash, positions := range resultMap {
			distance := hash.HammingDistance(similarHash)
			fmt.Printf("SimHash: %016x (distance: %d) - %d matches\n",
				similarHash,
				distance,
				len(positions))

			fmt.Printf("\nHash: %x (Hamming distance: %d)\n", similarHash, distance)

			for _, pos := range positions {
				showMatchContext(idx.SourceFile, pos, idx.ChunkSize, originalChunk)
			}
		}

		defer idx.Close()
		return nil
	case "hash":
		if *input == "" {
			return fmt.Errorf("input file must be specified")
		}

		// Generate hyperplanes
		hyperplanes := simhash.GenerateHyperplanes(simhash.VectorDimensions, simhash.NumHyperplanes)

		if *verbose {
			logger.Printf("Generated %d hyperplanes\n", len(hyperplanes))
		}

		// Read the file content
		content, err := os.ReadFile(*input)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		// Calculate hash
		hash := simhash.Calculate(string(content), hyperplanes)
		fmt.Printf("%x\n", hash) // Only output the hash
		return nil

	case "compare":
		if *input == "" || *secondInput == "" {
			return fmt.Errorf("first input file must be specified")
		}

		// secondInput := fs.String("i2", "", "Second input file path")
		if *secondInput == "" {
			return fmt.Errorf("second input file must be specified")
		}

		// Read first file
		content1, err := os.ReadFile(*input)
		if err != nil {
			return fmt.Errorf("error reading %s: %w", *input, err)
		}

		// Read second file
		content2, err := os.ReadFile(*secondInput)
		if err != nil {
			return fmt.Errorf("error reading %s: %w", *secondInput, err)
		}

		detector := simhash.NewDocumentSimilarity()
		// in this case the value ignored is the similarity number which is basically the level of similarity.
		_, details := detector.CompareDocuments(string(content1), string(content2))

		fmt.Println(details)

		if *output != "" {
			report := fmt.Sprintf("Comparison Report\n\nFile 1: %s\nFile 2: %s\n\n%s",
				*input, *secondInput, details)
			if err := os.WriteFile(*output, []byte(report), 0o644); err != nil {
				return fmt.Errorf("error writing report: %w", err)
			}
			fmt.Printf("Report saved to %s\n", *output)
		}

	case "moderate":
		if *input == "" {
			return fmt.Errorf("input file must be specified")
		}
		if *wordlistPath == "" {
			return fmt.Errorf("wordlist file must be specified")
		}

		matches, err := processModeration(*input, *wordlistPath, *modLevel, *contextSize, logger, *verbose)
		if err != nil {
			return err
		}

		if matches > 0 {
			return fmt.Errorf("found %d potentially inappropriate matches", matches)
		}
		return nil

	default:
		// TODO: Setup chunk options
		printUsage(fs)
		return fmt.Errorf("unknown command: %s", *cmd)
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func printUsage(fs *flag.FlagSet) {
	fmt.Println("JamText - A text indexing and similarity search tool")
	fmt.Println("\nUsage:")
	fmt.Println("  textindex -c <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  index  - Create index from text file")
	fmt.Println("  lookup - Exact lookup by SimHash")
	fmt.Println("  fuzzy  - Fuzzy lookup by SimHash with threshold")
	fmt.Println("  hash   - Calculate SimHash for a file")
	fmt.Println("  stats  - Show index statistics")
	fmt.Println("  compare - Compare two text files for similarity")
	fmt.Println("\nOptions:")
	fs.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  textindex -c index -i book.txt -o book.idx -s 4096")
	fmt.Println("  textindex -c lookup -i book.idx -h a1b2c3d4e5f6")
	fmt.Println("  textindex -c fuzzy -i book.idx -h a1b2c3d4e5f6 -threshold 5")
	fmt.Println("  textindex -c hash -i text.txt")
	fmt.Println("  textindex -c compare -i doc1.txt -i2 doc2.txt -o report.txt")
}

// Add this function to help verify matches
func showMatchContext(sourceFile string, position int64, chunkSize int, originalText string) {
	content, _, _, err := chunk.ReadChunk(sourceFile, position, chunkSize, 0, 0)
	if err != nil {
		return
	}

	commonText := findLongestCommonSubstring(content, originalText)
	if len(commonText) > 50 {
		// Limit output to first 200 characters with ellipsis if needed
		matchedText := commonText
		if len(matchedText) > 200 {
			matchedText = matchedText[:200] + "..."
		}
		fmt.Printf("\nMatched text:\n%s\n", matchedText)

		// Limit context to first 100 characters
		context := content
		if len(context) > 100 {
			context = context[:100] + "..."
		}
		fmt.Printf("\nOriginal context:\n%s\n", context)
		fmt.Printf("\nSimilarity: %.2f%%\n", float64(len(commonText))/float64(len(content))*100)
	}
}

func findLongestCommonSubstring(s1, s2 string) string {
	// Create DP table
	m, n := len(s1), len(s2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	// Track maximum length and ending position
	maxLength := 0
	endPos := 0
	startPos := 0

	// Fill DP table and track all matches above minimum length
	minMatchLength := 20 // Minimum length to consider as potential plagiarism

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
				if dp[i][j] > maxLength && dp[i][j] >= minMatchLength {
					maxLength = dp[i][j]
					endPos = i
					startPos = i - maxLength
				}
			}
		}
	}

	if maxLength < minMatchLength {
		return "" // No significant match found
	}

	// Extract the longest common substring
	return s1[startPos:endPos]
}

// Add this struct to store word occurrences
type WordOccurrence struct {
	Word      string
	Count     int
	Locations []struct {
		LineNum int
		Context string
	}
}

func processModeration(inputPath, wordlistPath, modLevel string, contextSize int, logger *log.Logger, verbose bool) (int, error) {
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read input file: %w", err)
	}

	wordlist, err := os.ReadFile(wordlistPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read wordlist: %w", err)
	}

	// Store offensive words with their severity in a map
	words := make(map[string]string) // word -> severity
	scanner := bufio.NewScanner(strings.NewReader(string(wordlist)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, ",")
		word := strings.ToLower(strings.TrimSpace(parts[0]))
		severity := "medium" // default severity
		if len(parts) > 1 {
			severity = strings.ToLower(strings.TrimSpace(parts[1]))
		}
		words[word] = severity
	}

	matches := 0
	lines := strings.Split(string(content), "\n")

	// Track word statistics with severity
	type occurrence struct {
		count    int
		severity string
		contexts []string
		lineNums []int
	}
	wordStats := make(map[string]*occurrence)

	fmt.Printf("\n📑 Content Moderation Report\n")
	fmt.Printf("========================\n\n")

	// Iterate over each line and check for offensive words
	for i, line := range lines {
		lineNum := i + 1
		foundWords := make(map[string]bool) // Use map to avoid duplicates per line

		for word, severity := range words {
			var found bool

			if modLevel == "strict" {
				// Strict mode: Match only whole words
				for _, token := range strings.Fields(line) {
					cleanedToken := strings.ToLower(strings.Trim(token, ".,!?\"'"))
					if cleanedToken == word {
						foundWords[word] = true
						found = true
						if _, exists := wordStats[word]; !exists {
							wordStats[word] = &occurrence{severity: severity}
						}
						wordStats[word].count++
						wordStats[word].lineNums = append(wordStats[word].lineNums, lineNum)
						wordStats[word].contexts = append(wordStats[word].contexts,
							truncateContext(line, contextSize))
					}
				}
			} else {
				// Lenient mode: Match if the word appears anywhere in the line
				if strings.Contains(strings.ToLower(line), word) {
					foundWords[word] = true
					found = true
					if _, exists := wordStats[word]; !exists {
						wordStats[word] = &occurrence{severity: severity}
					}
					wordStats[word].count++
					wordStats[word].lineNums = append(wordStats[word].lineNums, lineNum)
					wordStats[word].contexts = append(wordStats[word].contexts,
						truncateContext(line, contextSize))
				}
			}

			if found {
				matches++
			}
		}

		// If any offensive words were found, print them by severity
		if len(foundWords) > 0 {
			fmt.Printf("🚨 Line %d:\n", lineNum)
			fmt.Printf("   %s\n", line)

			// Group words by severity
			severityGroups := make(map[string][]string)
			for w := range foundWords {
				sev := words[w]
				severityGroups[sev] = append(severityGroups[sev], w)
			}

			// Print words by severity
			for _, sev := range []string{"high", "medium", "low"} {
				if words, ok := severityGroups[sev]; ok {
					sort.Strings(words)
					fmt.Printf("❌ %s severity: %s\n",
						strings.Title(sev),
						strings.Join(words, ", "))
				}
			}

			if verbose {
				fmt.Printf("   Context: %s\n", truncateContext(line, contextSize))
			}
			fmt.Println()
		}
	}

	// Print summary statistics
	if matches > 0 {
		fmt.Printf("\n📊 Summary Statistics\n")
		fmt.Printf("==================\n")
		fmt.Printf("Total matches found: %d\n\n", matches)

		// Print statistics by severity
		for _, severity := range []string{"high", "medium", "low"} {
			fmt.Printf("\n%s Severity Words:\n", strings.Title(severity))
			fmt.Printf("------------------\n")

			// Collect words of current severity
			var sevWords []struct {
				word  string
				stats *occurrence
			}
			for word, stats := range wordStats {
				if stats.severity == severity {
					sevWords = append(sevWords, struct {
						word  string
						stats *occurrence
					}{word, stats})
				}
			}

			// Sort by frequency
			sort.Slice(sevWords, func(i, j int) bool {
				return sevWords[i].stats.count > sevWords[j].stats.count
			})

			// Print details
			for _, w := range sevWords {
				fmt.Printf("'%s' (%d occurrences)\n", w.word, w.stats.count)
				if verbose {
					for i, context := range w.stats.contexts {
						fmt.Printf("  Line %d: %s\n",
							w.stats.lineNums[i], context)
					}
				}
			}
		}
	} else {
		fmt.Println("✅ No concerning content found")
	}

	return matches, nil
}

func truncateContext(text string, size int) string {
	if len(text) <= size {
		return text
	}
	return text[:size/2] + "..." + text[len(text)-size/2:]
}

func formatLookupOutput(sourceFile string, hash simhash.SimHash, matches map[simhash.SimHash][]int64, chunkSize int, contextBefore, contextAfter int) {
	fmt.Printf("Found matches for SimHash %x:\n\n", hash)
	
	for simHash, positions := range matches {
		fmt.Printf("\nHash: %x (Hamming distance: %d)\n", simHash, hash.HammingDistance(simHash))
		for _, pos := range positions {
			chunk, before, after, err := chunk.ReadChunk(sourceFile, pos, chunkSize, contextBefore, contextAfter)
			if err != nil {
				fmt.Printf("Error reading chunk at position %d: %v\n", pos, err)
				continue
			}
			
			fmt.Printf("Match at position %d:\n", pos)
			if len(before) > 0 {
				fmt.Printf("Before: %s\n", before)
			}
			fmt.Printf("Chunk:  %s\n", chunk)
			if len(after) > 0 {
				fmt.Printf("After:  %s\n", after)
			}
			fmt.Println("---")
		}
	}
}
