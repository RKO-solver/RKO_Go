package config

import (
	"testing"

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

func TestDecodeValidLogLevelAndType(t *testing.T) {
	for _, s := range []string{"", "silent", "INFO", "Verbose"} {
		assert.True(t, validLogLevel(s), s)
	}
	for _, s := range []string{"DEBUG", " ", "FILE"} {
		assert.False(t, validLogLevel(s), s)
	}
	for _, s := range []string{"", "print", "CHANNEL"} {
		assert.True(t, validLogType(s), s)
	}
	for _, s := range []string{"FILE", " ", "INFO"} {
		assert.False(t, validLogType(s), s)
	}
}
