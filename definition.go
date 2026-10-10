// Package rko runs several random-key metaheuristics in parallel over a shared solution pool.
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
)

// Solver runs several metaheuristics in parallel over a shared solution pool.
type Solver struct {
	l             logger.Logger
	env           definition.Environment
	solutionPool  *solution.Pool
	solvers       []definition.Solver
	timeLimit     time.Duration
	maxIterations int
}

// Solve runs the metaheuristics and returns the pool's best solution, decoded.
func (s *Solver) Solve(ctx context.Context) any {
	logLevel := s.l.GetLogLevel()
	var loggerWg sync.WaitGroup

	s.l.Reset()
	if l, ok := s.l.(*channel.Log); ok {
		l.Start(&loggerWg)
	}

	var wg sync.WaitGroup
	if s.maxIterations > 0 {
		ctx = definition.WithMaxIterations(ctx, s.maxIterations)
	}
	if s.timeLimit > 0 {
		ctx_, cancel := context.WithTimeout(ctx, s.timeLimit)
		ctx = ctx_
		defer cancel()
	}

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
			l.Finish()
			l.WorkersPrint()
		}
	} else {
		wg.Wait()
	}

	rk := s.solutionPool.BestSolution()
	return s.env.Decode(rk.RK)
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
