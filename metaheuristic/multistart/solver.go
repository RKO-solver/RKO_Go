package multistart

import (
	"context"
	"math"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/rk"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
)

func (m *MultiStart) solve(ctx context.Context, solutionPool *solution.Pool) (*metaheuristic.RandomKeyValue, time.Duration) {
	ctx, ctxAction := definition.NewSolverContext(ctx)
	defer ctxAction.Stop()

	rg := m.rg

	local := m.search

	localSolution := metaheuristic.NewRandomKeyValue(m.env.NumKeys())

	bestSolutionCost := math.MaxInt
	start := definition.StartTime(ctx)
	for ; ctx.Err() == nil; ctxAction.Step() {

		rk.Reset(localSolution.RK, rg)
		localSolution.Cost = m.env.Cost(localSolution.RK)

		local.Search(ctx, localSolution)

		if localSolution.Cost < bestSolutionCost {
			solutionPool.Add(localSolution, start)
			bestSolutionCost = localSolution.Cost
		}

		elapsedTime := time.Since(start)

		m.logger.Register(localSolution.Cost, bestSolutionCost, ctxAction.Iteration(), elapsedTime)
	}

	return localSolution, time.Since(start)
}
