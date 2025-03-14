package simhash

import (
	"fmt"
	"os"
	"path/filepath"
)

type DocumentSimilarity struct {
	hyperplanes [][]float64
	vectorizer  Vectorizer
}

func NewDocumentSimilarity() *DocumentSimilarity {
	// Initialize with standard parameters
	hyperplanes := GenerateHyperplanes(VectorDimensions, NumHyperplanes)

	// Use NGramVectorizer for better accuracy with documents
	vectorizer := NewNGramVectorizer(VectorDimensions, 3) // 3-gram vectorization

	return &DocumentSimilarity{
		hyperplanes: hyperplanes,
		vectorizer:  vectorizer,
	}
}

func (ds *DocumentSimilarity) CompareDocuments(doc1, doc2 string) (similarity float64, details string) {
	// Calculate SimHashes for both documents
	hash1 := CalculateWithVectorizer(doc1, ds.hyperplanes, ds.vectorizer)
	hash2 := CalculateWithVectorizer(doc2, ds.hyperplanes, ds.vectorizer)

	// Calculate Hamming distance
	distance := hash1.HammingDistance(hash2)

	// Convert distance to similarity percentage (0-100)
	similarity = 100.0 * (64.0 - float64(distance)) / 64.0

	// Generate detailed report
	var assessment string
	switch {
	case similarity >= 90:
		assessment = "Nearly identical"
	case similarity >= 80:
		assessment = "Very similar"
	case similarity >= 70:
		assessment = "Moderately similar"
	case similarity >= 50:
		assessment = "Somewhat similar"
	default:
		assessment = "Different"
	}

	details = fmt.Sprintf("Similarity: %.2f%%\nHamming Distance: %d\nAssessment: %s\n",
		similarity, distance, assessment)

	return similarity, details
}

func CompareFiles(file1, file2 string) error {
	// Read files
	content1, err := os.ReadFile(file1)
	if err != nil {
		return fmt.Errorf("error reading %s: %w", file1, err)
	}

	content2, err := os.ReadFile(file2)
	if err != nil {
		return fmt.Errorf("error reading %s: %w", file2, err)
	}

	detector := NewDocumentSimilarity()
	similarity, details := detector.CompareDocuments(string(content1), string(content2))

	fmt.Printf("\nComparison of %s and %s:\n%s",
		filepath.Base(file1), filepath.Base(file2), details)

	// Optional: Save detailed report
	if similarity >= 70 {
		report := fmt.Sprintf("High similarity detected!\n\nFile 1: %s\nFile 2: %s\n\n%s",
			file1, file2, details)
		err := os.WriteFile("similarity_report.txt", []byte(report), 0o644)
		if err != nil {
			return fmt.Errorf("error saving report: %w", err)
		}
	}

	return nil
}
