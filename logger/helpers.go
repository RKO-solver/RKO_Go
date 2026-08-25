package logger

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

// SaveCsvFile writes data as CSV rows to filename, exiting the program on error.
func SaveCsvFile(filename string, data [][]string) {
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

// SavePoolCsv writes poolData to a CSV file with cost and time columns.
func SavePoolCsv(poolData []SolutionData, problemName string, filename ...string) {
	var saveFile string
	if len(filename) > 0 {
		saveFile = filename[0]
	} else {
		saveFile = fmt.Sprintf("%s-pool.csv", problemName)
	}

	data := make([][]string, 0, len(poolData)+1)

	data = append(data, []string{"cost", "time"})

	for _, line := range poolData {
		cost := fmt.Sprintf("%d", line.Cost)
		elapsed := fmt.Sprintf("%.8f", line.Time)
		data = append(data, []string{cost, elapsed})
	}

	SaveCsvFile(saveFile, data)
}

// SaveSolverCSV writes solverInfo's performance history to a CSV file.
func SaveSolverCSV(solverInfo SolverInformation, problemName string, filename ...string) {
	var saveFile string
	if len(filename) > 0 {
		saveFile = filename[0]
	} else {
		saveFile = fmt.Sprintf("%s-%s-%d.csv", problemName, solverInfo.Name, solverInfo.Id)
	}

	data := make([][]string, 0, len(solverInfo.Performance)+1)

	data = append(data, []string{"best", "local", "time"})

	for _, line := range solverInfo.Performance {
		best := fmt.Sprintf("%d", line.BestCost)
		local := fmt.Sprintf("%d", line.LocalCost)
		elapsed := fmt.Sprintf("%.8f", line.Time)
		data = append(data, []string{best, local, elapsed})
	}

	SaveCsvFile(saveFile, data)
}
