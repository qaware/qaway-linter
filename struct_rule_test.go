package qawaylinter

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"os"
	"path/filepath"
	"testing"
)

func TestStructRule(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	testdata := filepath.Join(wd, "testdata")
	plugin := AnalyzerPlugin{Settings: Settings{
		Targets: []Rules{
			{
				Packages: []string{"struct"},
				StructRule: &StructRule[StructRuleResults]{
					Params: StructRuleParameters{
						RequireHeadlineComment: true,
						RequireFieldComment:    true,
					},
				},
			},
		}},
	}
	analyzers, err := plugin.BuildAnalyzers()
	analysistest.Run(t, testdata, analyzers[0], "struct")
}
