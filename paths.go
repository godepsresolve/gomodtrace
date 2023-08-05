package gomodtrace

import (
	"strings"
)

// Path represents a path between two nodes of the graph.
type Path []Node

// Paths represents collection of several paths between two nodes of the graph.
type Paths []Path

// String implements Stringer interface to print object in human-readable form.
func (p Path) String() string {
	builder := strings.Builder{}
	builder.WriteString("[")
	for i := range p {
		builder.WriteString(p[i].name)
		if i == len(p)-1 {
			continue
		}
		builder.WriteString(", ")
	}
	builder.WriteString("]")
	return builder.String()
}

// ListInvolvedLibraries returns list of all libraries occurring in paths
// between two nodes of a graph.
func (ps Paths) ListInvolvedLibraries() []Library {
	var result []Library
	for _, path := range ps {
		for _, node := range path {
			result = append(result, node.name)
		}
	}
	return unique(result)
}
