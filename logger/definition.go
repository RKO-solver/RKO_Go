// Package logger provides interfaces and types for logging the progress and results
// of metaheuristic optimization algorithms. It supports different log levels, reporting,
// and allows for custom logger implementations.
package logger

// Interface defines the logging methods required for tracking metaheuristic progress and results.
// Implement this interface to provide custom logging behavior.
type Interface interface {
	// Report logs the best and local solution costs and elapsed time.
	Report(bestSolutionCost, localSolutionCost int, elapsed float64)
	// Verbose logs verbose-level messages.
	Verbose(message string)
	// Debug logs debug-level messages.
	Debug(message string)
	// Info logs info-level messages.
	Info(message string)
	// SetIdWorker sets the worker ID for the logger.
	SetIdWorker(idWorker int)
	// CreateLogger creates a new logger for a specific method or context.
	CreateLogger(method string) Interface
	// Save persists the log data.
	Save()
	// SaveFileName sets the filename for saving logs.
	SaveFileName(fileName string)
}

// Log is a wrapper for a logger.Interface, providing log level and report control.
// It is used throughout the metaheuristic framework to manage logging behavior.
type Log struct {
	saveReport bool      // Whether to save reports (calls Report)
	handler    Interface // The underlying logger implementation
	LogLevel   Level     // The current log level (DEBUG, INFO, VERBOSE, etc.)
}

// GetLogger creates a new Log instance for a specific method or context.
// This is useful for distinguishing logs from different metaheuristics or components.
func (log *Log) GetLogger(method string) *Log {
	handler := log.handler.CreateLogger(method)
	return &Log{saveReport: log.saveReport, LogLevel: log.LogLevel, handler: handler}
}

// Debug logs a debug-level message if the log level allows.
func (log *Log) Debug(message string) {
	if log.LogLevel >= DEBUG {
		log.handler.Debug(message)
	}
}

// Report logs the best and local solution costs and elapsed time if reporting is enabled.
func (log *Log) Report(bestSolutionCost, localSolutionCost int, elapsed float64) {
	if log.saveReport {
		log.handler.Report(bestSolutionCost, localSolutionCost, elapsed)
	}
}

// Info logs an info-level message if the log level allows.
func (log *Log) Info(message string) {
	if log.LogLevel >= INFO {
		log.handler.Info(message)
	}
}

// Verbose logs a verbose-level message if the log level allows.
func (log *Log) Verbose(message string) {
	if log.LogLevel >= VERBOSE {
		log.handler.Verbose(message)
	}
}

// SetIdWorker sets the worker ID for the logger.
func (log *Log) SetIdWorker(idWorker int) {
	log.handler.SetIdWorker(idWorker)
}

// CreateLogger creates a new Log instance with the specified log level, report setting, and handler.
// This function is typically used to initialize a logger for the metaheuristic framework.
func CreateLogger(logLevel Level, saveReport bool, handler Interface) *Log {
	return &Log{saveReport: saveReport, handler: handler, LogLevel: logLevel}
}
