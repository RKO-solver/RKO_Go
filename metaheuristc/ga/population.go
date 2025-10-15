package ga

import (
	"math"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/random"
)

func initialPopulation(population []*metaheuristc.RandomKeyValue, env definition.Environment, sorted bool, rg *random.Generator) *metaheuristc.RandomKeyValue {

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

	return bestPerson
}

func tournament(population []*metaheuristc.RandomKeyValue, rg *random.Generator) *metaheuristc.RandomKeyValue {
	first := population[rg.IntN(len(population))]
	second := population[rg.IntN(len(population))]
	third := population[rg.IntN(len(population))]

	if first.Cost < second.Cost && first.Cost < third.Cost {
		return first
	}
	if second.Cost < third.Cost {
		return second
	}
	return third
}
