package rko

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/channel"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

// Solver runs several metaheuristics in parallel over a shared solution pool.
type Solver struct {
	l                logger.Logger
	rg               *random.Generator      // Random number generator
	env              definition.Environment // Problem environment (user implementation)
	solutionPool     *solution.Pool         // Shared pool of solutions among solvers
	solvers          []definition.Solver    // List of metaheuristic solvers to run
	timeLimitSeconds time.Duration
}

// SolveCtx is like Solve but stops when ctx is cancelled or expires.
func (s *Solver) SolveCtx(ctx context.Context) any {
	logLevel := s.l.GetLogLevel()
	var loggerWg sync.WaitGroup

	if l, ok := s.l.(*channel.Log); ok {
		l.Start(&loggerWg)
	}

	var wg sync.WaitGroup
	ctxStartTime := definition.WithStartTime(ctx, time.Now())
	for i, sv := range s.solvers {
		if logLevel > logger.SILENT {
			fmt.Printf("Running solver %s (%d)\n", sv.Name(), i)
		}
		wg.Add(1)
		go metaheuristc.Worker(ctxStartTime, sv, &metaheuristc.Configuration{Id: i}, s.l, &wg)
	}

	if l, ok := s.l.(*channel.Log); ok {
		if logLevel > logger.SILENT {
			ticker := time.NewTicker(l.GetTicker())
			defer ticker.Stop()

			workersDone := make(chan bool)
			go func() {
				wg.Wait()
				close(workersDone)
			}()

			displaying := true
			for displaying {
				select {
				case <-ticker.C:
					l.Print()
				case <-workersDone:
					displaying = false
				}
			}
		} else {
			wg.Wait()
		}

		l.Shutdown()
		loggerWg.Wait()
		if logLevel > logger.SILENT {
			l.WorkersPrint()
		}
	} else {
		wg.Wait()
	}

	rk := s.solutionPool.BestSolution()
	return s.env.Decode(rk.RK)
}

// Solve runs every metaheuristic under the configured time limit and returns the best solution.
func (s *Solver) Solve() any {
	if s.timeLimitSeconds > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), s.timeLimitSeconds)
		defer cancel()

		return s.SolveCtx(ctx)
	}
	ctx := context.Background()

	return s.SolveCtx(ctx)
}

// GetSolutionPool returns the pool shared by every running metaheuristic.
func (s *Solver) GetSolutionPool() *solution.Pool {
	return s.solutionPool
}

// Print prints the configuration of every metaheuristic in the solver.
func (s *Solver) Print() {
	for i, sv := range s.solvers {
		fmt.Printf("(%d) ", i)
		sv.Print()
	}
}
