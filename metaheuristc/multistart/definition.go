package multistart

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

// name is the solver name used in logs.
const name = "MultiStart"

// MultiStart repeatedly resets to a random key and improves it with a local search.
type MultiStart struct {
	env          definition.Environment
	logger       logger.SolverLogger
	search       search.Local
	RG           *random.Generator
	solutionPool *solution.Pool
}

// Print prints the solver name and its local search.
func (m *MultiStart) Print() {
	search.PrintSolver(name, m.search)
}
