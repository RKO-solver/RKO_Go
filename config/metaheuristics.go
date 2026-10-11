package config

import (
	"fmt"

	"github.com/RKO-solver/rko-go"
)

func implemented(t rko.MetaHeuristic) bool {
	switch t {
	case rko.MULTISTART, rko.SA, rko.GA, rko.VNS, rko.ILS, rko.BRKGA, rko.LNS:
		return true
	}
	return false
}

func metaheuristicError(i int, t rko.MetaHeuristic) error {
	if implemented(t) {
		return nil
	}
	name := rko.GetMetaHeuristicString(t)
	if name == "" {
		return fmt.Errorf("config: Metaheuristics[%d]: %w %d", i, ErrUnknownValue, t)
	}
	return fmt.Errorf("config: Metaheuristics[%d]: %s is %w", i, name, ErrNotImplemented)
}
