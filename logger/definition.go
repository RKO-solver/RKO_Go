package logger

type SolverInformation struct {
	Name        string
	Id          int
	Performance []Data
}

type Data struct {
	LocalCost int
	BestCost  int
	Time      float64
}

type Logger interface {
	AddSolutionPool(cost int)
	WorkerDone(message string)
	GetLogger(name string) SolverLogger
	GetLogLevel() Level
	GetReportData() []SolverInformation
}

type SolverLogger interface {
	Register(local int, localBest int, timeStamp float64, extra string)
	Verbose(message string, timeStamp float64)
}
