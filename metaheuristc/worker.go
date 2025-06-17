package metaheuristc

import (
	"fmt"
	"github.com/lucasmends/rko-go/definition"
	"sync"
)

type Configuration struct {
	Id int
}

func Worker(solver definition.Solver, configuration *Configuration, wg *sync.WaitGroup) {
	defer wg.Done()
	id := configuration.Id
	solver.SetIdWorker(id)

	result := solver.Solve()
	fmt.Printf("(%d) %s Local Solution:\n\tCost: %d\n\t Time spent: %.2fs\n", id, solver.Name(), result.Cost, result.TimeSpentSecond)
}
