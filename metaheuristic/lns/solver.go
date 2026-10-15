package lns

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

var fareySequence = []float64{
	0.00,
	0.142857,
	0.166667,
	0.20,
	0.25,
	0.285714,
	0.333333,
	0.40,
	0.428571,
	0.50,
	0.571429,
	0.60,
	0.666667,
	0.714286,
	0.75,
	0.80,
	0.833333,
	0.857143,
	1.0,
}

func (lns *LNS) solve(ctx context.Context, solutionPool *solution.Pool) (*metaheuristic.RandomKeyValue, time.Duration) {
	ctx, ctxAction := metaheuristic.NewSolverContext(ctx)
	defer ctxAction.Stop()

	configuration := lns.configuration
	rg := lns.rg
	local := lns.search
	env := lns.env
	tuner := lns.tuner
	betaMin, betaMax := configuration.BetaMin[0], configuration.BetaMax[0]

	var bestSolution, localSolution *metaheuristic.RandomKeyValue

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
		l := math.Log1p(float64(ctxAction.Iteration()) / float64(lns.states))
		return l / (1 + l)
	}
	if tuner != nil {
		defer func() {
			if r, ok := tuner.(tuning.PolicyReporter); ok {
				lns.logger.Policy(r.Policy())
			}
		}()
	}

	bestSolution = metaheuristic.NewRandomKeyValue(env.NumKeys())
	rk.Reset(bestSolution.RK, rg)
	bestSolution.Cost = env.Cost(bestSolution.RK)
	solutionPool.Add(bestSolution, start)

	localSolution = &metaheuristic.RandomKeyValue{
		RK:   make(problem.RandomKey, env.NumKeys()),
		Cost: 0,
	}

	for ; ctx.Err() == nil; ctxAction.Step() {
		previousBest := bestSolution.Cost
		if tuner != nil {
			state := tuner.Choose(progress())
			betaMin, betaMax = lns.betaMin.At(state), lns.betaMax.At(state)
		}
		intensityMin, intensityMax := int(betaMin*float64(env.NumKeys())), int(betaMax*float64(env.NumKeys()))

		rk.Copy(localSolution, bestSolution)

		intensity := rg.RangeInt(intensityMin, intensityMax)

		order := rg.Perm(env.NumKeys())

		for k := range intensity {
			pos := order[k]
			currentBest := math.MaxInt
			rkBest := 0.0

			for j := 0; j < len(fareySequence)-1; j++ {
				if ctx.Err() != nil {
					return bestSolution, time.Since(start)

				}

				localSolution.RK[pos] = rg.RangeFloat64(fareySequence[j], fareySequence[j+1])

				localSolution.Cost = env.Cost(localSolution.RK)

				if localSolution.Cost < currentBest {
					currentBest = localSolution.Cost
					rkBest = localSolution.RK[pos]
				}
			}

			localSolution.Cost = currentBest
			localSolution.RK[pos] = rkBest
		}

		local.Search(ctx, localSolution)

		if localSolution.Cost < bestSolution.Cost {
			rk.Copy(bestSolution, localSolution)

			solutionPool.Add(bestSolution, start)
		}

		if tuner != nil {
			reward := 1.0
			if bestSolution.Cost >= previousBest {
				reward = float64(bestSolution.Cost-localSolution.Cost) / math.Max(1, math.Abs(float64(bestSolution.Cost)))
			}
			tuner.Learn(reward, progress())
		}

		lns.logger.Register(localSolution.Cost, bestSolution.Cost, ctxAction.Iteration(), time.Since(start))
	}

	return bestSolution, time.Since(start)
}
