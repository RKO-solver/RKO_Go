// Package genetic holds helpers shared by the genetic algorithms.
package genetic

import (
	"math"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/random"
)

// InitialPopulation resets and evaluates the population and returns a clone of its best individual.
func InitialPopulation(population []*metaheuristc.RandomKeyValue, env definition.Environment, sorted bool, rg *random.Generator) *metaheuristc.RandomKeyValue {

	populationSize := len(population)

	bestPerson := &metaheuristc.RandomKeyValue{
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
		metaheuristc.Sort(population)
		bestPerson = population[0]
	}

	return bestPerson.Clone()
}
