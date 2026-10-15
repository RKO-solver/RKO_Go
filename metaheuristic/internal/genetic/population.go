// Package genetic holds helpers and defaults shared by the genetic algorithms.
package genetic

import (
	"math"

	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/rk"
	"github.com/RKO-solver/rko-go/problem"
	"github.com/RKO-solver/rko-go/random"
)

// InitialPopulation resets and evaluates the population and returns a clone of its best individual.
func InitialPopulation(population []*metaheuristic.RandomKeyValue, env problem.Environment, sorted bool, rg *random.Generator) *metaheuristic.RandomKeyValue {

	populationSize := len(population)

	bestPerson := &metaheuristic.RandomKeyValue{
		RK:   nil,
		Cost: math.MaxInt,
	}

	for i := 0; i < populationSize; i++ {
		rk.Reset(population[i].RK, rg)
		population[i].Cost = env.Cost(population[i].RK)

		if !sorted && population[i].Cost < bestPerson.Cost {
			bestPerson = population[i]
		}
	}

	if sorted {
		metaheuristic.Sort(population)
		bestPerson = population[0]
	}

	return bestPerson.Clone()
}
