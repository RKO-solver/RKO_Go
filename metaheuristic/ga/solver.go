package ga

import (
	"context"
	"log/slog"
	"math"
	"slices"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/internal/genetic"
	"github.com/RKO-solver/rko-go/metaheuristic/rk"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/tuning"
)

func (ga *GA) solve(ctx context.Context, solutionPool *solution.Pool) (*metaheuristic.RandomKeyValue, time.Duration) {
	ctx, ctxAction := definition.NewSolverContext(ctx)
	defer ctxAction.Stop()

	configuration := ga.configuration
	env := ga.env
	rg := ga.rg
	local := ga.search
	numKeys := env.NumKeys()
	tuner := ga.tuner
	populationSize := configuration.PopulationSize[0]
	crossoverAlpha := configuration.CrossoverAlpha[0]
	mutationAlpha := configuration.MutationAlpha[0]
	maxGenerationNoImprovement := configuration.MaxGenerationNoImprovement[0]
	maxSize := slices.Max(configuration.PopulationSize)

	children := make([]*metaheuristic.RandomKeyValue, maxSize)
	population := make([]*metaheuristic.RandomKeyValue, maxSize)
	for i := range maxSize {
		population[i] = metaheuristic.NewRandomKeyValue(numKeys)
		children[i] = metaheuristic.NewRandomKeyValue(numKeys)
	}

	hasImproved := false
	isOdd := populationSize%2 == 1
	generationNoImprovement := 0

	start := definition.StartTime(ctx)
	bestPerson := genetic.InitialPopulation(population[:populationSize], env, false, rg)
	solutionPool.Add(bestPerson, start)

	progress := func() float64 {
		if deadline, ok := ctx.Deadline(); ok {
			if total := deadline.Sub(start); total > 0 {
				return min(1, float64(time.Since(start))/float64(total))
			}
			return 1
		}
		l := math.Log1p(float64(ctxAction.Iteration()) / float64(ga.states))
		return l / (1 + l)
	}
	if tuner != nil {
		defer func() {
			if r, ok := tuner.(tuning.Reporter); ok {
				ga.logger.Policy(r.Policy())
			}
		}()
	}

	var key1, key2 float64

	for ; ctx.Err() == nil; ctxAction.Step() {
		bestGen := math.MaxInt
		if tuner != nil {
			state := tuner.Choose(progress())
			previousSize := populationSize
			populationSize = ga.populationSize.At(state)
			crossoverAlpha = ga.crossoverAlpha.At(state)
			mutationAlpha = ga.mutationAlpha.At(state)
			maxGenerationNoImprovement = ga.maxGenerationNoImprovement.At(state)
			for i := previousSize; i < populationSize; i++ {
				rk.Reset(population[i].RK, rg)
				population[i].Cost = env.Cost(population[i].RK)
			}
			isOdd = populationSize%2 == 1
		}

		for j := 0; j < populationSize-1; j += 2 {
			offspring1 := tournament(population[:populationSize], rg)
			offspring2 := tournament(population[:populationSize], rg)
			copy(children[j].RK, offspring1.RK)
			copy(children[j+1].RK, offspring2.RK)

			// crossover
			if rg.Float64() < crossoverAlpha {
				for i := range numKeys {
					key1 = children[j].RK[i]
					key2 = children[j+1].RK[i]
					if rg.Float64() < 0.5 {
						children[j].RK[i] = key2
					}

					if rg.Float64() < 0.5 {
						children[j+1].RK[i] = key1
					}

					if rg.Float64() < mutationAlpha {
						children[j].RK[i] = rg.Float64()
					}

					if rg.Float64() < mutationAlpha {
						children[j+1].RK[i] = rg.Float64()
					}
				}
			}
			children[j].Cost = env.Cost(children[j].RK)
			children[j+1].Cost = env.Cost(children[j+1].RK)
			bestGen = min(bestGen, children[j].Cost, children[j+1].Cost)
			// end crossover

			if children[j].Cost < bestPerson.Cost {
				rk.Copy(bestPerson, children[j])
				hasImproved = true
			}
			if children[j+1].Cost < bestPerson.Cost {
				rk.Copy(bestPerson, children[j+1])
				hasImproved = true
			}

			solutionPool.Add(bestPerson, start)
		}

		if isOdd {
			rk.Copy(children[populationSize-1], bestPerson)
		}

		k := rg.IntN(populationSize)
		local.Search(ctx, children[k])
		bestGen = min(bestGen, children[k].Cost)

		if children[k].Cost < bestPerson.Cost {
			rk.Copy(bestPerson, children[k])
			hasImproved = true
		}

		solutionPool.Add(bestPerson, start)

		elapsed := time.Since(start)
		if hasImproved {
			generationNoImprovement = 0
		} else {
			generationNoImprovement++
			ga.logger.Event(slog.Attr{Key: "stalled", Value: slog.GroupValue(slog.Int("generation", ctxAction.Iteration()), slog.Int("no-improvement", generationNoImprovement))}, elapsed)
		}
		if tuner != nil {
			reward := 1 + 1/float64(populationSize)
			if !hasImproved {
				reward = float64(bestPerson.Cost-bestGen) / math.Max(1, math.Abs(float64(bestPerson.Cost)))
			}
			tuner.Learn(reward, progress())
		}
		hasImproved = false
		ga.logger.Register(children[k].Cost, bestPerson.Cost, ctxAction.Iteration(), elapsed)

		if ctx.Err() != nil {
			break
		}

		if generationNoImprovement > maxGenerationNoImprovement {
			ga.logger.Event(slog.String("restart", "limit generations without improvement"), time.Since(start))
			restartBest := genetic.InitialPopulation(population[:populationSize], env, false, rg)
			if restartBest.Cost < bestPerson.Cost {
				rk.Copy(bestPerson, restartBest)
				solutionPool.Add(bestPerson, start)

			}
			generationNoImprovement = 0
		} else {
			population, children = children, population
		}

	}

	return bestPerson, time.Since(start)
}
