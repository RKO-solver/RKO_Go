package ga

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/random"
)

func crossover(child, offspring1, offspring2 *metaheuristc.RandomKeyValue, probabilityCrossover float64, probabilityMutation float64, rg *random.Generator) {
	copy(child.RK, offspring1.RK)

	if rg.Float64() < probabilityCrossover {
		for i := range child.RK {
			if rg.Float64() < 0.5 {
				child.RK[i] = offspring2.RK[i]
			}
		}
	}

	for i := range child.RK {
		if rg.Float64() < probabilityMutation {
			child.RK[i] = rg.Float64()
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
