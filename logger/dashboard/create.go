package dashboard

import (
	"os"
	"time"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/internal/spill"
)

// Options configures a Log; the zero value selects the defaults.
type Options struct {
	Level         logger.Level  // zero: LevelInfo
	Refresh       time.Duration // 0: 500ms; below 300ms raised to 300ms
	BufferSize    int           // update-channel slots; 0: 1024; below 1 raised to 1
	PoolLines     int           // 0: 3; below 1 raised to 1 (Verbose)
	HistoryBuffer int           // bytes per metaheuristic; 0: 128 KiB; negative: keep all history in memory
	HistoryDir    string        // "": system temp dir
}

// New creates a Log that draws to os.Stdout; see Start.
func New(problem string, opts *Options) *Log {
	var o Options
	if opts != nil {
		o = *opts
	}
	if o.Refresh == 0 {
		o.Refresh = defaultRefresh
	}
	if o.BufferSize == 0 {
		o.BufferSize = defaultBufferSize
	}
	if o.PoolLines == 0 {
		o.PoolLines = defaultPoolLines
	}
	if o.HistoryBuffer == 0 {
		o.HistoryBuffer = defaultHistoryBytes
	}
	o.BufferSize = max(o.BufferSize, 1)
	o.PoolLines = max(o.PoolLines, 1)
	if o.Level < logger.LevelVerbose {
		o.PoolLines = 0 // the pool is shown at Verbose only
	}

	// The final data store
	store := &information{
		numLinesPool:   o.PoolLines,
		pool:           make([]logger.PoolEntry, 0, 100),
		metaheuristics: make([]*metaheuristicHistory, 0),
		workers:        make([]logger.WorkerResult, 0),
	}

	return &Log{
		bufferSize: o.BufferSize,
		updateChan: make(chan message, o.BufferSize),
		data:       store,
		level:      o.Level,
		refresh:    max(o.Refresh, minimumRefresh),
		out:        os.Stdout,
		problem:    problem,
		hist:       spill.Config{Dir: o.HistoryDir, Bytes: o.HistoryBuffer, Problem: problem},
	}
}
