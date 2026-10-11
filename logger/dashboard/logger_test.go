package dashboard

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/stretchr/testify/require"
)

// run resets l, writes to buf, runs f between Start and stop and returns the output.
func run(l *Log, f func()) string {
	var buf bytes.Buffer
	l.out = &buf
	l.Reset()
	stop := l.Start()
	f()
	stop()
	return buf.String()
}

func TestOptionsDefaultsAndMinimums(t *testing.T) {
	l := New("p", nil)
	require.Equal(t, logger.LevelInfo, l.level)
	require.Equal(t, 500*time.Millisecond, l.refresh)
	require.Equal(t, 1024, cap(l.updateChan))
	require.Zero(t, l.data.numLinesPool) // the pool is shown at Verbose only
	require.Equal(t, 128*1024, l.hist.Bytes)
	require.Equal(t, os.Stdout, l.out)

	l = New("p", &Options{Level: logger.LevelVerbose})
	require.Equal(t, 3, l.data.numLinesPool)

	l = New("p", &Options{Level: logger.LevelVerbose, Refresh: time.Millisecond, BufferSize: -4, PoolLines: -1})
	require.Equal(t, 300*time.Millisecond, l.refresh)
	require.Equal(t, 1, cap(l.updateChan))
	require.Equal(t, 1, l.data.numLinesPool)
}

func TestResetStartStopTwiceKeepsLastRun(t *testing.T) {
	l := New("test", &Options{Level: logger.LevelVerbose})
	sl := l.SolverLogger("solver")

	for iter := range 2 {
		run(l, func() {
			sl.Register(10+iter, 5, 1, 100*time.Millisecond)
			sl.Event(slog.Int("stall", iter), 100*time.Millisecond)
			l.AddSolutionPool(5+iter, 100*time.Millisecond)
			l.WorkerDone(logger.WorkerResult{ID: 0, Name: "solver", Cost: 5})
		})

		report := l.Report()
		require.Len(t, report, 1)
		require.Equal(t, []logger.Sample{{Local: 10 + iter, Best: 5, Elapsed: 100 * time.Millisecond}}, report[0].Samples)
		require.Equal(t, []logger.PoolEntry{{Cost: 5 + iter, Elapsed: 100 * time.Millisecond}}, l.Pool())
		require.Len(t, l.data.workers, 1)
		require.Len(t, l.data.solvers[0].events, 1)
		require.Equal(t, int64(iter), l.data.solvers[0].events[0].attr.Value.Int64())
	}
}

func TestPrintShellEntersAlternateScreenOnce(t *testing.T) {
	l := New("test", nil)
	sl := l.SolverLogger("solver")
	l.Reset()
	sl.Register(10, 9, 3, 100*time.Millisecond)
	l.data.apply(<-l.updateChan)

	var buf bytes.Buffer
	l.data.printShell(&buf)
	l.data.printShell(&buf)
	out := buf.String()

	require.True(t, strings.HasPrefix(out, enterAltScreenCode+cursorHomeCode))
	require.Equal(t, 1, strings.Count(out, enterAltScreenCode))
	require.Equal(t, 2, strings.Count(out, cursorHomeCode))
	require.Equal(t, 2, strings.Count(out, clearToEndCode))
	require.Contains(t, out, "(0) solver: iteration 3"+cleanLineCode+"\n")
	require.Contains(t, out, "    Local 10 Best 9"+cleanLineCode+"\n")
	require.Contains(t, out, "    Time 0.100s"+cleanLineCode+"\n")
}

func TestStartStopDrawsFinalFrameAndDoneLines(t *testing.T) {
	l := New("test", &Options{Refresh: time.Hour})
	sl := l.SolverLogger("solver")
	out := run(l, func() {
		sl.Register(7, 7, 2, 200*time.Millisecond)
		l.WorkerDone(logger.WorkerResult{ID: 1, Name: "B", Cost: 8, Elapsed: 2 * time.Second})
		l.WorkerDone(logger.WorkerResult{ID: 0, Name: "A", Cost: 7, Elapsed: time.Second})
	})

	require.NotContains(t, out, enterAltScreenCode) // no tick: the alternate screen was never entered
	require.NotContains(t, out, leaveAltScreenCode)
	require.Contains(t, out, "    Local 7 Best 7"+cleanLineCode+"\n")
	require.True(t, strings.HasSuffix(out, "(1) B done cost=8 time=2.000s\n(0) A done cost=7 time=1.000s\n"))
}

func TestStopLeavesAlternateScreenWhenEntered(t *testing.T) {
	l := New("test", &Options{Refresh: minimumRefresh})
	sl := l.SolverLogger("solver")
	out := run(l, func() {
		sl.Register(7, 7, 1, time.Second)
		time.Sleep(2 * minimumRefresh)
	})

	require.Equal(t, 1, strings.Count(out, enterAltScreenCode))
	require.Equal(t, 1, strings.Count(out, leaveAltScreenCode))
}

func TestSilentWritesNothing(t *testing.T) {
	l := New("test", &Options{Level: logger.LevelSilent})
	sl := l.SolverLogger("solver")
	out := run(l, func() {
		sl.Register(10, 10, 1, time.Second)
		l.WorkerStarted(0, "solver")
		l.WorkerDone(logger.WorkerResult{})
	})

	require.Empty(t, out)
	require.Len(t, l.Report()[0].Samples, 1)
}

func TestWorkerStartedLine(t *testing.T) {
	var buf bytes.Buffer
	l := New("test", nil)
	l.out = &buf
	l.WorkerStarted(2, "GA")
	require.Equal(t, "(2) GA started\n", buf.String())
}

func TestInfoShowsNoEventsOrPool(t *testing.T) {
	l := New("test", nil)
	sl := l.SolverLogger("solver")
	out := run(l, func() {
		sl.Register(7, 7, 1, time.Second)
		sl.Event(slog.Int("stall", 3), time.Second)
		l.AddSolutionPool(7, time.Second)
	})

	require.NotContains(t, out, "stall=")
	require.NotContains(t, out, "pool cost=")
	require.Empty(t, l.data.solvers[0].events)
}

func TestVerboseShowsLatestEventPerKeyAndLastPoolLines(t *testing.T) {
	l := New("test", &Options{Level: logger.LevelVerbose, PoolLines: 2})
	sl := l.SolverLogger("solver")
	out := run(l, func() {
		sl.Register(7, 7, 1, time.Second)
		sl.Event(slog.Int("a", 1), time.Second)
		sl.Event(slog.Int("b", 2), 2*time.Second)
		sl.Event(slog.Int("a", 3), 3*time.Second)
		for i := range 3 {
			l.AddSolutionPool(10+i, time.Duration(i)*time.Second)
		}
	})

	require.Contains(t, out, "    a=3 time=3.000s")
	require.Contains(t, out, "    b=2 time=2.000s")
	require.NotContains(t, out, "a=1")
	require.NotContains(t, out, "pool cost=10")
	require.Contains(t, out, "pool cost=11 time=1.000s")
	require.Contains(t, out, "pool cost=12 time=2.000s")
}

func TestStopAppliesEveryQueuedMessage(t *testing.T) {
	const n = 500
	l := New("test", &Options{Level: logger.LevelSilent, BufferSize: n})
	sl := l.SolverLogger("solver")
	run(l, func() { register(sl, n) })
	require.Len(t, l.Report()[0].Samples, n)
}

func TestResetStartsFreshRun(t *testing.T) {
	l := New("test", &Options{Level: logger.LevelSilent})
	sl := l.SolverLogger("solver")
	run(l, func() {
		register(sl, 5)
		l.AddSolutionPool(5, 0)
	})

	l.Reset()
	require.Empty(t, l.Report()[0].Samples)
	require.Empty(t, l.Pool())
	require.NotNil(t, l.Pool())
}

func TestLongLinesAreTruncated(t *testing.T) {
	l := New("test", &Options{Refresh: time.Hour})
	sl := l.SolverLogger(strings.Repeat("x", 200))
	out := run(l, func() { sl.Register(1, 1, 1, time.Second) })

	require.Contains(t, out, "..."+cleanLineCode+"\n")
}

// spillLog returns a logger with one solver that flushes every 4 records into a temp dir, closed at test end.
func spillLog(t *testing.T, level logger.Level) (*Log, logger.SolverLogger, string) {
	t.Helper()
	dir := t.TempDir()
	l := New("test", &Options{Level: level, HistoryBuffer: 4 * 24, HistoryDir: dir})
	t.Cleanup(func() { require.NoError(t, l.Close()) })
	return l, l.SolverLogger("solver"), dir
}

// historyFiles returns the history files in dir.
func historyFiles(t *testing.T, dir string) []string {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(dir, "rko-history-*.bin"))
	require.NoError(t, err)
	return files
}

// register sends n records (local i, best 100-i, elapsed i/2 seconds) and returns them.
func register(sl logger.SolverLogger, n int) []logger.Sample {
	want := make([]logger.Sample, 0, n)
	for i := range n {
		s := logger.Sample{Local: i, Best: 100 - i, Elapsed: time.Duration(i) * 500 * time.Millisecond}
		sl.Register(s.Local, s.Best, i, s.Elapsed)
		want = append(want, s)
	}
	return want
}

func TestHistorySpillsWhenBufferIsFull(t *testing.T) {
	l, sl, dir := spillLog(t, logger.LevelSilent)
	var want []logger.Sample
	run(l, func() { want = register(sl, 10) }) // two flushes of 4, 2 records left in memory

	require.Len(t, historyFiles(t, dir), 1)
	require.Len(t, l.data.solvers[0].records, 2)
	require.Equal(t, 4, cap(l.data.solvers[0].records)) // the buffer is reused, never grown
	report := l.Report()
	require.Equal(t, "solver", report[0].Name)
	require.Equal(t, 0, report[0].ID)
	require.Equal(t, want, report[0].Samples)

	l.Reset()
	require.Empty(t, historyFiles(t, dir))
}

func TestFrameShowsLastRecordAfterFlush(t *testing.T) {
	l, sl, _ := spillLog(t, logger.LevelInfo)
	out := run(l, func() {
		for range 3 {
			sl.Register(10, 10, 1, 100*time.Millisecond)
		}
		sl.Register(7, 7, 4, 400*time.Millisecond) // fourth record: fills the buffer and flushes it
	})
	require.Empty(t, l.data.solvers[0].records)

	require.Contains(t, out, "(0) solver: iteration 4")
	require.Contains(t, out, "Local 7 Best 7")
	require.Contains(t, out, "Time 0.400s")
}

func TestCloseDeletesFilesAndIsIdempotent(t *testing.T) {
	l, sl, dir := spillLog(t, logger.LevelSilent)
	run(l, func() { register(sl, 5) }) // one flush of 4
	require.Len(t, historyFiles(t, dir), 1)

	require.NoError(t, l.Close())
	require.Empty(t, historyFiles(t, dir))
	require.NoError(t, l.Close()) // a second Close is a no-op
}

func TestNegativeHistoryBufferKeepsEverythingInMemory(t *testing.T) {
	dir := t.TempDir()
	l := New("test", &Options{Level: logger.LevelSilent, HistoryBuffer: -1, HistoryDir: dir})
	sl := l.SolverLogger("solver")
	var want []logger.Sample
	run(l, func() { want = register(sl, 50) })
	require.Empty(t, historyFiles(t, dir))
	require.Equal(t, want, l.Report()[0].Samples)
}

func TestCSVOutputs(t *testing.T) {
	l, sl, _ := spillLog(t, logger.LevelSilent)
	run(l, func() {
		register(sl, 3)
		l.AddSolutionPool(9, 1250*time.Millisecond)
	})

	prefix := filepath.Join(t.TempDir(), "out")
	require.NoError(t, l.SaveCSV(prefix))

	b, err := os.ReadFile(prefix + "-solver-0.csv")
	require.NoError(t, err)
	require.Equal(t, "best,local,time\n100,0,0.00000000\n99,1,0.50000000\n98,2,1.00000000\n", string(b))

	b, err = os.ReadFile(prefix + "-pool.csv")
	require.NoError(t, err)
	require.Equal(t, "cost,time\n9,1.25000000\n", string(b))
}

func TestSaveCSVDefaultsPrefixToProblem(t *testing.T) {
	tmp := t.TempDir()
	t.Chdir(tmp)
	l := New("inst/a:b", &Options{Level: logger.LevelSilent})
	sl := l.SolverLogger("solver")
	run(l, func() { register(sl, 2) })

	require.NoError(t, l.SaveCSV(""))
	require.FileExists(t, filepath.Join(tmp, "inst_a_b-solver-0.csv"))
	require.FileExists(t, filepath.Join(tmp, "inst_a_b-pool.csv"))
}

func TestSaveCSVReportsEveryError(t *testing.T) {
	l := New("test", &Options{Level: logger.LevelSilent})
	l.SolverLogger("solver")

	err := l.SaveCSV(filepath.Join(t.TempDir(), "missing", "out"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "-solver-0.csv")
	require.Contains(t, err.Error(), "-pool.csv")
}

func TestHotPathDoesNotAllocate(t *testing.T) {
	l := New("test", &Options{Level: logger.LevelSilent, HistoryDir: t.TempDir()})
	t.Cleanup(func() { l.Close() })
	sl := l.SolverLogger("solver")
	stop := l.Start()
	defer stop()
	a := slog.Int("a", 1)
	require.Zero(t, testing.AllocsPerRun(1000, func() { sl.Register(2, 1, 1, time.Second) }))
	require.Zero(t, testing.AllocsPerRun(1000, func() { sl.Event(a, time.Second) }))
}

func BenchmarkRegister(b *testing.B) {
	l := New("test", &Options{Level: logger.LevelSilent, HistoryDir: b.TempDir()})
	b.Cleanup(func() { l.Close() })
	sl := l.SolverLogger("solver")
	stop := l.Start()
	defer stop()
	b.ReportAllocs()
	for b.Loop() {
		sl.Register(2, 1, 1, time.Second)
	}
}
