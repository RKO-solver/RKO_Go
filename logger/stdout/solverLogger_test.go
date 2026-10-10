package stdout

import (
	"io"
	"math"
	"os"
	"path/filepath"
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

func TestResetStartsFreshRun(t *testing.T) {
	l := NewLogger("test", logger.INFO)
	sl := l.GetLogger("solver").(*SolverLogger)
	captureStdout(t, func() {
		sl.Register(5, 5, 0.0, "")
		l.AddSolutionPool(5, 0.0)
	})

	l.Reset()
	require.Empty(t, l.GetSolutionData())
	require.Len(t, l.GetReportData(), 1)
	require.Empty(t, l.GetReportData()[0].Performance)
	require.Equal(t, math.Inf(-1), sl.lastPrint)
	require.Equal(t, math.MaxInt, sl.lastBest)

	out := captureStdout(t, func() {
		sl.Register(9, 9, 0.1, "") // first Register of the run: printed
	})
	require.Equal(t, 1, strings.Count(out, "\tLocal "))
}

func TestSilentPrintsNothing(t *testing.T) {
	sl := NewLogger("test", logger.SILENT).GetLogger("solver").(*SolverLogger)
	out := captureStdout(t, func() {
		sl.Register(1, 1, 0.1, "")
		sl.Verbose("a", 0.1)
	})
	require.Empty(t, out)
	require.Equal(t, math.MaxInt, sl.lastBest)
	require.Len(t, sl.history, 1)
}

// spillLogger returns a SILENT logger whose solver flushes every 4 records into a temp dir, closed at test end.
func spillLogger(t *testing.T) (*Log, *SolverLogger, string) {
	t.Helper()
	dir := t.TempDir()
	l := NewLogger("test", logger.SILENT)
	l.SetHistoryDir(dir)
	t.Cleanup(func() { require.NoError(t, l.Close()) })
	l.SetHistoryBuffer(4 * 24)
	return l, l.GetLogger("solver").(*SolverLogger), dir
}

// historyFiles returns the history files in dir.
func historyFiles(t *testing.T, dir string) []string {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(dir, "rko-history-*.bin"))
	require.NoError(t, err)
	return files
}

// register sends n records (local i, best 100-i, time i/2) and returns them.
func register(sl logger.SolverLogger, n int) []logger.Data {
	want := make([]logger.Data, 0, n)
	for i := range n {
		d := logger.Data{LocalCost: i, BestCost: 100 - i, Time: float64(i) / 2}
		sl.Register(d.LocalCost, d.BestCost, d.Time, "")
		want = append(want, d)
	}
	return want
}

// csvLines returns the lines of the CSV file at path.
func csvLines(t *testing.T, path string) []string {
	t.Helper()
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	return strings.Split(strings.TrimSuffix(string(b), "\n"), "\n")
}

func TestHistorySpillsWhenBufferIsFull(t *testing.T) {
	l, sl, dir := spillLogger(t)
	want := register(sl, 10) // two flushes of 4, 2 records left in memory

	require.Len(t, historyFiles(t, dir), 1)
	require.Len(t, sl.history, 2)
	require.Equal(t, 4, cap(sl.history)) // the buffer is reused, never grown
	require.Equal(t, want, l.GetReportData()[0].Performance)

	out := filepath.Join(t.TempDir(), "out")
	l.SaveCsv(out)
	lines := csvLines(t, out+"-solver-0.csv")
	require.Len(t, lines, 11)
	require.Equal(t, "best,local,time", lines[0])
	require.Equal(t, "100,0,0.00000000", lines[1])
	require.Equal(t, "91,9,4.50000000", lines[10])
}

func TestResetAppliesBufferAndDeletesFiles(t *testing.T) {
	l, sl, dir := spillLogger(t)
	register(sl, 5)
	require.Len(t, historyFiles(t, dir), 1)

	l.SetHistoryBuffer(8 * 24) // after GetLogger: applied by Reset
	l.Reset()
	require.Empty(t, historyFiles(t, dir))
	require.Empty(t, l.GetReportData()[0].Performance)
	require.Equal(t, 8, cap(sl.history))

	register(sl, 9)
	require.Len(t, historyFiles(t, dir), 1)
	require.NoError(t, l.Close())
	require.Empty(t, historyFiles(t, dir))
}

func TestHistoryBufferZeroKeepsEverythingInMemory(t *testing.T) {
	dir := t.TempDir()
	l := NewLogger("test", logger.SILENT)
	l.SetHistoryDir(dir)
	t.Cleanup(func() { require.NoError(t, l.Close()) })
	sl := l.GetLogger("solver").(*SolverLogger)
	require.Equal(t, 512*1024/24, cap(sl.history))

	l.SetHistoryBuffer(0)
	l.Reset()
	want := register(sl, 100)
	require.Empty(t, historyFiles(t, dir))
	require.Len(t, sl.history, 100)
	require.Equal(t, want, l.GetReportData()[0].Performance)
}

func TestCloseDeletesFilesAndIsIdempotent(t *testing.T) {
	l, sl, dir := spillLogger(t)
	var lg logger.Logger = l
	register(sl, 5) // one flush of 4
	require.Len(t, historyFiles(t, dir), 1)

	require.NoError(t, lg.Close())
	require.Empty(t, historyFiles(t, dir))
	require.NoError(t, lg.Close()) // a second Close is a no-op
}

func TestRegisterAfterCloseCreatesNoFile(t *testing.T) {
	l, sl, dir := spillLogger(t)
	var lg logger.Logger = l
	register(sl, 5) // one flush of 4, one record left in memory
	require.NoError(t, lg.Close())

	register(sl, 10) // fills the buffer again: kept in memory
	require.Empty(t, historyFiles(t, dir))
	require.Len(t, sl.history, 11)

	lg.Reset() // the next run spills again
	register(sl, 5)
	require.Len(t, historyFiles(t, dir), 1)
}
