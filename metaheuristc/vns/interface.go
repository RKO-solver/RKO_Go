package vns

import (
	"context"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/random"
)

func (vns *VNS) SetRG(rg *random.Generator) {
	vns.RG = rg
}

func (vns *VNS) Name() string {
	return name
}

func (vns *VNS) Solve(ctx context.Context) definition.Result {
	rko, elapsed := vns.solve(ctx, vns.solutionPool)

	return definition.Result{
		Solution:        vns.env.Decode(rko.RK),
		Cost:            rko.Cost,
		TimeSpentSecond: elapsed,
	}
}

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ definition.Solver = (*VNS)(nil)
