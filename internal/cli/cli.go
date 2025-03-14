package cli

import (
	"flag"
	"fmt"
	"jamtext/internal/chunk"
	"jamtext/internal/index"
	"jamtext/internal/simhash"
	"log"
	"os"
	"time"
)

func Run(args []string) error {
	fs := flag.NewFlagSet("textindex", flag.ExitOnError)

	// TODO: Add more flags
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
		f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("error opening log file: %v", err)
		}
		defer f.Close()
		logger = log.New(f, "", log.LstdFlags)
	} else {
		logger = log.New(os.Stderr, "", log.LstdFlags)
	}

	switch *cmd {
	case "index":
		if *input == "" || *output == "" {
			return fmt.Errorf("input and output file paths must be specified")
		}

		hyperplanes := simhash.GenerateHyperplanes(simhash.VectorDimensions, simhash.NumHyperplanes)

		// TODO: Setup chunk options
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

		// TODO: Add a function to index a file
		idx, err := chunk.ProcessFile(*input, opts, hyperplanes, *indexDir)
		if err != nil {
			return err
		}

		if err := index.Save(idx, *output); err != nil {
			return err
		}

		// TODO: Add a function to get stats for output file
		stats := idx.Stats()
		fmt.Printf("Indexed %d unique hashes with %d total positions in %v\n", 
		                    stats["unique_hashes"], 
		                    stats["total_positions"], 
		                    time.Since(start))
		fmt.Printf("Created %d shards\n", stats["shards"])

		return nil
		
	case "lookup":
		// TODO: Add a function to lookup chunks
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

		var hash simhash.SimHash
		if _, err := fmt.Sscanf(*hashStr, "%x", &hash); err != nil {
			return fmt.Errorf("invalid hash: %w", err)
		}

		idx, err := index.Load(*input)
		if err != nil {
			return err
		}

		// Use fuzzy search to find similar hashes with threshold
		resultMap, exists := idx.FuzzyLookup(hash, *threshold)
		if !exists {
			return fmt.Errorf("No similar hashes found")
		}

		fmt.Printf("Fond %d similar hashes for SimHash %x\n", len(resultMap), hash)
		count := 0
		for similarHash, positions := range resultMap {
			distance := hash.HammingDistance(similarHash)
			fmt.Printf("SimHash: %016x (distance: %d) - %d matches\n",
		                    similarHash,
		                    distance,
		                    len(positions))

			// Show sample positions
			for i, pos := range positions[:min(2, len(positions))] {
				// TODO: Add a function to read chunk contents
				content, _, _, err := chunk.ReadChunk(idx.SourceFile, pos, idx.ChunkSize, *contextBefore, *contextAfter)
				if err != nil {
					return nil
				}

				preview := content
				if len(content) > 100 {
					preview = content[:100] + "..."
				}

				fmt.Printf(" %d.%d. Position: %d, Preview: %s\n", count+1, i+1, pos, preview)
			}

			count++
			if count >= 3 {
				fmt.Printf("... and %d more similar hashes\n", len(resultMap)-3)
				break
			}
		}

		defer idx.Close()
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
	fmt.Println("  index  - Create index from text file")
	fmt.Println("  lookup - Exact lookup by SimHash")
	fmt.Println("  fuzzy  - Fuzzy lookup by SimHash with threshold")
	fmt.Println("  stats  - Show index statistics")
	fmt.Println("\nOptions:")
	fs.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  textindex -c index -i book.txt -o book.idx -s 4096")
	fmt.Println("  textindex -c lookup -i book.idx -h a1b2c3d4e5f6")
	fmt.Println("  textindex -c fuzzy -i book.idx -h a1b2c3d4e5f6 -threshold 5")
}