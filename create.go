package rko

import (
	"fmt"

	"github.com/lucasmends/rko-go/definition"
	"github.com/lucasmends/rko-go/logger"
	"github.com/lucasmends/rko-go/metaheuristc/ga"
	"github.com/lucasmends/rko-go/metaheuristc/ils"
	"github.com/lucasmends/rko-go/metaheuristc/multistart"
	"github.com/lucasmends/rko-go/metaheuristc/sa"
	"github.com/lucasmends/rko-go/metaheuristc/solution"
	"github.com/lucasmends/rko-go/metaheuristc/vns"
	"github.com/lucasmends/rko-go/random"
)

func CreateDefaultSolver(mh []MetaHeuristic, env definition.Environment, logLevel logger.Level, saveReport bool, handler logger.Interface) *Solver {
	lo := logger.CreateLogger(logLevel, saveReport, handler)
	rg := random.GetGlobalInstance()
	solutionPool := solution.NewDefaultPool(env, rg, lo)

	solvers := make([]definition.Solver, 0, len(mh))

	for _, m := range mh {
		var solver definition.Solver
		switch m {
		case ILS:
			solver = ils.CreateDefaultILS(env, rg, solutionPool, lo)
		case VNS:
			solver = vns.CreateDefaultVNS(env, rg, solutionPool, lo)
		case MULTISTART:
			solver = multistart.CreateDefaultMultiStart(env, rg, solutionPool, lo)
		case SA:
			solver = sa.CreateDefaultSA(env, rg, solutionPool, lo)
		case GA:
			solver = ga.CreateDefaultGA(env, rg, solutionPool, lo)
		case BRKGA:
			solver = ga.CreateDefaultBRKGA(env, rg, solutionPool, lo)
		default:
			fmt.Printf("%s not implemented yet\n", GetMetaHeuristicString(m))
			continue
		}

		if solver != nil {
			solvers = append(solvers, solver)
		}
	}

	return &Solver{
		l:            lo,
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
func CreateDefaultSolverTimeLimitSecond(mh []MetaHeuristic, timeLimitSecond float64, env definition.Environment, logLevel logger.Level, saveReport bool, handler logger.Interface) *Solver {
	solver := CreateDefaultSolver(mh, env, logLevel, saveReport, handler)
	for _, sol := range solver.solvers {
		sol.SetTimeLimitSecond(timeLimitSecond)
	}

	return solver
}
