package multistart

import (
	"context"
	"math"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
)

func (m *MultiStart) solve(ctx context.Context, solutionPool *solution.Pool) (*metaheuristc.RandomKeyValue, time.Duration) {
	ctx, ctxAction := definition.GetMHContext(ctx)
	defer ctxAction.Stop()

	rg := m.rg

	local := m.search

	localSolution := metaheuristc.CreateEmptyRandomKeyValue(m.env.NumKeys())

	bestSolutionCost := math.MaxInt
	start := definition.StartTime(ctx)
	for ; ctx.Err() == nil; ctxAction.Step() {

		rk.Reset(localSolution.RK, rg)
		localSolution.Cost = m.env.Cost(localSolution.RK)

		local.Search(ctx, localSolution)

		if localSolution.Cost < bestSolutionCost {
			solutionPool.AddSolution(localSolution, time.Since(start))
			bestSolutionCost = localSolution.Cost
		}

		elapsedTime := time.Since(start)

		m.logger.Register(localSolution.Cost, bestSolutionCost, ctxAction.Iteration(), elapsedTime)
	}

	return localSolution, time.Since(start)
}
