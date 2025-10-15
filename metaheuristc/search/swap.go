package search

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/random"
)

func swapSearch(rko *metaheuristc.RandomKeyValue, environment definition.Environment) {
	for _, n := range environment.SwapSearch() {
		if len(n) < 2 {
			continue
		}
		start := n[0]
		end := n[1]
		for i := start; i < end-1; i++ {
			for j := i + 1; j < end; j++ {
				rko.RK[i], rko.RK[j] = rko.RK[j], rko.RK[i]
				cost := environment.Cost(rko.RK)
				if cost < rko.Cost {
					// maintain best solution
					rko.Cost = cost
				} else {
					// return to the best solution
					rko.RK[i], rko.RK[j] = rko.RK[j], rko.RK[i]
				}
			}
		}
	}

}

type swapLocalSearch struct {
	environment definition.Environment
}

func (s swapLocalSearch) String() string {
	return GetSearchString(Swap)
}

func (s swapLocalSearch) Search(rko *metaheuristc.RandomKeyValue) {
	swapSearch(rko, s.environment)
}

func (s swapLocalSearch) SetRG(rg *random.Generator) {}

func CreateSwapLocalSearch(environment definition.Environment) Local {
	return swapLocalSearch{environment}
}

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ Local = (*swapLocalSearch)(nil)
