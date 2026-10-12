package ga

import (
	"context"
	"log/slog"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
)

func (ga *GA) solve(ctx context.Context, solutionPool *solution.Pool) (*metaheuristc.RandomKeyValue, time.Duration) {
	ctx, ctxAction := definition.NewSolverContext(ctx)
	defer ctxAction.Stop()

	configuration := ga.configuration
	env := ga.env
	rg := ga.rg
	local := ga.search
	numKeys := env.NumKeys()

	children := make([]*metaheuristc.RandomKeyValue, configuration.PopulationSize)
	population := make([]*metaheuristc.RandomKeyValue, configuration.PopulationSize)
	for i := range configuration.PopulationSize {
		population[i] = metaheuristc.NewRandomKeyValue(numKeys)
		children[i] = metaheuristc.NewRandomKeyValue(numKeys)
	}

	hasImproved := false
	isOdd := configuration.PopulationSize%2 == 1
	generationNoImprovement := 0

	start := definition.StartTime(ctx)
	bestPerson := initialPopulation(population, env, false, rg)
	solutionPool.Add(bestPerson, start)

	var key1, key2 float64

	for ; ctx.Err() == nil; ctxAction.Step() {

		for j := 0; j < configuration.PopulationSize-1; j += 2 {
			offspring1 := tournament(population, rg)
			offspring2 := tournament(population, rg)
			copy(children[j].RK, offspring1.RK)
			copy(children[j+1].RK, offspring2.RK)

			// crossover
			if rg.Float64() < configuration.CrossoverAlpha {
				for i := range numKeys {
					key1 = children[j].RK[i]
					key2 = children[j+1].RK[i]
					if rg.Float64() < 0.5 {
						children[j].RK[i] = key2
					}

					if rg.Float64() < 0.5 {
						children[j+1].RK[i] = key1
					}

					if rg.Float64() < configuration.MutationAlpha {
						children[j].RK[i] = rg.Float64()
					}

					if rg.Float64() < configuration.MutationAlpha {
						children[j+1].RK[i] = rg.Float64()
					}
				}
			}
			children[j].Cost = env.Cost(children[j].RK)
			children[j+1].Cost = env.Cost(children[j+1].RK)
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
			rk.Copy(children[configuration.PopulationSize-1], bestPerson)
		}

		k := rg.IntN(configuration.PopulationSize)
		local.Search(ctx, children[k])

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
		hasImproved = false
		ga.logger.Register(children[k].Cost, bestPerson.Cost, ctxAction.Iteration(), elapsed)

		if ctx.Err() != nil {
			break
		}

		if generationNoImprovement > configuration.MaxGenerationNoImprovement {
			ga.logger.Event(slog.String("restart", "limit generations without improvement"), time.Since(start))
			restartBest := initialPopulation(population, env, false, rg)
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
