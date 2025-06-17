package multistart

import (
	"fmt"
	"github.com/lucasmends/rko-go/definition"
	"github.com/lucasmends/rko-go/metaheuristc"
	"github.com/lucasmends/rko-go/metaheuristc/rk"
	"github.com/lucasmends/rko-go/metaheuristc/solution"
	"time"
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
			solutionPool.AddSolution(localSolution.Clone())
			bestSolutionCost = localSolution.Cost
		}

		elapsedTime := time.Since(start).Seconds()
		message := fmt.Sprintf("Iteration: %d, best solution: %d, local solution %d, time %.2f", iteration, bestSolutionCost, localSolution.Cost, elapsedTime)
		m.logger.Debug(message)
		m.logger.Report(bestSolutionCost, localSolution.Cost, time.Since(start).Seconds())

	}

	return localSolution, time.Since(start).Seconds()
}
