package ga

import (
	"context"
	"fmt"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
)

func (ga *GA) solve(ctx context.Context, solutionPool *solution.Pool) (*metaheuristc.RandomKeyValue, float64) {
	configuration := ga.configuration
	env := ga.env
	rg := ga.RG
	local := ga.search
	numKeys := env.NumKeys()

	children := make([]*metaheuristc.RandomKeyValue, configuration.PopulationSize)
	population := make([]*metaheuristc.RandomKeyValue, configuration.PopulationSize)
	for i := range configuration.PopulationSize {
		population[i] = &metaheuristc.RandomKeyValue{
			RK:   make(definition.RandomKey, numKeys),
			Cost: 0,
		}
		children[i] = &metaheuristc.RandomKeyValue{
			RK:   make(definition.RandomKey, numKeys),
			Cost: 0,
		}
	}

	hasImproved := false
	generationNoImprovement := 0

	start := definition.StartTime(ctx)
	bestPerson := initialPopulation(population, env, false, rg)
	if bestPerson.Cost < solutionPool.BestSolutionCost() {
		solutionPool.AddSolution(bestPerson.Clone(), time.Since(start).Seconds())
	}

	for generation := 0; generation < configuration.MaxGenerations && ctx.Err() == nil; generation++ {

		for j := 0; j < configuration.PopulationSize; j++ {
			offspring1 := tournament(population, rg)
			offspring2 := tournament(population, rg)

			crossover(children[j], offspring1, offspring2, configuration.CrossoverAlpha, configuration.MutationAlpha, rg)
			children[j].Cost = env.Cost(children[j].RK)

			if children[j].Cost < bestPerson.Cost {
				bestPerson = children[j].Clone()
				hasImproved = true
			}

			bestSolutionCost := solutionPool.BestSolutionCost()
			if bestPerson.Cost < bestSolutionCost {
				solutionPool.AddSolution(bestPerson.Clone(), time.Since(start).Seconds())
			}
		}

		// apply local search in random element
		k := rg.IntN(configuration.PopulationSize)
		local.Search(ctx, children[k])

		if children[k].Cost < bestPerson.Cost {
			bestPerson = children[k].Clone()
			hasImproved = true
		}

		bestSolutionCost := solutionPool.BestSolutionCost()
		if bestPerson.Cost < bestSolutionCost {
			solutionPool.AddSolution(bestPerson.Clone(), time.Since(start).Seconds())
		}

		if ctx.Err() != nil {
			break
		}

		if hasImproved {
			generationNoImprovement = 0
		} else {
			generationNoImprovement++
		}
		hasImproved = false

		elapsed := time.Since(start).Seconds()

		ga.logger.Verbose(fmt.Sprintf("Generations without improvement: %d", generationNoImprovement), elapsed)
		ga.logger.Register(children[k].Cost, bestPerson.Cost, elapsed, fmt.Sprintf("Generation %d", generation))

		if generationNoImprovement > configuration.MaxGenerationNoImprovement {
			ga.logger.Verbose("Limit Generations without improvement exceeded. Resiting population", time.Since(start).Seconds())
			restartBest := initialPopulation(population, env, false, rg)
			if restartBest.Cost < bestPerson.Cost {
				bestPerson = restartBest

			}
			if restartBest.Cost < bestSolutionCost {
				solutionPool.AddSolution(bestPerson.Clone(), elapsed)
			}
			generationNoImprovement = 0
		} else {
			population, children = children, population
		}
	}

	return bestPerson, time.Since(start).Seconds()
}
