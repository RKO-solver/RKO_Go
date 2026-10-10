// Package vns implements Variable Neighbourhood Search (VNS).
package vns

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

const name = "VNS"

// Configuration holds the VNS parameters.
type Configuration struct {
	Rate float64 `yaml:"Rate"`
}

// VNS is a Variable Neighbourhood Search solver.
type VNS struct {
	env           definition.Environment
	configuration *Configuration
	logger        logger.SolverLogger
	search        search.Local
	RG            *random.Generator
	solutionPool  *solution.Pool
}

func (vns *VNS) Print() {
	search.PrintSolver(name, vns.search)
}
