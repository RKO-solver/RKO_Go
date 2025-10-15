package logger

import "strings"

type SolverInformation struct {
	Name        string
	Id          int
	Performance []Data
}

type SolutionData struct {
	Cost int
	Time float64
}
type Data struct {
	LocalCost int
	BestCost  int
	Time      float64
}

type Logger interface {
	AddSolutionPool(cost int, time float64)
	WorkerDone(message string)
	GetLogger(name string) SolverLogger
	GetLogLevel() Level
	GetReportData() []SolverInformation
	GetSolutionData() []SolutionData
}

type SolverLogger interface {
	Register(local int, localBest int, time float64, extra string)
	Verbose(message string, timeStamp float64)
}

type LogType = uint8

const (
	CHANNEL LogType = iota
	PRINT
)

func GetLogType(label string) LogType {
	label = strings.ToUpper(label)
	switch label {
	case "CHANNEL":
		return CHANNEL
	case "PRINT":
		return PRINT
	default:
		return PRINT
	}
}
