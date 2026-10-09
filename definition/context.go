package definition

import (
	"context"
	"time"
)

// contextKey is the type of the keys this package stores in a context.
type contextKey int

const (
	startTimeKey contextKey = iota
	maxIterationsKey
)

// IterationCounter counts a metaheuristic's iterations and cancels its context at the limit.
type IterationCounter struct {
	iteration, max int
	cancel         context.CancelFunc // nil when there is no limit
}

// Iteration returns the number of completed iterations.
func (l *IterationCounter) Iteration() int {
	return l.iteration
}

// Step counts one iteration and cancels the context when the limit is reached.
func (l *IterationCounter) Step() {
	l.iteration++
	if l.cancel != nil && l.iteration >= l.max {
		l.cancel()
	}
}

// Stop releases the context; call it with defer.
func (l *IterationCounter) Stop() {
	if l.cancel != nil {
		l.cancel()
	}
}

// WithStartTime returns a context carrying the given start time.
func WithStartTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, startTimeKey, t)
}

// StartTime returns the start time stored in ctx, or time.Now().
func StartTime(ctx context.Context) time.Time {
	if t, ok := ctx.Value(startTimeKey).(time.Time); ok {
		return t
	}
	return time.Now()
}

// WithMaxIterations returns a context carrying the iteration limit of each metaheuristic.
func WithMaxIterations(ctx context.Context, n int) context.Context {
	return context.WithValue(ctx, maxIterationsKey, n)
}

// maxIterations returns the limit stored in ctx and whether it is positive.
func maxIterations(ctx context.Context) (int, bool) {
	if n, ok := ctx.Value(maxIterationsKey).(int); ok && n > 0 {
		return n, true
	}
	return -1, false
}

// GetMHContext returns a context for one metaheuristic and its IterationCounter, which cancels that context at the limit set by WithMaxIterations.
func GetMHContext(ctx context.Context) (context.Context, *IterationCounter) {
	maxIteration, hasLimit := maxIterations(ctx)
	if !hasLimit {
		return ctx, &IterationCounter{}
	}
	ctx, cancel := context.WithCancel(ctx)
	return ctx, &IterationCounter{max: maxIteration, cancel: cancel}
}
