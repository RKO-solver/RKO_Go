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
	solutionPool := solution.NewDefaultPool(lo)
	rg := random.GetGlobalInstance()

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

func CreateDefaultSolverTimeLimitSecond(mh []MetaHeuristic, timeLimitSecond float64, env definition.Environment, logLevel logger.Level, saveReport bool, handler logger.Interface) *Solver {
	solver := CreateDefaultSolver(mh, env, logLevel, saveReport, handler)
	for _, sol := range solver.solvers {
		sol.SetTimeLimitSecond(timeLimitSecond)
	}

	return solver
}
