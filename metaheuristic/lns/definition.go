// Package lns implements Large Neighbourhood Search (LNS).
package lns

import (
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/tuning"
	"github.com/RKO-solver/rko-go/problem"
	"github.com/RKO-solver/rko-go/random"
)

const name = "LNS"

// Config holds the LNS parameters.
type Config struct {
	BetaMin metaheuristic.List[float64] `yaml:"BetaMin"`
	BetaMax metaheuristic.List[float64] `yaml:"BetaMax"`
}

// LNS is a Large Neighbourhood Search metaheuristic.
type LNS struct {
	env           problem.Environment
	configuration *Config
	logger        logger.MetaheuristicLogger
	search        metaheuristic.LocalSearch
	rg            *random.Generator
	solutionPool  *solution.Pool
	tuner         tuning.Tuner
	states        int
	betaMin       tuning.Param[float64]
	betaMax       tuning.Param[float64]
}

func (lns *LNS) Print() {
	metaheuristic.Print(name, lns.search)
}
