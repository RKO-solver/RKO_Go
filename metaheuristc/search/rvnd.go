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

func rvnd(ctx context.Context, rko *metaheuristc.RandomKeyValue, environment definition.Environment, s *solution.Pool, r *random.Generator, neighbourhood []Type) {

	initialNeighbourhood := make([]Type, len(neighbourhood))
	copy(initialNeighbourhood, neighbourhood)

	localSolutionCost := rko.Cost
	maxIterations := int(float64(environment.NumKeys()) * math.Exp(-2))

	for len(neighbourhood) > 0 {

		if ctx.Err() != nil {
			break
		}

		neighbourhoodId := r.IntN(len(neighbourhood))

		switch neighbourhood[neighbourhoodId] {
		case Swap:
			swapSearch(ctx, rko, environment)
		case Mirror:

			mirrorSearch(ctx, rko, environment)
		case Farey:
			fareySearch(ctx, rko, environment, r)
		case Nelder:
			nelderMeadSearch(ctx, rko, maxIterations, environment, s, r)
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

type rvndSearch struct {
	environment   definition.Environment
	rg            *random.Generator
	s             *solution.Pool
	neighbourhood []Type
}

func (s rvndSearch) String() string {
	composition := GetSearchString(RVND) + ": "
	for i, neighbour := range s.neighbourhood {
		composition += GetSearchString(neighbour)
		if i != len(s.neighbourhood)-1 {
			composition += ","
		}
	}
	return composition
}

func (s rvndSearch) Search(ctx context.Context, rko *metaheuristc.RandomKeyValue) {
	neighbourhood := make([]Type, len(s.neighbourhood))
	copy(neighbourhood, s.neighbourhood)
	rvnd(ctx, rko, s.environment, s.s, s.rg, neighbourhood)
}

// CreateRVND returns an RVND local search over neighbourhood.
func CreateRVND(environment definition.Environment, s *solution.Pool, rg *random.Generator, neighbourhood []Type) Local {
	n := 0
	for _, neighbour := range neighbourhood {
		if neighbour != RVND {
			neighbourhood[n] = neighbour
			n++
		}
	}

	neighbourhood = neighbourhood[:n]

	return &rvndSearch{environment, rg, s, neighbourhood}
}

var _ Local = (*rvndSearch)(nil)
