package lns

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

const name = "LNS"

type Configuration struct {
	MaxIterations    int     `yaml:"MaxIterations"`
	TimeLimitSeconds float64 `yaml:"TimeLimitSeconds"`
	BetaMin          float64 `yaml:"BetaMin"`
	BetaMax          float64 `yaml:"BetaMax"`
}

type LNS struct {
	env           definition.Environment
	configuration *Configuration
	logger        logger.SolverLogger
	search        search.Local
	RG            *random.Generator
	solutionPool  *solution.Pool
}

func (lns *LNS) Print() {
	search.PrintSolver(name, lns.search)
}
