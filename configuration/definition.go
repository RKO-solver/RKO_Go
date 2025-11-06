package configuration

import (
	"fmt"
	"os"
	"strings"

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

type MetaheuristicsConfiguration struct {
	MultiStart *multistart.Configuration
	BRKGA      *ga.ConfigurationBRKGA
	GA         *ga.ConfigurationGA
	ILS        *ils.Configuration
	SA         *sa.Configuration
	VNS        *vns.Configuration
	LNS        *lns.Configuration
}

type mhYamlConfiguration struct {
	TimeLimitSeconds float64                   `yaml:"TimeLimitSeconds"`
	MultiStart       *multistart.Configuration `yaml:"MultiStart"`
	BRKGA            *ga.ConfigurationBRKGA    `yaml:"BRKGA"`
	GA               *ga.ConfigurationGA       `yaml:"GA"`
	ILS              *ils.Configuration        `yaml:"ILS"`
	SA               *sa.Configuration         `yaml:"SA"`
	VNS              *vns.Configuration        `yaml:"VNS"`
	LNS              *lns.Configuration        `yaml:"LNS"`
}
type Option func(*MetaheuristicsConfiguration)

func newYamlConfiguration(opts ...Option) *MetaheuristicsConfiguration {
	// Start with a configuration struct populated with all default values.
	config := &MetaheuristicsConfiguration{
		MultiStart: multistart.DefaulConfigurationtMultiStart(),
		BRKGA:      ga.DefaultConfigurationBRKGA(),
		GA:         ga.DefaultConfigurationGA(),
		ILS:        ils.DefaultConfigurationILS(),
		SA:         sa.DefaultConfigurationSA(),
		VNS:        vns.DefaultConfigurationVNS(),
		LNS:        lns.DefaultConfigurationVNS(),
	}

	// Apply all provided options, which will overwrite the defaults if specified.
	for _, opt := range opts {
		opt(config)
	}

	return config
}

func CreateYamlMHConfiguration(filePath string) (*MetaheuristicsConfiguration, error) {

	data, err := os.ReadFile(filePath)
	if err != nil {

		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	configuration := mhYamlConfiguration{}

	if err = yaml.Unmarshal(data, &configuration); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)

	}

	var opts = []Option{
		withMultiStart(configuration.MultiStart, configuration.TimeLimitSeconds),
		withGA(configuration.GA, configuration.TimeLimitSeconds),
		withBRKGA(configuration.BRKGA, configuration.TimeLimitSeconds),
		withSA(configuration.SA, configuration.TimeLimitSeconds),
		withVNS(configuration.VNS, configuration.TimeLimitSeconds),
		withILS(configuration.ILS, configuration.TimeLimitSeconds),
		withLNS(configuration.LNS, configuration.TimeLimitSeconds),
	}

	return newYamlConfiguration(opts...), nil
}

type configurationSolver struct {
	LogLevel       string   `yaml:"logLevel"`
	LogType        string   `yaml:"logType"`
	Metaheuristics []string `yaml:"metaheuristics"`
}

type Solver struct {
	MetaHeuristic rko.MetaHeuristic
	Search        []search.Type
}
type SolverConfiguration struct {
	LoggerLevel logger.Level
	LoggerType  logger.LogType
	Solvers     []Solver
}

func processMetaheuristics(metaheuristics []string) []Solver {
	mhs := make([]Solver, 0)

	for _, env := range metaheuristics {
		// You can now process each string as needed
		parts := strings.SplitN(env, "=", 2)
		mhType := rko.GetMetaHeuristic(parts[0])
		// not valid
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
				searchType := search.GetSearchType(component)
				// ignore invalid and RVND for be circular
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

func CreateYamlSolverConfiguration(filePath string) (*SolverConfiguration, error) {

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	config := configurationSolver{}

	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	problemConfiguration := &SolverConfiguration{
		LoggerLevel: logger.GetLogLevel(config.LogLevel),
		LoggerType:  logger.GetLogType(config.LogType),
		Solvers:     processMetaheuristics(config.Metaheuristics),
	}

	return problemConfiguration, nil
}
