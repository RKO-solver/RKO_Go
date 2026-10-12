package config

import (
	"fmt"

	"github.com/RKO-solver/rko-go"
)

func implemented(t rko.Metaheuristic) bool {
	switch t {
	case rko.MultiStart, rko.SA, rko.GA, rko.VNS, rko.ILS, rko.BRKGA, rko.LNS:
		return true
	}
	return false
}

func metaheuristicError(i int, t rko.Metaheuristic) error {
	if implemented(t) {
		return nil
	}
	name := t.String()
	if name == "" {
		return fmt.Errorf("config: Metaheuristics[%d]: %w %d", i, ErrUnknownValue, t)
	}
	return fmt.Errorf("config: Metaheuristics[%d]: %s is %w", i, name, ErrNotImplemented)
}
