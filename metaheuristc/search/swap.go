package search

import (
	"context"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
)

func swapSearch(ctx context.Context, rko *metaheuristc.RandomKeyValue, environment definition.Environment) {
	for _, n := range environment.SwapSearch() {
		if ctx.Err() != nil {
			break
		}

		start := n[0]
		end := n[1]
		for i := start; i < end-1 && ctx.Err() == nil; i++ {
			for j := i + 1; j < end; j++ {
				rko.RK[i], rko.RK[j] = rko.RK[j], rko.RK[i]
				cost := environment.Cost(rko.RK)
				if cost < rko.Cost {
					rko.Cost = cost
				} else {
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

func (s swapLocalSearch) Search(ctx context.Context, rko *metaheuristc.RandomKeyValue) {
	swapSearch(ctx, rko, s.environment)
}

// CreateSwapLocalSearch returns the Swap local search.
func CreateSwapLocalSearch(environment definition.Environment) Local {
	return swapLocalSearch{environment}
}

var _ Local = (*swapLocalSearch)(nil)
