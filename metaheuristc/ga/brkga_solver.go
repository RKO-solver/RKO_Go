package ga

import (
	"context"
	"fmt"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
)

func (ga *BRKGA) solve(ctx context.Context, solutionPool *solution.Pool) (*metaheuristc.RandomKeyValue, float64) {
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
	solutionPool.AddSolution(bestPerson, time.Since(start).Seconds())

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

			solutionPool.AddSolution(bestPerson, time.Since(start).Seconds())

			break
		}

		if population[0].Cost < bestPerson.Cost {
			bestPerson = population[0].Clone()
			generationNoImprovement = 0
		} else {
			generationNoImprovement++
		}

		solutionPool.AddSolution(bestPerson, time.Since(start).Seconds())

		elapsed := time.Since(start).Seconds()

		ga.logger.Register(population[0].Cost, bestPerson.Cost, elapsed, fmt.Sprintf("Generation %d, without improvement %d", ctxAction.Iteration(), generationNoImprovement))

		if generationNoImprovement > configuration.MaxGenerationNoImprovement {
			ga.logger.Verbose("Limit Generations without improvement exceeded. Resetting population", time.Since(start).Seconds())
			restartBest := initialPopulation(population, env, true, rg)
			if restartBest.Cost < bestPerson.Cost {
				bestPerson = restartBest

			}
			solutionPool.AddSolution(bestPerson, time.Since(start).Seconds())

			generationNoImprovement = 0
		}

	}

	return bestPerson, time.Since(start).Seconds()
}
