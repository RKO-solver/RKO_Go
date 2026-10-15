package metaheuristic_test

import (
	"strings"
	"testing"

	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"
)

type listConfig struct {
	Size  metaheuristic.List[int]     `yaml:"Size"`
	Ratio metaheuristic.List[float64] `yaml:"Ratio"`
	On    metaheuristic.List[bool]    `yaml:"On"`
}

func decodeList(src string, strict bool) (listConfig, error) {
	var c listConfig
	d := yaml.NewDecoder(strings.NewReader(src))
	d.KnownFields(strict)
	return c, d.Decode(&c)
}

func TestListScalar(t *testing.T) {
	c, err := decodeList("Size: 7\nRatio: 0.5\nOn: true\n", true)
	require.NoError(t, err)
	assert.Equal(t, metaheuristic.List[int]{7}, c.Size)
	assert.Equal(t, metaheuristic.List[float64]{0.5}, c.Ratio)
	assert.Equal(t, metaheuristic.List[bool]{true}, c.On)
}

func TestListFlow(t *testing.T) {
	c, err := decodeList("Size: [1, 2, 3]\nRatio: [0.1, 0.2]\n", true)
	require.NoError(t, err)
	assert.Equal(t, metaheuristic.List[int]{1, 2, 3}, c.Size)
	assert.Equal(t, metaheuristic.List[float64]{0.1, 0.2}, c.Ratio)
	assert.Nil(t, c.On)
}

func TestListBlock(t *testing.T) {
	c, err := decodeList("Size:\n  - 10\n  - 20\nOn:\n  - true\n  - false\n", true)
	require.NoError(t, err)
	assert.Equal(t, metaheuristic.List[int]{10, 20}, c.Size)
	assert.Equal(t, metaheuristic.List[bool]{true, false}, c.On)
}

func TestListWrongType(t *testing.T) {
	for _, src := range []string{"Size: abc\n", "Size: [1, x]\n", "Size: {a: 1}\n"} {
		_, err := decodeList(src, false)
		var te *yaml.TypeError
		assert.ErrorAs(t, err, &te, src)
	}
}

func TestListStrictRejectsUnknownKey(t *testing.T) {
	src := "Size: [1, 2]\nOther: 3\n"

	_, err := decodeList(src, true)
	var te *yaml.TypeError
	require.ErrorAs(t, err, &te)
	assert.Contains(t, te.Errors[0], "Other")

	c, err := decodeList(src, false)
	require.NoError(t, err)
	assert.Equal(t, metaheuristic.List[int]{1, 2}, c.Size)
}
