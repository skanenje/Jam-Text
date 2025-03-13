package simhash

import (
	"math"
	"math/bits"
	"math/rand"
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
