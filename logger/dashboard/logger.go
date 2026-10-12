// Package dashboard implements a live-refreshing terminal Logger for interactive runs.
package dashboard

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/internal/spill"
)

// Log is the channel-based Logger.
type Log struct {
	updateChan chan message
	bufferSize int
	data       *information
	level      logger.Level
	refresh    time.Duration
	out        io.Writer
	problem    string
	hist       spill.Config // spill settings shared by every solver
}

func (l *Log) AddSolutionPool(cost int, elapsed time.Duration) {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	l.data.pool = append(l.data.pool, logger.PoolEntry{Cost: cost, Elapsed: elapsed})
}

func (l *Log) WorkerStarted(id int, name string) {
	if l.level > logger.LevelSilent {
		fmt.Fprintf(l.out, "(%d) %s started\n", id, name)
	}
}

func (l *Log) WorkerDone(r logger.WorkerResult) {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	l.data.workers = append(l.data.workers, r)
}

// Reset empties the previous run's data, deletes its history files and opens a fresh update channel.
func (l *Log) Reset() {
	l.data.reset(l.hist)
	l.updateChan = make(chan message, l.bufferSize)
}

// Start runs the aggregator and the display; call stop after the solvers return.
func (l *Log) Start() (stop func()) {
	ctx, cancel := context.WithCancel(context.Background())
	ch, store := l.updateChan, l.data
	var wg sync.WaitGroup

	// The single aggregator goroutine: the only one that writes solver data to store.
	wg.Go(func() {
		for {
			select {
			case m := <-ch:
				store.apply(m)
			case <-ctx.Done():
				for len(ch) > 0 { // stop runs after the solvers returned: apply what is left
					store.apply(<-ch)
				}
				return
			}
		}
	})

	if l.level > logger.LevelSilent {
		wg.Go(func() { // the display ticker
			t := time.NewTicker(l.refresh)
			defer t.Stop()
			for {
				select {
				case <-t.C:
					store.printShell(l.out)
				case <-ctx.Done():
					return
				}
			}
		})
	}

	return func() {
		cancel()
		wg.Wait()
		if l.level > logger.LevelSilent {
			store.printFinal(l.out)
		}
	}
}

func (l *Log) SolverLogger(name string) logger.SolverLogger {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	id := len(l.data.solvers)
	l.data.solvers = append(l.data.solvers, &solverHistory{id: id, name: name})
	l.data.solvers[id].reset(l.hist) // sizes the buffer; the file is created on the first spill

	return &SolverLogger{
		id:    id,
		name:  name,
		level: l.level,
		log:   l,
	}
}

// Close deletes the last run's history files; call it after reading the results.
// Spilling stays off until the next Reset.
func (l *Log) Close() error {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	var err error
	for _, s := range l.data.solvers {
		err = errors.Join(err, s.file.Remove())
	}
	return err
}

// Compile-time checks that *Log implements logger.Recorder and logger.Runner.
var (
	_ logger.Recorder = (*Log)(nil)
	_ logger.Runner   = (*Log)(nil)
)
