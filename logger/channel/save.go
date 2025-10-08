package channel

import (
	"fmt"

	"github.com/RKO-solver/rko-go/logger"
)

func (l *Log) SaveCsv(filename ...string) {
	for _, report := range l.GetReportData() {
		if len(filename) > 0 {
			reportFilename := fmt.Sprintf("%s-%s-%d.csv", filename[0], report.Name, report.Id)
			logger.SaveSolverCSV(report, l.problemName, reportFilename)
		} else {
			logger.SaveSolverCSV(report, l.problemName)
		}
	}

	var poolFilename string

	if len(filename) > 0 {
		poolFilename = fmt.Sprintf("%s-%s.csv", filename[0], "pool")
	} else {
		poolFilename = fmt.Sprintf("%s-%s.csv", l.problemName, "pool")
	}

	logger.SavePoolCsv(l.GetSolutionData(), l.problemName, poolFilename)
}
