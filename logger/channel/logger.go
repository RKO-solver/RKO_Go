package channel

import (
	"fmt"
	"sync"
	"time"

	"github.com/RKO-solver/rko-go/logger"
)

type Log struct {
	updateChan  chan channelMessage
	data        *information
	LogLevel    logger.Level
	ticker      time.Duration
	problemName string
}

func DefaultLogger(problemName string) *Log {
	return NewLogger(logger.DefaultLogLevel, problemName, defaultBufferSize)
}

func NewLoggerLevel(problemName string, level logger.Level) *Log {
	return NewLogger(level, problemName, defaultBufferSize)
}

func NewLogger(logLevel logger.Level, problemName string, bufferSize int) *Log {
	// The channel for communication
	progressChan := make(chan channelMessage, bufferSize)

	// The final data store
	store := &information{
		init:              false,
		previousLineCount: 0,
		solutionCost:      make([]int, 0, 100),
		solvers:           make([][]solverInfo, 0),
		extraMessages:     make([][]extraInfo, 0),
		workerMessages:    make([]string, 0),
	}

	return &Log{
		updateChan:  progressChan,
		data:        store,
		LogLevel:    logLevel,
		ticker:      defaultTickerMilliseconds * time.Millisecond,
		problemName: problemName,
	}
}

func (l *Log) AddSolutionPool(cost int) {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	l.data.solutionCost = append(l.data.solutionCost, cost)
}

func (l *Log) WorkerDone(message string) {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	l.data.workerMessages = append(l.data.workerMessages, message)
}

func (l *Log) Start(aggregatorWg *sync.WaitGroup) {
	progressChan := l.updateChan
	store := l.data

	// Start the single aggregator goroutine in the background.
	// It's the only one that writes to 'store'.
	aggregatorWg.Add(1)
	go func() {
		defer aggregatorWg.Done()
		for message := range progressChan {
			id := message.id
			switch message.t {
			case infoMessage:
				store.registerInfo(id, message.info)
			case verboseMessage:
				store.registerVerbose(id, message.extra)
			}

		}
	}()
}

// Shutdown closes the main channel to stop the aggregator
func (l *Log) Shutdown() {
	close(l.updateChan)
}

func (l *Log) Print() {
	l.data.printShell()
}

func (l *Log) WorkersPrint() {
	for _, message := range l.data.workerMessages {
		fmt.Println(message)
	}
}

func (l *Log) GetTicker() time.Duration {
	return l.ticker
}

func (l *Log) SetTicker(timeMilliseconds int) {
	if timeMilliseconds < minimumTickerMilliseconds {
		return
	}
	l.ticker = time.Duration(timeMilliseconds) * time.Millisecond
}

func (l *Log) GetLogger(name string) logger.SolverLogger {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	id := len(l.data.solvers)
	l.data.solvers = append(l.data.solvers, make([]solverInfo, 0))
	l.data.extraMessages = append(l.data.extraMessages, make([]extraInfo, 0))

	return &SolverLoggerChannel{
		id:         id,
		name:       name,
		logLevel:   l.LogLevel,
		updateChan: l.updateChan,
	}
}

func (l *Log) GetLogLevel() logger.Level {
	return l.LogLevel
}

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
				Time:      info.timeStamp,
			})
		}

		report = append(report, solverReport)
	}

	return report
}

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ logger.Logger = (*Log)(nil)
