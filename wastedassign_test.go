package wastedassign_test

import (
	"testing"

	"github.com/sanposhiho/wastedassign/v2"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, wastedassign.Analyzer, "a")
}
