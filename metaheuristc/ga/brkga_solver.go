package ga

import (
	"fmt"
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"time"
)

func (ga *BRKGA) solve(solutionPool *solution.Pool) (*metaheuristc.RandomKeyValue, float64) {
	configuration := ga.configuration
	env := ga.env
	rg := ga.RG
	local := ga.search
	numKeys := env.NumKeys()
	populationSize := configuration.PopulationSize

	population := make([]*metaheuristc.RandomKeyValue, populationSize)
	populationIntermediary := make([]*metaheuristc.RandomKeyValue, populationSize)
	for i := range populationSize {
		population[i] = &metaheuristc.RandomKeyValue{
			RK:   make(definition.RandomKey, numKeys),
			Cost: 0,
		}
		populationIntermediary[i] = &metaheuristc.RandomKeyValue{
			RK:   make(definition.RandomKey, numKeys),
			Cost: 0,
		}
	}

	hasImproved := false
	generationNoImprovement := 0

	start := time.Now()
	bestPerson := initialPopulation(population, env, true, rg)
	if bestPerson.Cost < solutionPool.BestSolutionCost() {
		solutionPool.AddSolution(bestPerson.Clone())
	}

	for generation := 0; generation < configuration.MaxGenerations && time.Since(start).Seconds() < configuration.TimeLimitSeconds; generation++ {
		eliteSize := int(configuration.EliteRatio * float64(populationSize))
		mutantSize := int(configuration.MutantRatio * float64(populationSize))

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
			populationIntermediary[i].Cost = population[i].Cost
		}

		copy(population, populationIntermediary)

		// apply local search in random element
		k := rg.IntN(populationSize)
		local.Search(population[k])

		metaheuristc.Sort(population)

		if population[0].Cost < bestPerson.Cost {
			bestPerson = population[0].Clone()
			hasImproved = true
		}

		bestSolutionCost := solutionPool.BestSolutionCost()
		if bestPerson.Cost < bestSolutionCost {
			solutionPool.AddSolution(bestPerson.Clone())
		}

		if hasImproved {
			generationNoImprovement = 0
		} else {
			generationNoImprovement++
		}
		hasImproved = false

		elapsed := time.Since(start).Seconds()

		ga.logger.Verbose(fmt.Sprintf("\tGenerations without improvement: %d", generationNoImprovement), elapsed)
		ga.logger.Register(bestPerson.Cost, population[0].Cost, elapsed, fmt.Sprintf("Generation %d", generation))

		if generationNoImprovement > configuration.MaxGenerationNoImprovement {
			ga.logger.Verbose("Limit Generations without improvement exceeded. Resiting population", time.Since(start).Seconds())
			bestPerson = initialPopulation(population, env, true, rg)

			if bestPerson.Cost < bestSolutionCost {
				solutionPool.AddSolution(bestPerson.Clone())
			} else {
				bestPool := solutionPool.BestSolution()
				population[populationSize-1] = bestPool
				metaheuristc.Sort(population)
			}
		}

	}

	return bestPerson, time.Since(start).Seconds()
}
