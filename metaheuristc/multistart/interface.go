package multistart

import (
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

func (m *MultiStart) SetIdWorker(id int) {
	if m.logger != nil {
		m.logger.SetIdWorker(id)
	}
}

func (m *MultiStart) Solve() definition.Result {
	rko, elapsed := m.solve(m.solutionPool)

	return definition.Result{
		Solution:        m.env.Decode(rko.RK),
		Cost:            rko.Cost,
		TimeSpentSecond: elapsed,
	}
}

func (m *MultiStart) SetTimeLimitSecond(timeLimitSecond float64) {
	m.configuration.TimeLimitSeconds = timeLimitSecond
}

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ definition.Solver = (*MultiStart)(nil)
