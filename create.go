package rko

import (
	"fmt"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/ga"
	"github.com/RKO-solver/rko-go/metaheuristc/ils"
	"github.com/RKO-solver/rko-go/metaheuristc/multistart"
	"github.com/RKO-solver/rko-go/metaheuristc/sa"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/metaheuristc/vns"
	"github.com/RKO-solver/rko-go/random"
)

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
		rg:           rg,
		env:          env,
		solutionPool: solutionPool,
		solvers:      solvers,
	}
}

// CreateDefaultSolverTimeLimitSecond creates a Solver as in CreateDefaultSolver, but also sets
// a time limit (in seconds) for all metaheuristics.
//
// Parameters:
//   - mh: slice of MetaHeuristic types to run
//   - timeLimitSecond: time limit in seconds for each metaheuristic
//   - env: user-implemented problem environment
//   - logLevel: logging level for all solvers
//   - saveReport: whether to save progress reports
//   - handler: logger implementation
//
// Returns:
//   - Pointer to a configured Solver with time limits set for all metaheuristics.
func CreateDefaultSolverTimeLimitSecond(mh []MetaHeuristic, timeLimitSecond float64, env definition.Environment, logger logger.Logger) *Solver {
	solver := CreateDefaultSolver(mh, env, logger)
	for _, sol := range solver.solvers {
		sol.SetTimeLimitSecond(timeLimitSecond)
	}

	return solver
}

func CreateFullSolver(logger logger.Logger, rg *random.Generator, env definition.Environment, solutionPool *solution.Pool, solvers []definition.Solver) *Solver {
	return &Solver{
		l:            logger,
		rg:           rg,
		env:          env,
		solutionPool: solutionPool,
		solvers:      solvers,
	}
}
