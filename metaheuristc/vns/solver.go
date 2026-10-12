package vns

import (
	"context"
	"log/slog"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
)

func (vns *VNS) solve(ctx context.Context, solutionPool *solution.Pool) (*metaheuristc.RandomKeyValue, time.Duration) {
	ctx, ctxAction := definition.GetMHContext(ctx)
	defer ctxAction.Stop()

	configuration := vns.configuration
	rg := vns.RG
	local := vns.search
	env := vns.env

	var bestSolution, localSolution, neighbour *metaheuristc.RandomKeyValue

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

		k := 1
		for k < rk.ShakeMax && ctx.Err() == nil {
			beta := rg.RangeFloat64(float64(k)*configuration.Rate, float64(k+1)*configuration.Rate)

			copy(neighbour.RK, localSolution.RK)
			rk.Shake(neighbour, beta, beta, rg, env)
			local.Search(ctx, neighbour)

			if neighbour.Cost < bestSolution.Cost {
				localSolution, neighbour = neighbour, localSolution

				rk.Copy(bestSolution, localSolution)

				solutionPool.AddSolution(bestSolution, time.Since(start))

				k = 1

				elapsedTime := time.Since(start)

				vns.logger.Event(slog.Attr{Key: "improved", Value: slog.GroupValue(slog.Int("iteration", ctxAction.Iteration()), slog.Int("local", localSolution.Cost))}, elapsedTime)
			} else {
				k++
			}
		}

		vns.logger.Register(localSolution.Cost, bestSolution.Cost, ctxAction.Iteration(), time.Since(start))
	}

	return bestSolution, time.Since(start)
}
