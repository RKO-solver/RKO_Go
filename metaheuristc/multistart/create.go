package multistart

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

// CreateDefaultMultiStart creates a MultiStart with the default local search.
func CreateDefaultMultiStart(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *MultiStart {
	local := search.NewDefault(env, solutionPool, rg)

	return &MultiStart{
		env:          env,
		logger:       logger.SolverLogger(name),
		search:       local,
		solutionPool: solutionPool,
		rg:           rg,
	}
}

// CreateMultiStart creates a MultiStart with the local search se.
func CreateMultiStart(env definition.Environment, se search.Local, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *MultiStart {
	return &MultiStart{
		env:          env,
		logger:       logger.SolverLogger(name),
		search:       se,
		solutionPool: solutionPool,
		rg:           rg,
	}
}
