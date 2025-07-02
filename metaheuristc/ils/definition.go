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
	MaxIterations       int     `yaml:"MaxIterations"`
	TimeLimitSeconds    float64 `yaml:"TimeLimitSeconds"`
	ShakeMin            float64 `yaml:"ShakeMin"`
	ShakeMax            float64 `yaml:"ShakeMax"`
	MetropolisCriterion bool    `yaml:"MetropolisCriterion"`
}

type ILS struct {
	env           definition.Environment
	configuration *Configuration
	logger        *logger.Log
	search        search.Local
	RG            *random.Generator
	solutionPool  *solution.Pool
}
