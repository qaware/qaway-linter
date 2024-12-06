package qawaylinter

import (
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
)

// StringSimilarity returns the similarity between two strings.
func StringSimilarity(a string, b string) float64 {
	return strutil.Similarity(a, b, metrics.NewLevenshtein())
}
