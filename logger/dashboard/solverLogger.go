package dashboard

import (
	"log/slog"
	"time"

	"github.com/RKO-solver/rko-go/logger"
)

type SolverLogger struct {
	id    int
	name  string
	level logger.Level
	log   *Log
}

func (sl *SolverLogger) Register(local, best, iteration int, elapsed time.Duration) {
	sl.log.updateChan <- message{
		t:    infoMessage,
		id:   sl.id,
		info: solverInfo{local: local, best: best, iteration: iteration, elapsed: elapsed},
	}
}

func (sl *SolverLogger) Event(a slog.Attr, elapsed time.Duration) {
	// ignore if log level is not verbose
	if sl.level < logger.LevelVerbose {
		return
	}
	sl.log.updateChan <- message{
		t:     eventMessage,
		id:    sl.id,
		event: event{attr: a, elapsed: elapsed},
	}
}

// Policy stores t directly, skipping the update queue.
func (sl *SolverLogger) Policy(t logger.Table) {
	sl.log.data.setPolicy(sl.id, t)
}

// Compile-time check that *SolverLogger implements logger.SolverLogger.
var _ logger.SolverLogger = (*SolverLogger)(nil)
