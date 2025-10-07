package channel

import "github.com/RKO-solver/rko-go/logger"

func (l *Log) GetReportData() []logger.SolverInformation {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	report := make([]logger.SolverInformation, 0, len(l.data.solvers))

	for _, solver := range l.data.solvers {
		solverReport := logger.SolverInformation{
			Name:        "",
			Id:          -1,
			Performance: make([]logger.Data, 0),
		}

		for _, info := range solver {
			if solverReport.Name == "" {
				solverReport.Name = info.name
			}
			if solverReport.Id < 0 {
				solverReport.Id = info.id
			}

			solverReport.Performance = append(solverReport.Performance, logger.Data{
				LocalCost: info.local,
				BestCost:  info.localBest,
				Time:      info.time,
			})
		}

		report = append(report, solverReport)
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
