package ils

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

const name = "ILS"

type Configuration struct {
	MaxIterations       int     `yaml:"MaxIterations"`
	TimeLimitSeconds    float64 `yaml:"TimeLimitSeconds"`
	ShakeMin            float64 `yaml:"ShakeMin"`
	ShakeMax            float64 `yaml:"ShakeMax"`
	MetropolisCriterion bool    `yaml:"MetropolisCriterion"`
}

type ILS struct {
	env           definition.Environment
	configuration *Configuration
	logger        logger.SolverLogger
	search        search.Local
	RG            *random.Generator
	solutionPool  *solution.Pool
}
