package stdout

import (
	"io"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/stretchr/testify/require"
)

func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	require.NoError(t, err)
	stdout := os.Stdout
	os.Stdout = w
	f()
	os.Stdout = stdout
	require.NoError(t, w.Close())
	out, err := io.ReadAll(r)
	require.NoError(t, err)
	return string(out)
}

func TestRegisterPrintsOnImprovementOrInterval(t *testing.T) {
	l := NewLogger("test", logger.INFO)
	sl := l.GetLogger("solver")
	l.SetPrintInterval(1000) // after GetLogger: must reach the existing solver

	out := captureStdout(t, func() {
		sl.Register(10, 10, 0.0, "") // first call: printed
		sl.Register(12, 10, 0.1, "") // skipped
		sl.Register(9, 9, 0.2, "")   // improvement: printed
		sl.Register(11, 9, 0.8, "")  // 0.6s < 1s: skipped
		sl.Register(11, 9, 1.3, "")  // interval elapsed: printed
	})
	require.Equal(t, 3, strings.Count(out, "\tLocal "))
	require.Len(t, l.GetReportData()[0].Performance, 5)
}

func TestVerboseAlwaysPrints(t *testing.T) {
	sl := NewLogger("test", logger.VERBOSE).GetLogger("solver")
	out := captureStdout(t, func() {
		sl.Verbose("a", 0.01)
		sl.Verbose("b", 0.02)
		sl.Verbose("c", 0.03)
	})
	require.Equal(t, 3, strings.Count(out, "Verbose:"))
}

func TestSilentPrintsNothing(t *testing.T) {
	sl := NewLogger("test", logger.SILENT).GetLogger("solver").(*SolverLogger)
	out := captureStdout(t, func() {
		sl.Register(1, 1, 0.1, "")
		sl.Verbose("a", 0.1)
	})
	require.Empty(t, out)
	require.Equal(t, math.MaxInt, sl.lastBest)
	require.Len(t, sl.solver, 1)
}
