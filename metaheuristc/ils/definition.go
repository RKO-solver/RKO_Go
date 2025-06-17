package ils

import (
	"github.com/lucasmends/rko-go/definition"
	"github.com/lucasmends/rko-go/logger"
	"github.com/lucasmends/rko-go/metaheuristc/search"
	"github.com/lucasmends/rko-go/metaheuristc/solution"
	"github.com/lucasmends/rko-go/random"
)

const name = "ILS"

type Configuration struct {
	MaxIterations       int
	TimeLimitSeconds    float64
	ShakeMin            float64
	ShakeMax            float64
	MetropolisCriterion bool
}

type ILS struct {
	env           definition.Environment
	configuration *Configuration
	logger        *logger.Log
	search        search.Local
	RG            *random.Generator
	solutionPool  *solution.Pool
}
