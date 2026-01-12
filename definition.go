package rko

import (
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

// Solver coordinates the execution of multiple metaheuristic solvers in parallel,
// sharing a solution pool and environment. It manages logging, random number generation,
// and provides a unified interface for running and retrieving the best solution.
type Solver struct {
	l            logger.Logger
	rg           *random.Generator      // Random number generator
	env          definition.Environment // Problem environment (user implementation)
	solutionPool *solution.Pool         // Shared pool of solutions among solvers
	solvers      []definition.Solver    // List of metaheuristic solvers to run
}

// Solve runs all configured metaheuristic solvers in parallel, waits for their completion,
// and returns the best solution decoded into the problem's representation.
//
// Returns:
//   - The best solution found, decoded using the Environment's Decode method.
func (s *Solver) Solve() any {
	logLevel := s.l.GetLogLevel()
	var loggerWg sync.WaitGroup

	if l, ok := s.l.(*channel.Log); ok {
		l.Start(&loggerWg)
	}

	var wg sync.WaitGroup
	for i, sv := range s.solvers {
		if logLevel > logger.SILENT {
			fmt.Printf("Running solver %s (%d)\n", sv.Name(), i)
		}
		wg.Add(1)
		go metaheuristc.Worker(sv, &metaheuristc.Configuration{Id: i}, s.l, &wg)
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

func (s *Solver) GetSolutionPool() *solution.Pool {
	return s.solutionPool
}

func (s *Solver) Print() {
	for i, sv := range s.solvers {
		fmt.Printf("(%d) ", i)
		sv.Print()
	}
}
