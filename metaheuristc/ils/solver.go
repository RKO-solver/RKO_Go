package ils

import (
	"context"
	"math"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
)

func (ils *ILS) solve(ctx context.Context, solutionPool *solution.Pool) (*metaheuristc.RandomKeyValue, time.Duration) {
	ctx, ctxAction := definition.GetMHContext(ctx)
	defer ctxAction.Stop()

	configuration := ils.configuration
	env := ils.env
	rg := ils.rg
	local := ils.search
	_, hasTimeOut := ctx.Deadline()
	metropolisCriterion := configuration.MetropolisCriterion && hasTimeOut

	historyInformation := &history{
		defaultMin:         configuration.ShakeMin,
		defaultMax:         configuration.ShakeMax,
		min:                configuration.ShakeMin,
		max:                configuration.ShakeMax,
		timesNoImprovement: 0,
	}

	var localSolution, bestSolution, neighbour *metaheuristc.RandomKeyValue

	start := definition.StartTime(ctx)

	bestSolution = metaheuristc.CreateEmptyRandomKeyValue(env.NumKeys())
	rk.Reset(bestSolution.RK, rg)
	bestSolution.Cost = env.Cost(bestSolution.RK)
	solutionPool.AddSolution(bestSolution, time.Since(start))

	localSolution = bestSolution.Clone()

	neighbour = &metaheuristc.RandomKeyValue{
		RK:   make(definition.RandomKey, env.NumKeys()),
		Cost: 0,
	}

	for ; ctx.Err() == nil; ctxAction.Step() {

		copy(neighbour.RK, localSolution.RK)
		shake(neighbour, historyInformation, rg, env)
		local.Search(ctx, neighbour)

		// acceptance criterion
		delta := neighbour.Cost - bestSolution.Cost
		if delta < 0 {
			localSolution, neighbour = neighbour, localSolution

			rk.Copy(bestSolution, localSolution)

			historyInformation.timesNoImprovement = 0

			bestSolutionCost := solutionPool.BestSolutionCost()
			if bestSolution.Cost < bestSolutionCost {
				solutionPool.AddSolution(bestSolution, time.Since(start))
			}
		} else {
			historyInformation.timesNoImprovement++
			if metropolisCriterion {
				deadline, _ := ctx.Deadline()
				elapsed := time.Since(start).Seconds()
				total := time.Until(deadline).Seconds() + elapsed
				prob := math.Exp(-(float64(delta) + 0.00001) / (1000.0 - 1000.0*(elapsed/total)))
				if rg.Float64() < prob {
					localSolution, neighbour = neighbour, localSolution
				}
			}

		}

		elapsedTime := time.Since(start)

		ils.logger.Register(localSolution.Cost, bestSolution.Cost, ctxAction.Iteration(), elapsedTime)
	}

	return bestSolution, time.Since(start)
}
