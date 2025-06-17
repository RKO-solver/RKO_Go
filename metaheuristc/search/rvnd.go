package search

import (
	"github.com/lucasmends/rko-go/definition"
	"github.com/lucasmends/rko-go/metaheuristc"
	"github.com/lucasmends/rko-go/random"
	"slices"
)

func rvnd(rko *metaheuristc.RandomKeyValue, environment definition.Environment, r *random.Generator, neighbourhood []Type) {

	localSolutionCost := rko.Cost

	for len(neighbourhood) > 0 {
		neighbourhoodId := r.IntN(len(neighbourhood))

		switch neighbourhood[neighbourhoodId] {
		case Swap:
			swapSearch(rko, environment)
		case Mirror:

			mirrorSearch(rko, environment)
		case Farey:
			fareySearch(rko, environment, r)
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
	neighbourhood []Type
}

func (s rvndseach) Search(rko *metaheuristc.RandomKeyValue) {
	neighbourhood := make([]Type, len(s.neighbourhood))
	copy(neighbourhood, s.neighbourhood)
	rvnd(rko, s.environment, s.rg, neighbourhood)
}

func (s rvndseach) SetRG(rg *random.Generator) {
	s.rg = rg
}

func CreateRVND(environment definition.Environment, rg *random.Generator, neighbourhood []Type) Local {
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

	return rvndseach{environment, rg, neighbourhood}
}
