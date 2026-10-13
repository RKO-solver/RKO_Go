package vns

import (
	"context"
	"log/slog"
	"math"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/rk"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/tuning"
)

func (vns *VNS) solve(ctx context.Context, solutionPool *solution.Pool) (*metaheuristic.RandomKeyValue, time.Duration) {
	ctx, ctxAction := definition.NewSolverContext(ctx)
	defer ctxAction.Stop()

	configuration := vns.configuration
	rg := vns.rg
	local := vns.search
	env := vns.env
	tuner := vns.tuner
	rate := configuration.Rate[0]

	var bestSolution, localSolution, neighbour *metaheuristic.RandomKeyValue

	start := definition.StartTime(ctx)

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
		l := math.Log1p(float64(ctxAction.Iteration()) / float64(vns.states))
		return l / (1 + l)
	}
	if tuner != nil {
		defer func() {
			if r, ok := tuner.(tuning.Reporter); ok {
				vns.logger.Policy(r.Policy())
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
			rate = vns.rate.At(tuner.Choose(progress()))
		}

		k := 1
		for k < rk.ShakeMax && ctx.Err() == nil {
			beta := rg.RangeFloat64(float64(k)*rate, float64(k+1)*rate)

			copy(neighbour.RK, localSolution.RK)
			rk.Shake(neighbour, beta, beta, rg, env)
			local.Search(ctx, neighbour)

			if neighbour.Cost < bestSolution.Cost {
				localSolution, neighbour = neighbour, localSolution

				rk.Copy(bestSolution, localSolution)

				solutionPool.Add(bestSolution, start)

				k = 1

				elapsedTime := time.Since(start)

				vns.logger.Event(slog.Attr{Key: "improved", Value: slog.GroupValue(slog.Int("iteration", ctxAction.Iteration()), slog.Int("local", localSolution.Cost))}, elapsedTime)
			} else {
				k++
			}
		}

		if tuner != nil {
			reward := 1.0
			if bestSolution.Cost >= previousBest {
				reward = float64(bestSolution.Cost-neighbour.Cost) / math.Max(1, math.Abs(float64(bestSolution.Cost)))
			}
			tuner.Learn(reward, progress())
		}

		vns.logger.Register(localSolution.Cost, bestSolution.Cost, ctxAction.Iteration(), time.Since(start))
	}

	return bestSolution, time.Since(start)
}
