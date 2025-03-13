package simhash

import (
	"crypto/md5"
	"encoding/binary"
	"math"
	"math/bits"
	"math/rand"
	"strings"
	"sync"
)

const (
	VectorDimensions = 128
	NumHyperplanes   = 64
)

// SimHash represents a 64-bit fingerprint
type SimHash uint64

// PermutationTable implements locality-sensitive hashing for faster similarity search
type PermutationTable struct {
	permutations [][]int
	bandSize     int
	bands        int
}

// Vectorizer defines an interface for converting text to vectors
type Vectorizer interface {
	TextToVector(text string) []float64
}

// FrequencyVectorizer implements simple frequency-based vectorization
type FrequencyVectorizer struct {
	dimensions int
}

// NewFrequencyVectorizer creates a new frequency vectorizer
func NewFrequencyVectorizer(dimensions int) *FrequencyVectorizer {
	return &FrequencyVectorizer{dimensions: dimensions}
}

// CalculateWithVectorizer computes SimHash using a custom vectorizer
func CalculateWithVectorizer(text string, hyperplanes [][]float64, vectorizer Vectorizer) SimHash {
	vector := vectorizer.TextToVector(text)
	var hash SimHash

	// Compute dot products with hyperplanes
	for i, hyperplane := range hyperplanes {
		dotProduct := 0.0
		for j := range vector {
			dotProduct += vector[j] * hyperplane[j]
		}
		if dotProduct >= 0 {
			hash |= 1 << i
		}
	}

	return hash
}

// Calculate computes SimHash for text using hyperplanes
func Calculate(text string, hyperplanes [][]float64) SimHash {
	vectorizer := NewFrequencyVectorizer(VectorDimensions)
	return CalculateWithVectorizer(text, hyperplanes, vectorizer)
}

// NGramVectorizer implements n-gram based vectorization
type NGramVectorizer struct {
	dimensions int
	ngramSize  int
}

// HammingDistance calculates the number of bit positions where two SimHashes differ
func (s SimHash) HammingDistance(other SimHash) int {
	return bits.OnesCount64(uint64(s ^ other))
}

// IsSimilar determines if two SimHashes are similar based on a threshold
func (s SimHash) IsSimilar(other SimHash, threshold int) bool {
	return s.HammingDistance(other) <= threshold
}
func GenerateHyperplanes(dimensions, count int) [][]float64 {
	// Create a deterministic source for reproducibility
	source := rand.NewSource(42)
	r := rand.New(source)

	hyperplanes := make([][]float64, count)

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Parallel hyperplane generation
	for i := 0; i < count; i += 4 {
		wg.Add(1)
		go func(startIdx int) {
			defer wg.Done()

			localPlanes := make([][]float64, 0, 4)
			for j := 0; j < 4 && startIdx+j < count; j++ {
				hyperplane := make([]float64, dimensions)
				sumSquared := 0.0

				// Generate random values using Box-Muller transform
				for k := range hyperplane {
					u1, u2 := r.Float64(), r.Float64()
					z := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2.0*math.Pi*u2)
					hyperplane[k] = z
					sumSquared += z * z
				}

				// Normalize to unit vector
				magnitude := math.Sqrt(sumSquared)
				for k := range hyperplane {
					hyperplane[k] /= magnitude
				}

				localPlanes = append(localPlanes, hyperplane)
			}

			mu.Lock()
			for j, plane := range localPlanes {
				if startIdx+j < count {
					hyperplanes[startIdx+j] = plane
				}
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	return hyperplanes
}

// NewNGramVectorizer creates a new n-gram vectorizer
func NewNGramVectorizer(dimensions, ngramSize int) *NGramVectorizer {
	return &NGramVectorizer{
		dimensions: dimensions,
		ngramSize:  ngramSize,
	}
}

// TextToVector converts text to a normalized n-gram vector
func (nv *NGramVectorizer) TextToVector(text string) []float64 {
	if len(text) < nv.ngramSize {
		// Handle edge case for very short texts
		return NewFrequencyVectorizer(nv.dimensions).TextToVector(text)
	}

	ngramFreq := make(map[string]int)

	// Generate n-grams
	for i := 0; i <= len(text)-nv.ngramSize; i++ {
		ngram := text[i : i+nv.ngramSize]
		ngramFreq[ngram]++
	}

	vector := make([]float64, nv.dimensions)

	// Distribute frequencies to dimensions using hashing
	for ngram, freq := range ngramFreq {
		hash := md5.Sum([]byte(ngram))
		dim := int(binary.BigEndian.Uint32(hash[:4]) % uint32(nv.dimensions))
		vector[dim] += float64(freq)
	}

	// Normalize vector
	magnitude := 0.0
	for _, v := range vector {
		magnitude += v * v
	}
	magnitude = math.Sqrt(magnitude)

	if magnitude > 0 {
		for i := range vector {
			vector[i] /= magnitude
		}
	}

	return vector
}
func (fv *FrequencyVectorizer) TextToVector(text string) []float64 {
	wordFreq := make(map[string]int)

	// Count word frequencies
	for _, word := range strings.Fields(text) {
		word = strings.ToLower(strings.Trim(word, ".,!?:;\"'()[]{}"))
		if word != "" {
			wordFreq[word]++
		}
	}

	if len(wordFreq) == 0 {
		return make([]float64, fv.dimensions)
	}

	vector := make([]float64, fv.dimensions)

	// Distribute frequencies to dimensions using hashing
	for word, freq := range wordFreq {
		hash := md5.Sum([]byte(word))
		dim := int(binary.BigEndian.Uint32(hash[:4]) % uint32(fv.dimensions))
		vector[dim] += float64(freq)
	}

	// Normalize vector
	magnitude := 0.0
	for _, v := range vector {
		magnitude += v * v
	}
	magnitude = math.Sqrt(magnitude)

	if magnitude > 0 {
		for i := range vector {
			vector[i] /= magnitude
		}
	}

	return vector
}
