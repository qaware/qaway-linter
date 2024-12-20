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
			b:        "downloadArtifacts downloads the artifacts",
			expected: 0.41463414634146345,
		},
		{
			name:     "Similar strings",
			a:        "getValue",
			b:        "getValue returns the value",
			expected: 0.3076923076923077,
		},
		{
			name:     "Different strings",
			a:        "Calculate",
			b:        "Calculate uses a special algorithm to determine the business metric by taking into account the current state of the system",
			expected: 0.07377049180327866,
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
