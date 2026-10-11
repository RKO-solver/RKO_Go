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

func (ga *BRKGA) solve(ctx context.Context, solutionPool *solution.Pool) (*metaheuristc.RandomKeyValue, time.Duration) {
	ctx, ctxAction := definition.GetMHContext(ctx)
	defer ctxAction.Stop()

	configuration := ga.configuration
	env := ga.env
	rg := ga.RG
	local := ga.search
	numKeys := env.NumKeys()
	populationSize := configuration.PopulationSize

	population := make([]*metaheuristc.RandomKeyValue, populationSize)
	populationIntermediary := make([]*metaheuristc.RandomKeyValue, populationSize)
	for i := range populationSize {
		population[i] = metaheuristc.CreateEmptyRandomKeyValue(numKeys)
		populationIntermediary[i] = metaheuristc.CreateEmptyRandomKeyValue(numKeys)
	}

	generationNoImprovement := 0

	start := definition.StartTime(ctx)
	bestPerson := initialPopulation(population, env, true, rg)
	solutionPool.AddSolution(bestPerson, time.Since(start))

	eliteSize := int(configuration.EliteRatio * float64(populationSize))
	mutantSize := int(configuration.MutantRatio * float64(populationSize))
	for ; ctx.Err() == nil; ctxAction.Step() {
		for i := range eliteSize {
			copy(populationIntermediary[i].RK, population[i].RK)
			populationIntermediary[i].Cost = population[i].Cost
		}

		for i := eliteSize; i < (populationSize - mutantSize); i++ {
			populationIntermediary[i].RK = crossoverUniformElite(population, eliteSize, populationSize, configuration.MutationAlpha, configuration.CrossoverAlpha, rg)
			populationIntermediary[i].Cost = env.Cost(populationIntermediary[i].RK)

		}

		for i := populationSize - mutantSize; i < populationSize; i++ {
			rk.Reset(populationIntermediary[i].RK, rg)
			populationIntermediary[i].Cost = env.Cost(populationIntermediary[i].RK)
		}

		population, populationIntermediary = populationIntermediary, population

		k := rg.IntN(populationSize)
		local.Search(ctx, population[k])

		metaheuristc.Sort(population)

		if ctx.Err() != nil {
			if population[0].Cost < bestPerson.Cost {
				bestPerson = population[0].Clone()
			}

			solutionPool.AddSolution(bestPerson, time.Since(start))

			break
		}

		elapsed := time.Since(start)
		if population[0].Cost < bestPerson.Cost {
			bestPerson = population[0].Clone()
			generationNoImprovement = 0
		} else {
			generationNoImprovement++
			ga.logger.Event(slog.Attr{Key: "stalled", Value: slog.GroupValue(slog.Int("generation", ctxAction.Iteration()), slog.Int("no-improvement", generationNoImprovement))}, elapsed)

		}
		solutionPool.AddSolution(bestPerson, time.Since(start))
		ga.logger.Register(population[0].Cost, bestPerson.Cost, ctxAction.Iteration(), elapsed)

		if generationNoImprovement > configuration.MaxGenerationNoImprovement {
			ga.logger.Event(slog.String("restart", "limit generations without improvement"), time.Since(start))
			restartBest := initialPopulation(population, env, true, rg)
			if restartBest.Cost < bestPerson.Cost {
				bestPerson = restartBest

			}
			solutionPool.AddSolution(bestPerson, time.Since(start))

			generationNoImprovement = 0
		}

	}

	return bestPerson, time.Since(start)
}
