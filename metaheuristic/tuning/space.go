package tuning

import "fmt"

// Space is the set of configurations a tuner chooses from.
type Space struct {
	// Every combination of the values of the parameters registered with [Add],
	// each numbered by a state in [0, States()). Create it with [NewSpace] and
	// add every parameter before creating a tuner on it.

	// A state is a mixed-radix number with one digit per parameter, parameter
	// 0 least significant: state = sum(index[i] * strides[i]).

	// counts[i] is the number of values of parameter i (its digit's base).
	counts []int
	// strides[i] is the product of counts[0:i], the weight of digit i.
	strides []int
	// states is the product of all counts.
	states int
	// formats[i] returns the text of value index of parameter i, for reports.
	formats []func(index int) string
}

// NewSpace returns an empty Space, with no parameters and a single state.
func NewSpace() *Space {
	// states starts at 1, the empty product, so the first parameter gets
	// stride 1 and each Add multiplies in its count.
	return &Space{states: 1}
}

// States returns the number of configurations, or 1 if no parameter was added.
func (s *Space) States() int {
	return s.states
}

// Param is a typed handle to one parameter of a [Space], returned by [Add].
type Param[T comparable] struct {
	// It is cheap to copy and stays valid as more parameters are added. The zero
	// Param is not usable.

	values []T // candidate values; a state stores an index into it
	stride int // weight of this parameter in the state number
}

// At returns the parameter's value in state, in [0, States()) of its Space.
func (p Param[T]) At(state int) T {
	// state is not checked.
	// Dividing by the stride drops the less significant digits; the remainder
	// by the count drops the more significant ones.
	return p.values[state/p.stride%len(p.values)]
}

// Add registers a parameter with the given values in space and returns its handle.
func Add[T comparable](space *Space, values []T) Param[T] {
	// Value index i in a state stands for values[i]. Add keeps values without
	// copying it, so do not modify it afterwards.
	//
	// Add panics if values is empty or holds duplicates (compared with ==).
	// param is the index the new parameter gets; it is used in panic messages.
	param := len(space.counts)
	if len(values) == 0 {
		panic(fmt.Sprintf("tuning: parameter %d has no values", param))
	}
	for i := range values {
		for j := i + 1; j < len(values); j++ {
			if values[i] == values[j] {
				panic(fmt.Sprintf("tuning: parameter %d has duplicate value %v at positions %d and %d", param, values[i], i, j))
			}
		}
	}
	// The new digit is the most significant: it advances once every
	// space.states states, then the state count grows by its base.
	stride := space.states
	space.counts = append(space.counts, len(values))
	space.strides = append(space.strides, stride)
	space.formats = append(space.formats, func(index int) string { return fmt.Sprint(values[index]) })
	space.states *= len(values)
	return Param[T]{values: values, stride: stride}
}

// digit returns the value index of parameter d in state, as in Param.At.
func (s *Space) digit(state, d int) int {
	return state / s.strides[d] % s.counts[d]
}
