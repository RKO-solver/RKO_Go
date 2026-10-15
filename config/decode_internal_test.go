package config

import (
	"testing"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/stretchr/testify/assert"
	"go.yaml.in/yaml/v3"
)

func TestDecodeTypeErrorsTranslation(t *testing.T) {
	te := &yaml.TypeError{Errors: []string{
		"line 4: field SeedPoolMaxSize not found in type config.yamlConfig",
		"line 2: cannot unmarshal !!str `abc` into int64",
	}}

	errs := typeErrors(te)
	assert.Len(t, errs, 2)
	assert.ErrorIs(t, errs[0], ErrUnknownField)
	assert.Equal(t, `config: line 4: unknown field "SeedPoolMaxSize"`, errs[0].Error())
	assert.NotErrorIs(t, errs[1], ErrUnknownField)
	assert.Equal(t, "config: line 2: cannot unmarshal !!str `abc` into int64", errs[1].Error())
}

func TestDecodeConvertLogLevelAndType(t *testing.T) {
	cfg, errs := convert(&yamlConfig{LogLevel: "silent", LogType: "dashboard"})
	assert.Empty(t, errs)
	assert.Equal(t, logger.LevelSilent, cfg.LogLevel)
	assert.Equal(t, logger.TypeDashboard, cfg.LogType)

	cfg, errs = convert(&yamlConfig{LogLevel: " ", LogType: "CHANNEL"})
	assert.Len(t, errs, 2)
	assert.Equal(t, logger.LevelInfo, cfg.LogLevel)
	assert.Equal(t, logger.TypePrint, cfg.LogType)
}
