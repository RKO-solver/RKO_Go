package ga

import (
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/random"
)

func tournament(population []*metaheuristic.RandomKeyValue, rg *random.Generator) *metaheuristic.RandomKeyValue {
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
