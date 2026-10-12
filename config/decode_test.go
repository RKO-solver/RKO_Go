package config_test

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/RKO-solver/rko-go"
	"github.com/RKO-solver/rko-go/config"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// decodeAll decodes every document in src until io.EOF or the first error.
func decodeAll(src string, strict bool) ([]config.Config, error) {
	d := config.NewDecoder(strings.NewReader(src))
	d.Strict(strict)
	var cfgs []config.Config
	for {
		var c config.Config
		err := d.Decode(&c)
		if err == io.EOF {
			return cfgs, nil
		}
		if err != nil {
			return cfgs, err
		}
		cfgs = append(cfgs, c)
	}
}

// decodeOne decodes the first document in src.
func decodeOne(src string, strict bool) (config.Config, error) {
	d := config.NewDecoder(strings.NewReader(src))
	d.Strict(strict)
	var c config.Config
	err := d.Decode(&c)
	return c, err
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

func defaultDecoded() config.Config {
	return config.Config{
		LogLevel:   logger.LevelInfo,
		LogType:    logger.TypePrint,
		SeedPool:   config.SeedPool{MaxSize: solution.DefaultSize, InitialSize: solution.DefaultSize},
		Parameters: config.DefaultParameters(),
	}
}

func mustDecode(t *testing.T, src string, strict bool) config.Config {
	t.Helper()
	c, err := decodeOne(src, strict)
	require.NoError(t, err)
	return c
}

func TestDecodeFull(t *testing.T) {
	src := `LogLevel: VERBOSE
LogType: DASHBOARD
TimeLimitSeconds: 60
MaxIterations: 500
SeedPool:
  MaxSize: 50
  InitialSize: 20
Metaheuristics:
  - "ILS=Swap, Mirror"
  - "BRKGA"
  - "SA=Farey"
Parameters:
  BRKGA:
    PopulationSize: 300
    EliteRatio: 0.25
    MutantRatio: 0.1
    CrossoverAlpha: 0.6
    MutationAlpha: 0.15
    MaxGenerationNoImprovement: 20
  GA:
    PopulationSize: 123
    CrossoverAlpha: 0.7
    MutationAlpha: 0.2
    MaxGenerationNoImprovement: 10
  ILS:
    ShakeMin: 2.0
    ShakeMax: 8.0
    MetropolisCriterion: true
  SA:
    Alpha: 0.95
    TemperatureInitial: 1000
    TemperatureGoal: 10
    TemperatureReheat: 200
    ShakeMin: 2.0
    ShakeMax: 5.0
    QtdReheat: 3
    Iterations: 100
  VNS:
    Rate: 0.33
  LNS:
    BetaMin: 0.5
    BetaMax: 0.9
`
	for _, strict := range []bool{false, true} {
		cfg := mustDecode(t, src, strict)

		assert.Equal(t, logger.LevelVerbose, cfg.LogLevel)
		assert.Equal(t, logger.TypeDashboard, cfg.LogType)
		assert.Equal(t, 60*time.Second, cfg.TimeLimit)
		assert.Equal(t, 500, cfg.MaxIterations)
		assert.Equal(t, config.SeedPool{MaxSize: 50, InitialSize: 20}, cfg.SeedPool)
		assert.Equal(t, []config.Metaheuristic{
			{Type: rko.ILS, Searches: []search.Type{search.Swap, search.Mirror}},
			{Type: rko.BRKGA},
			{Type: rko.SA, Searches: []search.Type{search.Farey}},
		}, cfg.Metaheuristics)

		p := cfg.Parameters
		require.NotNil(t, p)

		assert.Equal(t, 300, p.BRKGA.PopulationSize)
		assert.Equal(t, 0.25, p.BRKGA.EliteRatio)
		assert.Equal(t, 0.1, p.BRKGA.MutantRatio)
		assert.Equal(t, 0.6, p.BRKGA.CrossoverAlpha)
		assert.Equal(t, 0.15, p.BRKGA.MutationAlpha)
		assert.Equal(t, 20, p.BRKGA.MaxGenerationNoImprovement)

		assert.Equal(t, 123, p.GA.PopulationSize)
		assert.Equal(t, 0.7, p.GA.CrossoverAlpha)
		assert.Equal(t, 0.2, p.GA.MutationAlpha)
		assert.Equal(t, 10, p.GA.MaxGenerationNoImprovement)

		assert.Equal(t, 2.0, p.ILS.ShakeMin)
		assert.Equal(t, 8.0, p.ILS.ShakeMax)
		assert.True(t, p.ILS.MetropolisCriterion)

		assert.Equal(t, 0.95, p.SA.Alpha)
		assert.Equal(t, 1000.0, p.SA.TemperatureInitial)
		assert.Equal(t, 10.0, p.SA.TemperatureGoal)
		assert.Equal(t, 200.0, p.SA.TemperatureReheat)
		assert.Equal(t, 2.0, p.SA.ShakeMin)
		assert.Equal(t, 5.0, p.SA.ShakeMax)
		assert.Equal(t, uint8(3), p.SA.QtdReheat)
		assert.Equal(t, 100, p.SA.Iterations)

		assert.Equal(t, 0.33, p.VNS.Rate)

		assert.Equal(t, 0.5, p.LNS.BetaMin)
		assert.Equal(t, 0.9, p.LNS.BetaMax)
	}
}

func TestDecodeParametersSections(t *testing.T) {
	tests := []struct {
		name  string
		src   string
		check func(t *testing.T, p *config.Parameters)
	}{
		{
			name: "GA",
			src:  "Parameters:\n  GA:\n    PopulationSize: 123\n    CrossoverAlpha: 0.7\n    MutationAlpha: 0.2\n    MaxGenerationNoImprovement: 10\n",
			check: func(t *testing.T, p *config.Parameters) {
				assert.Equal(t, 123, p.GA.PopulationSize)
				assert.Equal(t, 0.7, p.GA.CrossoverAlpha)
				assert.Equal(t, 0.2, p.GA.MutationAlpha)
				assert.Equal(t, 10, p.GA.MaxGenerationNoImprovement)
			},
		},
		{
			name: "BRKGA",
			src:  "Parameters:\n  BRKGA:\n    PopulationSize: 300\n    EliteRatio: 0.25\n    MutantRatio: 0.1\n    CrossoverAlpha: 0.6\n    MutationAlpha: 0.15\n    MaxGenerationNoImprovement: 20\n",
			check: func(t *testing.T, p *config.Parameters) {
				assert.Equal(t, 300, p.BRKGA.PopulationSize)
				assert.Equal(t, 0.25, p.BRKGA.EliteRatio)
				assert.Equal(t, 0.1, p.BRKGA.MutantRatio)
				assert.Equal(t, 0.6, p.BRKGA.CrossoverAlpha)
				assert.Equal(t, 0.15, p.BRKGA.MutationAlpha)
				assert.Equal(t, 20, p.BRKGA.MaxGenerationNoImprovement)
			},
		},
		{
			name: "SA",
			src:  "Parameters:\n  SA:\n    Alpha: 0.95\n    TemperatureInitial: 1000\n    TemperatureGoal: 10\n    TemperatureReheat: 200\n    ShakeMin: 2.0\n    ShakeMax: 5.0\n    QtdReheat: 3\n    Iterations: 100\n",
			check: func(t *testing.T, p *config.Parameters) {
				assert.Equal(t, 0.95, p.SA.Alpha)
				assert.Equal(t, 1000.0, p.SA.TemperatureInitial)
				assert.Equal(t, 10.0, p.SA.TemperatureGoal)
				assert.Equal(t, 200.0, p.SA.TemperatureReheat)
				assert.Equal(t, 2.0, p.SA.ShakeMin)
				assert.Equal(t, 5.0, p.SA.ShakeMax)
				assert.Equal(t, uint8(3), p.SA.QtdReheat)
				assert.Equal(t, 100, p.SA.Iterations)
			},
		},
		{
			name: "VNS",
			src:  "Parameters:\n  VNS:\n    Rate: 0.33\n",
			check: func(t *testing.T, p *config.Parameters) {
				assert.Equal(t, 0.33, p.VNS.Rate)
			},
		},
		{
			name: "ILS",
			src:  "Parameters:\n  ILS:\n    ShakeMin: 2.0\n    ShakeMax: 8.0\n    MetropolisCriterion: true\n",
			check: func(t *testing.T, p *config.Parameters) {
				assert.Equal(t, 2.0, p.ILS.ShakeMin)
				assert.Equal(t, 8.0, p.ILS.ShakeMax)
				assert.True(t, p.ILS.MetropolisCriterion)
			},
		},
		{
			name: "LNS",
			src:  "Parameters:\n  LNS:\n    BetaMin: 2.0\n    BetaMax: 8.0\n",
			check: func(t *testing.T, p *config.Parameters) {
				assert.Equal(t, 2.0, p.LNS.BetaMin)
				assert.Equal(t, 8.0, p.LNS.BetaMax)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := mustDecode(t, tt.src, true)
			require.NotNil(t, cfg.Parameters)
			tt.check(t, cfg.Parameters)
		})
	}
}

func TestDecodePartialSectionKeepsDefaults(t *testing.T) {
	cfg := mustDecode(t, "Parameters:\n  SA:\n    Alpha: 0.9\n    Iterations: 7\n  ILS:\n    ShakeMax: 0.4\n", false)
	defaults := config.DefaultParameters()

	expectedSA := *defaults.SA
	expectedSA.Alpha = 0.9
	expectedSA.Iterations = 7
	assert.Equal(t, expectedSA, *cfg.Parameters.SA)

	expectedILS := *defaults.ILS
	expectedILS.ShakeMax = 0.4
	assert.Equal(t, expectedILS, *cfg.Parameters.ILS)

	assert.Equal(t, defaults.BRKGA, cfg.Parameters.BRKGA)
	assert.Equal(t, defaults.GA, cfg.Parameters.GA)
	assert.Equal(t, defaults.VNS, cfg.Parameters.VNS)
	assert.Equal(t, defaults.LNS, cfg.Parameters.LNS)
}

func TestDecodeNonPositiveValuesKeepDefaults(t *testing.T) {
	cfg := mustDecode(t, "Parameters:\n  VNS:\n    Rate: -1\n  ILS:\n    MetropolisCriterion: false\n", false)
	defaults := config.DefaultParameters()

	assert.Equal(t, defaults.VNS, cfg.Parameters.VNS)
	assert.Equal(t, defaults.ILS, cfg.Parameters.ILS)
}

func TestDecodeMissingSectionsUseDefaults(t *testing.T) {
	cfg := mustDecode(t, "LogLevel: SILENT\nTimeLimitSeconds: 5\nMetaheuristics:\n  - \"GA\"\n", false)
	assert.Equal(t, config.DefaultParameters(), cfg.Parameters)
}

func TestDecodeEmptyInputReturnsEOF(t *testing.T) {
	for _, src := range []string{"", "\n\n", "# a\n# b\n"} {
		for _, strict := range []bool{false, true} {
			var c config.Config
			d := config.NewDecoder(strings.NewReader(src))
			d.Strict(strict)
			assert.Equal(t, io.EOF, d.Decode(&c), "%q", src)
		}
	}
}

func TestDecodeNullDocumentGivesDefaults(t *testing.T) {
	for _, src := range []string{"---\n", "~\n", "{}\n"} {
		t.Run(strings.TrimSpace(src), func(t *testing.T) {
			c := config.Config{
				LogLevel:       logger.LevelVerbose,
				LogType:        logger.TypeDashboard,
				TimeLimit:      time.Hour,
				MaxIterations:  9,
				SeedPool:       config.SeedPool{MaxSize: 1, InitialSize: 1},
				Metaheuristics: []config.Metaheuristic{{Type: rko.GA}},
				Parameters:     &config.Parameters{},
			}
			d := config.NewDecoder(strings.NewReader(src))
			require.NoError(t, d.Decode(&c))
			assert.Equal(t, defaultDecoded(), c)
		})
	}
}

func TestDecodeMultipleDocuments(t *testing.T) {
	d := config.NewDecoder(strings.NewReader("LogLevel: SILENT\n---\nLogLevel: VERBOSE\n"))

	var a, b, c config.Config
	require.NoError(t, d.Decode(&a))
	require.NoError(t, d.Decode(&b))
	assert.Equal(t, logger.LevelSilent, a.LogLevel)
	assert.Equal(t, logger.LevelVerbose, b.LogLevel)

	assert.Equal(t, io.EOF, d.Decode(&c))
	assert.Equal(t, io.EOF, d.Decode(&c))
}

func TestDecodeLenientIgnoresUnknown(t *testing.T) {
	tests := map[string]string{
		"unknown top-level key":   "Foo: 1\n",
		"top-level GA":            "GA:\n  PopulationSize: 123\n",
		"old Configuration key":   "Configuration:\n  GA:\n    PopulationSize: 123\n",
		"old SeedPoolMaxSize":     "SeedPoolMaxSize: 3\n",
		"lowercase keys":          "loglevel: SILENT\nlogtype: DASHBOARD\ntimelimitseconds: 10\nmaxiterations: 10\n",
		"MaxIterations in ILS":    "Parameters:\n  ILS:\n    MaxIterations: 88\n",
		"misspelled nested key":   "SeedPool:\n  Max: 4\n",
		"unknown section":         "Parameters:\n  TABU:\n    Size: 3\n",
		"unknown LogLevel":        "LogLevel: DEBUG\n",
		"unknown Metaheuristic":   "Metaheuristics:\n  - TABU\n",
		"unimplemented heuristic": "Metaheuristics:\n  - GRASP\n",
	}

	for name, src := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, defaultDecoded(), mustDecode(t, src, false))
		})
	}
}

func TestDecodeLogLevel(t *testing.T) {
	tests := []struct {
		value    string
		expected logger.Level
	}{
		{"SILENT", logger.LevelSilent},
		{"silent", logger.LevelSilent},
		{"INFO", logger.LevelInfo},
		{"Info", logger.LevelInfo},
		{"VERBOSE", logger.LevelVerbose},
		{"verbose", logger.LevelVerbose},
		{"DEBUG", logger.LevelInfo},
		{`""`, logger.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			assert.Equal(t, tt.expected, mustDecode(t, "LogLevel: "+tt.value+"\n", false).LogLevel)
		})
	}
}

func TestDecodeLogType(t *testing.T) {
	tests := []struct {
		value    string
		expected logger.Type
	}{
		{"DASHBOARD", logger.TypeDashboard},
		{"dashboard", logger.TypeDashboard},
		{"PRINT", logger.TypePrint},
		{"Print", logger.TypePrint},
		{"FILE", logger.TypePrint},
		{"CHANNEL", logger.TypePrint},
		{`""`, logger.TypePrint},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			assert.Equal(t, tt.expected, mustDecode(t, "LogType: "+tt.value+"\n", false).LogType)
		})
	}
}

func TestDecodeMetaheuristics(t *testing.T) {
	tests := []struct {
		name     string
		entries  []string
		expected []config.Metaheuristic
	}{
		{
			name:     "name without equals gives nil searches",
			entries:  []string{"VNS"},
			expected: []config.Metaheuristic{{Type: rko.VNS}},
		},
		{
			name:     "whitespace is trimmed",
			entries:  []string{"  ils  =  swap ,  MIRROR  , Nelder "},
			expected: []config.Metaheuristic{{Type: rko.ILS, Searches: []search.Type{search.Swap, search.Mirror, search.NelderMead}}},
		},
		{
			name:    "names are case-insensitive",
			entries: []string{"multistart", "Lns=farey"},
			expected: []config.Metaheuristic{
				{Type: rko.MultiStart},
				{Type: rko.LNS, Searches: []search.Type{search.Farey}},
			},
		},
		{
			name:     "unknown and unimplemented names are dropped",
			entries:  []string{"TABU=Swap", "GA", "", "NOPE", "GRASP"},
			expected: []config.Metaheuristic{{Type: rko.GA}},
		},
		{
			name:     "unknown searches and RVND are filtered out",
			entries:  []string{"SA=Swap, Unknown, RVND, rvnd, Mirror"},
			expected: []config.Metaheuristic{{Type: rko.SA, Searches: []search.Type{search.Swap, search.Mirror}}},
		},
		{
			name:    "only invalid searches leave nil searches",
			entries: []string{"BRKGA=RVND, Bogus", "VNS="},
			expected: []config.Metaheuristic{
				{Type: rko.BRKGA},
				{Type: rko.VNS},
			},
		},
		{
			name:     "no entries gives nil",
			entries:  nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := "Metaheuristics:\n"
			for _, e := range tt.entries {
				src += "  - \"" + e + "\"\n"
			}
			cfg := mustDecode(t, src, false)
			assert.Equal(t, tt.expected, cfg.Metaheuristics)
		})
	}
}

func TestDecodeMaxIterations(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected int
	}{
		{"positive", "MaxIterations: 250\n", 250},
		{"zero", "MaxIterations: 0\n", 0},
		{"negative is clamped to zero", "MaxIterations: -5\n", 0},
		{"unset", "LogLevel: INFO\n", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, mustDecode(t, tt.src, false).MaxIterations)
		})
	}
}

func TestDecodeTimeLimitSeconds(t *testing.T) {
	assert.Equal(t, 90*time.Second, mustDecode(t, "TimeLimitSeconds: 90\n", false).TimeLimit)
	assert.Equal(t, time.Duration(0), mustDecode(t, "LogLevel: INFO\n", false).TimeLimit)

	cfg := mustDecode(t, "TimeLimitSeconds: 10\nMaxIterations: 7\n", false)
	assert.Equal(t, 10*time.Second, cfg.TimeLimit)
	assert.Equal(t, 7, cfg.MaxIterations)
}

func TestDecodeSeedPool(t *testing.T) {
	def := solution.DefaultSize
	tests := []struct {
		name     string
		src      string
		expected config.SeedPool
	}{
		{"unset", "LogLevel: INFO\n", config.SeedPool{MaxSize: def, InitialSize: def}},
		{"both set", "SeedPool:\n  MaxSize: 50\n  InitialSize: 20\n", config.SeedPool{MaxSize: 50, InitialSize: 20}},
		{"max raised to initial", "SeedPool:\n  MaxSize: 5\n  InitialSize: 20\n", config.SeedPool{MaxSize: 20, InitialSize: 20}},
		{"negative max means no cap", "SeedPool:\n  MaxSize: -1\n", config.SeedPool{MaxSize: -1, InitialSize: def}},
		{"non-positive initial uses default", "SeedPool:\n  InitialSize: -3\n", config.SeedPool{MaxSize: def, InitialSize: def}},
		{"only initial", "SeedPool:\n  InitialSize: 30\n", config.SeedPool{MaxSize: 30, InitialSize: 30}},
		{"misspelled nested key is ignored", "SeedPool:\n  Max: 4\n", config.SeedPool{MaxSize: def, InitialSize: def}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, mustDecode(t, tt.src, false).SeedPool)
		})
	}
}

func TestDecodeStrictUnknownField(t *testing.T) {
	tests := []struct {
		src     string
		message string
	}{
		{"Foo: 1\n", `config: line 1: unknown field "Foo"`},
		{"Parameters:\n  ILS:\n    MaxIterations: 88\n", `config: line 3: unknown field "MaxIterations"`},
		{"GA:\n  PopulationSize: 1\n", `config: line 1: unknown field "GA"`},
		{"Configuration: {}\n", `config: line 1: unknown field "Configuration"`},
		{"SeedPoolMaxSize: 3\n", `config: line 1: unknown field "SeedPoolMaxSize"`},
		{"SeedPool:\n  Max: 4\n", `config: line 2: unknown field "Max"`},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			_, err := decodeOne(tt.src, true)
			require.Error(t, err)
			assert.ErrorIs(t, err, config.ErrUnknownField)
			assert.NotErrorIs(t, err, config.ErrUnknownValue)
			assert.Equal(t, tt.message, err.Error())
		})
	}
}

func TestDecodeStrictUnknownValue(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		message  string
		sentinel bool
	}{
		{"LogLevel", "LogLevel: DEBUG\n", `config: LogLevel: unknown value "DEBUG"`, true},
		{"LogType", "LogType: FILE\n", `config: LogType: unknown value "FILE"`, true},
		{"LogType CHANNEL", "LogType: CHANNEL\n", `config: LogType: unknown value "CHANNEL"`, true},
		{"metaheuristic", "Metaheuristics:\n  - \"TABU=Swap\"\n", `config: Metaheuristics[0]: unknown value "TABU"`, true},
		{"search", "Metaheuristics:\n  - GA\n  - \"SA=Swap, Bogus\"\n", `config: Metaheuristics[1]: unknown value "Bogus"`, true},
		{"empty entry", "Metaheuristics:\n  - \"\"\n", `config: Metaheuristics[0]: unknown value ""`, true},
		{"unimplemented", "Metaheuristics:\n  - GRASP\n", `config: Metaheuristics[0]: Greedy Randomized Adaptive Search Procedure is not implemented`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decodeOne(tt.src, true)
			require.Error(t, err)
			assert.Equal(t, tt.sentinel, errors.Is(err, config.ErrUnknownValue))
			assert.Equal(t, !tt.sentinel, errors.Is(err, config.ErrNotImplemented))
			assert.NotErrorIs(t, err, config.ErrUnknownField)
			assert.Equal(t, tt.message, err.Error())
		})
	}
}

func TestDecodeLenientDropsUnimplemented(t *testing.T) {
	cfg := mustDecode(t, "Metaheuristics:\n  - GA\n  - GRASP\n", false)
	assert.Equal(t, []config.Metaheuristic{{Type: rko.GA}}, cfg.Metaheuristics)
}

func TestDecodeStrictAcceptsMissingAndDefaults(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected []config.Metaheuristic
	}{
		{"empty mapping", "{}\n", nil},
		{"empty LogLevel", `LogLevel: ""` + "\n", nil},
		{"null LogLevel", "LogLevel:\n", nil},
		{"lowercase values", "LogLevel: info\nLogType: dashboard\n", nil},
		{"empty searches", "Metaheuristics:\n  - \"VNS=\"\n", []config.Metaheuristic{{Type: rko.VNS}}},
		{"RVND is accepted and skipped", "Metaheuristics:\n  - \"SA=Swap, RVND\"\n", []config.Metaheuristic{{Type: rko.SA, Searches: []search.Type{search.Swap}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := mustDecode(t, tt.src, true)
			assert.Equal(t, tt.expected, cfg.Metaheuristics)
		})
	}
}

func TestDecodeStrictCollectsAllErrors(t *testing.T) {
	src := "Foo: 1\nLogLevel: DEBUG\nBar: 2\nLogType: X\nMetaheuristics:\n  - TABU\n"
	_, err := decodeOne(src, true)
	require.Error(t, err)

	assert.ErrorIs(t, err, config.ErrUnknownField)
	assert.ErrorIs(t, err, config.ErrUnknownValue)

	assert.Equal(t, []string{
		`config: line 1: unknown field "Foo"`,
		`config: line 3: unknown field "Bar"`,
		`config: LogLevel: unknown value "DEBUG"`,
		`config: LogType: unknown value "X"`,
		`config: Metaheuristics[0]: unknown value "TABU"`,
	}, strings.Split(err.Error(), "\n"))

	joined, ok := err.(interface{ Unwrap() []error })
	require.True(t, ok)
	for _, e := range joined.Unwrap() {
		field := errors.Is(e, config.ErrUnknownField)
		value := errors.Is(e, config.ErrUnknownValue)
		assert.NotEqual(t, field, value, e.Error())
	}
}

func TestDecodeStrictToggle(t *testing.T) {
	d := config.NewDecoder(strings.NewReader("Foo: 1\n---\nFoo: 1\n"))

	var c config.Config
	d.Strict(true)
	assert.ErrorIs(t, d.Decode(&c), config.ErrUnknownField)

	d.Strict(false)
	require.NoError(t, d.Decode(&c))
	assert.Equal(t, defaultDecoded(), c)
}

func TestDecodeErrorLeavesConfigUntouched(t *testing.T) {
	c := config.Config{MaxIterations: 7}
	d := config.NewDecoder(strings.NewReader("Foo: 1\nMaxIterations: 3\n"))
	d.Strict(true)
	require.Error(t, d.Decode(&c))
	assert.Equal(t, config.Config{MaxIterations: 7}, c)
}

func TestDecodeTypeError(t *testing.T) {
	for _, strict := range []bool{false, true} {
		_, err := decodeOne("TimeLimitSeconds: abc\n", strict)
		require.Error(t, err)
		assert.NotErrorIs(t, err, config.ErrUnknownField)
		assert.True(t, strings.HasPrefix(err.Error(), "config: line 1: cannot unmarshal"), err.Error())
	}
}

func TestDecodeSyntaxError(t *testing.T) {
	for _, strict := range []bool{false, true} {
		_, err := decodeOne("LogLevel: [INFO\n", strict)
		require.Error(t, err)
		assert.True(t, strings.HasPrefix(err.Error(), "config: yaml:"), err.Error())
	}
}

func TestDecodeNilConfig(t *testing.T) {
	d := config.NewDecoder(strings.NewReader("LogLevel: SILENT\n"))
	err := d.Decode(nil)
	require.Error(t, err)
	assert.Equal(t, "config: Decode called with nil *Config", err.Error())

	var c config.Config
	require.NoError(t, d.Decode(&c))
	assert.Equal(t, logger.LevelSilent, c.LogLevel)
}

func TestReadFile(t *testing.T) {
	path := writeTempFile(t, "LogLevel: SILENT\nMaxIterations: 1\n---\nLogLevel: VERBOSE\nMaxIterations: 2\n")

	cfgs, err := config.ReadFile(path)
	require.NoError(t, err)
	require.Len(t, cfgs, 2)
	assert.Equal(t, logger.LevelSilent, cfgs[0].LogLevel)
	assert.Equal(t, 1, cfgs[0].MaxIterations)
	assert.Equal(t, logger.LevelVerbose, cfgs[1].LogLevel)
	assert.Equal(t, 2, cfgs[1].MaxIterations)
}

func TestReadFileStrict(t *testing.T) {
	path := writeTempFile(t, "LogLevel: SILENT\n---\nFoo: 1\n")

	cfgs, err := config.ReadFile(path)
	require.Error(t, err)
	assert.Nil(t, cfgs)
	assert.ErrorIs(t, err, config.ErrUnknownField)
	assert.True(t, strings.HasPrefix(err.Error(), path+": document 2: config: line "), err.Error())
}

func TestReadFileNoConfig(t *testing.T) {
	for _, src := range []string{"", "\n\n", "# only\n# comments\n"} {
		path := writeTempFile(t, src)
		cfgs, err := config.ReadFile(path)
		require.Error(t, err, "%q", src)
		assert.Nil(t, cfgs)
		assert.ErrorIs(t, err, config.ErrNoConfig)
		assert.Equal(t, path+": config: no configuration found", err.Error())
	}

	cfgs, err := config.ReadFile(writeTempFile(t, "---\n"))
	require.NoError(t, err)
	require.Len(t, cfgs, 1)
	assert.Equal(t, defaultDecoded(), cfgs[0])
}

func TestReadFileMissing(t *testing.T) {
	cfgs, err := config.ReadFile(filepath.Join(t.TempDir(), "does-not-exist.yaml"))
	require.Error(t, err)
	assert.Nil(t, cfgs)
	assert.ErrorIs(t, err, fs.ErrNotExist)
}
