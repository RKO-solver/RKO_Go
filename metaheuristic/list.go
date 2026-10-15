package metaheuristic

import "go.yaml.in/yaml/v3"

// List holds a parameter's values; YAML accepts a single value or a list.
type List[T any] []T

// UnmarshalYAML reads a single value or a YAML list.
func (l *List[T]) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind == yaml.SequenceNode {
		return n.Decode((*[]T)(l))
	}
	*l = List[T]{*new(T)}
	return n.Decode(&(*l)[0])
}
