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

	saveCsvFile(saveFile, data)
}

func (l *Log) savePool(filename ...string) {
	var saveFile string
	if len(filename) > 0 {
		saveFile = filename[0]
	} else {
		saveFile = fmt.Sprintf("%s-pool.csv", l.problemName)
	}

	rawData := l.GetSolutionData()

	data := make([][]string, 0, len(rawData)+1)

	data = append(data, []string{"cost", "time"})

	for _, line := range rawData {
		cost := fmt.Sprintf("%d", line.Cost)
		elapsed := fmt.Sprintf("%.3f", line.Time)
		data = append(data, []string{cost, elapsed})
	}

	saveCsvFile(saveFile, data)
}

func saveCsvFile(filename string, data [][]string) {
	file, err := os.Create(filename)
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

	var poolFilename string

	if len(filename) > 0 {
		poolFilename = fmt.Sprintf("%s-%s.csv", filename[0], "pool")
	} else {
		poolFilename = fmt.Sprintf("%s-%s.csv", l.problemName, "pool")
	}

	l.savePool(poolFilename)
}
