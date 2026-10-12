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
	BRKGA *brkga.ConfigurationBRKGA `yaml:"BRKGA"`
	GA    *ga.ConfigurationGA       `yaml:"GA"`
	ILS   *ils.Configuration        `yaml:"ILS"`
	SA    *sa.Configuration         `yaml:"SA"`
	VNS   *vns.Configuration        `yaml:"VNS"`
	LNS   *lns.Configuration        `yaml:"LNS"`
}
