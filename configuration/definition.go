// Package configuration builds solvers from Go structs or YAML files.
package configuration

import (
	"time"

	"github.com/RKO-solver/rko-go"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/ga"
	"github.com/RKO-solver/rko-go/metaheuristc/ils"
	"github.com/RKO-solver/rko-go/metaheuristc/lns"
	"github.com/RKO-solver/rko-go/metaheuristc/sa"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/vns"
)

// MetaheuristicsConfiguration holds the tuning parameters of every metaheuristic.
type MetaheuristicsConfiguration struct {
	BRKGA *ga.ConfigurationBRKGA
	GA    *ga.ConfigurationGA
	ILS   *ils.Configuration
	SA    *sa.Configuration
	VNS   *vns.Configuration
	LNS   *lns.Configuration
}

// Solver pairs a metaheuristic with the local searches it uses.
type Solver struct {
	MetaHeuristic rko.MetaHeuristic
	Search        []search.Type
}

// SolverConfiguration holds the logger, time and iteration limits, metaheuristics and tuning parameters of a solver.
type SolverConfiguration struct {
	LoggerLevel   logger.Level
	LoggerType    logger.LogType
	TimeLimit     time.Duration // 0 or negative means no limit
	MaxIterations int           // iterations (generations for GA/BRKGA) per metaheuristic; 0 or negative means no limit
	Solvers       []Solver
	Configuration *MetaheuristicsConfiguration
}

// solverYamlConfiguration mirrors the top level of the YAML file.
type solverYamlConfiguration struct {
	LogLevel         string              `yaml:"LogLevel"`
	LogType          string              `yaml:"LogType"`
	TimeLimitSeconds int64               `yaml:"TimeLimitSeconds"`
	MaxIterations    int                 `yaml:"MaxIterations"`
	Metaheuristics   []string            `yaml:"Metaheuristics"`
	Configuration    mhYamlConfiguration `yaml:"Configuration"`
}

// mhYamlConfiguration mirrors the Configuration section of the YAML file.
type mhYamlConfiguration struct {
	BRKGA *ga.ConfigurationBRKGA `yaml:"BRKGA"`
	GA    *ga.ConfigurationGA    `yaml:"GA"`
	ILS   *ils.Configuration     `yaml:"ILS"`
	SA    *sa.Configuration      `yaml:"SA"`
	VNS   *vns.Configuration     `yaml:"VNS"`
	LNS   *lns.Configuration     `yaml:"LNS"`
}

// Option modifies a MetaheuristicsConfiguration.
type Option func(*MetaheuristicsConfiguration)
