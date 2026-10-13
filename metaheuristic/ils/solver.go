package ils

import (
	"context"
	"math"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/rk"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/tuning"
)

func (ils *ILS) solve(ctx context.Context, solutionPool *solution.Pool) (*metaheuristic.RandomKeyValue, time.Duration) {
	ctx, ctxAction := definition.NewSolverContext(ctx)
	defer ctxAction.Stop()

	configuration := ils.configuration
	env := ils.env
	rg := ils.rg
	local := ils.search
	tuner := ils.tuner
	shakeMin, shakeMax := configuration.ShakeMin[0], configuration.ShakeMax[0]
	_, hasTimeOut := ctx.Deadline()
	metropolisCriterion := configuration.MetropolisCriterion[0] && hasTimeOut

	historyInformation := &history{
		defaultMin:         shakeMin,
		defaultMax:         shakeMax,
		min:                shakeMin,
		max:                shakeMax,
		timesNoImprovement: 0,
	}

	var localSolution, bestSolution, neighbour *metaheuristic.RandomKeyValue

	start := definition.StartTime(ctx)

	progress := func() float64 {
		if deadline, ok := ctx.Deadline(); ok {
			if total := deadline.Sub(start); total > 0 {
				return min(1, float64(time.Since(start))/float64(total))
			}
			return 1
		}
		l := math.Log1p(float64(ctxAction.Iteration()) / float64(ils.states))
		return l / (1 + l)
	}
	if tuner != nil {
		defer func() {
			if r, ok := tuner.(tuning.Reporter); ok {
				ils.logger.Policy(r.Policy())
			}
		}()
	}

	bestSolution = metaheuristic.NewRandomKeyValue(env.NumKeys())
	rk.Reset(bestSolution.RK, rg)
	bestSolution.Cost = env.Cost(bestSolution.RK)
	solutionPool.Add(bestSolution, start)

	localSolution = bestSolution.Clone()

	neighbour = &metaheuristic.RandomKeyValue{
		RK:   make(definition.RandomKey, env.NumKeys()),
		Cost: 0,
	}

	for ; ctx.Err() == nil; ctxAction.Step() {
		previousBest := bestSolution.Cost
		if tuner != nil {
			state := tuner.Choose(progress())
			historyInformation.defaultMin, historyInformation.defaultMax = ils.shakeMin.At(state), ils.shakeMax.At(state)
		}

		copy(neighbour.RK, localSolution.RK)
		shake(neighbour, historyInformation, rg, env)
		local.Search(ctx, neighbour)

		// acceptance criterion
		delta := neighbour.Cost - bestSolution.Cost
		if delta < 0 {
			localSolution, neighbour = neighbour, localSolution

			rk.Copy(bestSolution, localSolution)

			historyInformation.timesNoImprovement = 0

			solutionPool.Add(bestSolution, start)
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

		if tuner != nil {
			reward := 1.0
			if bestSolution.Cost >= previousBest {
				reward = float64(bestSolution.Cost-localSolution.Cost) / math.Max(1, math.Abs(float64(bestSolution.Cost)))
			}
			tuner.Learn(reward, progress())
		}

		elapsedTime := time.Since(start)

		ils.logger.Register(localSolution.Cost, bestSolution.Cost, ctxAction.Iteration(), elapsedTime)
	}

	return bestSolution, time.Since(start)
}
