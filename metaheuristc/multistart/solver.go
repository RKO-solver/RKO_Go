package multistart

import (
	"context"
	"fmt"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
)

// solve restarts from random keys, adding improvements to the pool, and returns the last local solution and the elapsed seconds.
func (m *MultiStart) solve(ctx context.Context, solutionPool *solution.Pool) (*metaheuristc.RandomKeyValue, float64) {
	ctx, ctxAction := definition.GetMHContext(ctx)
	defer ctxAction.Stop()

	rg := m.RG

	local := m.search

	localSolution := &metaheuristc.RandomKeyValue{
		RK:   make(definition.RandomKey, m.env.NumKeys()),
		Cost: 0,
	}

	start := definition.StartTime(ctx)
	for ; ctx.Err() == nil; ctxAction.Step() {

		rk.Reset(localSolution.RK, rg)
		localSolution.Cost = m.env.Cost(localSolution.RK)

		local.Search(ctx, localSolution)

		bestSolutionCost := solutionPool.BestSolutionCost()
		if localSolution.Cost < bestSolutionCost {
			solutionPool.AddSolution(localSolution, time.Since(start).Seconds())
			bestSolutionCost = localSolution.Cost
		}

		elapsedTime := time.Since(start).Seconds()

		m.logger.Register(localSolution.Cost, bestSolutionCost, elapsedTime, fmt.Sprintf("Iteration: %d", ctxAction.Iteration()))
	}

	return localSolution, time.Since(start).Seconds()
}
