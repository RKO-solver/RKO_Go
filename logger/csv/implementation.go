package csv

import (
	"encoding/csv"
	"fmt"
	"github.com/RKO-solver/rko-go/logger"
	"log"
	"os"
)

type Logger struct {
	reportName string
	method     string
	idWorker   int
	lines      []*performance
	subLogs    []*Logger
}

type performance struct {
	best  int
	local int
	time  float64
}

func (logger *Logger) SaveCsv(filename ...string) {
	hasSaved := false
	if len(logger.lines) > 0 {
		logger.saveCsv(filename...)
		hasSaved = true
	}

	for _, subLogger := range logger.subLogs {
		if !hasSaved {
			subLogger.SaveCsv(filename...)
		} else if len(filename) > 1 {
			subLogger.SaveCsv(filename[1:]...)
		} else {
			subLogger.SaveCsv()
		}
	}
}

func (logger *Logger) saveCsv(filename ...string) {
	var saveFile string
	if len(filename) > 0 {
		saveFile = filename[0]
	} else if logger.idWorker < 0 {
		saveFile = fmt.Sprintf("%s-%s.csv", logger.reportName, logger.method)
	} else {
		saveFile = fmt.Sprintf("%s-%s-(%d).csv", logger.reportName, logger.method, logger.idWorker)
	}

	data := make([][]string, 0, len(logger.lines)+1)

	data = append(data, []string{"best", "local", "time"})

	for _, line := range logger.lines {
		best := fmt.Sprintf("%d", line.best)
		local := fmt.Sprintf("%d", line.local)
		elapsed := fmt.Sprintf("%.3f", line.time)
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

func (logger *Logger) print(message string) {
	if len(logger.method) < 1 {
		fmt.Printf("%s\n", message)
		return
	}

	if logger.idWorker > 0 {
		fmt.Printf("%s (%d): %s\n", logger.method, logger.idWorker, message)
		return
	}

	fmt.Printf("%s: %s\n", logger.method, message)

}

func (logger *Logger) Report(bestSolutionCost, localSolutionCost int, elapsed float64) {
	logger.lines = append(logger.lines, &performance{best: bestSolutionCost, local: localSolutionCost, time: elapsed})
}

func (logger *Logger) Debug(message string) {
	logger.print(message)
}

func (logger *Logger) Info(message string) {
	logger.print(message)
}

func (logger *Logger) Verbose(message string) {
	logger.print(message)
}

func (logger *Logger) CreateLogger(method string) logger.Interface {
	newLogger := &Logger{
		reportName: logger.reportName,
		method:     method,
		idWorker:   -1,
		lines:      []*performance{},
		subLogs:    []*Logger{},
	}

	logger.subLogs = append(logger.subLogs, newLogger)
	return newLogger
}

func (logger *Logger) SetIdWorker(idWorker int) {
	logger.idWorker = idWorker
}

func (logger *Logger) Save() {
	logger.SaveCsv()
}

func (logger *Logger) SaveFileName(fileName string) {
	logger.SaveCsv(fileName)
}

func CreateLogger(reportName string) *Logger {
	return &Logger{
		reportName: reportName,
		method:     "",
		idWorker:   -1,
		lines:      []*performance{},
		subLogs:    []*Logger{},
	}
}

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ logger.Interface = (*Logger)(nil)
