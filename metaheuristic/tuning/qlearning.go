package tuning

import (
	"fmt"
	"math"
	"strconv"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/random"
)

// DefaultDiscount is the gamma used by [NewQLearning] unless [WithDiscount] is given.
const DefaultDiscount = 0.8

const (
	// Initial Q-values are uniform in [initialQMin, initialQMax): slightly
	// above the non-improving rewards (at most 0), so untried actions get
	// tried, and random, so each state's first best action is random.
	initialQMin = 0.005 // lower bound of initial Q-values
	initialQMax = 0.01  // upper bound of initial Q-values

	// Exploration schedule; see epsilon.
	epsilonMin      = 0.1 // exploration floor
	epsilonPeak     = 1.0 // exploration at the first restart
	epsilonPeakStep = 0.1 // peak decrease per restart
	restartPeriod   = 0.1 // progress between restarts

	// Learning rate schedule; see learningRate.
	learningRateStart = 1.0 // learning rate at progress 0
	learningRateDecay = 0.9 // learning rate decrease up to progress 1
)

// Option configures a [QLearning] created by [NewQLearning].
type Option func(*QLearning)

// WithDiscount sets the discount factor gamma; it panics outside [0, 1) or NaN.
func WithDiscount(discount float64) Option {
	// It panics when called, not when the option is applied.
	// Negated range check, so NaN, which fails every comparison, is rejected.
	if !(discount >= 0 && discount < 1) {
		panic(fmt.Sprintf("tuning: discount must be in [0, 1), got %v", discount))
	}
	return func(q *QLearning) { q.discount = discount }
}

// QLearning is a [Tuner] using tabular Q-learning over a [Space].
type QLearning struct {
	// It uses epsilon-greedy exploration over the configurations of the Space.
	// Each action keeps the configuration or changes one parameter to another of
	// its values.
	//
	// A QLearning is not safe for concurrent use: use one per goroutine, each
	// with its own [random.Generator].

	// transitions[state*actions+action] is the state reached by taking
	// action in state, precomputed so moving is a lookup.
	transitions []int
	// qValues[state*actions+action] is Q(state, action), row-major.
	qValues []float64
	// bestActions[state] caches the action with the largest Q-value in row
	// state, the lowest one on ties.
	bestActions []int
	// actions is the row length: 1 + sum(count-1), the same for every state.
	actions int
	// space is a copy of the Space at creation, so later Adds don't change
	// it; Policy reads its formats, counts and strides.
	space Space
	// discount is the discount factor gamma.
	discount float64
	// state is the current state, from which the next Choose acts.
	state int
	// action and next record the last Choose: the action taken from state
	// and the state it reached. Learn updates Q(state, action) and then
	// moves state to next.
	action int
	next   int
	// generator draws the initial Q-values, the start state and the
	// epsilon-greedy choices.
	generator *random.Generator
}

// Compile-time checks that *QLearning implements Tuner and Reporter.
var (
	_ Tuner    = (*QLearning)(nil)
	_ Reporter = (*QLearning)(nil)
)

// NewQLearning creates a QLearning over the configurations of space, drawing from generator.
func NewQLearning(space *Space, generator *random.Generator, options ...Option) *QLearning {
	// The tables are sized from space now, so add every parameter first.
	states := space.States()
	// Action 0 (stay), plus one action per other value of each parameter.
	actions := 1
	for _, count := range space.counts {
		actions += count - 1
	}
	q := &QLearning{
		transitions: make([]int, states*actions),
		qValues:     make([]float64, states*actions),
		bestActions: make([]int, states),
		actions:     actions,
		space:       *space,
		discount:    DefaultDiscount,
		generator:   generator,
	}
	// Row layout: row[0] is stay, then one block per parameter, in the order
	// they were added, listing its other values by increasing index.
	for state := range states {
		row := q.transitions[state*actions : (state+1)*actions]
		row[0] = state
		// action is the next free column of the row.
		action := 1
		for d, count := range space.counts {
			stride := space.strides[d]
			// current is the value index of parameter d in state, as in
			// Param.At.
			current := state / stride % count
			for value := range count {
				if value == current {
					continue
				}
				// Replace digit d: remove current*stride, add value*stride.
				row[action] = state + (value-current)*stride
				action++
			}
		}
	}
	// Keep this draw order (Q-values, then start state) so a seed reproduces.
	for i := range q.qValues {
		q.qValues[i] = generator.RangeFloat64(initialQMin, initialQMax)
	}
	for state := range states {
		q.bestActions[state] = q.scan(state)
	}
	q.state = generator.IntN(states)
	// Options are applied in order, after the random draws, so they do not
	// change the random sequence.
	for _, option := range options {
		option(q)
	}
	return q
}

// Choose implements [Tuner] with an epsilon-greedy action and returns the state it reaches.
func (q *QLearning) Choose(progress float64) int {
	// It takes a random action with probability epsilon(progress), otherwise
	// the best one. It does not move the current state; Learn does.
	if q.generator.Float64() < epsilon(progress) {
		q.action = q.generator.IntN(q.actions)
	} else {
		q.action = q.bestActions[q.state]
	}
	q.next = q.transitions[q.state*q.actions+q.action]
	return q.next
}

// Learn implements [Tuner] by updating the estimate of the last choice.
func (q *QLearning) Learn(reward float64, progress float64) {
	// With s the current state, a the chosen action and s' the state it reached,
	//
	//	Q(s, a) += alpha * (reward + gamma*max Q(s', .) - Q(s, a))
	//
	// where alpha = learningRate(progress) and gamma is the discount; then s'
	// becomes the current state.
	// max Q(s', .), read through the cached best action of s'. When s' is s
	// (stay), this is the value before the update below.
	future := q.qValues[q.next*q.actions+q.bestActions[q.next]]
	i := q.state*q.actions + q.action
	q.qValues[i] += learningRate(progress) * (reward + q.discount*future - q.qValues[i])
	// Only row s changed. Rescan all of it instead of comparing the new value
	// with the old best: if Q(s, a) was the maximum and decreased, another
	// action may now be best.
	q.bestActions[q.state] = q.scan(q.state)
	q.state = q.next
}

// scan returns the action with the largest Q-value in row state. The strict >
// keeps the lowest index on ties, so stay (action 0) wins them.
func (q *QLearning) scan(state int) int {
	row := q.qValues[state*q.actions : (state+1)*q.actions]
	best := 0
	for action, value := range row {
		if value > row[best] {
			best = action
		}
	}
	return best
}

// Policy implements [Reporter], returning one row per state with its best move.
func (q *QLearning) Policy() logger.Table {
	// Each row holds the state, its value of each parameter, its best action as
	// a move, the state that move reaches and the action's Q-value.
	params := len(q.space.counts)
	header := make([]string, 0, params+4)
	header = append(header, "state")
	for d := range params {
		header = append(header, "p"+strconv.Itoa(d))
	}
	header = append(header, "move", "next", "q")

	rows := make([][]string, len(q.bestActions))
	for state := range rows {
		i := state*q.actions + q.bestActions[state]
		next := q.transitions[i]
		row := make([]string, 0, len(header))
		row = append(row, strconv.Itoa(state))
		for d := range params {
			row = append(row, q.space.formats[d](q.space.digit(state, d)))
		}
		// 'g' with -1 is the shortest exact text: lossless in the CSV and
		// independent of the reward scale.
		row = append(row, q.move(state, next), strconv.Itoa(next), strconv.FormatFloat(q.qValues[i], 'g', -1, 64))
		rows[state] = row
	}
	return logger.Table{Header: header, Rows: rows}
}

// move describes going from state to next: "stay", or "p<d> -> <value>" for
// the one parameter that differs.
func (q *QLearning) move(state, next int) string {
	for d := range q.space.counts {
		if to := q.space.digit(next, d); to != q.space.digit(state, d) {
			return "p" + strconv.Itoa(d) + " -> " + q.space.formats[d](to)
		}
	}
	return "stay"
}

// epsilon returns the exploration rate at progress: cosine annealing from a
// peak down to epsilonMin within each slice of restartPeriod, with the peak
// lowered by epsilonPeakStep at every restart. It depends on progress alone.
func epsilon(progress float64) float64 {
	// Index of the current slice. At some exact boundaries rounding gives the
	// previous slice (0.3/0.1 is 2.9999999999999996), so epsilon is about
	// epsilonMin there instead of the new peak.
	cycle := math.Floor(progress / restartPeriod)
	// Peak of this slice, never below the floor.
	peak := max(epsilonMin, epsilonPeak-epsilonPeakStep*cycle)
	// Position inside the slice: 0 at its start, 1 at its end.
	phase := (progress - cycle*restartPeriod) / restartPeriod
	// 1+cos(pi*phase) goes from 2 to 0, so epsilon goes from peak to
	// epsilonMin.
	return epsilonMin + 0.5*(peak-epsilonMin)*(1+math.Cos(math.Pi*phase))
}

// learningRate returns the learning rate alpha at progress: linear from 1 at
// progress 0 to 0.1 at progress 1, with no clamping.
func learningRate(progress float64) float64 {
	return learningRateStart - learningRateDecay*progress
}
