package sa

import (
	"context"
	"log/slog"
	"math"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
)

func (sa *SimulatedAnnealing) solve(ctx context.Context, solutionPool *solution.Pool) (*metaheuristc.RandomKeyValue, time.Duration) {
	ctx, ctxAction := definition.NewSolverContext(ctx)
	defer ctxAction.Stop()

	configuration := sa.configuration
	env := sa.env
	rg := sa.rg
	local := sa.search
	alpha := 1.0 - configuration.Alpha
	temperatureLocal := configuration.TemperatureInitial
	timesReHeat := 1
	heatedLeft := configuration.QtdReheat

	var localSolution, neighbour *metaheuristc.RandomKeyValue

	localSolution = solutionPool.Best()

	start := definition.StartTime(ctx)

	if localSolution == nil {
		localSolution = metaheuristc.NewRandomKeyValue(env.NumKeys())
		rk.Reset(localSolution.RK, rg)
		localSolution.Cost = env.Cost(localSolution.RK)
		solutionPool.Add(localSolution, start)
	}

	neighbour = &metaheuristc.RandomKeyValue{
		RK:   make(definition.RandomKey, env.NumKeys()),
		Cost: 0,
	}

	bestSolutionCost := localSolution.Cost

	for ; temperatureLocal > configuration.TemperatureGoal && ctx.Err() == nil; ctxAction.Step() {

		for iteration := 0; iteration < configuration.Iterations && ctx.Err() == nil; iteration++ {

			copy(neighbour.RK, localSolution.RK)
			rk.Shake(neighbour, configuration.ShakeMin, configuration.ShakeMax, rg, env)
			local.Search(ctx, neighbour)

			delta := neighbour.Cost - localSolution.Cost
			if delta < 0 {
				localSolution, neighbour = neighbour, localSolution

				if localSolution.Cost < bestSolutionCost {
					bestSolutionCost = localSolution.Cost
				}

				solutionPool.Add(localSolution, start)
			} else {
				prob := math.Exp(-(float64(delta) + 0.00001) / temperatureLocal)
				if rg.Float64() < prob {
					localSolution, neighbour = neighbour, localSolution
				}
			}
		}

		elapsedTime := time.Since(start)

		if heatedLeft > 0 && temperatureLocal < configuration.TemperatureReheat {
			heatedLeft--
			sa.logger.Event(slog.Int("reheat", timesReHeat), time.Since(start))
			temperatureLocal = configuration.TemperatureInitial / math.Pow(2.0, float64(timesReHeat))
			timesReHeat++
		} else {
			temperatureLocal = temperatureLocal * alpha
		}

		sa.logger.Register(localSolution.Cost, bestSolutionCost, ctxAction.Iteration(), elapsedTime)
		sa.logger.Event(slog.Float64("temperature", temperatureLocal), elapsedTime)

	}

	return localSolution, time.Since(start)
}
