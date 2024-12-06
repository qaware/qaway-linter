package qawaylinter

import (
	"go/ast"
	"go/scanner"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"log"
	"os"
	"regexp"
	"strings"
)

// patterns for determining logger calls. the (?i) in the regex makes the regex case-insensitive.
var loggerPattern = regexp.MustCompile("(?i)(log|logger)")

// the method pattern also covers calls like Printf etc. as (?i)print also matches Printf.
var loggerMethodPattern = regexp.MustCompile("(?i)(debug|info|warn|error|fatal|print|panic|trace|log)")

type FunctionFilters struct {
	// MinLinesOfCode determines the minimum number of lines of code that a function must have to be considered.
	MinLinesOfCode int `json:"minLinesOfCode"`
}

type FunctionRuleParameters struct {
	// RequireHeadlineComment determines if a comment must be placed on top of the function.
	RequireHeadlineComment bool `json:"requireHeadlineComment"`
	// MinHeadlineCommentDensity determines the minimum percentage of comments in the headline of the function compared to the body length.
	MinHeadlineCommentDensity float64 `json:"minHeadlineCommentDensity"`
	// MinCommentDensity determines the minimum percentage of comments in the body of the function compared to the body length.
	MinCommentDensity       float64 `json:"minCommentDensity"`
	TrivialCommentThreshold float64 `json:"trivialCommentThreshold"`
	MinLoggingDensity       float64 `json:"minLoggingDensity"`
}

type FunctionRuleResults struct {
	// Number of lines of comments in the headline of the function.
	HeadlineComments int
	// Number of lines of code in the body of the function (excludes comments).
	BodyLinesOfCode int
	// Number of lines of comments in the body of the function.
	BodyComments int
	// Indicates the similarity between the method name and the headline comments.
	CommentSimilarity float64
	// Number of logging statements in the function.
	LoggingStatements int
}

type FunctionRule[ResultType FunctionRuleResults] struct {
	Filters FunctionFilters        `json:"filters"`
	Params  FunctionRuleParameters `json:"params"`

	analysisResults *FunctionRuleResults
}

func (f FunctionRule[ResultType]) IsApplicable(node ast.Node, pass *analysis.Pass, file *ast.File) bool {
	if _, ok := node.(*ast.FuncDecl); !ok {
		return false
	}

	// an analysis must be done in this step already as the lines of code are
	// relevant for the target filter `minLinesOfCode`.
	// the result is cached and the analysis is not executed twice.
	f.analysisResults = f.Analyse(node, pass, file)

	if f.analysisResults.BodyLinesOfCode < f.Filters.MinLinesOfCode {
		return false
	}

	return true
}

func (f FunctionRule[ResultType]) Analyse(node ast.Node, pass *analysis.Pass, file *ast.File) *FunctionRuleResults {
	if f.analysisResults != nil {
		// return cached results determined in IsApplicable method.
		return f.analysisResults
	}

	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return nil
	}

	linesInFunction := countLinesInFunction(funcDecl, pass.Fset)
	linesOfCommentsInMethodBody := countInlineCommentsInFunction(funcDecl, file.Comments, pass.Fset)
	loggingStatements := countLoggingStatementsInFunction(funcDecl)

	linesOfHeadlineComments := 0
	if funcDecl.Doc != nil {
		linesOfHeadlineComments = countCommentLines(funcDecl.Doc, pass.Fset)
	}

	commentSimilarity := StringSimilarity(funcDecl.Name.Name, funcDecl.Doc.Text())

	return &FunctionRuleResults{
		HeadlineComments:  linesOfHeadlineComments,
		BodyLinesOfCode:   linesInFunction,
		BodyComments:      linesOfCommentsInMethodBody,
		CommentSimilarity: commentSimilarity,
		LoggingStatements: loggingStatements,
	}
}

func (f FunctionRule[ResultType]) Apply(analysis *FunctionRuleResults, node ast.Node, pass *analysis.Pass) {
	if analysis == nil {
		return
	}
	funcDecl := node.(*ast.FuncDecl)
	if analysis.HeadlineComments == 0 && f.Params.RequireHeadlineComment {
		pass.Reportf(node.Pos(), "Method '%s' is missing required headline comment", funcDecl.Name.Name)
	}
	if analysis.CommentDensity() < f.Params.MinCommentDensity {
		pass.Reportf(node.Pos(), "Method '%s' has less than %.0f%% comment density. Actual: %.0f%%", funcDecl.Name.Name, f.Params.MinCommentDensity*100, analysis.CommentDensity()*100)
	}
	if analysis.HeadlineCommentDensity() < f.Params.MinHeadlineCommentDensity {
		pass.Reportf(node.Pos(), "Method '%s' has less than %.0f%% headline comment density. Actual: %.0f%%", funcDecl.Name.Name, f.Params.MinHeadlineCommentDensity*100, analysis.HeadlineCommentDensity()*100)
	}
	if f.Params.TrivialCommentThreshold > 0 && analysis.CommentSimilarity > f.Params.TrivialCommentThreshold {
		pass.Reportf(node.Pos(), "Method '%s' has a trivial comment. Similarity to method name: %.0f%%", funcDecl.Name.Name, analysis.CommentSimilarity*100)
	}
	if f.Params.MinLoggingDensity > 0 && analysis.LoggingDensity() < f.Params.MinLoggingDensity {
		pass.Reportf(node.Pos(), "Method '%s' has less than %.0f%% logging density. Actual: %.0f%%", funcDecl.Name.Name, f.Params.MinLoggingDensity*100, analysis.LoggingDensity()*100)
	}
}

func (r FunctionRuleResults) CommentDensity() float64 {
	if r.BodyLinesOfCode == 0 {
		return 0
	}
	return (float64(r.BodyComments) + float64(r.HeadlineComments)) / float64(r.BodyLinesOfCode)
}

func (r FunctionRuleResults) HeadlineCommentDensity() float64 {
	if r.BodyLinesOfCode == 0 {
		return 0
	}
	return float64(r.HeadlineComments) / float64(r.BodyLinesOfCode)
}

func (r FunctionRuleResults) LoggingDensity() float64 {
	if r.BodyLinesOfCode == 0 {
		return 0
	}
	return float64(r.LoggingStatements) / float64(r.BodyLinesOfCode)
}

// countInlineCommentsInFunction determines the number of lines of comments that are part of the method body.
// These comments are not returned as part of the AST of a FuncDecl.
// But all comments within a given file are available in the file's comments.
// This function determines the number of lines of comment within a method body by checking the comments in the file.
func countInlineCommentsInFunction(f *ast.FuncDecl, commentsInFile []*ast.CommentGroup, fset *token.FileSet) int {
	commentLines := 0
	for _, comment := range commentsInFile {
		if (comment.Pos() >= f.Pos()) && (comment.End() <= f.End()) {
			commentLines += countCommentLines(comment, fset)
		}
	}
	return commentLines
}

func countLoggingStatementsInFunction(f *ast.FuncDecl) int {
	loggingStatements := 0
	ast.Inspect(f, func(n ast.Node) bool {
		if isLoggingStatement(n) {
			loggingStatements++
		}
		return true
	})
	return loggingStatements
}

// isLoggingStatements determines whether a given node contains a logging statement.
// It first ensures that the node is of the correct types.
func isLoggingStatement(n ast.Node) bool {
	callExpr, ok := n.(*ast.CallExpr)
	if !ok {
		return false
	}

	// selExpr.Sel.Name contains the method that is called on the logger, e.g. `Warnf`
	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	// selExpr.X.Name contains struct that is called (e.g.) `log`
	x, ok := selExpr.X.(*ast.Ident)
	if !ok {
		return false
	}

	return loggerPattern.MatchString(x.Name) && loggerMethodPattern.MatchString(selExpr.Sel.Name)
}

// countLinesInFunction counts the lines between the start and end of a given function declaration.
// Note that this method also includes comments.
func countLinesInFunction(funcDecl *ast.FuncDecl, fset *token.FileSet) int {
	// Get the start and end positions
	start := fset.Position(funcDecl.Body.Pos())
	end := fset.Position(funcDecl.Body.End())

	// Extract the function body as source code
	sourceCode := extractSource(fset.File(funcDecl.Pos()).Name(), start.Offset, end.Offset)

	// Tokenize the source code
	return countMeaningfulLines(sourceCode)
}

// Extract source code from the file for a given offset range
func extractSource(filePath string, start, end int) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	return string(data[start:end])
}

// Count meaningful lines in the given source code
func countMeaningfulLines(source string) int {
	var loc int
	var s scanner.Scanner
	fset := token.NewFileSet()
	s.Init(fset.AddFile("", fset.Base(), len(source)), []byte(source), nil, scanner.ScanComments)

	linesEncountered := make(map[int]bool)

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		position := fset.Position(pos)

		// Skip comments and empty lines
		if tok == token.COMMENT || strings.TrimSpace(lit) == "" {
			continue
		}

		// Ensure each line is only counted once
		if !linesEncountered[position.Line] {
			linesEncountered[position.Line] = true
			loc++
		}
	}
	return loc
}

// countCommentLines counts the lines covered by comments in a given AST node.
// This method takes into account that a command can span multiple lines using the /* */ syntax.
// It can count the number of comments in both the headline and within a method's body.
func countCommentLines(node ast.Node, fset *token.FileSet) int {
	commentLines := 0

	// Traverse the node to find all comment groups
	ast.Inspect(node, func(n ast.Node) bool {
		if commentGroup, ok := n.(*ast.CommentGroup); ok {
			for _, comment := range commentGroup.List {
				if strings.HasPrefix(comment.Text, "// want `") {
					// filter comments that are used for testing as these comments for analysistest
					// would otherwise increase the logging density.
					continue
				}
				start := fset.Position(comment.Pos()).Line
				end := fset.Position(comment.End()).Line
				commentLines += end - start + 1
			}
		}
		return true
	})

	return commentLines
}
