package dashboard

import (
	"slices"

	"github.com/RKO-solver/rko-go/logger"
)

func (l *Log) Report() []logger.SolverReport {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	report := make([]logger.SolverReport, 0, len(l.data.solvers))

	for id, s := range l.data.solvers {
		report = append(report, logger.SolverReport{
			Name:    s.name,
			ID:      id,
			Samples: s.file.All(s.records), // spilled records, then the in-memory tail
		})
	}
	return report
}

func (l *Log) Pool() []logger.PoolEntry {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	return slices.Clone(l.data.pool)
}
