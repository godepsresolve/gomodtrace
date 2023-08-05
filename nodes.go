package gomodtrace

import (
	"fmt"
	"strings"
)

// Node represents a node of a graph.
type Node struct {
	name Library
	from []*Node
	to   []*Node
}

// String implements Stringer interface to print object in human-readable form.
func (n Node) String() string {
	builder := strings.Builder{}
	for i := range n.from {
		builder.WriteString(n.from[i].name)
		if i == len(n.from)-1 {
			continue
		}
		builder.WriteString(", ")
	}
	from := builder.String()
	builder.Reset()
	for i := range n.to {
		builder.WriteString(n.to[i].name)
		if i == len(n.to)-1 {
			continue
		}
		builder.WriteString(", ")
	}
	to := builder.String()
	return fmt.Sprintf("{%s from: [%s] to: [%s]}", n.name, from, to)
}
