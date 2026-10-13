// Package lns implements Large Neighbourhood Search (LNS).
package lns

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/search"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/tuning"
	"github.com/RKO-solver/rko-go/random"
)

const name = "LNS"

// Config holds the LNS parameters.
type Config struct {
	BetaMin metaheuristic.List[float64] `yaml:"BetaMin"`
	BetaMax metaheuristic.List[float64] `yaml:"BetaMax"`
}

// LNS is a Large Neighbourhood Search solver.
type LNS struct {
	env           definition.Environment
	configuration *Config
	logger        logger.SolverLogger
	search        search.Local
	rg            *random.Generator
	solutionPool  *solution.Pool
	tuner         tuning.Tuner
	states        int
	betaMin       tuning.Param[float64]
	betaMax       tuning.Param[float64]
}

func (lns *LNS) Print() {
	search.PrintSolver(name, lns.search)
}
