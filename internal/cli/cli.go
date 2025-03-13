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
	cmd := fs.String("cmd", "", "Command to run")
	input := fs.String("i", "", "Input file path")
	output := fs.String("o", "", "Output file path")
	size := fs.Int("s", 4096, "Chunk size in bytes")
	hashStr := fs.String("h", "", "SimHash algorithm to use")
	contextBefore := fs.Int("context-before", 100, "Number of bytes to include before hash")
	contextAfter := fs.Int("context-after", 100, "Number of bytes to include after hash")

	// Advanced commands to be added
	overlapSize := fs.Int("overlap", 256, "Overlap size in bytes")
	splitBoundary := fs.Bool("boundary", true, "Split on text boundaries")
	boundaryChars := fs.String("boundary-chars", ".!?\n", "Text boundaries to split on")
	maxChunkSize := fs.Int("max-size", 6144, "Maximum chunk size in bytes")
	preserveNewlines := fs.Bool("preserve-nl", true, "Preserve newlines in chunks")
	indexDir := fs.String("index-dir", "", "Directory to store index shards")

	

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
			return fmt.Errorf("index file and hash must be specified")
		}

		var hash simhash.SimHash
		if _, err := fmt.Sscanf(*hashStr, "%x", &hash); err != nil {
			return fmt.Errorf("invalid hash: %w", err)
		}

		idx, err := index.Load(*input)
		if err != nil {
			return err
		}

		positions, exists := idx.Lookup(hash)
		if !exists {
			return fmt.Errorf("SimHash not found in index")
		}

		fmt.Printf("Found %d matches: \n", len(positions))
		for i, pos := range positions[:min(3, len(positions))] {
			// TODO: Add a function to read chunk contents
			content, err := chunk.ReadChunk(idx.SourceFile, pos, idx.ChunkSize, *contextBefore, *contextAfter)
			if err != nil {
				return err
			}

			preview := content
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}

			fmt.Printf("%d. Position: %d, Preview: %s\n", i+1, pos, preview)
		}

		defer idx.Close()
		return nil

	// case "stats":
	// 	if *input == "" {
	// 		return fmt.Errorf("index file must be specified")
	// 	}

	// 	idx, err := index.Load(*input)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	stats := idx.Stats()


	default:
		// TODO: Setup chunk options

		
		return fmt.Errorf("unknown command: %s", *cmd)
	}

	return nil
}