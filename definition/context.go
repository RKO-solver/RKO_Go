package definition

import (
	"context"
	"time"
)

type contextKey int

const (
	startTimeKey contextKey = iota
	maxIterationsKey
)

// IterationCounter counts a metaheuristic's iterations and cancels its context at the limit.
type IterationCounter struct {
	iteration, max int
	cancel         context.CancelFunc
}

// Iteration returns the number of completed iterations.
func (l *IterationCounter) Iteration() int {
	return l.iteration
}

// Limit returns the iteration limit, or -1 if there is none.
func (l *IterationCounter) Limit() int {
	return l.max
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

// WithStartTime stores the start time and, if timeLimit > 0, a deadline timeLimit later.
func WithStartTime(ctx context.Context, timeLimit time.Duration) (context.Context, context.CancelFunc) {
	start := time.Now()
	ctx = context.WithValue(ctx, startTimeKey, start)
	if timeLimit > 0 {
		return context.WithDeadline(ctx, start.Add(timeLimit))
	}
	return ctx, func() {}
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

func maxIterations(ctx context.Context) (int, bool) {
	if n, ok := ctx.Value(maxIterationsKey).(int); ok && n > 0 {
		return n, true
	}
	return -1, false
}

// NewSolverContext returns a context for one metaheuristic and its IterationCounter.
func NewSolverContext(ctx context.Context) (context.Context, *IterationCounter) {
	maxIteration, hasLimit := maxIterations(ctx)
	if !hasLimit {
		return ctx, &IterationCounter{max: maxIteration}
	}
	ctx, cancel := context.WithCancel(ctx)
	return ctx, &IterationCounter{max: maxIteration, cancel: cancel}
}
