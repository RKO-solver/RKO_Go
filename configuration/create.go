package configuration

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/RKO-solver/rko-go"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/ga"
	"github.com/RKO-solver/rko-go/metaheuristc/ils"
	"github.com/RKO-solver/rko-go/metaheuristc/lns"
	"github.com/RKO-solver/rko-go/metaheuristc/multistart"
	"github.com/RKO-solver/rko-go/metaheuristc/sa"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/vns"
	"go.yaml.in/yaml/v3"
)

// DefaultConfiguration returns the default parameters for every metaheuristic.
func DefaultConfiguration() *MetaheuristicsConfiguration {
	return &MetaheuristicsConfiguration{
		MultiStart: multistart.DefaulConfigurationtMultiStart(),
		BRKGA:      ga.DefaultConfigurationBRKGA(),
		GA:         ga.DefaultConfigurationGA(),
		ILS:        ils.DefaultConfigurationILS(),
		SA:         sa.DefaultConfigurationSA(),
		VNS:        vns.DefaultConfigurationVNS(),
		LNS:        lns.DefaultConfigurationLNS(),
	}
}

// createMHConfiguration applies opts on top of the default configuration.
func createMHConfiguration(opts ...Option) *MetaheuristicsConfiguration {
	config := DefaultConfiguration()

	for _, opt := range opts {
		opt(config)
	}

	return config
}

// processMetaheuristics parses "NAME[=search,...]" entries, dropping unknown names, unknown searches and RVND.
func processMetaheuristics(metaheuristics []string) []Solver {
	mhs := make([]Solver, 0)

	for _, env := range metaheuristics {
		parts := strings.SplitN(env, "=", 2)
		mhType := rko.GetMetaHeuristic(strings.TrimSpace(parts[0]))
		if mhType < 0 {
			continue
		}

		mh := Solver{
			MetaHeuristic: mhType,
			Search:        make([]search.Type, 0),
		}

		if len(parts) == 2 {
			components := strings.Split(parts[1], ",")

			for _, component := range components {
				searchType := search.GetSearchType(strings.TrimSpace(component))
				// RVND cannot be a component of itself.
				if searchType < 0 || searchType == search.RVND {
					continue
				}
				mh.Search = append(mh.Search, searchType)
			}
		}

		mhs = append(mhs, mh)
	}

	return mhs
}

// CreateSolverConfiguration loads the YAML file at filePath, keeping defaults for unset or non-positive parameters.
func CreateSolverConfiguration(filePath string) (*SolverConfiguration, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	config := solverYamlConfiguration{}

	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	mhConfiguration := config.Configuration

	var opts = []Option{
		withMultiStart(mhConfiguration.MultiStart),
		withGA(mhConfiguration.GA),
		withBRKGA(mhConfiguration.BRKGA),
		withSA(mhConfiguration.SA),
		withVNS(mhConfiguration.VNS),
		withILS(mhConfiguration.ILS),
		withLNS(mhConfiguration.LNS),
	}

	problemConfiguration := &SolverConfiguration{
		LoggerLevel:   logger.GetLogLevel(config.LogLevel),
		LoggerType:    logger.GetLogType(config.LogType),
		TimeLimit:     time.Duration(config.TimeLimitSeconds * int64(time.Second)),
		Solvers:       processMetaheuristics(config.Metaheuristics),
		Configuration: createMHConfiguration(opts...),
	}

	return problemConfiguration, nil
}
