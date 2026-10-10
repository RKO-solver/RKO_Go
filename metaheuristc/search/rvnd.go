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
		// there was improvement
		if rko.Cost < localSolutionCost {
			localSolutionCost = rko.Cost
		} else {
			// there wasn't improvement
			// remove neighborhood
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

func (s *rvndSearch) SetRG(rg *random.Generator) {
	s.rg = rg
}

func CreateRVND(environment definition.Environment, s *solution.Pool, rg *random.Generator, neighbourhood []Type) Local {
	n := 0
	for _, neighbour := range neighbourhood {
		// Filter, Keep only
		if neighbour != RVND {
			// Move the kept element to the front of the slice.
			neighbourhood[n] = neighbour
			n++
		}
	}

	// Truncate the slice. This is the crucial step.
	// It updates the slice's length to only include the elements we kept.
	neighbourhood = neighbourhood[:n]

	return &rvndSearch{environment, rg, s, neighbourhood}
}

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ Local = (*rvndSearch)(nil)
