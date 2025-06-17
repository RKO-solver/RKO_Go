package ils

import (
	"github.com/lucasmends/rko-go/definition"
	"github.com/lucasmends/rko-go/random"
)

func (ils *ILS) SetRG(rg *random.Generator) {
	ils.RG = rg
	ils.search.SetRG(ils.RG)
}

func (ils *ILS) Name() string {
	return name
}

func (ils *ILS) SetIdWorker(id int) {
	if ils.logger != nil {
		ils.logger.SetIdWorker(id)
	}
}

func (ils *ILS) Solve() definition.Result {
	rko, elapsed := ils.solve(ils.solutionPool)

	return definition.Result{
		Solution:        ils.env.Decode(rko.RK),
		Cost:            rko.Cost,
		TimeSpentSecond: elapsed,
	}
}

func (ils *ILS) SetTimeLimitSecond(timeLimitSecond float64) {
	ils.configuration.TimeLimitSeconds = timeLimitSecond
}

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ definition.Solver = (*ILS)(nil)
