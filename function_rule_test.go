package qawaylinter

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"os"
	"path/filepath"
	"testing"
)

func TestFunctionRule(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	testdata := filepath.Join(wd, "testdata")
	plugin := AnalyzerPlugin{Settings: Settings{
		Targets: []Rules{
			{
				Packages: []string{"functions"},
				FunctionRule: &FunctionRule[FunctionRuleResults]{
					Params: FunctionRuleParameters{
						RequireHeadlineComment:  true,
						MinCommentDensity:       0.1,
						TrivialCommentThreshold: 0.3,
						MinLoggingDensity:       0.1,
					},
				},
			},
		},
	}}

	analyzers, err := plugin.BuildAnalyzers()
	analysistest.Run(t, testdata, analyzers[0], "functions")
}
