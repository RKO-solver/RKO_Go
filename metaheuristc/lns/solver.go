package lns

import (
	"math"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
)

var fareySequence = []float64{
	0.00,
	0.142857,
	0.166667,
	0.20,
	0.25,
	0.285714,
	0.333333,
}

func (lns *LNS) solve(solutionPool *solution.Pool) (*metaheuristc.RandomKeyValue, float64) {
	configuration := lns.configuration
	rg := lns.RG
	local := lns.search
	env := lns.env

	var bestSolution, localSolution *metaheuristc.RandomKeyValue

	start := time.Now()

	bestSolution = &metaheuristc.RandomKeyValue{
		RK:   make(definition.RandomKey, env.NumKeys()),
		Cost: 0,
	}
	rk.Reset(bestSolution.RK, rg)
	bestSolution.Cost = env.Cost(bestSolution.RK)
	solutionPool.AddSolution(bestSolution.Clone(), time.Since(start).Seconds())

	localSolution = &metaheuristc.RandomKeyValue{
		RK:   make(definition.RandomKey, env.NumKeys()),
		Cost: 0,
	}

	for iteration := 0; iteration < configuration.MaxIterations && time.Since(start).Seconds() < configuration.TimeLimitSeconds; iteration++ {
		// Future Q-Learning
		intensityMin, intensityMax := int(configuration.BetaMin*float64(env.NumKeys())), int(configuration.BetaMax*float64(env.NumKeys()))

		copy(localSolution.RK, bestSolution.RK)
		localSolution.Cost = bestSolution.Cost

		intensity := rg.RangeInt(intensityMin, intensityMax)

		order := rg.Permutation(env.NumKeys())

		for k := 0; k < intensity; k++ {
			pos := order[k]
			currentBest := math.MaxInt
			rkBest := 0.0

			for j := 0; j < len(fareySequence)-1; j++ {
				elapsedTime := time.Since(start).Seconds()
				if elapsedTime >= configuration.TimeLimitSeconds {
					return bestSolution, elapsedTime
				}

				localSolution.RK[pos] = rg.RangeFloat64(fareySequence[j], fareySequence[j+1])

				localSolution.Cost = env.Cost(localSolution.RK)

				if localSolution.Cost < currentBest {
					currentBest = localSolution.Cost
					rkBest = localSolution.RK[pos]
				}
			}

			// Keep the best found key
			localSolution.Cost = currentBest
			localSolution.RK[pos] = rkBest
		}

		local.Search(localSolution)

		if localSolution.Cost < bestSolution.Cost {
			copy(bestSolution.RK, localSolution.RK)
			bestSolution.Cost = localSolution.Cost

			bestSolutionCost := solutionPool.BestSolutionCost()
			if bestSolution.Cost < bestSolutionCost {
				solutionPool.AddSolution(bestSolution.Clone(), time.Since(start).Seconds())
			}
		}

		lns.logger.Register(localSolution.Cost, bestSolution.Cost, time.Since(start).Seconds(), "")
	}

	return bestSolution, time.Since(start).Seconds()
}
