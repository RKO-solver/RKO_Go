package config_test

import (
	"context"
	"math"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/RKO-solver/rko-go"
	"github.com/RKO-solver/rko-go/config"
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testEnv is a cheap Environment with n keys.
type testEnv struct{ n int }

func (e testEnv) NumKeys() int        { return e.n }
func (e testEnv) ProblemName() string { return "test" }
func (e testEnv) Cost(r definition.RandomKey) int {
	c := 0
	for i, v := range r {
		c += int(v*1000) * (i + 1)
	}
	return c
}
func (e testEnv) Decode(r definition.RandomKey) any { return r.Clone() }
func (e testEnv) SwapSearch() [][2]int              { return [][2]int{{0, e.n}} }

// solveWithin runs solver.Solve and fails the test if it does not return within d.
func solveWithin(t *testing.T, solver *rko.Solver, d time.Duration) {
	t.Helper()
	done := make(chan struct{})
	go func() { solver.Solve(context.Background()); close(done) }()
	select {
	case <-done:
	case <-time.After(d):
		t.Fatal("solver did not stop")
	}
}

func silentConfig() *config.Config {
	return &config.Config{LogLevel: logger.LevelSilent, LogType: logger.TypePrint, MaxIterations: 3}
}

func cloneConfig(cfg *config.Config) config.Config {
	c := *cfg
	c.Metaheuristics = make([]config.Metaheuristic, len(cfg.Metaheuristics))
	for i, m := range cfg.Metaheuristics {
		c.Metaheuristics[i] = config.Metaheuristic{Type: m.Type, Searches: append([]search.Type(nil), m.Searches...)}
	}
	if cfg.Parameters != nil {
		p := *cfg.Parameters
		if p.BRKGA != nil {
			v := *p.BRKGA
			p.BRKGA = &v
		}
		if p.GA != nil {
			v := *p.GA
			p.GA = &v
		}
		if p.ILS != nil {
			v := *p.ILS
			p.ILS = &v
		}
		if p.SA != nil {
			v := *p.SA
			p.SA = &v
		}
		if p.VNS != nil {
			v := *p.VNS
			p.VNS = &v
		}
		if p.LNS != nil {
			v := *p.LNS
			p.LNS = &v
		}
		c.Parameters = &p
	}
	return c
}

// mhConfigs returns the configuration pointer held by each solver inside s.
func mhConfigs(t *testing.T, s *rko.Solver) []reflect.Value {
	t.Helper()
	solvers := reflect.ValueOf(s).Elem().FieldByName("solvers")
	require.True(t, solvers.IsValid())
	out := make([]reflect.Value, solvers.Len())
	for i := range out {
		c := solvers.Index(i).Elem().Elem().FieldByName("configuration")
		require.True(t, c.IsValid())
		out[i] = c
	}
	return out
}

func newSolver(t *testing.T, cfg *config.Config, opts ...config.Option) (*rko.Solver, logger.Recorder) {
	t.Helper()
	solver, log, err := cfg.NewSolver(testEnv{n: 30}, opts...)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, log.Close()) })
	return solver, log
}

func allMetaheuristics() []config.Metaheuristic {
	var out []config.Metaheuristic
	for _, mh := range []rko.MetaHeuristic{rko.MULTISTART, rko.SA, rko.GA, rko.VNS, rko.ILS, rko.BRKGA, rko.LNS} {
		out = append(out, config.Metaheuristic{Type: mh, Searches: []search.Type{search.Swap}})
	}
	return out
}

func TestNewSolverMaxIterationsStopsEveryMetaheuristic(t *testing.T) {
	const maxIterations = 5
	params := config.DefaultParameters()
	params.SA.Iterations = 10
	cfg := silentConfig()
	cfg.MaxIterations = maxIterations
	cfg.Metaheuristics = allMetaheuristics()
	cfg.Parameters = params

	solver, log := newSolver(t, cfg, config.WithSeed(42))
	solveWithin(t, solver, 30*time.Second)

	report := log.Report()
	require.Len(t, report, len(cfg.Metaheuristics))
	for _, info := range report {
		assert.Len(t, info.Samples, maxIterations, info.Name)
	}
}

func TestNewSolverTimeLimitStopsBeforeMaxIterations(t *testing.T) {
	cfg := silentConfig()
	cfg.TimeLimit = 200 * time.Millisecond
	cfg.MaxIterations = math.MaxInt
	cfg.Metaheuristics = []config.Metaheuristic{{Type: rko.MULTISTART, Searches: []search.Type{search.Swap}}}

	solver, _ := newSolver(t, cfg, config.WithSeed(42))
	solveWithin(t, solver, 10*time.Second)
}

func TestNewSolverSeedPoolInitialSize(t *testing.T) {
	tests := []struct {
		name string
		pool config.SeedPool
		want int
	}{
		{"zero pool uses default", config.SeedPool{}, solution.DefaultSize},
		{"larger than default", config.SeedPool{MaxSize: -1, InitialSize: 25}, 25},
		{"smaller than default", config.SeedPool{MaxSize: -1, InitialSize: 3}, 3},
		{"max raised to initial", config.SeedPool{MaxSize: 4, InitialSize: 15}, 15},
		{"negative initial uses default", config.SeedPool{MaxSize: 50, InitialSize: -3}, solution.DefaultSize},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := silentConfig()
			cfg.SeedPool = tt.pool
			solver, _ := newSolver(t, cfg, config.WithSeed(1))
			assert.Equal(t, tt.want, solver.GetSolutionPool().Len())
		})
	}
}

func TestNewSolverSeedPoolMaxSizeCap(t *testing.T) {
	tests := []struct {
		name string
		pool config.SeedPool
		cap  int
	}{
		{"zero pool capped at default", config.SeedPool{}, solution.DefaultSize},
		{"explicit cap", config.SeedPool{MaxSize: 12, InitialSize: 5}, 12},
		{"cap equal to initial", config.SeedPool{MaxSize: 6, InitialSize: 6}, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := silentConfig()
			cfg.MaxIterations = 200
			cfg.SeedPool = tt.pool
			cfg.Metaheuristics = []config.Metaheuristic{{Type: rko.MULTISTART, Searches: []search.Type{search.Swap}}}
			solver, _ := newSolver(t, cfg, config.WithSeed(7))
			solveWithin(t, solver, 30*time.Second)
			assert.LessOrEqual(t, solver.GetSolutionPool().Len(), tt.cap)
		})
	}
}

func TestNewSolverInvalidLogType(t *testing.T) {
	cfg := silentConfig()
	cfg.LogType = 7
	solver, log, err := cfg.NewSolver(testEnv{n: 30})
	assert.Nil(t, solver)
	assert.Nil(t, log)
	require.ErrorIs(t, err, config.ErrUnknownValue)
	assert.EqualError(t, err, "config: LogType: unknown value 7")
}

func TestNewSolverInvalidLogLevel(t *testing.T) {
	cfg := silentConfig()
	cfg.LogLevel = 9
	solver, log, err := cfg.NewSolver(testEnv{n: 30})
	assert.Nil(t, solver)
	assert.Nil(t, log)
	require.ErrorIs(t, err, config.ErrUnknownValue)
	assert.EqualError(t, err, "config: LogLevel: unknown value 9")
}

func TestNewSolverNilParameterSection(t *testing.T) {
	sections := map[string]func(p *config.Parameters){
		"BRKGA": func(p *config.Parameters) { p.BRKGA = nil },
		"GA":    func(p *config.Parameters) { p.GA = nil },
		"ILS":   func(p *config.Parameters) { p.ILS = nil },
		"SA":    func(p *config.Parameters) { p.SA = nil },
		"VNS":   func(p *config.Parameters) { p.VNS = nil },
		"LNS":   func(p *config.Parameters) { p.LNS = nil },
	}
	for name, clear := range sections {
		t.Run(name, func(t *testing.T) {
			cfg := silentConfig()
			cfg.Parameters = config.DefaultParameters()
			clear(cfg.Parameters)
			solver, log, err := cfg.NewSolver(testEnv{n: 30})
			assert.Nil(t, solver)
			assert.Nil(t, log)
			assert.EqualError(t, err, "config: Parameters."+name+" is nil")
		})
	}

	t.Run("combined", func(t *testing.T) {
		cfg := silentConfig()
		cfg.LogType = 7
		cfg.Metaheuristics = []config.Metaheuristic{{Type: rko.GA}, {Type: rko.IPR}}
		cfg.Parameters = config.DefaultParameters()
		cfg.Parameters.GA = nil
		cfg.Parameters.LNS = nil
		_, _, err := cfg.NewSolver(testEnv{n: 30})
		assert.EqualError(t, err, "config: LogType: unknown value 7\n"+
			"config: Metaheuristics[1]: Implicit Path Relinking is not implemented\n"+
			"config: Parameters.GA is nil\n"+
			"config: Parameters.LNS is nil")
	})
}

func TestNewSolverUnimplementedMetaheuristic(t *testing.T) {
	cfg := silentConfig()
	cfg.Metaheuristics = []config.Metaheuristic{{Type: rko.GA}, {Type: rko.GRASP}}
	solver, log, err := cfg.NewSolver(testEnv{n: 30})
	assert.Nil(t, solver)
	assert.Nil(t, log)
	require.ErrorIs(t, err, config.ErrNotImplemented)
	assert.EqualError(t, err, "config: Metaheuristics[1]: Greedy Randomized Adaptive Search Procedure is not implemented")

	cfg.Metaheuristics = []config.Metaheuristic{{Type: 99}}
	solver, log, err = cfg.NewSolver(testEnv{n: 30})
	assert.Nil(t, solver)
	assert.Nil(t, log)
	require.ErrorIs(t, err, config.ErrUnknownValue)
	assert.EqualError(t, err, "config: Metaheuristics[0]: unknown value 99")
}

func TestNewSolverNilParametersUsesDefaults(t *testing.T) {
	cfg := silentConfig()
	cfg.Metaheuristics = allMetaheuristics()
	solver, _ := newSolver(t, cfg, config.WithSeed(3))
	solveWithin(t, solver, 30*time.Second)
	assert.Nil(t, cfg.Parameters)
}

func TestNewSolverDoesNotMutateConfig(t *testing.T) {
	for _, name := range []string{"explicit parameters", "nil parameters"} {
		t.Run(name, func(t *testing.T) {
			cfg := silentConfig()
			cfg.Metaheuristics = []config.Metaheuristic{
				{Type: rko.ILS, Searches: []search.Type{search.Swap, search.Mirror}},
				{Type: rko.SA, Searches: []search.Type{search.Farey}},
			}
			if name == "explicit parameters" {
				cfg.Parameters = config.DefaultParameters()
			}
			before := cloneConfig(cfg)
			p := cfg.Parameters
			var ilsPtr, saPtr any
			if p != nil {
				ilsPtr, saPtr = p.ILS, p.SA
			}

			solver, _ := newSolver(t, cfg, config.WithSeed(5))
			solveWithin(t, solver, 30*time.Second)

			assert.Equal(t, before, *cfg)
			if p == nil {
				assert.Nil(t, cfg.Parameters)
				return
			}
			assert.Same(t, p, cfg.Parameters)
			assert.Same(t, ilsPtr, cfg.Parameters.ILS)
			assert.Same(t, saPtr, cfg.Parameters.SA)
		})
	}
}

func TestNewSolverCopiesParameters(t *testing.T) {
	cfg := silentConfig()
	cfg.Parameters = config.DefaultParameters()
	cfg.Metaheuristics = []config.Metaheuristic{{Type: rko.ILS}, {Type: rko.ILS}, {Type: rko.SA}}
	solver, _ := newSolver(t, cfg, config.WithSeed(9))

	confs := mhConfigs(t, solver)
	require.Len(t, confs, 3)

	ilsPtr := reflect.ValueOf(cfg.Parameters.ILS).Pointer()
	saPtr := reflect.ValueOf(cfg.Parameters.SA).Pointer()
	assert.NotEqual(t, ilsPtr, confs[0].Pointer())
	assert.NotEqual(t, ilsPtr, confs[1].Pointer())
	assert.NotEqual(t, confs[0].Pointer(), confs[1].Pointer())
	assert.NotEqual(t, saPtr, confs[2].Pointer())

	want := config.DefaultParameters().ILS.ShakeMax
	cfg.Parameters.ILS.ShakeMax = 99
	assert.Equal(t, want, confs[0].Elem().FieldByName("ShakeMax").Float())
	assert.Equal(t, want, confs[1].Elem().FieldByName("ShakeMax").Float())
}

func TestNewSolverConcurrent(t *testing.T) {
	for _, name := range []string{"explicit parameters", "nil parameters"} {
		t.Run(name, func(t *testing.T) {
			cfg := silentConfig()
			cfg.MaxIterations = 2
			cfg.Metaheuristics = allMetaheuristics()
			if name == "explicit parameters" {
				cfg.Parameters = config.DefaultParameters()
			}

			const n = 4
			solvers := make([]*rko.Solver, 2*n)
			logs := make([]logger.Recorder, 2*n)
			errs := make([]error, 2*n)
			var wg sync.WaitGroup
			for i := range solvers {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var opts []config.Option
					if i < n {
						opts = append(opts, config.WithSeed(uint64(i)))
					}
					solvers[i], logs[i], errs[i] = cfg.NewSolver(testEnv{n: 30}, opts...)
				}()
			}
			wg.Wait()
			for i := range solvers {
				require.NoError(t, errs[i])
				log := logs[i]
				t.Cleanup(func() { require.NoError(t, log.Close()) })
			}

			for _, s := range solvers {
				wg.Add(1)
				go func() {
					defer wg.Done()
					s.Solve(context.Background())
				}()
			}
			if cfg.Parameters != nil {
				for i := range 50 {
					cfg.Parameters.ILS.ShakeMax = float64(i)
					cfg.Parameters.SA.Alpha = float64(i)
				}
			}
			wg.Wait()
		})
	}
}

func TestWithSeedSamePool(t *testing.T) {
	opt := config.WithSeed(42)
	costs := make([]int, 2)
	for i := range costs {
		cfg := silentConfig()
		cfg.Metaheuristics = []config.Metaheuristic{{Type: rko.ILS, Searches: []search.Type{search.Swap}}}
		solver, _ := newSolver(t, cfg, opt)
		solveWithin(t, solver, 30*time.Second)
		costs[i] = solver.GetSolutionPool().BestCost()
	}
	assert.Equal(t, costs[0], costs[1])
}

func TestNewSolverNoOptions(t *testing.T) {
	cfg := silentConfig()
	cfg.Metaheuristics = allMetaheuristics()
	solver, log := newSolver(t, cfg)
	solveWithin(t, solver, 30*time.Second)
	assert.Len(t, log.Report(), len(cfg.Metaheuristics))
}
