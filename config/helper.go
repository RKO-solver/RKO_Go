package config

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristic/search"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/random"
)

func getSearch(sea []search.Type, environment definition.Environment, s *solution.Pool, rg *random.Generator) search.Local {
	if len(sea) == 1 {
		switch sea[0] {
		case search.Swap:
			return search.NewSwap(environment)
		case search.Mirror:
			return search.NewMirror(environment)
		case search.Farey:
			return search.NewFarey(environment, rg)
		case search.NelderMead:
			return search.NewNelderMead(environment, s, rg)
		case search.RVND:
			return search.NewDefault(environment, s, rg)
		}
	}

	neighbourhood := make([]search.Type, 0)

	for _, t := range sea {
		switch t {
		case search.Swap:
			neighbourhood = append(neighbourhood, search.Swap)
		case search.Mirror:
			neighbourhood = append(neighbourhood, search.Mirror)
		case search.Farey:
			neighbourhood = append(neighbourhood, search.Farey)
		case search.NelderMead:
			neighbourhood = append(neighbourhood, search.NelderMead)
		case search.RVND:
			continue

		}
	}

	if len(neighbourhood) == 0 {
		return search.NewDefault(environment, s, rg)
	}

	return search.NewRVND(environment, s, rg, neighbourhood)
}
