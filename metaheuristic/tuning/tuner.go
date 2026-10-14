// Package tuning chooses metaheuristic parameters online.
package tuning

import "github.com/RKO-solver/rko-go/logger"

// Before each cycle of its run, a metaheuristic asks a [Tuner] for a
// configuration of a [Space]; after the cycle it reports a reward, from which
// the tuner learns.
//
// See README.md for a usage example, how states are numbered and how
// [QLearning] explores and learns.

// Tuner chooses a configuration for each cycle and learns from its reward.
type Tuner interface {
	// Choose must come before every Learn: Learn rewards the state returned by
	// the most recent Choose, and a Learn before the first Choose goes
	// undetected. progress is the fraction of the run done, in [0, 1]; it is not
	// checked.

	// Choose returns the state, in [0, Space.States()), for the next cycle.
	Choose(progress float64) int
	// Learn rewards the last chosen state; higher is better. progress is
	// measured after the cycle.
	Learn(reward float64, progress float64)
}

// PolicyReporter is implemented by Tuners that can describe what they learned.
type PolicyReporter interface{ Policy() logger.Table }
