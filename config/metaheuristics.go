package config

import (
	"fmt"

	"github.com/RKO-solver/rko-go/metaheuristic"
)

func metaheuristicError(i int, t metaheuristic.Type) error {
	if t.Implemented() {
		return nil
	}
	name := t.String()
	if name == "" {
		return fmt.Errorf("config: Metaheuristics[%d]: %w %d", i, ErrUnknownValue, t)
	}
	return fmt.Errorf("config: Metaheuristics[%d]: %s is %w", i, name, ErrNotImplemented)
}
