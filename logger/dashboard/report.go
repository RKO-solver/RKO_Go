package dashboard

import (
	"slices"

	"github.com/RKO-solver/rko-go/logger"
)

// Report returns the history of every registered metaheuristic.
func (l *Log) Report() []logger.MetaheuristicReport {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	report := make([]logger.MetaheuristicReport, 0, len(l.data.metaheuristics))

	for id, s := range l.data.metaheuristics {
		report = append(report, logger.MetaheuristicReport{
			Name:    s.name,
			ID:      id,
			Samples: s.file.All(s.records), // spilled records, then the in-memory tail
		})
	}
	return report
}

// Pool returns a copy of the recorded pool solutions.
func (l *Log) Pool() []logger.PoolEntry {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	return slices.Clone(l.data.pool)
}
