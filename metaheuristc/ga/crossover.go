package ga

import (
	"github.com/lucasmends/rko-go/definition"
	"github.com/lucasmends/rko-go/metaheuristc"
	"github.com/lucasmends/rko-go/random"
)

func crossover(offspring1 *metaheuristc.RandomKeyValue, offspring2 *metaheuristc.RandomKeyValue, probabilityCrossover float64, probabilityMutation float64, rg *random.Generator) {

	if rg.Float64() < probabilityCrossover {
		n := offspring1.RK.Len()

		for i := 0; i < n; i++ {
			// swap alleles of offspring1
			if rg.Float64() < 0.5 {
				offspring1.RK[i] = offspring2.RK[i]
			} else {
				offspring2.RK[i] = offspring1.RK[i]
			}

			// mutation
			if rg.Float64() < probabilityMutation {
				offspring1.RK[i] = rg.Float64()
			}

			if rg.Float64() < probabilityMutation {
				offspring2.RK[i] = rg.Float64()
			}
		}
	}
}

func crossoverUniformElite(population []*metaheuristc.RandomKeyValue, eliteSize, populationSize int, alphaMutation, alphaCrossover float64, rg *random.Generator) definition.RandomKey {

	eliteParentId := rg.IntN(eliteSize)
	nonEliteParentId := rg.RangeInt(eliteSize, populationSize)

	childRK := population[eliteParentId].RK.Clone()
	keySize := childRK.Len()

	for i := 0; i < keySize; i++ {
		if rg.Float64() < alphaMutation {
			childRK[i] = rg.Float64()
		} else {
			if rg.Float64() >= alphaCrossover {
				childRK[i] = population[nonEliteParentId].RK[i]
			}
		}
	}

	return childRK
}
