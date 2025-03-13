package simhash

import (
	"testing"
)

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
func TestHammingDistance(t *testing.T) {
	tests := []struct {
		a, b     SimHash
		expected int
	}{
		{0x0000, 0x0000, 0},
		{0xFFFF, 0x0000, 16},
		{0xFF00, 0x0F00, 4},
	}

	for _, tt := range tests {
		if got := tt.a.HammingDistance(tt.b); got != tt.expected {
			t.Errorf("HammmingDistance(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.expected)
		}
	}
}

func TestGenerateHyperplanes(t *testing.T) {
	dims, count := 128, 64
	planes := GenerateHyperplanes(dims, count)

	if len(planes) != count {
		t.Errorf("Expected %d planes, got %d", count, len(planes))
	}

	// Test plane normalization
	for i, plane := range planes {
		if len(plane) != dims {
			t.Errorf("Plane %d: expected %d dimensions, got %d", i, dims, len(plane))
		}

		// Verifying unit vector
		sumSquared := 0.0
		for _, v := range plane {
			sumSquared += v * v
		}
		if abs(sumSquared-1.0) > 1e-10 {
			t.Errorf("Plane %d not normalized: magnitude = %f", i, sumSquared)
		}

	}
}
