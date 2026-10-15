package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLevelNamesAndOrder(t *testing.T) {
	require.Equal(t, "SILENT", LevelSilent.String())
	require.Equal(t, "INFO", LevelInfo.String())
	require.Equal(t, "VERBOSE", LevelVerbose.String())
	require.Equal(t, "Level(7)", Level(7).String())
	require.Equal(t, "Level(-5)", Level(-5).String())
	require.True(t, LevelSilent < LevelInfo && LevelInfo < LevelVerbose)
	var zero Level
	require.Equal(t, LevelInfo, zero)
	require.Equal(t, LevelInfo, DefaultLevel)
}

func TestParseLevel(t *testing.T) {
	for _, l := range []Level{LevelSilent, LevelInfo, LevelVerbose} {
		got, err := ParseLevel(l.String())
		require.NoError(t, err)
		require.Equal(t, l, got)
	}
	got, err := ParseLevel("vErBoSe")
	require.NoError(t, err)
	require.Equal(t, LevelVerbose, got)
	for _, s := range []string{"", "DEBUG", " info"} {
		_, err := ParseLevel(s)
		require.Error(t, err, s)
	}
}
