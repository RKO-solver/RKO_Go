// Package ils implements Iterated Local Search (ILS).
package ils

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

const name = "ILS"

// Configuration holds the ILS parameters.
type Configuration struct {
	ShakeMin            float64 `yaml:"ShakeMin"`
	ShakeMax            float64 `yaml:"ShakeMax"`
	MetropolisCriterion bool    `yaml:"MetropolisCriterion"`
}

// ILS is an Iterated Local Search solver.
type ILS struct {
	env           definition.Environment
	configuration *Configuration
	logger        logger.SolverLogger
	search        search.Local
	rg            *random.Generator
	solutionPool  *solution.Pool
}

func (ils *ILS) Print() {
	search.PrintSolver(name, ils.search)
}
