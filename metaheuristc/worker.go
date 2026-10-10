package metaheuristc

import (
	"context"
	"fmt"
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

	result := solver.Solve(ctx)
	log.WorkerDone(fmt.Sprintf("(%d) %s\n\tBest Local Solution Cost: %d\n\tTime spent: %.2fs\n", id, solver.Name(), result.Cost, result.TimeSpentSecond))
}
