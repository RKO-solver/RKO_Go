package metaheuristc

import (
	"fmt"
	"sync"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
)

type Configuration struct {
	Id int
}

// Worker runs a metaheuristic solver in a goroutine, setting its worker Id and printing the result.
// It is intended to be used with sync.WaitGroup for parallel execution.
//
// Parameters:
//   - solver: the metaheuristic solver implementing definition.Solver
//   - configuration: pointer to Configuration containing the worker Id
//   - wg: pointer to sync.WaitGroup for goroutine synchronization
func Worker(solver definition.Solver, configuration *Configuration, log logger.Logger, wg *sync.WaitGroup) {
	defer wg.Done()
	id := configuration.Id

	result := solver.Solve()
	log.WorkerDone(fmt.Sprintf("(%d) %s\n\tLocal Solution Cost: %d\n\tTime spent: %.2fs\n", id, solver.Name(), result.Cost, result.TimeSpentSecond))
}
