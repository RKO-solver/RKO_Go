package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypeNames(t *testing.T) {
	require.Equal(t, "PRINT", TypePrint.String())
	require.Equal(t, "DASHBOARD", TypeDashboard.String())
	require.Equal(t, "Type(7)", Type(7).String())
	var zero Type
	require.Equal(t, TypePrint, zero)
}

func TestParseType(t *testing.T) {
	for _, v := range []Type{TypePrint, TypeDashboard} {
		got, err := ParseType(v.String())
		require.NoError(t, err)
		require.Equal(t, v, got)
	}
	got, err := ParseType("dashboard")
	require.NoError(t, err)
	require.Equal(t, TypeDashboard, got)
	for _, s := range []string{"", "CHANNEL", "FILE"} {
		_, err := ParseType(s)
		require.Error(t, err, s)
	}
}
