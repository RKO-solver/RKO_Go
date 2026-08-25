package rko

import (
	"fmt"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
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

// CreateDefaultSolver creates a Solver running the given metaheuristics with default parameters.
func CreateDefaultSolver(mh []MetaHeuristic, env definition.Environment, logger logger.Logger) *Solver {
	rg := random.GetGlobalInstance()
	solutionPool := solution.NewDefaultPool(env, rg, logger)

	solvers := make([]definition.Solver, 0, len(mh))

	for _, m := range mh {
		var solver definition.Solver
		switch m {
		case ILS:
			solver = ils.CreateDefaultILS(env, rg, solutionPool, logger)
		case VNS:
			solver = vns.CreateDefaultVNS(env, rg, solutionPool, logger)
		case MULTISTART:
			solver = multistart.CreateDefaultMultiStart(env, rg, solutionPool, logger)
		case SA:
			solver = sa.CreateDefaultSA(env, rg, solutionPool, logger)
		case GA:
			solver = ga.CreateDefaultGA(env, rg, solutionPool, logger)
		case BRKGA:
			solver = ga.CreateDefaultBRKGA(env, rg, solutionPool, logger)
		case LNS:
			solver = lns.CreateDefaultLNS(env, rg, solutionPool, logger)
		default:
			fmt.Printf("%s not implemented yet\n", GetMetaHeuristicString(m))
			continue
		}

		if solver != nil {
			solvers = append(solvers, solver)
		}
	}

	return &Solver{
		l:                logger,
		rg:               rg,
		env:              env,
		solutionPool:     solutionPool,
		solvers:          solvers,
		timeLimitSeconds: constants.DefaultTimeLimit,
	}
}

// CreateDefaultSolverTimeLimitSecond is like CreateDefaultSolver with an explicit time limit.
func CreateDefaultSolverTimeLimitSecond(mh []MetaHeuristic, timeLimitSecond time.Duration, env definition.Environment, logger logger.Logger) *Solver {
	solver := CreateDefaultSolver(mh, env, logger)

	solver.timeLimitSeconds = timeLimitSecond

	return solver
}

// CreateFullSolver creates a Solver from already-built solvers, pool and generator.
func CreateFullSolver(logger logger.Logger, rg *random.Generator, env definition.Environment, timeLimitSeconds int64, solutionPool *solution.Pool, solvers []definition.Solver) *Solver {
	return &Solver{
		l:                logger,
		rg:               rg,
		env:              env,
		solutionPool:     solutionPool,
		solvers:          solvers,
		timeLimitSeconds: time.Duration(timeLimitSeconds),
	}
}
