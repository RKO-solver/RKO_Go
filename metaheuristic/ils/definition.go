// Package ils implements Iterated Local Search (ILS).
package ils

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/search"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/tuning"
	"github.com/RKO-solver/rko-go/random"
)

const name = "ILS"

// Config holds the ILS parameters.
type Config struct {
	ShakeMin            metaheuristic.List[float64] `yaml:"ShakeMin"`
	ShakeMax            metaheuristic.List[float64] `yaml:"ShakeMax"`
	MetropolisCriterion *bool                       `yaml:"MetropolisCriterion"`
}

// ILS is an Iterated Local Search solver.
type ILS struct {
	env           definition.Environment
	configuration *Config
	logger        logger.SolverLogger
	search        search.Local
	rg            *random.Generator
	solutionPool  *solution.Pool
	tuner         tuning.Tuner
	states        int
	shakeMin      tuning.Param[float64]
	shakeMax      tuning.Param[float64]
}

func (ils *ILS) Print() {
	search.PrintSolver(name, ils.search)
}
