package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
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

	fs.Parse(args[1:])

	// Setup logger
	var logger *log.Logger
	if *logFile != "" {
		f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
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
			// TODO: Add a function to read chunk contents
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
			}
		}

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
			return fmt.Errorf("No similar hashes found")
		}

		fmt.Printf("Fond %d similar hashes for SimHash %x\n", len(resultMap), hash)
		count := 0
			fmt.Printf("No exact matches found. Trying with increased threshold...\n")
			// Try with a higher threshold
			resultMap, exists = idx.FuzzyLookup(hash, *threshold+2)
			if !exists {
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
		if *input == "" || *secondInput == ""{
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
		//in this case the value ignored is the similarity number which is basically the level of similarity.
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
	fmt.Println("  ./jamtext -c <command> [options]")
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
}
	fmt.Println("  textindex -c hash -i text.txt")
	fmt.Println("  textindex -c compare -i doc1.txt -i2 doc2.txt -o report.txt")
	
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
