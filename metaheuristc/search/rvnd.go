package search

import (
	"math"
	"slices"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

func rvnd(rko *metaheuristc.RandomKeyValue, environment definition.Environment, s *solution.Pool, r *random.Generator, neighbourhood []Type) {

	localSolutionCost := rko.Cost
	maxIterations := int(float64(environment.NumKeys()) * math.Exp(-2))

	for len(neighbourhood) > 0 {
		neighbourhoodId := r.IntN(len(neighbourhood))

		switch neighbourhood[neighbourhoodId] {
		case Swap:
			swapSearch(rko, environment)
		case Mirror:

			mirrorSearch(rko, environment)
		case Farey:
			fareySearch(rko, environment, r)
		case Nelder:
			nelderMeadSearch(rko, maxIterations, environment, s, r)
		}
		// there was improvement
		if localSolutionCost < rko.Cost {
			localSolutionCost = rko.Cost
		} else {
			// there wasn't improvement
			// remove neighborhood
			neighbourhood = slices.Delete(neighbourhood, neighbourhoodId, neighbourhoodId+1)
		}
	}
}

type rvndseach struct {
	environment   definition.Environment
	rg            *random.Generator
	s             *solution.Pool
	neighbourhood []Type
}

func (s rvndseach) String() string {
	composition := GetSearchString(RVND) + ": "
	for i, neighbour := range s.neighbourhood {
		composition += GetSearchString(neighbour)
		if i != len(s.neighbourhood)-1 {
			composition += ","
		}
	}
	return composition
}

func (s rvndseach) Search(rko *metaheuristc.RandomKeyValue) {
	neighbourhood := make([]Type, len(s.neighbourhood))
	copy(neighbourhood, s.neighbourhood)
	rvnd(rko, s.environment, s.s, s.rg, neighbourhood)
}

func (s rvndseach) SetRG(rg *random.Generator) {
	s.rg = rg
}

func CreateRVND(environment definition.Environment, s *solution.Pool, rg *random.Generator, neighbourhood []Type) Local {
	n := 0
	for _, neighboor := range neighbourhood {
		// Filter, Keep only
		if neighboor != RVND {
			// Move the kept element to the front of the slice.
			neighbourhood[n] = neighboor
			n++
		}
	}

	// Truncate the slice. This is the crucial step.
	// It updates the slice's length to only include the elements we kept.
	neighbourhood = neighbourhood[:n]

	return rvndseach{environment, rg, s, neighbourhood}
}

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ Local = (*rvndseach)(nil)
