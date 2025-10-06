package logger

type SolverLogger struct {
	id         int
	name       string
	logLevel   Level
	updateChan chan channelMessage
	extraChan  chan string
}

func (sl *SolverLogger) Register(local int, localBest int, timeStamp float64, extra string) {
	sl.updateChan <- channelMessage{
		t:  infoMessage,
		id: sl.id,
		info: solverInfo{
			name:      sl.name,
			localBest: localBest,
			local:     local,
			timeStamp: timeStamp,
			extra:     extra,
		},
		message: "",
	}
}

func (sl *SolverLogger) Verbose(message string) {
	// ignore if log level is not verbose
	if sl.logLevel >= VERBOSE {
		sl.updateChan <- channelMessage{
			t:       verboseMessage,
			info:    solverInfo{},
			message: message,
		}
	}
}

func (l *Logger) GetLogger(name string) *SolverLogger {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	id := len(l.data.solvers)
	l.data.solvers = append(l.data.solvers, make([]solverInfo, 0, 50))
	l.data.extraMessages = append(l.data.extraMessages, []string{})

	return &SolverLogger{
		id:         id,
		name:       name,
		logLevel:   l.LogLevel,
		updateChan: l.updateChan,
	}
}
