package search

import (
	"context"
	"math"
	"slices"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

func rvnd(ctx context.Context, rko *metaheuristc.RandomKeyValue, env definition.Environment, pool *solution.Pool, rg *random.Generator, neighbourhood []Type) {

	initialNeighbourhood := make([]Type, len(neighbourhood))
	copy(initialNeighbourhood, neighbourhood)

	localSolutionCost := rko.Cost
	maxIterations := int(float64(env.NumKeys()) * math.Exp(-2))

	for len(neighbourhood) > 0 {

		if ctx.Err() != nil {
			break
		}

		neighbourhoodId := rg.IntN(len(neighbourhood))

		switch neighbourhood[neighbourhoodId] {
		case Swap:
			swapSearch(ctx, rko, env)
		case Mirror:

			mirrorSearch(ctx, rko, env)
		case Farey:
			fareySearch(ctx, rko, env, rg)
		case NelderMead:
			nelderMeadSearch(ctx, rko, maxIterations, env, pool, rg)
		}
		if rko.Cost < localSolutionCost {
			localSolutionCost = rko.Cost
			neighbourhood = neighbourhood[:len(initialNeighbourhood)]
			copy(neighbourhood, initialNeighbourhood)
		} else {
			neighbourhood = slices.Delete(neighbourhood, neighbourhoodId, neighbourhoodId+1)

		}
	}
}

type rvndLocalSearch struct {
	env           definition.Environment
	rg            *random.Generator
	pool          *solution.Pool
	neighbourhood []Type
}

func (s rvndLocalSearch) String() string {
	composition := RVND.String() + ": "
	for i, neighbour := range s.neighbourhood {
		composition += neighbour.String()
		if i != len(s.neighbourhood)-1 {
			composition += ","
		}
	}
	return composition
}

func (s rvndLocalSearch) Search(ctx context.Context, rko *metaheuristc.RandomKeyValue) {
	neighbourhood := make([]Type, len(s.neighbourhood))
	copy(neighbourhood, s.neighbourhood)
	rvnd(ctx, rko, s.env, s.pool, s.rg, neighbourhood)
}

// NewRVND returns an RVND local search over neighbourhood.
func NewRVND(env definition.Environment, pool *solution.Pool, rg *random.Generator, neighbourhood []Type) Local {
	n := 0
	for _, neighbour := range neighbourhood {
		if neighbour != RVND {
			neighbourhood[n] = neighbour
			n++
		}
	}

	neighbourhood = neighbourhood[:n]

	return &rvndLocalSearch{env, rg, pool, neighbourhood}
}

var _ Local = (*rvndLocalSearch)(nil)
