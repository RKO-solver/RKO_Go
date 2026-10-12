package search

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/random"
)

// New returns the local search of searchType.
func New(searchType Type, env definition.Environment, pool *solution.Pool, rg *random.Generator) Local {
	switch searchType {
	case Swap:
		return NewSwap(env)
	case Mirror:
		return NewMirror(env)
	case Farey:
		return NewFarey(env, rg)
	case RVND:
		neighbourhood := make([]Type, 3)
		neighbourhood[0] = Swap
		neighbourhood[1] = Mirror
		neighbourhood[2] = Farey
		return NewRVND(env, pool, rg, neighbourhood)
	case NelderMead:
		return NewNelderMead(env, pool, rg)
	default:
		return NewDefault(env, pool, rg)
	}
}

// NewDefault returns an RVND over Swap, Mirror, Farey and Nelder-Mead.
func NewDefault(env definition.Environment, pool *solution.Pool, rg *random.Generator) Local {
	neighbourhood := make([]Type, 4)
	neighbourhood[0] = Swap
	neighbourhood[1] = Mirror
	neighbourhood[2] = Farey
	neighbourhood[3] = NelderMead

	return NewRVND(env, pool, rg, neighbourhood)
}
