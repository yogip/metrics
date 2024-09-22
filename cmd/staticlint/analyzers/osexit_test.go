package analyzers

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	dir, err := filepath.Abs("./testdata")
	require.NoError(t, err)

	analysistest.Run(t, dir, OsExitAnalyzer, "./...")
}
