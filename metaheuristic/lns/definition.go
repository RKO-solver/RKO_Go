// Package lns implements Large Neighbourhood Search (LNS).
package lns

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic/search"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/random"
)

const name = "LNS"

// Config holds the LNS parameters.
type Config struct {
	BetaMin float64 `yaml:"BetaMin"`
	BetaMax float64 `yaml:"BetaMax"`
}

// LNS is a Large Neighbourhood Search solver.
type LNS struct {
	env           definition.Environment
	configuration *Config
	logger        logger.SolverLogger
	search        search.Local
	rg            *random.Generator
	solutionPool  *solution.Pool
}

func (lns *LNS) Print() {
	search.PrintSolver(name, lns.search)
}
