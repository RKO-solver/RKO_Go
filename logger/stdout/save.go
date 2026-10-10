package stdout

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

	for _, s := range l.solversLogger {
		if err := s.file.SaveCSV(fmt.Sprintf("%s-%s-%d.csv", prefix, s.name, s.id), s.history); err != nil {
			log.Fatal(err)
		}
	}

	logger.SavePoolCsv(l.GetSolutionData(), l.problemName, fmt.Sprintf("%s-pool.csv", prefix))
}
