package qawaylinter

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

type StructRuleParameters struct {
	// RequireHeadlineComment determines if a comment must be placed on top of the interface.
	RequireHeadlineComment bool `json:"requireHeadlineComment"`
	// RequireFieldComment determines if a comment must be placed on top of each field in the struct.
	RequireFieldComment bool `json:"requireFieldComment"`
}

type StructRuleResults struct {
	// Number of lines of comments in the headline of the interface.
	HeadlineComments int
	// Number of lines of comments on top of each field. Key: field name, value: number of comment lines
	FieldComments map[string]int
}

type StructRule[ResultType StructRuleResults] struct {
	Params StructRuleParameters `json:"params"`
}

func (i StructRule[ResultType]) IsApplicable(node ast.Node, pass *analysis.Pass, _ *ast.File) bool {
	n, ok := node.(*ast.GenDecl)
	if !ok {
		return false
	}
	spec, ok := n.Specs[0].(*ast.TypeSpec)
	if !ok || spec.Type == nil {
		return false
	}
	if _, ok := spec.Type.(*ast.StructType); !ok {
		return false
	}

	return true
}

func (i StructRule[ResultType]) Analyse(node ast.Node, pass *analysis.Pass, _ *ast.File) *StructRuleResults {
	typespec, ok := node.(*ast.GenDecl)
	if !ok {
		return nil
	}

	typeComments := countHeadlineComments(typespec.Doc, pass.Fset)

	var fieldComments = make(map[string]int)
	ast.Inspect(node, func(n ast.Node) bool {
		if stru, ok := n.(*ast.StructType); ok {
			for _, field := range stru.Fields.List {
				if _, ok := field.Type.(*ast.Ident); ok {
					if len(field.Names) == 0 {
						continue
					}
					fieldComments[field.Names[0].Name] = countHeadlineComments(field.Doc, pass.Fset)
				}
			}
		}
		return true
	})

	return &StructRuleResults{
		HeadlineComments: typeComments,
		FieldComments:    fieldComments,
	}
}

func (i StructRule[ResultType]) Apply(analysis *StructRuleResults, node ast.Node, pass *analysis.Pass) {
	if analysis == nil {
		return
	}
	if analysis.HeadlineComments == 0 && i.Params.RequireHeadlineComment {
		pass.Reportf(node.Pos(), "Struct '%s' is missing required headline comment", node.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name)
	}
	for name, comments := range analysis.FieldComments {
		if comments == 0 && i.Params.RequireFieldComment {
			pass.Reportf(node.Pos(), "Field '%s' is missing required comment", name)
		}
	}
}
