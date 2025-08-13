package configuration

import (
	"github.com/lucasmends/rko-go/metaheuristc/ga"
	"github.com/lucasmends/rko-go/metaheuristc/ils"
	"github.com/lucasmends/rko-go/metaheuristc/multistart"
	"github.com/lucasmends/rko-go/metaheuristc/sa"
	"github.com/lucasmends/rko-go/metaheuristc/vns"
	"go.yaml.in/yaml/v3"
	"os"
)

type YamlConfiguration struct {
	TimeLimitSeconds float64                   `yaml:"TimeLimitSeconds"`
	MultiStart       *multistart.Configuration `yaml:"MultiStart"`
	BRKGA            *ga.ConfigurationBRKGA    `yaml:"BRKGA"`
	GA               *ga.ConfigurationGA       `yaml:"GA"`
	ILS              *ils.Configuration        `yaml:"ILS"`
	SA               *sa.Configuration         `yaml:"SA"`
	VNS              *vns.Configuration        `yaml:"VNS"`
}

type Option func(*YamlConfiguration)

func newYamlConfiguration(opts ...Option) *YamlConfiguration {
	// Start with a configuration struct populated with all default values.
	config := &YamlConfiguration{
		MultiStart: DefaultMultiStart(),
		BRKGA:      DefaultBRKGA(),
		GA:         DefaultGA(),
		ILS:        DefaultILS(),
		SA:         DefaultSA(),
		VNS:        DefaultVNS(),
	}

	// Apply all provided options, which will overwrite the defaults if specified.
	for _, opt := range opts {
		opt(config)
	}

	return config
}

func CreateYamlConfiguration(filePath string) (*YamlConfiguration, error) {

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	information := ymlStructure{}

	err = yaml.Unmarshal(data, &information)

	if err != nil {
		return nil, err
	}

	return &YamlConfiguration{&information}, nil
}
