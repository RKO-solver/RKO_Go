// Package vns implements Variable Neighbourhood Search (VNS).
package vns

import (
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/tuning"
	"github.com/RKO-solver/rko-go/problem"
	"github.com/RKO-solver/rko-go/random"
)

const name = "VNS"

// Config holds the VNS parameters.
type Config struct {
	Rate metaheuristic.List[float64] `yaml:"Rate"`
}

// VNS is a Variable Neighbourhood Search metaheuristic.
type VNS struct {
	env           problem.Environment
	configuration *Config
	logger        logger.MetaheuristicLogger
	search        metaheuristic.LocalSearch
	rg            *random.Generator
	solutionPool  *solution.Pool
	tuner         tuning.Tuner
	states        int
	rate          tuning.Param[float64]
}

func (vns *VNS) Print() {
	metaheuristic.Print(name, vns.search)
}
