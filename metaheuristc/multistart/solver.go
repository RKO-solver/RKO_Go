package multistart

import (
	"fmt"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
)

func (m *MultiStart) solve(solutionPool *solution.Pool) (*metaheuristc.RandomKeyValue, float64) {
	configuration := m.configuration
	rg := m.RG

	local := m.search

	localSolution := &metaheuristc.RandomKeyValue{
		RK:   make(definition.RandomKey, m.env.NumKeys()),
		Cost: 0,
	}

	start := time.Now()
	for iteration := 0; iteration < configuration.MaxIterations && time.Since(start).Seconds() < configuration.TimeLimitSeconds; iteration++ {
		rk.Reset(localSolution.RK, rg)
		localSolution.Cost = m.env.Cost(localSolution.RK)

		local.Search(localSolution)

		bestSolutionCost := solutionPool.BestSolutionCost()
		if localSolution.Cost < bestSolutionCost {
			solutionPool.AddSolution(localSolution.Clone(), time.Since(start).Seconds())
			bestSolutionCost = localSolution.Cost
		}

		elapsedTime := time.Since(start).Seconds()

		m.logger.Register(localSolution.Cost, bestSolutionCost, elapsedTime, fmt.Sprintf("Iteration: %d", iteration))

	}

	return localSolution, time.Since(start).Seconds()
}
