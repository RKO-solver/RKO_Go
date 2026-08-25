package multistart

import (
	"context"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/random"
)

func (m *MultiStart) SetRG(rg *random.Generator) {
	m.RG = rg
	m.search.SetRG(m.RG)
}

func (m *MultiStart) Name() string {
	return name
}

func (m *MultiStart) Solve(ctx context.Context) definition.Result {
	rko, elapsed := m.solve(ctx, m.solutionPool)

	return definition.Result{
		Solution:        m.env.Decode(rko.RK),
		Cost:            rko.Cost,
		TimeSpentSecond: elapsed,
	}
}

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ definition.Solver = (*MultiStart)(nil)
