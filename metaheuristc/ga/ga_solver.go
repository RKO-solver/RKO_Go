package ga

import (
	"fmt"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
)

func (ga *GA) solve(solutionPool *solution.Pool) (*metaheuristc.RandomKeyValue, float64) {
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

	start := time.Now()
	bestPerson := initialPopulation(population, env, false, rg)
	if bestPerson.Cost < solutionPool.BestSolutionCost() {
		solutionPool.AddSolution(bestPerson.Clone())
	}

	for generation := 0; generation < configuration.MaxGenerations && time.Since(start).Seconds() < configuration.TimeLimitSeconds; generation++ {
		for j := 0; j < configuration.PopulationSize-1; j++ {
			offspring1 := tournament(population, rg)
			copy(children[j].RK, offspring1.RK)
			children[j].Cost = offspring1.Cost
			offspring1 = children[j]

			offspring2 := tournament(population, rg)
			copy(children[j+1].RK, offspring2.RK)
			children[j+1].Cost = offspring2.Cost
			offspring2 = children[j+1]

			crossover(offspring1, offspring2, configuration.CrossoverAlpha, configuration.MutationAlpha, rg)

			if offspring1.Cost < bestPerson.Cost {
				bestPerson = offspring1.Clone()
				hasImproved = true
			}

			if offspring2.Cost < bestPerson.Cost {
				bestPerson = offspring2.Clone()
				hasImproved = true
			}

			bestSolutionCost := solutionPool.BestSolutionCost()
			if bestPerson.Cost < bestSolutionCost {
				solutionPool.AddSolution(bestPerson.Clone())
			}
		}

		// apply local search in random element
		k := rg.IntN(configuration.PopulationSize)
		local.Search(children[k])

		if children[k].Cost < bestPerson.Cost {
			bestPerson = children[k].Clone()
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

		ga.logger.Verbose(fmt.Sprintf("Generations without improvement: %d", generationNoImprovement), elapsed)
		ga.logger.Register(children[k].Cost, bestPerson.Cost, elapsed, fmt.Sprintf("Generation %d", generation))

		if generationNoImprovement > configuration.MaxGenerationNoImprovement {
			ga.logger.Verbose("Limit Generations without improvement exceeded. Resiting population", time.Since(start).Seconds())
			bestPerson = initialPopulation(population, env, false, rg)
			if bestPerson.Cost < bestSolutionCost {
				solutionPool.AddSolution(bestPerson.Clone())
			}
		} else {
			copy(population, children)
		}

	}

	return bestPerson, time.Since(start).Seconds()
}
