package vns

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/random"
)

func (vns *VNS) SetRG(rg *random.Generator) {
	vns.RG = rg
}

func (vns *VNS) Name() string {
	return name
}

func (vns *VNS) SetIdWorker(id int) {
	if vns.logger != nil {
		vns.logger.SetIdWorker(id)
	}
}

func (vns *VNS) Solve() definition.Result {
	rko, elapsed := vns.solve(vns.solutionPool)

	return definition.Result{
		Solution:        vns.env.Decode(rko.RK),
		Cost:            rko.Cost,
		TimeSpentSecond: elapsed,
	}
}

func (vns *VNS) SetTimeLimitSecond(timeLimitSecond float64) {
	vns.configuration.TimeLimitSeconds = timeLimitSecond
}

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ definition.Solver = (*VNS)(nil)
