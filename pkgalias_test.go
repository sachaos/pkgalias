package pkgalias

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"testing"

	"github.com/gostaticanalysis/testutil"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	var cnf = &config{
		settings: []*expectedPackage{
			{
				alias:    "m",
				fullpath: "math",
			},
			{
				alias:    "format",
				fullpath: "fmt",
			},
			{
				alias:    "opesys",
				fullpath: "os",
			},
			{
				alias:    "",
				fullpath: "net",
			},
		},
	}


	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analyzer := &analysis.Analyzer{
		Run: runWithConfig(cnf),
		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
		},
	}
	analysistest.Run(t, testdata, analyzer, "a")
}
