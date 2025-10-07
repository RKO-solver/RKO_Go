package channel

import "github.com/RKO-solver/rko-go/logger"

type SolverLoggerChannel struct {
	id         int
	name       string
	logLevel   logger.Level
	updateChan chan channelMessage
	extraChan  chan string
}

func (sl *SolverLoggerChannel) Register(local int, localBest int, timeStamp float64, extra string) {
	sl.updateChan <- channelMessage{
		t:  infoMessage,
		id: sl.id,
		info: solverInfo{
			id:        sl.id,
			name:      sl.name,
			localBest: localBest,
			local:     local,
			time:      timeStamp,
			extra:     extra,
		},
		extra: extraInfo{},
	}
}

func (sl *SolverLoggerChannel) Verbose(message string, timeStamp float64) {
	// ignore if log level is not verbose
	if sl.logLevel >= logger.VERBOSE {
		sl.updateChan <- channelMessage{
			t:    verboseMessage,
			id:   sl.id,
			info: solverInfo{},
			extra: extraInfo{
				message:   message,
				timeStamp: timeStamp,
			},
		}
	}
}

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ logger.SolverLogger = (*SolverLoggerChannel)(nil)
