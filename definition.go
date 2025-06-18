package rko

import (
	"fmt"
	"github.com/lucasmends/rko-go/definition"
	"github.com/lucasmends/rko-go/logger"
	"github.com/lucasmends/rko-go/metaheuristc"
	"github.com/lucasmends/rko-go/metaheuristc/solution"
	"github.com/lucasmends/rko-go/random"
	"sync"
)

// Solver coordinates the execution of multiple metaheuristic solvers in parallel,
// sharing a solution pool and environment. It manages logging, random number generation,
// and provides a unified interface for running and retrieving the best solution.
type Solver struct {
	l            *logger.Log                // Logger for reporting progress and results
	rg           *random.Generator          // Random number generator
	env          definition.Environment     // Problem environment (user implementation)
	solutionPool *solution.Pool             // Shared pool of solutions among solvers
	solvers      []definition.Solver        // List of metaheuristic solvers to run
}

// Solve runs all configured metaheuristic solvers in parallel, waits for their completion,
// and returns the best solution decoded into the problem's representation.
//
// Returns:
//   - The best solution found, decoded using the Environment's Decode method.
func (s *Solver) Solve() any {
	var wg sync.WaitGroup

	for i, sv := range s.solvers {
		fmt.Printf("Running solver %s (%d)\n", sv.Name(), i+1)
		wg.Add(1)
		go metaheuristc.Worker(sv, &metaheuristc.Configuration{Id: i + 1}, &wg)
	}

	wg.Wait()

	rk := s.solutionPool.BestSolution()
	return s.env.Decode(rk.RK)
}
