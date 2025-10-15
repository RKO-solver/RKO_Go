package vns

import (
	"fmt"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
)

func (vns *VNS) solve(solutionPool *solution.Pool) (*metaheuristc.RandomKeyValue, float64) {
	configuration := vns.configuration
	rg := vns.RG
	local := vns.search
	env := vns.env

	var bestSolution, localSolution, neighbour *metaheuristc.RandomKeyValue

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
		k := 1
		for k < rk.ShakeMax && time.Since(start).Seconds() < configuration.TimeLimitSeconds {
			beta := rg.RangeFloat64(float64(k)*configuration.Rate, float64(k+1)*configuration.Rate)

			copy(neighbour.RK, localSolution.RK)
			rk.Shake(neighbour, beta, beta, rg, env)
			local.Search(neighbour)

			if neighbour.Cost < bestSolution.Cost {
				localSolution = neighbour
				copy(bestSolution.RK, localSolution.RK)
				bestSolution.Cost = localSolution.Cost

				bestSolutionCost := solutionPool.BestSolutionCost()
				if bestSolution.Cost < bestSolutionCost {
					solutionPool.AddSolution(bestSolution.Clone(), time.Since(start).Seconds())
				}

				k = 1
			} else {
				k++
			}

			elapsedTime := time.Since(start).Seconds()

			message := fmt.Sprintf("Iteration: %d, best solution: %d, local solution %d", iteration, bestSolution.Cost, localSolution.Cost)
			vns.logger.Verbose(message, elapsedTime)
		}

		vns.logger.Register(localSolution.Cost, bestSolution.Cost, time.Since(start).Seconds(), "")
	}

	return bestSolution, time.Since(start).Seconds()
}
