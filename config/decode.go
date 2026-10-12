package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/RKO-solver/rko-go"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic/search"
	"go.yaml.in/yaml/v3"
)

var unknownFieldRE = regexp.MustCompile(`^line (\d+): field (.+) not found in type \S+$`)

// Decoder reads Config documents from a YAML stream.
type Decoder struct {
	dec    *yaml.Decoder
	strict bool
}

// NewDecoder returns a lenient Decoder reading YAML documents from r.
func NewDecoder(r io.Reader) *Decoder { return &Decoder{dec: yaml.NewDecoder(r)} }

// Strict sets whether Decode rejects unknown keys and values.
func (d *Decoder) Strict(enable bool) { d.strict = enable }

// Decode reads the next document into cfg, returning io.EOF at the end.
func (d *Decoder) Decode(cfg *Config) error {
	if cfg == nil {
		return errors.New("config: Decode called with nil *Config")
	}

	var y yamlConfig
	d.dec.KnownFields(d.strict)
	err := d.dec.Decode(&y)
	if err == io.EOF {
		return io.EOF
	}

	var errs []error
	var te *yaml.TypeError
	if err != nil {
		if !errors.As(err, &te) {
			return fmt.Errorf("config: %w", err)
		}
		errs = typeErrors(te)
	}

	c, valueErrs := convert(&y)
	if d.strict {
		errs = append(errs, valueErrs...)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	*cfg = c
	return nil
}

func typeErrors(te *yaml.TypeError) []error {
	errs := make([]error, 0, len(te.Errors))
	for _, s := range te.Errors {
		if m := unknownFieldRE.FindStringSubmatch(s); m != nil {
			errs = append(errs, fmt.Errorf("config: line %s: %w %q", m[1], ErrUnknownField, m[2]))
			continue
		}
		errs = append(errs, errors.New("config: "+s))
	}
	return errs
}

func convert(y *yamlConfig) (Config, []error) {
	var errs []error

	level := logger.DefaultLevel
	if y.LogLevel != "" {
		if v, err := logger.ParseLevel(y.LogLevel); err != nil {
			errs = append(errs, fmt.Errorf("config: LogLevel: %w %q", ErrUnknownValue, y.LogLevel))
		} else {
			level = v
		}
	}

	logType := logger.TypePrint
	if y.LogType != "" {
		if v, err := logger.ParseType(y.LogType); err != nil {
			errs = append(errs, fmt.Errorf("config: LogType: %w %q", ErrUnknownValue, y.LogType))
		} else {
			logType = v
		}
	}

	maxSize, initialSize := seedPoolSizes(y.SeedPool.MaxSize, y.SeedPool.InitialSize)

	mhs, mhErrs := parseMetaheuristics(y.Metaheuristics)
	errs = append(errs, mhErrs...)

	params := DefaultParameters()
	for _, opt := range []parameterOption{
		withGA(y.Parameters.GA),
		withBRKGA(y.Parameters.BRKGA),
		withSA(y.Parameters.SA),
		withVNS(y.Parameters.VNS),
		withILS(y.Parameters.ILS),
		withLNS(y.Parameters.LNS),
	} {
		opt(params)
	}

	return Config{
		LogLevel:       level,
		LogType:        logType,
		TimeLimit:      time.Duration(y.TimeLimitSeconds * int64(time.Second)),
		MaxIterations:  max(y.MaxIterations, 0),
		SeedPool:       SeedPool{MaxSize: maxSize, InitialSize: initialSize},
		Metaheuristics: mhs,
		Parameters:     params,
	}, errs
}

func parseMetaheuristics(entries []string) ([]Metaheuristic, []error) {
	var out []Metaheuristic
	var errs []error

	for i, e := range entries {
		parts := strings.SplitN(e, "=", 2)
		name := strings.TrimSpace(parts[0])
		mt := rko.ParseMetaheuristic(name)

		keep := true
		if mt < 0 {
			errs = append(errs, fmt.Errorf("config: Metaheuristics[%d]: %w %q", i, ErrUnknownValue, name))
			keep = false
		} else if err := metaheuristicError(i, mt); err != nil {
			errs = append(errs, err)
			keep = false
		}

		var searches []search.Type
		if len(parts) == 2 {
			for _, comp := range strings.Split(parts[1], ",") {
				comp = strings.TrimSpace(comp)
				if comp == "" {
					continue
				}
				st := search.ParseType(comp)
				if st < 0 {
					errs = append(errs, fmt.Errorf("config: Metaheuristics[%d]: %w %q", i, ErrUnknownValue, comp))
					continue
				}
				if st == search.RVND {
					continue
				}
				searches = append(searches, st)
			}
		}

		if keep {
			out = append(out, Metaheuristic{Type: mt, Searches: searches})
		}
	}

	return out, errs
}

// ReadFile reads every configuration in the file at path, in strict mode.
func ReadFile(path string) ([]Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d := NewDecoder(f)
	d.Strict(true)

	var cfgs []Config
	for {
		var c Config
		err = d.Decode(&c)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%s: document %d: %w", path, len(cfgs)+1, err)
		}
		cfgs = append(cfgs, c)
	}

	if len(cfgs) == 0 {
		return nil, fmt.Errorf("%s: config: %w", path, ErrNoConfig)
	}
	return cfgs, nil
}
