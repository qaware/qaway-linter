package qawaylinter

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
)

type InterfaceRuleParameters struct {
	// RequireHeadlineComment determines if a comment must be placed on top of the interface.
	RequireHeadlineComment bool `json:"requireHeadlineComment"`
	// RequireMethodComment determines if a comment must be placed on top of each method in the interface.
	RequireMethodComment bool `json:"requireMethodComment"`
}

type InterfaceRuleResults struct {
	// Number of lines of comments in the headline of the interface.
	HeadlineComments int
	// Number of lines of comments in the body of the function. Key: function name, value: number of comment lines
	FunctionComments map[string]int
}

type InterfaceRule[ResultType InterfaceRuleResults] struct {
	Params InterfaceRuleParameters `json:"params"`
}

func (i InterfaceRule[ResultType]) IsApplicable(node ast.Node, pass *analysis.Pass, _ *ast.File) bool {
	n, ok := node.(*ast.GenDecl)
	if !ok {
		return false
	}
	spec, ok := n.Specs[0].(*ast.TypeSpec)
	if !ok || spec.Type == nil {
		return false
	}
	if _, ok := spec.Type.(*ast.InterfaceType); !ok {
		return false
	}

	return true
}

func (i InterfaceRule[ResultType]) Analyse(node ast.Node, pass *analysis.Pass, _ *ast.File) *InterfaceRuleResults {
	typespec, ok := node.(*ast.GenDecl)
	if !ok {
		return nil
	}

	typeComments := countHeadlineComments(typespec.Doc, pass.Fset)

	var methodComments = make(map[string]int)
	ast.Inspect(node, func(n ast.Node) bool {
		if field, ok := n.(*ast.Field); ok {
			if _, ok := field.Type.(*ast.FuncType); ok {
				methodComments[field.Names[0].Name] = countHeadlineComments(field.Doc, pass.Fset)
			}
		}
		return true
	})

	return &InterfaceRuleResults{
		HeadlineComments: typeComments,
		FunctionComments: methodComments,
	}
}

func countHeadlineComments(comments *ast.CommentGroup, fset *token.FileSet) int {
	if comments == nil {
		return 0
	}
	start := fset.Position(comments.Pos()).Line
	end := fset.Position(comments.End()).Line
	return end - start + 1
}

func (i InterfaceRule[ResultType]) Apply(analysis *InterfaceRuleResults, node ast.Node, pass *analysis.Pass) {
	if analysis == nil {
		return
	}
	if analysis.HeadlineComments == 0 && i.Params.RequireHeadlineComment {
		pass.Reportf(node.Pos(), "Interface '%s' is missing required headline comment", node.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name)
	}
	for name, comments := range analysis.FunctionComments {
		if comments == 0 && i.Params.RequireMethodComment {
			pass.Reportf(node.Pos(), "Method '%s' is missing required comment", name)
		}
	}
}
