// Package multistart implements a multi-start local search.
package multistart

import (
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/problem"
	"github.com/RKO-solver/rko-go/random"
)

const name = "MultiStart"

// MultiStart is a multi-start local search metaheuristic.
type MultiStart struct {
	env          problem.Environment
	logger       logger.MetaheuristicLogger
	search       metaheuristic.LocalSearch
	rg           *random.Generator
	solutionPool *solution.Pool
}

func (m *MultiStart) Print() {
	metaheuristic.Print(name, m.search)
}
