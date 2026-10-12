package rko

import (
	"fmt"
	"os"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/brkga"
	"github.com/RKO-solver/rko-go/metaheuristc/constants"
	"github.com/RKO-solver/rko-go/metaheuristc/ga"
	"github.com/RKO-solver/rko-go/metaheuristc/ils"
	"github.com/RKO-solver/rko-go/metaheuristc/lns"
	"github.com/RKO-solver/rko-go/metaheuristc/multistart"
	"github.com/RKO-solver/rko-go/metaheuristc/sa"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/metaheuristc/vns"
	"github.com/RKO-solver/rko-go/random"
)

// CreateDefaultSolver creates a Solver that runs the metaheuristics in mh with default parameters.
func CreateDefaultSolver(mh []MetaHeuristic, env definition.Environment, logger logger.Logger) *Solver {
	rg := random.NewRandom()
	solutionPool := solution.NewPool(env, rg, logger)

	solvers := make([]definition.Solver, 0, len(mh))

	for _, m := range mh {
		var solver definition.Solver
		srg := rg.Split()
		switch m {
		case ILS:
			solver = ils.CreateDefaultILS(env, srg, solutionPool, logger)
		case VNS:
			solver = vns.CreateDefaultVNS(env, srg, solutionPool, logger)
		case MULTISTART:
			solver = multistart.CreateDefaultMultiStart(env, srg, solutionPool, logger)
		case SA:
			solver = sa.CreateDefaultSA(env, srg, solutionPool, logger)
		case GA:
			solver = ga.CreateDefaultGA(env, srg, solutionPool, logger)
		case BRKGA:
			solver = brkga.CreateDefaultBRKGA(env, srg, solutionPool, logger)
		case LNS:
			solver = lns.CreateDefaultLNS(env, srg, solutionPool, logger)
		default:
			fmt.Printf("%s not implemented yet\n", GetMetaHeuristicString(m))
			continue
		}

		if solver != nil {
			solvers = append(solvers, solver)
		}
	}

	return &Solver{
		l:            logger,
		env:          env,
		solutionPool: solutionPool,
		solvers:      solvers,
		timeLimit:    constants.DefaultTimeLimit,
	}
}

// CreateDefaultSolverTimeLimit is like CreateDefaultSolver with an explicit time limit.
func CreateDefaultSolverTimeLimit(mh []MetaHeuristic, timeLimit time.Duration, env definition.Environment, logger logger.Logger) *Solver {
	solver := CreateDefaultSolver(mh, env, logger)

	solver.timeLimit = timeLimit

	return solver
}

// CreateFullSolver creates a Solver from already-built solvers and pool; non-positive limits mean no limit.
func CreateFullSolver(logger logger.Logger, env definition.Environment, timeLimit time.Duration, maxIterations int, solutionPool *solution.Pool, solvers []definition.Solver) *Solver {
	if timeLimit <= 0 && maxIterations <= 0 {
		fmt.Fprintln(os.Stdout, NoLimitWarning)
	}
	return &Solver{
		l:             logger,
		env:           env,
		solutionPool:  solutionPool,
		solvers:       solvers,
		timeLimit:     timeLimit,
		maxIterations: maxIterations,
	}
}
