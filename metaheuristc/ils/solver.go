package ils

import (
	"fmt"
	"math"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
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

	var localSolution, bestSolution, neighbour *metaheuristc.RandomKeyValue

	start := time.Now()

	bestSolution = &metaheuristc.RandomKeyValue{
		RK:   make(definition.RandomKey, env.NumKeys()),
		Cost: 0,
	}
	rk.Reset(bestSolution.RK, rg)
	bestSolution.Cost = env.Cost(bestSolution.RK)
	solutionPool.AddSolution(bestSolution.Clone(), time.Since(start).Seconds())

	localSolution = bestSolution.Clone()

	neighbour = &metaheuristc.RandomKeyValue{
		RK:   make(definition.RandomKey, env.NumKeys()),
		Cost: 0,
	}

	for iteration := 0; iteration < configuration.MaxIterations && time.Since(start).Seconds() < configuration.TimeLimitSeconds; iteration++ {
		copy(neighbour.RK, localSolution.RK)
		shake(neighbour, historyInformation, rg, env)
		local.Search(neighbour)

		// acceptance criterion
		delta := neighbour.Cost - bestSolution.Cost
		if delta < 0 {
			localSolution = neighbour
			copy(bestSolution.RK, localSolution.RK)
			bestSolution.Cost = localSolution.Cost
			historyInformation.timesNoImprovement = 0

			bestSolutionCost := solutionPool.BestSolutionCost()
			if bestSolution.Cost < bestSolutionCost {
				solutionPool.AddSolution(bestSolution.Clone(), time.Since(start).Seconds())
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
			}

		}

		elapsedTime := time.Since(start).Seconds()

		ils.logger.Register(localSolution.Cost, bestSolution.Cost, elapsedTime, fmt.Sprintf("Iteration: %d", iteration))
	}

	return bestSolution, time.Since(start).Seconds()
}
