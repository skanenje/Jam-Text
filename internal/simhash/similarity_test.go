package simhash

import (
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func TestNewDocumentSimilarity(t *testing.T) {
    ds := NewDocumentSimilarity()
    
    if ds == nil {
        t.Fatal("NewDocumentSimilarity returned nil")
    }
    
    if ds.hyperplanes == nil {
        t.Error("hyperplanes not initialized")
    }
    
    if len(ds.hyperplanes) != NumHyperplanes {
        t.Errorf("expected %d hyperplanes, got %d", NumHyperplanes, len(ds.hyperplanes))
    }
    
    if ds.vectorizer == nil {
        t.Error("vectorizer not initialized")
    }
}

func TestCompareDocuments(t *testing.T) {
    tests := []struct {
        name           string
        doc1           string
        doc2           string
        minSimilarity float64
        maxSimilarity float64
        wantAssessment string
    }{
        {
            name:           "identical documents",
            doc1:           "This is a test document",
            doc2:           "This is a test document",
            minSimilarity: 95.0,
            maxSimilarity: 100.0,
            wantAssessment: "Nearly identical",
        },
        {
            name:           "very similar documents",
            doc1:           "This is a test document",
            doc2:           "This is a test document!",
            minSimilarity: 90.0,
            maxSimilarity: 100.0,
            wantAssessment: "Nearly identical",
        },
        {
            name:           "moderately similar documents",
            doc1:           "This is a test document",
            doc2:           "This is another test file",
            minSimilarity: 60.0,
            maxSimilarity: 90.0,
            wantAssessment: "Moderately similar",
        },
        {
            name:           "somewhat similar documents",
            doc1:           "This is a test document",
            doc2:           "This is something completely different",
            minSimilarity: 40.0,
            maxSimilarity: 70.0,
            wantAssessment: "Somewhat similar",
        },
        {
            name:           "different documents",
            doc1:           "This is a test document",
            doc2:           "Something entirely different here",
            minSimilarity: 40.0,
            maxSimilarity: 60.0,
            wantAssessment: "Somewhat similar",
        },
        {
            name:           "empty documents",
            doc1:           "",
            doc2:           "",
            minSimilarity: 95.0,
            maxSimilarity: 100.0,
            wantAssessment: "Nearly identical",
        },
    }

    ds := NewDocumentSimilarity()
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            similarity, details := ds.CompareDocuments(tt.doc1, tt.doc2)
            
            if similarity < tt.minSimilarity || similarity > tt.maxSimilarity {
                t.Errorf("similarity = %.2f, want between %.2f and %.2f", 
                    similarity, tt.minSimilarity, tt.maxSimilarity)
            }
            
            if !strings.Contains(details, tt.wantAssessment) {
                t.Errorf("details = %q, want to contain %q", details, tt.wantAssessment)
            }
        })
    }
}

func TestCompareFiles(t *testing.T) {
    // Create temporary test files
    tmpDir := t.TempDir()
    
    file1Path := filepath.Join(tmpDir, "test1.txt")
    file2Path := filepath.Join(tmpDir, "test2.txt")
    
    content1 := "This is the first test document."
    content2 := "This is the second test document."
    
    if err := os.WriteFile(file1Path, []byte(content1), 0644); err != nil {
        t.Fatalf("failed to create test file 1: %v", err)
    }
    if err := os.WriteFile(file2Path, []byte(content2), 0644); err != nil {
        t.Fatalf("failed to create test file 2: %v", err)
    }
    
    tests := []struct {
        name    string
        file1   string
        file2   string
        wantErr bool
    }{
        {
            name:    "existing files",
            file1:   file1Path,
            file2:   file2Path,
            wantErr: false,
        },
        {
            name:    "non-existent first file",
            file1:   "nonexistent.txt",
            file2:   file2Path,
            wantErr: true,
        },
        {
            name:    "non-existent second file",
            file1:   file1Path,
            file2:   "nonexistent.txt",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := CompareFiles(tt.file1, tt.file2)
            if (err != nil) != tt.wantErr {
                t.Errorf("CompareFiles() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestCompareFilesWithLongContent(t *testing.T) {
    // Create temporary test files
    tmpDir := t.TempDir()
    
    file1Path := filepath.Join(tmpDir, "long1.txt")
    file2Path := filepath.Join(tmpDir, "long2.txt")
    
    // Create long content with some similarities
    content1 := strings.Repeat("This is a test paragraph with some unique content. ", 100)
    content2 := strings.Repeat("This is a test paragraph with some different content. ", 100)
    
    if err := os.WriteFile(file1Path, []byte(content1), 0644); err != nil {
        t.Fatalf("failed to create long test file 1: %v", err)
    }
    if err := os.WriteFile(file2Path, []byte(content2), 0644); err != nil {
        t.Fatalf("failed to create long test file 2: %v", err)
    }
    
    err := CompareFiles(file1Path, file2Path)
    if err != nil {
        t.Errorf("CompareFiles() failed with long content: %v", err)
    }
}