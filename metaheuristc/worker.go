package metaheuristc

import (
	"context"
	"sync"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
)

// Configuration holds the settings of a Worker.
type Configuration struct {
	Id int // worker number shown in the log
}

// Worker runs solver.Solve and reports the result to log.
func Worker(ctx context.Context, solver definition.Solver, configuration *Configuration, log logger.Logger, wg *sync.WaitGroup) {
	defer wg.Done()
	id := configuration.Id

	log.WorkerStarted(id, solver.Name())
	result := solver.Solve(ctx)
	log.WorkerDone(logger.WorkerResult{ID: id, Name: solver.Name(), Cost: result.Cost, Elapsed: result.Elapsed})
}
