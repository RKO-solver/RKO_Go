package ga

import (
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/random"
)

func crossover(offspring1, offspring2 *metaheuristc.RandomKeyValue, probabilityCrossover float64, probabilityMutation float64, rg *random.Generator) {

	var key1, key2 float64

	if rg.Float64() < probabilityCrossover {
		for i := range offspring1.RK {
			key1 = offspring1.RK[i]
			key2 = offspring2.RK[i]
			if rg.Float64() < 0.5 {
				offspring1.RK[i] = key2
			}

			if rg.Float64() < 0.5 {
				offspring2.RK[i] = key1
			}

			if rg.Float64() < probabilityMutation {
				offspring1.RK[i] = rg.Float64()
			}

			if rg.Float64() < probabilityMutation {
				offspring2.RK[i] = rg.Float64()
			}
		}
	}

}
