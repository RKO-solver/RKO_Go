package search

import (
	"context"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristic"
)

func mirrorSearch(ctx context.Context, rko *metaheuristic.RandomKeyValue, env definition.Environment) {
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
	env definition.Environment
}

func (s mirrorLocalSearch) String() string {
	return Mirror.String()
}

func (s mirrorLocalSearch) Search(ctx context.Context, rko *metaheuristic.RandomKeyValue) {
	mirrorSearch(ctx, rko, s.env)
}

// NewMirror returns the Mirror local search.
func NewMirror(env definition.Environment) Local {
	return mirrorLocalSearch{env}
}

var _ Local = (*mirrorLocalSearch)(nil)
