package search

import (
	"context"

	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/problem"
)

func mirrorSearch(ctx context.Context, rko *metaheuristic.RandomKeyValue, env problem.Environment) {
	n := rko.RK.Len()

	for i := 0; i < n; i++ {
		if ctx.Err() != nil {
			break
		}

		rko.RK[i] = 1.0 - rko.RK[i]
		cost := env.Cost(rko.RK)
		if cost < rko.Cost {
			rko.Cost = cost
		} else {
			rko.RK[i] = 1.0 - rko.RK[i]
		}
	}
}

type mirrorLocalSearch struct {
	env problem.Environment
}

func (s mirrorLocalSearch) String() string {
	return Mirror.String()
}

func (s mirrorLocalSearch) Search(ctx context.Context, rko *metaheuristic.RandomKeyValue) {
	mirrorSearch(ctx, rko, s.env)
}

// NewMirror returns the Mirror local search.
func NewMirror(env problem.Environment) metaheuristic.LocalSearch {
	return mirrorLocalSearch{env}
}

var _ metaheuristic.LocalSearch = (*mirrorLocalSearch)(nil)
