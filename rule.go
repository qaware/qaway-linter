package qawaylinter

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

// Rule is an interface that is used to define checks for different types of nodes in the AST.
type Rule[ResultType any] interface {
	// IsApplicable determines if the rule is applicable to the given node.
	IsApplicable(node ast.Node, pass *analysis.Pass, file *ast.File) bool
	// Analyse is used to analyse the node and return the results of the analysis.
	Analyse(node ast.Node, pass *analysis.Pass, file *ast.File) *ResultType
	// Apply is used to compare the analysis results to the input parameters of node.
	// It returns an error if the analysis results do not match the input parameters.
	Apply(analysis *ResultType, node ast.Node, pass *analysis.Pass)
}
