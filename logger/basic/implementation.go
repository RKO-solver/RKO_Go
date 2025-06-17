package basic

import (
	"fmt"
	"github.com/lucasmends/rko-go/logger"
)

type Logger struct {
	method   string
	idWorker int
}

func (logger *Logger) Save() {

}

func (logger *Logger) SaveFileName(fileName string) {

}

func (logger *Logger) Report(bestSolutionCost, localSolutionCost int, elapsed float64) {
	fmt.Printf("best solution cost: %d, local solution cost: %d, elapsed %.3f\n", bestSolutionCost, localSolutionCost, elapsed)
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
	return &Logger{
		method: method,
	}
}

func (logger *Logger) SetIdWorker(idWorker int) {
	logger.idWorker = idWorker
}

func CreateLogger() logger.Interface {
	return &Logger{}
}

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ logger.Interface = (*Logger)(nil)
