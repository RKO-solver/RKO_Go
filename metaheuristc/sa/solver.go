package sa

import (
	"fmt"
	"math"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
)

func (sa *SimulatedAnnealing) solve(solutionPool *solution.Pool) (*metaheuristc.RandomKeyValue, float64) {
	configuration := sa.configuration
	env := sa.env
	rg := sa.RG
	local := sa.search
	alpha := 1.0 - configuration.Alpha
	temperatureLocal := configuration.TemperatureInitial
	timesReHeat := 1
	heatedLeft := configuration.QtdReheat

	var localSolution, neighbour *metaheuristc.RandomKeyValue

	localSolution = solutionPool.BestSolution()

	start := time.Now()

	if localSolution == nil {
		localSolution = &metaheuristc.RandomKeyValue{
			RK:   make(definition.RandomKey, env.NumKeys()),
			Cost: 0,
		}
		rk.Reset(localSolution.RK, rg)
		localSolution.Cost = env.Cost(localSolution.RK)
		solutionPool.AddSolution(localSolution.Clone(), time.Since(start).Seconds())
	}

	neighbour = &metaheuristc.RandomKeyValue{
		RK:   make(definition.RandomKey, env.NumKeys()),
		Cost: 0,
	}

	var bestSolutionCost int

	for time.Since(start).Seconds() < configuration.TimeLimitSeconds && temperatureLocal > configuration.TemperatureGoal {

		for iteration := 0; iteration < configuration.Iterations && time.Since(start).Seconds() < configuration.TimeLimitSeconds; iteration++ {
			copy(neighbour.RK, localSolution.RK)
			rk.Shake(neighbour, configuration.ShakeMin, configuration.ShakeMax, rg, env)
			local.Search(neighbour)

			delta := neighbour.Cost - localSolution.Cost
			if delta < 0 {
				localSolution.Cost = neighbour.Cost
				copy(localSolution.RK, neighbour.RK)

				bestSolutionCost = solutionPool.BestSolutionCost()
				if neighbour.Cost < bestSolutionCost {
					solutionPool.AddSolution(neighbour.Clone(), time.Since(start).Seconds())
				}
			} else {
				prob := math.Exp(-(float64(delta) + 0.00001) / temperatureLocal)
				if rg.Float64() < prob {
					localSolution.Cost = neighbour.Cost
					copy(localSolution.RK, neighbour.RK)
				}
			}
		}

		elapsedTime := time.Since(start).Seconds()

		if heatedLeft > 0 && temperatureLocal < configuration.TemperatureReheat {
			heatedLeft--
			sa.logger.Verbose(fmt.Sprintf("Re-heat %d", timesReHeat), time.Since(start).Seconds())
			temperatureLocal = configuration.TemperatureInitial / math.Pow(2.0, float64(timesReHeat))
			timesReHeat++
		} else {
			temperatureLocal = temperatureLocal * alpha
		}

		sa.logger.Register(localSolution.Cost, bestSolutionCost, elapsedTime, fmt.Sprintf("Temperature %.4f", temperatureLocal))

	}

	return localSolution, time.Since(start).Seconds()
}
