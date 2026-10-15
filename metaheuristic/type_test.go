package metaheuristic_test

import (
	"testing"

	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/stretchr/testify/assert"
)

func TestTypeImplemented(t *testing.T) {
	assert.True(t, metaheuristic.GA.Implemented())
	assert.False(t, metaheuristic.GRASP.Implemented())
	assert.False(t, metaheuristic.Type(99).Implemented())
}
