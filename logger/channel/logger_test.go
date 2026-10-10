package channel

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/stretchr/testify/require"
)

func TestResetStartShutdownTwiceKeepsLastRun(t *testing.T) {
	l := NewLoggerLevel("test", logger.VERBOSE)
	sl := l.GetLogger("solver")

	for run := range 2 {
		var wg sync.WaitGroup
		l.Reset()
		l.Start(&wg)
		sl.Register(10+run, 5, 0.1, "")
		sl.Verbose("event", 0.1)
		l.AddSolutionPool(5+run, 0.1)
		l.WorkerDone("done")
		l.Shutdown()
		wg.Wait()

		report := l.GetReportData()
		require.Len(t, report, 1)
		require.Len(t, report[0].Performance, 1)
		require.Equal(t, 10+run, report[0].Performance[0].LocalCost)
		require.Equal(t, []logger.SolutionData{{Cost: 5 + run, Time: 0.1}}, l.GetSolutionData())
		require.Len(t, l.data.extraMessages[0], 1)
		require.Len(t, l.data.workerMessages, 1)
	}
}

func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	require.NoError(t, err)
	out := make(chan string)
	go func() {
		b, _ := io.ReadAll(r)
		out <- string(b)
	}()
	stdout := os.Stdout
	os.Stdout = w
	f()
	os.Stdout = stdout
	require.NoError(t, w.Close())
	return <-out
}

func TestDashboardDrawsInAlternateScreenAndPrintsFinalFrame(t *testing.T) {
	l := NewLoggerLevel("test", logger.INFO)
	sl := l.GetLogger("solver")

	for range 2 { // the second run must enter the alternate screen again
		var wg sync.WaitGroup
		l.Reset()
		l.Start(&wg)
		out := captureStdout(t, func() {
			sl.Register(10, 10, 0.1, "")
			l.Print()
			l.Print()
			sl.Register(7, 7, 0.2, "") // after the last frame: only the final frame shows it
			l.Shutdown()
			wg.Wait()
			l.Finish()
		})

		require.True(t, strings.HasPrefix(out, enterAltScreenCode+cursorHomeCode))
		require.Equal(t, 1, strings.Count(out, enterAltScreenCode))
		require.Equal(t, 2, strings.Count(out, cursorHomeCode))
		require.Equal(t, 2, strings.Count(out, clearToEndCode))
		require.Equal(t, 1, strings.Count(out, leaveAltScreenCode))

		alternate, final, found := strings.Cut(out, leaveAltScreenCode)
		require.True(t, found)
		require.True(t, strings.HasSuffix(alternate, clearToEndCode))
		require.NotContains(t, alternate, "Local 7 Best 7")
		require.Contains(t, final, "Local 7 Best 7")
		require.NotContains(t, final, cursorHomeCode)
		require.NotContains(t, final, clearToEndCode)
	}
}

func TestFinishWithoutPrintSkipsAlternateScreen(t *testing.T) {
	l := NewLoggerLevel("test", logger.INFO)
	sl := l.GetLogger("solver")
	var wg sync.WaitGroup
	l.Reset()
	l.Start(&wg)

	out := captureStdout(t, func() {
		sl.Register(7, 7, 0.2, "")
		l.Shutdown()
		wg.Wait()
		l.Finish()
	})

	require.NotContains(t, out, enterAltScreenCode)
	require.NotContains(t, out, leaveAltScreenCode)
	require.NotContains(t, out, cursorHomeCode)
	require.NotContains(t, out, clearToEndCode)
	require.Contains(t, out, "Local 7 Best 7")
}

func TestSilentRunWritesNothing(t *testing.T) {
	l := NewLoggerLevel("test", logger.SILENT)
	sl := l.GetLogger("solver")
	var wg sync.WaitGroup

	// Same calls as Solve's SILENT path, which never calls Print or Finish.
	out := captureStdout(t, func() {
		l.Reset()
		l.Start(&wg)
		sl.Register(10, 10, 0.1, "")
		l.Shutdown()
		wg.Wait()
		l.CleanScreen()
	})

	require.Empty(t, out)
}

// spillLog returns a logger with one solver that flushes every 4 records into a temp dir, closed at test end.
func spillLog(t *testing.T, level logger.Level) (*Log, logger.SolverLogger, string) {
	t.Helper()
	dir := t.TempDir()
	l := NewLoggerLevel("test", level)
	l.SetHistoryDir(dir)
	t.Cleanup(func() { require.NoError(t, l.Close()) })
	l.SetHistoryBuffer(4 * 24)
	return l, l.GetLogger("solver"), dir
}

// runOnce does one Reset, Start, f, Shutdown cycle, as Solve does.
func runOnce(l *Log, f func()) {
	var wg sync.WaitGroup
	l.Reset()
	l.Start(&wg)
	f()
	l.Shutdown()
	wg.Wait()
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
	l, sl, dir := spillLog(t, logger.SILENT)
	var want []logger.Data
	runOnce(l, func() { want = register(sl, 10) }) // two flushes of 4, 2 records left in memory

	require.Len(t, historyFiles(t, dir), 1)
	require.Len(t, l.data.solvers[0].records, 2)
	require.Equal(t, 4, cap(l.data.solvers[0].records)) // the buffer is reused, never grown
	report := l.GetReportData()
	require.Equal(t, "solver", report[0].Name)
	require.Equal(t, 0, report[0].Id)
	require.Equal(t, want, report[0].Performance)

	out := filepath.Join(t.TempDir(), "out")
	l.SaveCsv(out)
	lines := csvLines(t, out+"-solver-0.csv")
	require.Len(t, lines, 11)
	require.Equal(t, "best,local,time", lines[0])
	require.Equal(t, "100,0,0.00000000", lines[1])
	require.Equal(t, "91,9,4.50000000", lines[10])

	l.Reset()
	require.Empty(t, historyFiles(t, dir))
}

func TestDashboardShowsLastRecordAfterFlush(t *testing.T) {
	l, sl, _ := spillLog(t, logger.INFO)
	runOnce(l, func() {
		for range 3 {
			sl.Register(10, 10, 0.1, "")
		}
		sl.Register(7, 7, 0.4, "last") // fourth record: fills the buffer and flushes it
	})
	require.Empty(t, l.data.solvers[0].records)

	out := captureStdout(t, l.Finish)
	require.Contains(t, out, "(0) solver: last")
	require.Contains(t, out, "Local 7 Best 7")
	require.Contains(t, out, "Time 0.400s")
}
