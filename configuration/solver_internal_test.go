package configuration

import (
	"bytes"
	"testing"
	"time"

	rko "github.com/RKO-solver/rko-go"
	"github.com/stretchr/testify/assert"
)

func TestWarnIfUnbounded(t *testing.T) {
	tests := []struct {
		name          string
		timeLimit     time.Duration
		maxIterations int
		warn          bool
	}{
		{"neither set", 0, 0, true},
		{"negative values", -time.Second, -1, true},
		{"time limit only", time.Second, 0, false},
		{"max iterations only", 0, 10, false},
		{"both set", time.Second, 10, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			warnIfUnbounded(&buf, &SolverConfiguration{TimeLimit: tt.timeLimit, MaxIterations: tt.maxIterations})
			if tt.warn {
				assert.Equal(t, rko.NoLimitWarning+"\n", buf.String())
			} else {
				assert.Empty(t, buf.String())
			}
		})
	}
}
