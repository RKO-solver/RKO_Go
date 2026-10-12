// Package multistart implements a multi-start local search.
package multistart

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

const name = "MultiStart"

// MultiStart is a multi-start local search solver.
type MultiStart struct {
	env          definition.Environment
	logger       logger.SolverLogger
	search       search.Local
	rg           *random.Generator
	solutionPool *solution.Pool
}

func (m *MultiStart) Print() {
	search.PrintSolver(name, m.search)
}
