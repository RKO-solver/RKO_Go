package sa

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/random"
)

func (sa *SimulatedAnnealing) SetRG(rg *random.Generator) {
	sa.RG = rg
	sa.search.SetRG(sa.RG)
}

func (sa *SimulatedAnnealing) Name() string {
	return name
}

func (sa *SimulatedAnnealing) Solve() definition.Result {
	rko, elapsed := sa.solve(sa.solutionPool)

	return definition.Result{
		Solution:        sa.env.Decode(rko.RK),
		Cost:            rko.Cost,
		TimeSpentSecond: elapsed,
	}
}

func (sa *SimulatedAnnealing) SetTimeLimitSecond(timeLimitSecond float64) {
	sa.configuration.TimeLimitSeconds = timeLimitSecond
}

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ definition.Solver = (*SimulatedAnnealing)(nil)
