package search

import (
	"context"

	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/problem"
)

func swapSearch(ctx context.Context, rko *metaheuristic.RandomKeyValue, env problem.Environment) {
	for _, n := range env.SwapSearch() {
		if ctx.Err() != nil {
			break
		}

		start := n[0]
		end := n[1]
		for i := start; i < end-1 && ctx.Err() == nil; i++ {
			for j := i + 1; j < end; j++ {
				rko.RK[i], rko.RK[j] = rko.RK[j], rko.RK[i]
				cost := env.Cost(rko.RK)
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
	env problem.Environment
}

func (s swapLocalSearch) String() string {
	return Swap.String()
}

func (s swapLocalSearch) Search(ctx context.Context, rko *metaheuristic.RandomKeyValue) {
	swapSearch(ctx, rko, s.env)
}

// NewSwap returns the Swap local search.
func NewSwap(env problem.Environment) metaheuristic.LocalSearch {
	return swapLocalSearch{env}
}

var _ metaheuristic.LocalSearch = (*swapLocalSearch)(nil)
