package channel

import "github.com/RKO-solver/rko-go/logger"

func (l *Log) GetReportData() []logger.SolverInformation {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	report := make([]logger.SolverInformation, 0, len(l.data.solvers))

	for id, s := range l.data.solvers {
		report = append(report, logger.SolverInformation{
			Name:        s.name,
			Id:          id,
			Performance: s.file.All(s.records), // spilled records, then the in-memory tail
		})
	}
	return report
}

func (l *Log) GetSolutionData() []logger.SolutionData {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	report := make([]logger.SolutionData, 0, len(l.data.pool))

	for _, pool := range l.data.pool {
		report = append(report, logger.SolutionData{
			Cost: pool.cost,
			Time: pool.time,
		})
	}

	return report
}
