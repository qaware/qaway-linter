package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"
	"qaway-linter/pkg/golinters/commentdensity"
)

func main() {
	singlechecker.Main(commentdensity.Analyzer)
}
