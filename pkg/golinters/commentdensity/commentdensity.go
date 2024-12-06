package commentdensity

import (
	"flag"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

//nolint:gochecknoglobals
var flagSet flag.FlagSet

//nolint:gochecknoglobals
var (
	minCommentDensity int
	minLinesOfCode    int
)

const (
	defaultMinCommentDensity = 10
	defaultMinLinesOfCode    = 10
)

//nolint:gochecknoinits
func init() {
	flagSet.IntVar(&minCommentDensity, "minCommentDensity", defaultMinCommentDensity, "percentage of comments required for functions in relation to the method length")
	flagSet.IntVar(&minLinesOfCode, "minLinesOfCode", defaultMinLinesOfCode, "minimum lines of codes for methods to require a comment")
}

var Analyzer = &analysis.Analyzer{
	Name:     "gocommentdensity",
	Doc:      "Checks that a given function has an appropriate amount of comments.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Flags:    flagSet,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := func(node ast.Node) bool {
		funcDecl, ok := node.(*ast.FuncDecl)
		if !ok {
			return true
		}

		linesOfCode := pass.Fset.Position(funcDecl.End()).Line - pass.Fset.Position(funcDecl.Pos()).Line
		if linesOfCode < 10 {
			// no comments necessary for short functions
			return true
		}

		linesOfComment := 0
		if funcDecl.Doc == nil {
			pass.Reportf(node.Pos(), "function '%s' should have a comment explaining what it does", funcDecl.Name.Name)
			return true
		}

		linesOfComment = len(funcDecl.Doc.List)

		// Inspect the function body for comments
		linesOfComment += countCommentLines(funcDecl.Body, pass.Fset)
		linesOfCode -= countCommentLines(funcDecl.Body, pass.Fset)

		// Traverse the function body to count inline comments
		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			if _, ok := n.(*ast.Comment); ok {
				linesOfComment++
			}
			return true
		})

		if linesOfComment/linesOfCode*100 < minCommentDensity/100 {
			pass.Reportf(node.Pos(), "function '%s' should not have enough coverage. Lines of code: %d, expected lines of comment: %d", funcDecl.Name.Name, linesOfCode, linesOfComment)
		}

		return true
	}

	for _, f := range pass.Files {
		ast.Inspect(f, inspect)
	}
	return nil, nil
}

// countCommentLines counts the lines covered by comments in a given AST node.
func countCommentLines(node ast.Node, fset *token.FileSet) int {
	commentLines := 0

	// Traverse the node to find all comment groups
	ast.Inspect(node, func(n ast.Node) bool {
		if commentGroup, ok := n.(*ast.CommentGroup); ok {
			for _, comment := range commentGroup.List {
				start := fset.Position(comment.Pos()).Line
				end := fset.Position(comment.End()).Line
				commentLines += (end - start + 1)
			}
		}
		return true
	})

	return commentLines
}
