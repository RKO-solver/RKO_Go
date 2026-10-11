package printer

import (
	"bytes"
	"io"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/stretchr/testify/require"
)

func ms(n int) time.Duration { return time.Duration(n) * time.Millisecond }

func newLog(level logger.Level) (*Log, *bytes.Buffer) {
	var out bytes.Buffer
	return New(&out, "test", &Options{Level: level, PrintInterval: time.Second}), &out
}

func TestRegisterPrintsOnImprovementOrInterval(t *testing.T) {
	l, out := newLog(logger.LevelInfo)
	sl := l.SolverLogger("solver")

	sl.Register(10, 10, 1, 0)       // first call: printed
	sl.Register(12, 10, 2, ms(100)) // skipped
	sl.Register(9, 9, 3, ms(200))   // improvement: printed
	sl.Register(11, 9, 4, ms(800))  // 0.6s < 1s: skipped
	sl.Register(11, 9, 5, ms(1300)) // interval elapsed: printed

	require.Equal(t, 3, strings.Count(out.String(), "\n"))
	require.Len(t, l.Report()[0].Samples, 5)
}

func TestRegisterLineFormat(t *testing.T) {
	l, out := newLog(logger.LevelInfo)
	l.SolverLogger("solver").Register(10, 10, 3, 0)
	require.Equal(t, "(0) solver iteration=3 local=10 best=10 time=0.000s\n", out.String())
}

func TestEventPrintsAtVerboseOnly(t *testing.T) {
	l, out := newLog(logger.LevelVerbose)
	l.SolverLogger("solver").Event(slog.Int("stall", 3), ms(10))
	require.Equal(t, "(0) solver stall=3 time=0.010s\n", out.String())

	out.Reset()
	l.SolverLogger("solver").Event(slog.Attr{Key: "improved", Value: slog.GroupValue(slog.Int("iteration", 152), slog.Int("local", 10))}, ms(10))
	require.Equal(t, "(1) solver improved=[iteration=152 local=10] time=0.010s\n", out.String())

	l, out = newLog(logger.LevelInfo)
	l.SolverLogger("solver").Event(slog.Int("stall", 3), ms(10))
	require.Empty(t, out.String())
}

func TestWorkerLines(t *testing.T) {
	for _, level := range []logger.Level{logger.LevelInfo, logger.LevelVerbose} {
		l, out := newLog(level)
		l.WorkerStarted(0, "GA")
		l.WorkerDone(logger.WorkerResult{ID: 0, Name: "GA", Cost: 9, Elapsed: 12345 * time.Millisecond})
		require.Equal(t, "(0) GA started\n(0) GA done cost=9 time=12.345s\n", out.String())
	}

	l, out := newLog(logger.LevelSilent)
	l.WorkerStarted(0, "GA")
	l.WorkerDone(logger.WorkerResult{ID: 0, Name: "GA"})
	require.Empty(t, out.String())
}

func TestPoolLineOnlyAtVerbose(t *testing.T) {
	l, out := newLog(logger.LevelVerbose)
	l.AddSolutionPool(9, 1234*time.Millisecond)
	require.Equal(t, "pool cost=9 time=1.234s\n", out.String())

	for _, level := range []logger.Level{logger.LevelSilent, logger.LevelInfo} {
		l, out = newLog(level)
		l.AddSolutionPool(9, ms(5))
		require.Empty(t, out.String())
		require.Equal(t, []logger.PoolEntry{{Cost: 9, Elapsed: ms(5)}}, l.Pool())
	}
}

func TestSilentPrintsNothing(t *testing.T) {
	l, out := newLog(logger.LevelSilent)
	sl := l.SolverLogger("solver").(*SolverLogger)
	sl.Register(1, 1, 1, ms(100))
	sl.Event(slog.Int("a", 1), ms(100))
	require.Empty(t, out.String())
	require.Equal(t, math.MaxInt, sl.lastBest)
	require.Len(t, sl.history, 1)
}

func TestResetStartsFreshRun(t *testing.T) {
	l, out := newLog(logger.LevelInfo)
	sl := l.SolverLogger("solver").(*SolverLogger)
	sl.Register(5, 5, 1, 0)
	l.AddSolutionPool(5, 0)

	l.Reset()
	require.Empty(t, l.Pool())
	require.NotNil(t, l.Pool())
	require.Len(t, l.Report(), 1)
	require.Empty(t, l.Report()[0].Samples)
	require.Equal(t, time.Duration(0), sl.next)
	require.Equal(t, math.MaxInt, sl.lastBest)

	out.Reset()
	sl.Register(9, 9, 1, ms(100)) // first Register of the run: printed
	require.Equal(t, 1, strings.Count(out.String(), "\n"))
}

func TestOptionsDefaults(t *testing.T) {
	l := New(nil, "test", nil)
	require.Equal(t, logger.LevelInfo, l.level)
	require.Equal(t, 500*time.Millisecond, l.interval)
	require.Equal(t, os.Stdout, l.w)
	sl := l.SolverLogger("solver").(*SolverLogger)
	require.Equal(t, 512*1024/24, cap(sl.history))

	require.Equal(t, 10*time.Millisecond, New(io.Discard, "t", &Options{PrintInterval: 1}).interval)
}

func TestHistoryBufferNegativeKeepsEverythingInMemory(t *testing.T) {
	dir := t.TempDir()
	l := New(io.Discard, "test", &Options{Level: logger.LevelSilent, HistoryBuffer: -1, HistoryDir: dir})
	t.Cleanup(func() { require.NoError(t, l.Close()) })
	sl := l.SolverLogger("solver").(*SolverLogger)
	require.Equal(t, 0, cap(sl.history))

	want := register(sl, 100)
	require.Empty(t, historyFiles(t, dir))
	require.Len(t, sl.history, 100)
	require.Equal(t, want, l.Report()[0].Samples)
}

// spillLogger returns a Silent Log whose solver flushes every 4 records into a temp dir, closed at test end.
func spillLogger(t *testing.T) (*Log, *SolverLogger, string) {
	t.Helper()
	dir := t.TempDir()
	l := New(io.Discard, "test", &Options{Level: logger.LevelSilent, HistoryBuffer: 4 * 24, HistoryDir: dir})
	t.Cleanup(func() { require.NoError(t, l.Close()) })
	return l, l.SolverLogger("solver").(*SolverLogger), dir
}

// historyFiles returns the history files in dir.
func historyFiles(t *testing.T, dir string) []string {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(dir, "rko-history-*.bin"))
	require.NoError(t, err)
	return files
}

// register sends n records (local i, best 100-i, elapsed i*500ms) and returns them.
func register(sl logger.SolverLogger, n int) []logger.Sample {
	want := make([]logger.Sample, 0, n)
	for i := range n {
		d := logger.Sample{Local: i, Best: 100 - i, Elapsed: time.Duration(i) * 500 * time.Millisecond}
		sl.Register(d.Local, d.Best, i, d.Elapsed)
		want = append(want, d)
	}
	return want
}

func TestHistorySpillsWhenBufferIsFull(t *testing.T) {
	l, sl, dir := spillLogger(t)
	want := register(sl, 10) // two flushes of 4, 2 records left in memory

	require.Len(t, historyFiles(t, dir), 1)
	require.Len(t, sl.history, 2)
	require.Equal(t, 4, cap(sl.history)) // the buffer is reused, never grown
	require.Equal(t, want, l.Report()[0].Samples)
}

func TestCloseDeletesFilesAndIsIdempotent(t *testing.T) {
	l, sl, dir := spillLogger(t)
	register(sl, 5) // one flush of 4
	require.Len(t, historyFiles(t, dir), 1)

	require.NoError(t, l.Close())
	require.Empty(t, historyFiles(t, dir))
	require.NoError(t, l.Close()) // a second Close is a no-op
}

func TestRegisterAfterCloseCreatesNoFile(t *testing.T) {
	l, sl, dir := spillLogger(t)
	register(sl, 5) // one flush of 4, one record left in memory
	require.NoError(t, l.Close())

	register(sl, 10) // fills the buffer again: kept in memory
	require.Empty(t, historyFiles(t, dir))
	require.Len(t, sl.history, 11)

	l.Reset() // the next run spills again
	register(sl, 5)
	require.Len(t, historyFiles(t, dir), 1)
}

func TestSaveCSVContent(t *testing.T) {
	l, sl, _ := spillLogger(t)
	register(sl, 10)
	l.AddSolutionPool(9, 1250*time.Millisecond)
	l.AddSolutionPool(7, 2*time.Second)

	prefix := filepath.Join(t.TempDir(), "out")
	require.NoError(t, l.SaveCSV(prefix))

	b, err := os.ReadFile(prefix + "-solver-0.csv")
	require.NoError(t, err)
	lines := strings.Split(strings.TrimSuffix(string(b), "\n"), "\n")
	require.Len(t, lines, 11)
	require.Equal(t, "best,local,time", lines[0])
	require.Equal(t, "100,0,0.00000000", lines[1])
	require.Equal(t, "91,9,4.50000000", lines[10])

	b, err = os.ReadFile(prefix + "-pool.csv")
	require.NoError(t, err)
	require.Equal(t, "cost,time\n9,1.25000000\n7,2.00000000\n", string(b))
}

func TestSaveCSV(t *testing.T) {
	l, _ := newLog(logger.LevelSilent)
	register(l.SolverLogger("solver"), 3)
	l.AddSolutionPool(5, ms(1))

	prefix := filepath.Join(t.TempDir(), "out")
	require.NoError(t, l.SaveCSV(prefix))
	require.FileExists(t, prefix+"-solver-0.csv")
	require.FileExists(t, prefix+"-pool.csv")
}

func TestSaveCSVDefaultPrefixIsProblem(t *testing.T) {
	tmp := t.TempDir()
	t.Chdir(tmp)
	l := New(io.Discard, "inst/a:b", &Options{Level: logger.LevelSilent})
	l.SolverLogger("solver").Register(1, 1, 1, 0)

	require.NoError(t, l.SaveCSV(""))
	require.FileExists(t, filepath.Join(tmp, "inst_a_b-solver-0.csv"))
	require.FileExists(t, filepath.Join(tmp, "inst_a_b-pool.csv"))
}

func TestSaveCSVReportsEveryError(t *testing.T) {
	l, _ := newLog(logger.LevelSilent)
	l.SolverLogger("solver")

	err := l.SaveCSV(filepath.Join(t.TempDir(), "missing", "out"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "-solver-0.csv")
	require.Contains(t, err.Error(), "-pool.csv")
}

type lockedWriter struct {
	mu sync.Mutex
	w  bytes.Buffer
}

func (lw *lockedWriter) Write(p []byte) (int, error) {
	lw.mu.Lock()
	defer lw.mu.Unlock()
	return lw.w.Write(p)
}

func TestConcurrentUse(t *testing.T) {
	var out lockedWriter
	l := New(&out, "test", &Options{Level: logger.LevelVerbose, PrintInterval: time.Millisecond})
	a, b := l.SolverLogger("a"), l.SolverLogger("b")

	var wg sync.WaitGroup
	for _, sl := range []logger.SolverLogger{a, b} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range 200 {
				sl.Register(i, 1000-i, i, time.Duration(i)*time.Millisecond)
				sl.Event(slog.Int("i", i), time.Duration(i)*time.Millisecond)
			}
		}()
	}
	for i := range 100 {
		l.AddSolutionPool(i, time.Duration(i))
	}
	wg.Wait()

	require.Len(t, l.Pool(), 100)
	require.Len(t, l.Report(), 2)
	require.Len(t, l.Report()[0].Samples, 200)
}

func TestHotPathDoesNotAllocate(t *testing.T) {
	l := New(io.Discard, "test", &Options{Level: logger.LevelSilent, HistoryDir: t.TempDir()})
	t.Cleanup(func() { l.Close() })
	sl := l.SolverLogger("solver")
	a := slog.Int("a", 1)
	require.Zero(t, testing.AllocsPerRun(1000, func() { sl.Register(2, 1, 1, time.Second) }))
	require.Zero(t, testing.AllocsPerRun(1000, func() { sl.Event(a, time.Second) }))
}

func BenchmarkRegister(b *testing.B) {
	l := New(io.Discard, "test", &Options{Level: logger.LevelSilent, HistoryDir: b.TempDir()})
	b.Cleanup(func() { l.Close() })
	sl := l.SolverLogger("solver")
	b.ReportAllocs()
	for b.Loop() {
		sl.Register(2, 1, 1, time.Second)
	}
}
