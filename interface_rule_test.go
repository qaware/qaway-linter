package qawaylinter

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"os"
	"path/filepath"
	"testing"
)

func TestInterfaceRule(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	testdata := filepath.Join(wd, "testdata")
	plugin := AnalyzerPlugin{Settings: Settings{
		Targets: []Rules{
			{
				Packages: []string{"interfaces"},
				InterfaceRule: &InterfaceRule[InterfaceRuleResults]{
					Params: InterfaceRuleParameters{
						RequireHeadlineComment: true,
						RequireMethodComment:   true,
					},
				},
			},
		},
	}}
	analyzers, err := plugin.BuildAnalyzers()
	analysistest.Run(t, testdata, analyzers[0], "interfaces")
}
