package cli

import (
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

	// Advanced commands to be added
	overlapSize := fs.Int("overlap", 256, "Overlap size in bytes")
	splitBoundary := fs.Bool("boundary", true, "Split on text boundaries")
	boundaryChars := fs.String("boundary-chars", ".!?\n", "Text boundaries to split on")
	maxChunkSize := fs.Int("max-size", 6144, "Maximum chunk size in bytes")
	preserveNewlines := fs.Bool("preserve-nl", true, "Preserve newlines in chunks")
	indexDir := fs.String("index-dir", "", "Directory to store index shards")
	contextBefore := fs.Int("context-before", 100, "Number of bytes to include before chunk")
	contextAfter := fs.Int("context-after", 100, "Number of bytes to include after chunk")
	threshold := fs.Int("threshold", 3, "Threshold for fuzzy lookup")

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

		var hash simhash.SimHash
		if _, err := fmt.Sscanf(*hashStr, "%x", &hash); err != nil {
			return fmt.Errorf("invalid hash: %w", err)
		}

		idx, err := index.Load(*input)
		if err != nil {
			return err
		}

		positions, err := idx.Lookup(hash)
		if err != nil {
			return fmt.Errorf("SimHash not found")
		}

		fmt.Printf("Found %d positions for SimHash %x\n", len(positions), hash)
		for i, pos := range positions[:min(3, len(positions))] {
			content, contextBeforeStr, contextAfterStr, err := chunk.ReadChunk(idx.SourceFile, pos, idx.ChunkSize, *contextBefore, *contextAfter)
			if err != nil {
				return nil
			}

			preview := content
			if len(content) > 100 {
				preview = content[:100] + "..."
			}

			if contextBeforeStr == "" && contextAfterStr == "" {
				fmt.Printf("%d. Position: %d\n    %s\n", i+1, pos, preview)
			} else if contextBeforeStr == "" && contextAfterStr != "" {
				fmt.Printf("%d. Position: %d\n\n    %s\n\n Context after: %s\n", i+1, pos, preview, contextAfterStr)
			} else if contextBeforeStr != "" && contextAfterStr == "" {
				fmt.Printf("%d. Position: %d\nContext before: %s\n\n    %s\n", i+1, pos, contextBeforeStr, preview)
			} else {
				fmt.Printf("%d. Position: %d\nContext before: %s\n\n    %s\n\n Context after: %s\n", i+1, pos, contextBeforeStr, preview, contextAfterStr)
				fmt.Printf("%d. Position: %d\nContext before: %s\n\n    %s\n\n Context after: %s\n", i+1, pos, contextBeforeStr, preview, contextAfterStr)
			}
		}

		defer idx.Close()
		return nil

	case "dedup":
		fmt.Println("Here")
		if *input == "" {
			return fmt.Errorf("input file must be specified")
		}

		idx, err := index.Load(*input)
		if err != nil {
			fmt.Printf("Error loading index: %v\n", err)
			return err
		}

		start := time.Now()
		deduped := make(map[simhash.SimHash][]int64)
		totalDupes := 0

		// TODO: Add a function to deduplicate the index
		for _, shard := range idx.Shards {
			fmt.Printf("Deduplicating shard %d\n", shard.ShardID)
			if shard == nil {
				continue
			}

			for hash, positions := range shard.SimHashToPos {
				// Keep only unique positions
				seen := make(map[int64]bool)
				unique := []int64{}

				for _, pos := range positions {
					if !seen[pos] {
						seen[pos] = true
						unique = append(unique, int64(pos))
					}
				}

				if len(unique) < len(positions) {
					totalDupes += len(positions) - len(unique)
				}

				deduped[hash] = unique
			}
		}

		// Update index with deduplicated data
		idx.Shards = []*index.IndexShard{{
			SimHashToPos: deduped,
			ShardID:      0,
		}}

		if err := index.Save(idx, *output); err != nil {
			return err
		}

		fmt.Printf("Removed %d duplicate positions in %v\n", totalDupes, time.Since(start))

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
			resultMap, exists = idx.FuzzyLookup(hash, *threshold+2)
			if !exists {
				return fmt.Errorf("no similar hashes found within threshold %d", *threshold+2)
				return fmt.Errorf("no similar hashes found within threshold %d", *threshold+2)
			}
		}

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
		fmt.Printf("%x\n", hash)
		return nil

	case "moderate":
		if *input == "" {
			return fmt.Errorf("input file must be specified")
		}

		opts := struct {
			wordlist    string
			modLevel    string
			contextSize int
		}{
			wordlist:    *fs.String("wordlist", "offensive_words.txt", "Path to offensive words list"),
			modLevel:    *fs.String("level", "strict", "Moderation level (strict/lenient)"),
			contextSize: *fs.Int("context", 50, "Context size in characters"),
		}

		start := time.Now()

		matches, err := processModeration(*input, opts.wordlist, opts.modLevel, opts.contextSize, logger, *verbose)
		if err != nil {
			return fmt.Errorf("moderation failed: %w", err)
		}

		if *verbose {
			fmt.Printf("Completed moderation in %v\n", time.Since(start))
		}

		if matches == 0 {
			fmt.Printf("No offensive content found\n")
		} else {
			fmt.Printf("Found %d instances of offensive content\n", matches)
		}

		return nil

	default:
		// TODO: Setup chunk options
		printUsage(fs)
		return fmt.Errorf("unknown command: %s", *cmd)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func printUsage(fs *flag.FlagSet) {
	fmt.Println("TextIndex - A text indexing and similarity search tool")
	fmt.Println("\nUsage:")
	fmt.Println("  textindex -c <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  index    - Create index from text file")
	fmt.Println("  lookup   - Exact lookup by SimHash")
	fmt.Println("  fuzzy    - Fuzzy lookup by SimHash with threshold")
	fmt.Println("  hash     - Calculate SimHash for a file")
	fmt.Println("  stats    - Show index statistics")
	fmt.Println("  moderate - Scan text for offensive content")
	fmt.Println("\nOptions:")
	fs.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  textindex -c index -i book.txt -o book.idx -s 4096")
	fmt.Println("  textindex -c lookup -i book.idx -h a1b2c3d4e5f6")
	fmt.Println("  textindex -c moderate -i document.txt -wordlist words.txt -level strict")
}


// Add this function to help verify matches
func showMatchContext(sourceFile string, position int64, chunkSize int, originalText string) {
	// Read the chunk from position
	content, _, _, err := chunk.ReadChunk(sourceFile, position, chunkSize, 0, 0)
	if err != nil {
		return
	}


	// Find the common substring
	commonText := findLongestCommonSubstring(content, originalText)
	if len(commonText) > 50 {
		fmt.Printf("\nMatched text:\n%s\n", commonText)
		fmt.Printf("\nOriginal context:\n%s\n", content[:min(200, len(content))])
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

	// Store offensive words in a map for quick lookup
	words := make(map[string]bool)
	for _, word := range strings.Fields(string(wordlist)) {
		words[strings.ToLower(strings.TrimSpace(word))] = true
	}

	matches := 0
	lines := strings.Split(string(content), "\n")

	// Track word statistics
	wordStats := make(map[string]int)

	fmt.Printf("\n📑 Content Moderation Report\n")
	fmt.Printf("========================\n\n")

	// Iterate over each line and check for offensive words
	for i, line := range lines {
		lineNum := i + 1
		foundWords := make(map[string]bool) // Use map to avoid duplicates per line

		for word := range words {
			var found bool

			if modLevel == "strict" {
				// Strict mode: Match only whole words
				for _, token := range strings.Fields(line) {
					cleanedToken := strings.ToLower(strings.Trim(token, ".,!?\"'"))
					if cleanedToken == word {
						foundWords[word] = true
						wordStats[word]++
						found = true
					}
				}
			} else {
				// Lenient mode: Match if the word appears anywhere in the line
				if strings.Contains(strings.ToLower(line), word) {
					foundWords[word] = true
					wordStats[word]++
					found = true
				}
			}

			if found {
				matches++
			}
		}

		// If any offensive words were found, print them
		if len(foundWords) > 0 {
			fmt.Printf("🚨 Line %d:\n", lineNum)
			fmt.Printf("   %s\n", line)

			// Convert map keys to slice for sorting
			var words []string
			for w := range foundWords {
				words = append(words, w)
			}
			sort.Strings(words)

			fmt.Printf("❌ Found: %s\n", strings.Join(words, ", "))
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

		// Sort words by frequency
		type wordCount struct {
			word  string
			count int
		}
		var sorted []wordCount
		for word, count := range wordStats {
			sorted = append(sorted, wordCount{word, count})
		}
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].count > sorted[j].count
		})

		fmt.Printf("Word Frequency:\n")
		for _, wc := range sorted {
			fmt.Printf("- '%s': %d occurrence(s)\n", wc.word, wc.count)
		}
	} else {
		fmt.Printf("\n✅ No offensive content found\n")
	}

	return matches, nil
}

func truncateContext(text string, size int) string {
	if len(text) <= size {
		return text
	}
	return text[:size] + "..."
}
