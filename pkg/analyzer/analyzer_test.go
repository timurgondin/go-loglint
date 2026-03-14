package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/timurgondin/go-loglint/pkg/analyzer"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "basic")
}

func TestLowercase(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "lowercase")
}

func TestEnglish(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "english")
}

func TestSpecialChars(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "specialchars")
}

func TestSensitive(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "sensitive")
}

func TestZap(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "zap")
}

func TestSlogContext(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "slogctx")
}
