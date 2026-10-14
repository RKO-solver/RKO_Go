package dashboard

import (
	"log/slog"
	"time"

	"github.com/RKO-solver/rko-go/logger"
)

// MetaheuristicLogger sends the progress of one metaheuristic to the dashboard.
type MetaheuristicLogger struct {
	id    int
	name  string
	level logger.Level
	log   *Log
}

func (ml *MetaheuristicLogger) Register(local, best, iteration int, elapsed time.Duration) {
	ml.log.updateChan <- message{
		t:    infoMessage,
		id:   ml.id,
		info: metaheuristicInfo{local: local, best: best, iteration: iteration, elapsed: elapsed},
	}
}

func (ml *MetaheuristicLogger) Event(a slog.Attr, elapsed time.Duration) {
	// ignore if log level is not verbose
	if ml.level < logger.LevelVerbose {
		return
	}
	ml.log.updateChan <- message{
		t:     eventMessage,
		id:    ml.id,
		event: event{attr: a, elapsed: elapsed},
	}
}

// Policy stores t directly, skipping the update queue.
func (ml *MetaheuristicLogger) Policy(t logger.Table) {
	ml.log.data.setPolicy(ml.id, t)
}

// Compile-time check that *MetaheuristicLogger implements logger.MetaheuristicLogger.
var _ logger.MetaheuristicLogger = (*MetaheuristicLogger)(nil)
