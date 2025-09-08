package ils

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/random"
)

type history struct {
	defaultMin         float64
	defaultMax         float64
	min                float64
	max                float64
	timesNoImprovement uint
}

func shake(sol *metaheuristc.RandomKeyValue, history *history, rg *random.Generator, env definition.Environment) {
	if history.timesNoImprovement == 0 {
		history.min = history.defaultMin
		history.max = history.defaultMax
	} else if history.timesNoImprovement > 10 && history.timesNoImprovement < 1000 {
		history.min *= 1.1
		history.max *= 1.1
	}

	rk.Shake(sol, history.min, history.max, rg, env)
}
