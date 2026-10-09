package channel

import (
	"io"
	"os"
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
