package configuration

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

func getSearch(sea []search.Type, environment definition.Environment, s *solution.Pool, rg *random.Generator) search.Local {
	if len(sea) == 1 {
		switch sea[0] {
		case search.Swap:
			return search.CreateSwapLocalSearch(environment)
		case search.Mirror:
			return search.CreateMirrorLocalSearch(environment)
		case search.Farey:
			return search.CreateFareyLocalSearch(environment, rg)
		case search.Nelder:
			return search.CreateNelderMeadLocalSearch(environment, s, rg)
		case search.RVND:
			return search.CreateDefault(environment, s, rg)
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
		case search.Nelder:
			neighbourhood = append(neighbourhood, search.Nelder)
		case search.RVND:
			continue

		}
	}

	if len(sea) < 0 || len(neighbourhood) == 0 {
		return search.CreateDefault(environment, s, rg)
	}

	return search.CreateRVND(environment, s, rg, neighbourhood)
}
