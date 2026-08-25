package ga

import (
	"context"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/random"
)

func (ga *GA) SetRG(rg *random.Generator) {
	ga.RG = rg
	ga.search.SetRG(ga.RG)
}

func (ga *GA) Name() string {
	return nameGA
}

func (ga *GA) Solve(ctx context.Context) definition.Result {
	rko, elapsed := ga.solve(ctx, ga.solutionPool)

	return definition.Result{
		Solution:        ga.env.Decode(rko.RK),
		Cost:            rko.Cost,
		TimeSpentSecond: elapsed,
	}
}

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ definition.Solver = (*GA)(nil)
