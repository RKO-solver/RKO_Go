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

type ymlStructure struct {
	MultiStart *multistart.Configuration `yaml:"MultiStart"`
	BRKGA      *ga.ConfigurationBRKGA    `yaml:"BRKGA"`
	GA         *ga.ConfigurationGA       `yaml:"GA"`
	ILS        *ils.Configuration        `yaml:"ILS"`
	SA         *sa.Configuration         `yaml:"SA"`
	VNS        *vns.Configuration        `yaml:"VNS"`
}

type YamlConfiguration struct {
	information *ymlStructure
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
