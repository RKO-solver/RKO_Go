package channel

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/RKO-solver/rko-go/logger"
)

func (l *Log) saveIndividualCsv(info logger.SolverInformation, filename ...string) {
	var saveFile string
	if len(filename) > 0 {
		saveFile = filename[0]
	} else {
		saveFile = fmt.Sprintf("%s-%s-%d.csv", l.problemName, info.Name, info.Id)
	}

	data := make([][]string, 0, len(info.Performance)+1)

	data = append(data, []string{"best", "local", "time"})

	for _, line := range info.Performance {
		best := fmt.Sprintf("%d", line.BestCost)
		local := fmt.Sprintf("%d", line.LocalCost)
		elapsed := fmt.Sprintf("%.3f", line.Time)
		data = append(data, []string{best, local, elapsed})
	}

	file, err := os.Create(saveFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(data)
	if err != nil {
		return
	}
}

func (l *Log) SaveCsv(filename ...string) {
	for _, report := range l.GetReportData() {
		if len(filename) > 0 {
			reportFilename := fmt.Sprintf("%s-%s-%d.csv", filename[0], report.Name, report.Id)
			l.saveIndividualCsv(report, reportFilename)
		} else {
			l.saveIndividualCsv(report)
		}
	}
}
