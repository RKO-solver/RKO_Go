package channel

import (
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
