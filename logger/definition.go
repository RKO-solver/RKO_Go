package logger

type Interface interface {
	Report(bestSolutionCost, localSolutionCost int, elapsed float64)
	Verbose(message string)
	Debug(message string)
	Info(message string)
	SetIdWorker(idWorker int)
	CreateLogger(method string) Interface
	Save()
	SaveFileName(fileName string)
}

type Log struct {
	saveReport bool
	handler    Interface
	LogLevel   Level
}

func (log *Log) GetLogger(method string) *Log {
	handler := log.handler.CreateLogger(method)
	return &Log{saveReport: log.saveReport, LogLevel: log.LogLevel, handler: handler}
}

func (log *Log) Debug(message string) {
	if log.LogLevel >= DEBUG {
		log.handler.Debug(message)
	}
}

func (log *Log) Report(bestSolutionCost, localSolutionCost int, elapsed float64) {
	if log.saveReport {
		log.handler.Report(bestSolutionCost, localSolutionCost, elapsed)
	}
}

func (log *Log) Info(message string) {
	if log.LogLevel >= INFO {
		log.handler.Info(message)
	}
}

func (log *Log) Verbose(message string) {
	if log.LogLevel >= VERBOSE {
		log.handler.Verbose(message)
	}
}

func (log *Log) SetIdWorker(idWorker int) {
	log.handler.SetIdWorker(idWorker)
}

func CreateLogger(logLevel Level, saveReport bool, handler Interface) *Log {
	return &Log{saveReport: saveReport, handler: handler, LogLevel: logLevel}
}
