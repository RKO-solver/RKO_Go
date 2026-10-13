package printer

import (
	"io"
	"os"
	"time"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/internal/spill"
)

// Options configures a Log; the zero value is valid.
type Options struct {
	Level         logger.Level  // zero: LevelInfo
	PrintInterval time.Duration // 0: 500ms; below 10ms raised to 10ms
	HistoryBuffer int           // bytes per solver; 0: 512 KiB; negative: keep all history in memory
	HistoryDir    string        // "": system temp dir
}

// Default creates a Log that writes to os.Stdout with default options.
func Default(problem string) *Log {
	return New(os.Stdout, problem, nil)
}

// New returns a Log writing to w (safe for concurrent use; nil is os.Stdout).
func New(w io.Writer, problem string, opts *Options) *Log {
	var o Options
	if opts != nil {
		o = *opts
	}
	if w == nil {
		w = os.Stdout
	}
	if o.PrintInterval == 0 {
		o.PrintInterval = defaultPrintInterval
	}
	if o.HistoryBuffer == 0 {
		o.HistoryBuffer = defaultHistoryBytes
	}
	return &Log{
		level:    o.Level,
		interval: max(o.PrintInterval, minPrintInterval),
		w:        w,
		problem:  problem,
		hist:     spill.Config{Dir: o.HistoryDir, Bytes: o.HistoryBuffer, Problem: problem},
		pool:     make([]logger.PoolEntry, 0),
	}
}
