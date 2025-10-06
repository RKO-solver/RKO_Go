package metaheuristc

import (
	"fmt"
	"sync"

	"github.com/RKO-solver/rko-go/definition"
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
func Worker(solver definition.Solver, configuration *Configuration, wg *sync.WaitGroup) {
	defer wg.Done()
	id := configuration.Id

	result := solver.Solve()
	fmt.Printf("(%d) %s Local Solution:\n\tCost: %d\n\t Time spent: %.2fs\n", id, solver.Name(), result.Cost, result.TimeSpentSecond)
}
