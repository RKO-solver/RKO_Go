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
	ctx, ctxAction := definition.GetMHContext(ctx)
	defer ctxAction.Stop()

	configuration := ga.configuration
	env := ga.env
	rg := ga.RG
	local := ga.search
	numKeys := env.NumKeys()

	children := make([]*metaheuristc.RandomKeyValue, configuration.PopulationSize)
	population := make([]*metaheuristc.RandomKeyValue, configuration.PopulationSize)
	for i := range configuration.PopulationSize {
		population[i] = metaheuristc.CreateEmptyRandomKeyValue(numKeys)
		children[i] = metaheuristc.CreateEmptyRandomKeyValue(numKeys)
	}

	hasImproved := false
	isOdd := configuration.PopulationSize%2 == 1
	generationNoImprovement := 0

	start := definition.StartTime(ctx)
	bestPerson := initialPopulation(population, env, false, rg)
	solutionPool.AddSolution(bestPerson, time.Since(start).Seconds())

	for ; ctx.Err() == nil; ctxAction.Step() {

		for j := 0; j < configuration.PopulationSize-1; j += 2 {
			offspring1 := tournament(population, rg)
			offspring2 := tournament(population, rg)
			copy(children[j].RK, offspring1.RK)
			copy(children[j+1].RK, offspring2.RK)

			crossover(children[j], children[j+1], configuration.CrossoverAlpha, configuration.MutationAlpha, rg)
			children[j].Cost = env.Cost(children[j].RK)
			children[j+1].Cost = env.Cost(children[j+1].RK)

			if children[j].Cost < bestPerson.Cost {
				bestPerson = children[j].Clone()
				hasImproved = true
			}
			if children[j+1].Cost < bestPerson.Cost {
				bestPerson = children[j+1].Clone()
				hasImproved = true
			}

			solutionPool.AddSolution(bestPerson, time.Since(start).Seconds())
		}

		if isOdd {
			copy(children[configuration.PopulationSize-1].RK, bestPerson.RK)
			children[configuration.PopulationSize-1].Cost = bestPerson.Cost
		}

		k := rg.IntN(configuration.PopulationSize)
		local.Search(ctx, children[k])

		if children[k].Cost < bestPerson.Cost {
			bestPerson = children[k].Clone()
			hasImproved = true
		}

		solutionPool.AddSolution(bestPerson, time.Since(start).Seconds())

		if hasImproved {
			generationNoImprovement = 0
		} else {
			generationNoImprovement++
		}
		hasImproved = false

		elapsed := time.Since(start).Seconds()

		ga.logger.Register(children[k].Cost, bestPerson.Cost, elapsed, fmt.Sprintf("Generation %d, without improvement %d", ctxAction.Iteration(), generationNoImprovement))

		if ctx.Err() != nil {
			break
		}

		if generationNoImprovement > configuration.MaxGenerationNoImprovement {
			ga.logger.Verbose("Limit Generations without improvement exceeded. Resetting population", time.Since(start).Seconds())
			restartBest := initialPopulation(population, env, false, rg)
			if restartBest.Cost < bestPerson.Cost {
				bestPerson = restartBest
				solutionPool.AddSolution(bestPerson, elapsed)

			}
			generationNoImprovement = 0
		} else {
			population, children = children, population
		}

	}

	return bestPerson, time.Since(start).Seconds()
}
