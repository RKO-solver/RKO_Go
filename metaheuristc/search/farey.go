package search

import (
	"context"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/random"
)

var fareySequence = []float64{
	0.00,
	0.142857,
	0.166667,
	0.20,
	0.25,
	0.285714,
	0.333333,
	0.40,
	0.428571,
	0.50,
	0.571429,
	0.60,
	0.666667,
	0.714286,
	0.75,
	0.80,
	0.833333,
	0.857143,
	1.0,
}

func fareySearch(ctx context.Context, rko *metaheuristc.RandomKeyValue, environment definition.Environment, rg *random.Generator) {

	fareyLen := len(fareySequence)
	n := rko.RK.Len()

	for i := 0; i < n; i++ {
		if ctx.Err() != nil {
			break
		}

		for j := 0; j < fareyLen-1; j++ {
			oldValue := rko.RK[i]
			value := rg.RangeFloat64(fareySequence[j], fareySequence[j+1])
			rko.RK[i] = value
			cost := environment.Cost(rko.RK)
			if cost < rko.Cost {
				rko.Cost = cost
			} else {
				rko.RK[i] = oldValue
			}
		}
	}

}

type fareyLocalSearch struct {
	environment definition.Environment
	rg          *random.Generator
}

func (s fareyLocalSearch) String() string {
	return GetSearchString(Farey)
}

func (s fareyLocalSearch) Search(ctx context.Context, rko *metaheuristc.RandomKeyValue) {
	fareySearch(ctx, rko, s.environment, s.rg)
}

// CreateFareyLocalSearch returns the Farey local search.
func CreateFareyLocalSearch(environment definition.Environment, rg *random.Generator) Local {
	return &fareyLocalSearch{environment, rg}
}

var _ Local = (*fareyLocalSearch)(nil)
