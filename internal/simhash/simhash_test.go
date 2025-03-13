package simhash

import (
	"testing"
)

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
