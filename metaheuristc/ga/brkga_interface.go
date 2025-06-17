package ga

import (
	"github.com/lucasmends/rko-go/definition"
	"github.com/lucasmends/rko-go/random"
)

func (ga *BRKGA) SetRG(rg *random.Generator) {
	ga.RG = rg
	ga.search.SetRG(ga.RG)
}

func (ga *BRKGA) Name() string {
	return nameBRKGA
}

func (ga *BRKGA) SetIdWorker(id int) {
	if ga.logger != nil {
		ga.logger.SetIdWorker(id)
	}
}

func (ga *BRKGA) Solve() definition.Result {
	rko, elapsed := ga.solve(ga.solutionPool)

	return definition.Result{
		Solution:        ga.env.Decode(rko.RK),
		Cost:            rko.Cost,
		TimeSpentSecond: elapsed,
	}
}

func (ga *BRKGA) SetTimeLimitSecond(timeLimitSecond float64) {
	ga.configuration.TimeLimitSeconds = timeLimitSecond
}

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ definition.Solver = (*BRKGA)(nil)
