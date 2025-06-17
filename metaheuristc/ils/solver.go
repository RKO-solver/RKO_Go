package ils

import (
	"fmt"
	"github.com/lucasmends/rko-go/definition"
	"github.com/lucasmends/rko-go/metaheuristc"
	"github.com/lucasmends/rko-go/metaheuristc/rk"
	"github.com/lucasmends/rko-go/metaheuristc/solution"
	"math"
	"time"
)

func (ils *ILS) solve(solutionPool *solution.Pool) (*metaheuristc.RandomKeyValue, float64) {
	configuration := ils.configuration
	env := ils.env
	rg := ils.RG
	local := ils.search
	metropolisCriterion := configuration.MetropolisCriterion && (configuration.TimeLimitSeconds < math.MaxInt)

	historyInformation := &history{
		defaultMin:         configuration.ShakeMin,
		defaultMax:         configuration.ShakeMax,
		min:                configuration.ShakeMin,
		max:                configuration.ShakeMax,
		timesNoImprovement: 0,
	}

	var localSolution, neighbour *metaheuristc.RandomKeyValue

	localSolution = solutionPool.BestSolution()

	if localSolution == nil {
		localSolution = &metaheuristc.RandomKeyValue{
			RK:   make(definition.RandomKey, env.NumKeys()),
			Cost: 0,
		}
		rk.Reset(localSolution.RK, rg)
		localSolution.Cost = env.Cost(localSolution.RK)
		solutionPool.AddSolution(localSolution.Clone())
	}

	neighbour = &metaheuristc.RandomKeyValue{
		RK:   make(definition.RandomKey, env.NumKeys()),
		Cost: 0,
	}

	start := time.Now()
	for iteration := 0; iteration < configuration.MaxIterations && time.Since(start).Seconds() < configuration.TimeLimitSeconds; iteration++ {
		copy(neighbour.RK, localSolution.RK)

		shake(neighbour, historyInformation, rg, env)
		local.Search(neighbour)

		// acceptance criterion
		bestSolutionCost := solutionPool.BestSolutionCost()
		delta := neighbour.Cost - localSolution.Cost
		if delta < 0 {
			localSolution.Cost = neighbour.Cost
			copy(localSolution.RK, neighbour.RK)
			historyInformation.timesNoImprovement = 0

			if neighbour.Cost < bestSolutionCost {
				solutionPool.AddSolution(neighbour.Clone())
			}
		} else {
			historyInformation.timesNoImprovement++
			if metropolisCriterion {
				prob := math.Exp(-(float64(delta) + 0.00001) / (1000.0 - 1000.0*(time.Since(start).Seconds()/(configuration.TimeLimitSeconds+0.5))))
				if rg.Float64() < prob {
					localSolution.Cost = neighbour.Cost
					copy(localSolution.RK, neighbour.RK)
				}
				// Get the best solution for local solution
			} else if historyInformation.timesNoImprovement > 100 {
				localSolution = solutionPool.BestSolution()
			}

		}

		elapsedTime := time.Since(start).Seconds()
		message := fmt.Sprintf("Iteration: %d, best solution: %d, local solution %d, time %.2f", iteration, bestSolutionCost, localSolution.Cost, elapsedTime)
		ils.logger.Debug(message)
		ils.logger.Report(bestSolutionCost, localSolution.Cost, elapsedTime)
	}

	return localSolution, time.Since(start).Seconds()
}
