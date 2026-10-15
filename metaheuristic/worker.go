package metaheuristic

import (
	"context"
	"sync"
	"time"

	"github.com/RKO-solver/rko-go/logger"
)

// Result is the outcome of one Algorithm Solve call.
type Result struct {
	Solution any           // Solution decoded by Environment.Decode
	Cost     int           // Cost of the solution
	Elapsed  time.Duration // from the start of the run until Solve returned
}

// Algorithm is a metaheuristic run by rko.Solver.
type Algorithm interface {
	// Solve runs the metaheuristic until it finishes or ctx is done and returns its Result.
	Solve(ctx context.Context) Result
	// Name returns the metaheuristic name.
	Name() string
	// Print prints the metaheuristic configuration.
	Print()
}

// Worker runs a metaheuristic and reports its result to log.
func Worker(ctx context.Context, metaheuristic Algorithm, id int, log logger.Logger, wg *sync.WaitGroup) {
	defer wg.Done()

	log.WorkerStarted(id, metaheuristic.Name())
	result := metaheuristic.Solve(ctx)
	log.WorkerDone(logger.WorkerResult{ID: id, Name: metaheuristic.Name(), Cost: result.Cost, Elapsed: result.Elapsed})
}
