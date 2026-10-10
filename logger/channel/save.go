package channel

import (
	"fmt"
	"log"

	"github.com/RKO-solver/rko-go/logger"
)

// SaveCsv writes a per-solver CSV, streaming its history file, and a pool CSV.
func (l *Log) SaveCsv(filename ...string) {
	prefix := l.problemName
	if len(filename) > 0 {
		prefix = filename[0]
	}

	// The lock keeps the aggregator from spilling while a solver's file and tail are read.
	l.data.mu.Lock()
	for id, s := range l.data.solvers {
		if err := s.file.SaveCSV(fmt.Sprintf("%s-%s-%d.csv", prefix, s.name, id), s.records); err != nil {
			log.Fatal(err)
		}
	}
	l.data.mu.Unlock() // GetSolutionData takes the lock itself

	logger.SavePoolCsv(l.GetSolutionData(), l.problemName, fmt.Sprintf("%s-pool.csv", prefix))
}
