package config

import (
	"github.com/RKO-solver/rko-go/metaheuristc/brkga"
	"github.com/RKO-solver/rko-go/metaheuristc/ga"
	"github.com/RKO-solver/rko-go/metaheuristc/ils"
	"github.com/RKO-solver/rko-go/metaheuristc/lns"
	"github.com/RKO-solver/rko-go/metaheuristc/sa"
	"github.com/RKO-solver/rko-go/metaheuristc/vns"
)

type yamlConfig struct {
	LogLevel         string         `yaml:"LogLevel"`
	LogType          string         `yaml:"LogType"`
	TimeLimitSeconds int64          `yaml:"TimeLimitSeconds"`
	MaxIterations    int            `yaml:"MaxIterations"`
	SeedPool         yamlSeedPool   `yaml:"SeedPool"`
	Metaheuristics   []string       `yaml:"Metaheuristics"`
	Parameters       yamlParameters `yaml:"Parameters"`
}

type yamlSeedPool struct {
	MaxSize     int `yaml:"MaxSize"`
	InitialSize int `yaml:"InitialSize"`
}

type yamlParameters struct {
	BRKGA *brkga.Config `yaml:"BRKGA"`
	GA    *ga.Config    `yaml:"GA"`
	ILS   *ils.Config   `yaml:"ILS"`
	SA    *sa.Config    `yaml:"SA"`
	VNS   *vns.Config   `yaml:"VNS"`
	LNS   *lns.Config   `yaml:"LNS"`
}
