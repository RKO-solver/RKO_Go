package logger

type Logger interface {
	AddSolutionPool(cost int)
	WorkerDone(message string)
	GetLogger(name string) SolverLogger
	GetLogLevel() Level
}

type SolverLogger interface {
	Register(local int, localBest int, timeStamp float64, extra string)
	Verbose(message string, timeStamp float64)
}
