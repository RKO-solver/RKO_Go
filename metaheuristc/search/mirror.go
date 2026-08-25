package search

import (
	"context"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/random"
)

func mirrorSearch(ctx context.Context, rko *metaheuristc.RandomKeyValue, environment definition.Environment) {
	n := rko.RK.Len()

	for i := 0; i < n; i++ {
		if ctx.Err() != nil {
			break
		}

		rko.RK[i] = 1.0 - rko.RK[i]
		cost := environment.Cost(rko.RK)
		if cost < rko.Cost {
			// mantain best solution
			rko.Cost = cost
		} else {
			// return to the best solution
			rko.RK[i] = 1.0 - rko.RK[i]
		}
	}
}

type mirrorLocalSearch struct {
	environment definition.Environment
}

func (s mirrorLocalSearch) String() string {
	return GetSearchString(Mirror)
}

func (s mirrorLocalSearch) Search(ctx context.Context, rko *metaheuristc.RandomKeyValue) {
	mirrorSearch(ctx, rko, s.environment)
}

func (s mirrorLocalSearch) SetRG(rg *random.Generator) {}

func CreateMirrorLocalSearch(environment definition.Environment) Local {
	return mirrorLocalSearch{environment}
}

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ Local = (*mirrorLocalSearch)(nil)
