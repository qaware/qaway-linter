package qawaylinter

import (
	"testing"
)

func TestStringSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected float64
	}{
		{
			name:     "Identical strings",
			a:        "downloadArtifacts",
			b:        "downloads the artifacts",
			expected: 0.6956521739130435,
		},
		{
			name:     "Similar strings",
			a:        "getValue",
			b:        "returns the value",
			expected: 0.3529411764705882,
		},
		{
			name:     "Different strings",
			a:        "Calculate",
			b:        "Calculate uses a special algorithm to determine the business metric",
			expected: 0.13432835820895528,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			similarity := StringSimilarity(tt.a, tt.b)
			if similarity != tt.expected {
				t.Errorf("StringSimilarity(%q, %q) = %v; want %v", tt.a, tt.b, similarity, tt.expected)
			}
		})
	}
}
