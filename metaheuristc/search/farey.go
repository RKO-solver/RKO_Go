package search

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/random"
)

func fareySearch(rko *metaheuristc.RandomKeyValue, environment definition.Environment, rg *random.Generator) {
	fareySequence := []float64{
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
	fareyLen := len(fareySequence)
	n := rko.RK.Len()

	for i := 0; i < n; i++ {
		for j := 1; j < fareyLen-1; j++ {
			oldValue := rko.RK[i]
			value := rg.RangeFloat64(fareySequence[j], fareySequence[j+1])
			rko.RK[i] = value
			cost := environment.Cost(rko.RK)
			if cost < rko.Cost {
				// mantain best solution
				rko.Cost = cost
			} else {
				// return to the best solution
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

func (s fareyLocalSearch) Search(rko *metaheuristc.RandomKeyValue) {
	fareySearch(rko, s.environment, s.rg)
}

func (s fareyLocalSearch) SetRG(rg *random.Generator) {
	s.rg = rg
}

func CreateFareyLocalSearch(environment definition.Environment, rg *random.Generator) Local {
	return fareyLocalSearch{environment, rg}
}

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ Local = (*fareyLocalSearch)(nil)
