package metaheuristic

import (
	"context"
	"sync"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
)

// Worker runs solver.Solve and reports the result to log.
func Worker(ctx context.Context, solver definition.Solver, id int, log logger.Logger, wg *sync.WaitGroup) {
	defer wg.Done()

	log.WorkerStarted(id, solver.Name())
	result := solver.Solve(ctx)
	log.WorkerDone(logger.WorkerResult{ID: id, Name: solver.Name(), Cost: result.Cost, Elapsed: result.Elapsed})
}
