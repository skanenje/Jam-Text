package simhash


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
