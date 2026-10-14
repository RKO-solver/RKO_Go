package rko

import "time"

// NoLimitWarning is printed when a Solver has no time or iteration limit.
const NoLimitWarning = "Warning: neither a time limit nor an iteration limit is set; the solver will run until its context is cancelled"

// DefaultTimeLimit is the time limit of New.
const DefaultTimeLimit = 150 * time.Second
