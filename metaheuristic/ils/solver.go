package ils

import (
	"context"
	"math"
	"time"

	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/rk"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/tuning"
	"github.com/RKO-solver/rko-go/problem"
)

func (ils *ILS) solve(ctx context.Context, solutionPool *solution.Pool) (*metaheuristic.RandomKeyValue, time.Duration) {
	ctx, ctxAction := metaheuristic.NewSolverContext(ctx)
	defer ctxAction.Stop()

	configuration := ils.configuration
	env := ils.env
	rg := ils.rg
	local := ils.search
	tuner := ils.tuner
	shakeMin, shakeMax := configuration.ShakeMin[0], configuration.ShakeMax[0]
	_, hasTimeOut := ctx.Deadline()
	metropolisCriterion := *configuration.MetropolisCriterion && hasTimeOut

	historyInformation := &history{
		defaultMin:         shakeMin,
		defaultMax:         shakeMax,
		min:                shakeMin,
		max:                shakeMax,
		timesNoImprovement: 0,
	}

	var localSolution, bestSolution, neighbour *metaheuristic.RandomKeyValue

	start := metaheuristic.StartTime(ctx)

	deadline, ok := ctx.Deadline()
	total := deadline.Sub(start)
	ok = ok && total > 0
	n := ctxAction.Limit()
	progress := func() float64 {
		if ok {
			return min(1, float64(time.Since(start))/float64(total))
		}
		if n > 0 {
			return min(1, float64(ctxAction.Iteration())/float64(n))
		}
		l := math.Log1p(float64(ctxAction.Iteration()) / float64(ils.states))
		return l / (1 + l)
	}
	if tuner != nil {
		defer func() {
			if r, ok := tuner.(tuning.PolicyReporter); ok {
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
		RK:   make(problem.RandomKey, env.NumKeys()),
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
				elapsed := time.Since(start).Seconds()
				prob := math.Exp(-(float64(delta) + 0.00001) / (1000.0 - 1000.0*(elapsed/total.Seconds())))
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
