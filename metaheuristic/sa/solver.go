package sa

import (
	"context"
	"log/slog"
	"math"
	"time"

	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/rk"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/tuning"
	"github.com/RKO-solver/rko-go/problem"
)

func (sa *SimulatedAnnealing) solve(ctx context.Context, solutionPool *solution.Pool) (*metaheuristic.RandomKeyValue, time.Duration) {
	ctx, ctxAction := metaheuristic.NewSolverContext(ctx)
	defer ctxAction.Stop()

	configuration := sa.configuration
	env := sa.env
	rg := sa.rg
	local := sa.search
	tuner := sa.tuner
	alpha := 1.0 - configuration.Alpha[0]
	temperatureInitial := configuration.TemperatureInitial[0]
	temperatureGoal := configuration.TemperatureGoal[0]
	temperatureReheat := configuration.TemperatureReheat[0]
	shakeMin, shakeMax := configuration.ShakeMin[0], configuration.ShakeMax[0]
	qtdReheat := configuration.QtdReheat[0]
	iterations := configuration.Iterations[0]
	restart := *configuration.Restart
	temperatureLocal := temperatureInitial
	timesReHeat := 1

	var localSolution, neighbour *metaheuristic.RandomKeyValue

	localSolution = solutionPool.Best()

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
		l := math.Log1p(float64(ctxAction.Iteration()) / float64(sa.states))
		return l / (1 + l)
	}
	if tuner != nil {
		defer func() {
			if r, ok := tuner.(tuning.PolicyReporter); ok {
				sa.logger.Policy(r.Policy())
			}
		}()
	}

	if localSolution == nil {
		localSolution = metaheuristic.NewRandomKeyValue(env.NumKeys())
		rk.Reset(localSolution.RK, rg)
		localSolution.Cost = env.Cost(localSolution.RK)
		solutionPool.Add(localSolution, start)
	}

	neighbour = &metaheuristic.RandomKeyValue{
		RK:   make(problem.RandomKey, env.NumKeys()),
		Cost: 0,
	}

	bestSolutionCost := localSolution.Cost

	for ; ctx.Err() == nil && (restart || temperatureLocal > temperatureGoal); ctxAction.Step() {
		if temperatureLocal <= temperatureGoal {
			temperatureLocal = temperatureInitial
			rk.Reset(localSolution.RK, rg)
			localSolution.Cost = env.Cost(localSolution.RK)
			solutionPool.Add(localSolution, start)

			timesReHeat = 1
			bestSolutionCost = localSolution.Cost
			sa.logger.Event(slog.Attr{Key: "restart", Value: slog.GroupValue(slog.Int("iteration", ctxAction.Iteration()))}, time.Since(start))
		}

		previousBest := bestSolutionCost
		if tuner != nil {
			state := tuner.Choose(progress())
			alpha = 1.0 - sa.alpha.At(state)
			temperatureInitial = sa.temperatureInitial.At(state)
			temperatureGoal = sa.temperatureGoal.At(state)
			temperatureReheat = sa.temperatureReheat.At(state)
			shakeMin, shakeMax = sa.shakeMin.At(state), sa.shakeMax.At(state)
			qtdReheat = sa.qtdReheat.At(state)
			iterations = sa.iterations.At(state)
		}

		for iteration := 0; iteration < iterations && ctx.Err() == nil; iteration++ {

			copy(neighbour.RK, localSolution.RK)
			rk.Shake(neighbour, shakeMin, shakeMax, rg, env)
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

		if tuner != nil {
			reward := 1.0
			if bestSolutionCost >= previousBest {
				reward = float64(bestSolutionCost-localSolution.Cost) / math.Max(1, math.Abs(float64(bestSolutionCost)))
			}
			tuner.Learn(reward, progress())
		}

		elapsedTime := time.Since(start)

		if int(qtdReheat) >= timesReHeat && temperatureLocal < temperatureReheat {
			sa.logger.Event(slog.Int("reheat", timesReHeat), time.Since(start))
			temperatureLocal = temperatureInitial / math.Pow(2.0, float64(timesReHeat))
			timesReHeat++
		} else {
			temperatureLocal = temperatureLocal * alpha
		}

		sa.logger.Register(localSolution.Cost, bestSolutionCost, ctxAction.Iteration(), elapsedTime)
		sa.logger.Event(slog.Float64("temperature", temperatureLocal), elapsedTime)

	}

	return localSolution, time.Since(start)
}
