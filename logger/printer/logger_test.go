package printer

import (
	"bytes"
	"io"
	"log/slog"
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
	return New(&out, "test", &Options{Level: level, PrintInterval: 1000}), &out
}

func TestRegisterPrintsEveryIntervalPlusOne(t *testing.T) {
	var out bytes.Buffer
	l := New(&out, "test", &Options{Level: logger.LevelInfo, PrintInterval: 2})
	ml := l.MetaheuristicLogger("metaheuristic")

	for i := 1; i <= 7; i++ {
		ml.Register(i, i, i, 0) // printed on calls 1, 4 and 7
	}

	require.Equal(t, 3, strings.Count(out.String(), "\n"))
	require.Contains(t, out.String(), "iteration=4 ")
	require.Len(t, l.Report()[0].Samples, 7)
}

func TestRegisterLineFormat(t *testing.T) {
	l, out := newLog(logger.LevelInfo)
	l.MetaheuristicLogger("metaheuristic").Register(10, 10, 3, 0)
	require.Equal(t, "(0) metaheuristic iteration=3 local=10 best=10 time=0.000s\n", out.String())
}

func TestEventPrintsAtVerboseOnly(t *testing.T) {
	l, out := newLog(logger.LevelVerbose)
	l.MetaheuristicLogger("metaheuristic").Event(slog.Int("stall", 3), ms(10))
	require.Equal(t, "(0) metaheuristic stall=3 time=0.010s\n", out.String())

	out.Reset()
	l.MetaheuristicLogger("metaheuristic").Event(slog.Attr{Key: "improved", Value: slog.GroupValue(slog.Int("iteration", 152), slog.Int("local", 10))}, ms(10))
	require.Equal(t, "(1) metaheuristic improved=[iteration=152 local=10] time=0.010s\n", out.String())

	l, out = newLog(logger.LevelInfo)
	l.MetaheuristicLogger("metaheuristic").Event(slog.Int("stall", 3), ms(10))
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
	ml := l.MetaheuristicLogger("metaheuristic").(*MetaheuristicLogger)
	ml.Register(1, 1, 1, ms(100))
	ml.Event(slog.Int("a", 1), ms(100))
	require.Empty(t, out.String())
	require.Equal(t, -1, ml.count) // never armed by a print
	require.Len(t, ml.history, 1)
}

func TestResetStartsFreshRun(t *testing.T) {
	l, out := newLog(logger.LevelInfo)
	ml := l.MetaheuristicLogger("metaheuristic").(*MetaheuristicLogger)
	ml.Register(5, 5, 1, 0)
	l.AddSolutionPool(5, 0)

	l.Reset()
	require.Empty(t, l.Pool())
	require.NotNil(t, l.Pool())
	require.Len(t, l.Report(), 1)
	require.Empty(t, l.Report()[0].Samples)
	require.Zero(t, ml.count)

	out.Reset()
	ml.Register(9, 9, 1, ms(100)) // first Register of the run: printed
	require.Equal(t, 1, strings.Count(out.String(), "\n"))
}

func TestOptionsDefaults(t *testing.T) {
	l := New(nil, "test", nil)
	require.Equal(t, logger.LevelInfo, l.level)
	require.Equal(t, defaultPrintInterval, l.interval)
	require.Equal(t, os.Stdout, l.w)
	ml := l.MetaheuristicLogger("metaheuristic").(*MetaheuristicLogger)
	require.Equal(t, 512*1024/24, cap(ml.history))

	require.Zero(t, New(io.Discard, "t", &Options{PrintInterval: -3}).interval)
}

func TestHistoryBufferNegativeKeepsEverythingInMemory(t *testing.T) {
	dir := t.TempDir()
	l := New(io.Discard, "test", &Options{Level: logger.LevelSilent, HistoryBuffer: -1, HistoryDir: dir})
	t.Cleanup(func() { require.NoError(t, l.Close()) })
	ml := l.MetaheuristicLogger("metaheuristic").(*MetaheuristicLogger)
	require.Equal(t, 0, cap(ml.history))

	want := register(ml, 100)
	require.Empty(t, historyFiles(t, dir))
	require.Len(t, ml.history, 100)
	require.Equal(t, want, l.Report()[0].Samples)
}

// spillLogger returns a Silent Log whose metaheuristic flushes every 4 records into a temp dir, closed at test end.
func spillLogger(t *testing.T) (*Log, *MetaheuristicLogger, string) {
	t.Helper()
	dir := t.TempDir()
	l := New(io.Discard, "test", &Options{Level: logger.LevelSilent, HistoryBuffer: 4 * 24, HistoryDir: dir})
	t.Cleanup(func() { require.NoError(t, l.Close()) })
	return l, l.MetaheuristicLogger("metaheuristic").(*MetaheuristicLogger), dir
}

// historyFiles returns the history files in dir.
func historyFiles(t *testing.T, dir string) []string {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(dir, "rko-history-*.bin"))
	require.NoError(t, err)
	return files
}

// register sends n records (local i, best 100-i, elapsed i*500ms) and returns them.
func register(ml logger.MetaheuristicLogger, n int) []logger.Sample {
	want := make([]logger.Sample, 0, n)
	for i := range n {
		d := logger.Sample{Local: i, Best: 100 - i, Elapsed: time.Duration(i) * 500 * time.Millisecond}
		ml.Register(d.Local, d.Best, i, d.Elapsed)
		want = append(want, d)
	}
	return want
}

func TestHistorySpillsWhenBufferIsFull(t *testing.T) {
	l, ml, dir := spillLogger(t)
	want := register(ml, 10) // two flushes of 4, 2 records left in memory

	require.Len(t, historyFiles(t, dir), 1)
	require.Len(t, ml.history, 2)
	require.Equal(t, 4, cap(ml.history)) // the buffer is reused, never grown
	require.Equal(t, want, l.Report()[0].Samples)
}

func TestCloseDeletesFilesAndIsIdempotent(t *testing.T) {
	l, ml, dir := spillLogger(t)
	register(ml, 5) // one flush of 4
	require.Len(t, historyFiles(t, dir), 1)

	require.NoError(t, l.Close())
	require.Empty(t, historyFiles(t, dir))
	require.NoError(t, l.Close()) // a second Close is a no-op
}

func TestRegisterAfterCloseCreatesNoFile(t *testing.T) {
	l, ml, dir := spillLogger(t)
	register(ml, 5) // one flush of 4, one record left in memory
	require.NoError(t, l.Close())

	register(ml, 10) // fills the buffer again: kept in memory
	require.Empty(t, historyFiles(t, dir))
	require.Len(t, ml.history, 11)

	l.Reset() // the next run spills again
	register(ml, 5)
	require.Len(t, historyFiles(t, dir), 1)
}

func TestSaveCSVContent(t *testing.T) {
	l, ml, _ := spillLogger(t)
	register(ml, 10)
	l.AddSolutionPool(9, 1250*time.Millisecond)
	l.AddSolutionPool(7, 2*time.Second)

	prefix := filepath.Join(t.TempDir(), "out")
	require.NoError(t, l.SaveCSV(prefix))

	b, err := os.ReadFile(prefix + "-metaheuristic-0.csv")
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
	register(l.MetaheuristicLogger("metaheuristic"), 3)
	l.AddSolutionPool(5, ms(1))

	prefix := filepath.Join(t.TempDir(), "out")
	require.NoError(t, l.SaveCSV(prefix))
	require.FileExists(t, prefix+"-metaheuristic-0.csv")
	require.FileExists(t, prefix+"-pool.csv")
}

func TestSaveCSVDefaultPrefixIsProblem(t *testing.T) {
	tmp := t.TempDir()
	t.Chdir(tmp)
	l := New(io.Discard, "inst/a:b", &Options{Level: logger.LevelSilent})
	l.MetaheuristicLogger("metaheuristic").Register(1, 1, 1, 0)

	require.NoError(t, l.SaveCSV(""))
	require.FileExists(t, filepath.Join(tmp, "inst_a_b-metaheuristic-0.csv"))
	require.FileExists(t, filepath.Join(tmp, "inst_a_b-pool.csv"))
}

func TestSaveCSVReportsEveryError(t *testing.T) {
	l, _ := newLog(logger.LevelSilent)
	l.MetaheuristicLogger("metaheuristic")

	err := l.SaveCSV(filepath.Join(t.TempDir(), "missing", "out"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "-metaheuristic-0.csv")
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
	l := New(&out, "test", &Options{Level: logger.LevelVerbose, PrintInterval: 0})
	a, b := l.MetaheuristicLogger("a"), l.MetaheuristicLogger("b")

	var wg sync.WaitGroup
	for _, ml := range []logger.MetaheuristicLogger{a, b} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range 200 {
				ml.Register(i, 1000-i, i, time.Duration(i)*time.Millisecond)
				ml.Event(slog.Int("i", i), time.Duration(i)*time.Millisecond)
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
	ml := l.MetaheuristicLogger("metaheuristic")
	a := slog.Int("a", 1)
	require.Zero(t, testing.AllocsPerRun(1000, func() { ml.Register(2, 1, 1, time.Second) }))
	require.Zero(t, testing.AllocsPerRun(1000, func() { ml.Event(a, time.Second) }))
}

func BenchmarkRegister(b *testing.B) {
	l := New(io.Discard, "test", &Options{Level: logger.LevelSilent, HistoryDir: b.TempDir()})
	b.Cleanup(func() { l.Close() })
	ml := l.MetaheuristicLogger("metaheuristic")
	b.ReportAllocs()
	for b.Loop() {
		ml.Register(2, 1, 1, time.Second)
	}
}
