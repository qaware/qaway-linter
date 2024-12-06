package qawaylinter

import (
	"github.com/golangci/plugin-module-register/register"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"strings"
)

// AnalyzerPlugin is the entry point for the linter.
// Please see https://golangci-lint.run/plugins/module-plugins/ for instructions on how to integrate custom linters
// into golangci-lint.
// There is also an example linter at https://github.com/golangci/example-plugin-module-linter which was used
// as baseline for this implementation.
type AnalyzerPlugin struct {
	Settings Settings
}

func (a *AnalyzerPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name:     "qawaylinter",
			Doc:      "Checks that a given function has an appropriate amount of documetation.",
			Run:      a.Run,
			Requires: []*analysis.Analyzer{inspect.Analyzer},
		},
	}, nil
}

// Run executes the analysis step of the linter.
// The method iterates over all files and applies all rules to the nodes in the file.
// Please refer to the Rule interface for more information on how to implement rules.
// Do to limitations in Go generics, the rules are split, e.g. into FunctionRules and InterfaceRules.
// It would be better if they would all be part of a single list as they all implement the same interface,
// but this was not possible.
func (a *AnalyzerPlugin) Run(pass *analysis.Pass) (interface{}, error) {
	var file *ast.File
	inspect := func(node ast.Node) bool {
		if node == nil {
			return true
		}

		target := a.Settings.GetMatchingTarget(pass.Pkg)
		if target == nil {
			return true
		}

		if target.FunctionRule != nil && target.FunctionRule.IsApplicable(node, pass, file) {
			results := target.FunctionRule.Analyse(node, pass, file)
			target.FunctionRule.Apply(results, node, pass)
		}

		if target.InterfaceRule != nil && target.InterfaceRule.IsApplicable(node, pass, file) {
			results := target.InterfaceRule.Analyse(node, pass, file)
			target.InterfaceRule.Apply(results, node, pass)
		}

		if target.StructRule != nil && target.StructRule.IsApplicable(node, pass, file) {
			results := target.StructRule.Analyse(node, pass, file)
			target.StructRule.Apply(results, node, pass)
		}

		return true

	}

	for _, f := range pass.Files {
		filename := pass.Fset.Position(f.Pos()).Filename

		// skip all tests fails as documenting them is not as important
		if strings.HasSuffix(filename, "_test.go") {
			continue
		}
		file = f
		ast.Inspect(f, inspect)
	}
	return nil, nil
}

func (a *AnalyzerPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
