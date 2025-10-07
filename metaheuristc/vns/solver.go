package vns

import (
	"fmt"
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"time"
)

func (vns *VNS) solve(solutionPool *solution.Pool) (*metaheuristc.RandomKeyValue, float64) {
	configuration := vns.configuration
	rg := vns.RG
	env := vns.env

	var bestSolution, localSolution *metaheuristc.RandomKeyValue

	bestSolution = solutionPool.BestSolution()

	if bestSolution == nil {
		bestSolution = &metaheuristc.RandomKeyValue{
			RK:   make(definition.RandomKey, env.NumKeys()),
			Cost: 0,
		}
		rk.Reset(bestSolution.RK, rg)
		bestSolution.Cost = env.Cost(bestSolution.RK)
		solutionPool.AddSolution(bestSolution.Clone())
	}

	localSolution = &metaheuristc.RandomKeyValue{
		RK:   make(definition.RandomKey, env.NumKeys()),
		Cost: 0,
	}

	start := time.Now()
	for iteration := 0; iteration < configuration.MaxIterations && time.Since(start).Seconds() < configuration.TimeLimitSeconds; iteration++ {
		k := 0
		for k < rk.ShakeMax && time.Since(start).Seconds() < configuration.TimeLimitSeconds {
			beta := rg.RangeFloat64(float64(k)*configuration.Rate, float64(k+1)*configuration.Rate)

			copy(localSolution.RK, bestSolution.RK)
			rk.Shake(localSolution, beta, beta, rg, env)

			poolSolutionCost := solutionPool.BestSolutionCost()
			if localSolution.Cost < bestSolution.Cost {
				copy(bestSolution.RK, localSolution.RK)
				bestSolution.Cost = localSolution.Cost

				if bestSolution.Cost < poolSolutionCost {
					solutionPool.AddSolution(bestSolution.Clone())
				}

				k = 0
			} else {
				k++
			}

			elapsedTime := time.Since(start).Seconds()
			message := fmt.Sprintf("Iteration: %d, best solution: %d, local solution %d, time %.2f", iteration, poolSolutionCost, bestSolution.Cost, elapsedTime)
			vns.logger.Verbose(message)
			message := fmt.Sprintf("Iteration: %d, best solution: %d, local solution %d", iteration, poolSolutionCost, bestSolution.Cost)
			vns.logger.Verbose(message, elapsedTime)
		}

		vns.logger.Register(localSolution.Cost, bestSolution.Cost, time.Since(start).Seconds(), "")
		if solutionPool.BestSolutionCost() < bestSolution.Cost {
			bestSolution = solutionPool.BestSolution()
		}
	}

	return bestSolution, time.Since(start).Seconds()
}
