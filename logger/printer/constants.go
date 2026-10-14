package printer

import "time"

const (
	defaultPrintInterval = 500 * time.Millisecond
	minPrintInterval     = 10 * time.Millisecond
)

// defaultHistoryBytes is each metaheuristic's default history buffer, in bytes, before it spills to a file.
const defaultHistoryBytes = 512 * 1024
