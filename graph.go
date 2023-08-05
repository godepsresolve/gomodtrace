package gomodtrace

import (
	"log"
	"strings"
)

var logger = log.Default()

// Library is just an alias for string representation of dependency lib.
type Library = string

// AdjacencyListItem is an element of AdjacencyList.
type AdjacencyListItem = [2]Library

// AdjacencyList is a graph representation form.
type AdjacencyList []AdjacencyListItem

// NodeIndex is an index for quick access to nodes of a graph.
// It also more human understandable form of graph representation form.
type NodeIndex map[string]*Node

// ParseGraph parses graph from basic string input.
// A graph is represented as adjacency list.
func ParseGraph(input []string) AdjacencyList {
	adjacencyList := AdjacencyList{}
	for i := range input {
		left, right, found := strings.Cut(input[i], " ")
		if !found {
			panic("invalid format of dependency graph," +
				"probably no adjacency list form was provided")
		}
		item := AdjacencyListItem{left, right}
		adjacencyList = append(adjacencyList, item)
	}
	return adjacencyList
}

// String implements Stringer interface to print object in human-readable form.
func (al AdjacencyList) String() string {
	builder := strings.Builder{}
	for i := range al {
		builder.WriteString(al[i][0])
		builder.WriteString(" ")
		builder.WriteString(al[i][1])
		builder.WriteString("\n")
	}
	return builder.String()
}

// WithOnly makes a copy of graph,
// but includes ONLY provided libraries and links between them.
func (al AdjacencyList) WithOnly(libraries []Library) AdjacencyList {
	var result AdjacencyList
	libIndex := toIndex(libraries)
	for _, item := range al {
		if !libIndex[item[0]] || !libIndex[item[1]] {
			continue
		}
		result = append(result, item)
	}
	return result
}

// BuildGraphIndex makes graph index to have fast and direct access to all nodes of a graph.
func BuildGraphIndex(adjacencyList AdjacencyList) NodeIndex {
	nodeIndex := make(NodeIndex)
	for _, item := range adjacencyList {
		from, to := item[0], item[1]
		nodeFrom, ok := nodeIndex[from]
		if !ok {
			nodeFrom = &Node{name: from}
			nodeIndex[from] = nodeFrom
		}
		nodeTo, ok := nodeIndex[to]
		if !ok {
			nodeTo = &Node{name: to}
			nodeIndex[to] = nodeTo
		}
		nodeFrom.to = append(nodeFrom.to, nodeTo)
		nodeTo.from = append(nodeTo.from, nodeFrom)
	}
	return nodeIndex
}

// FindPaths is intended to find all paths from parent library to children one.
func (index NodeIndex) FindPaths(
	parent Library, child Library, seen map[string]bool,
) Paths {
	parentNode, ok := index[parent]
	if !ok {
		panic("cannot find `parent` node, check input")
	}
	childNode, ok := index[child]
	if !ok {
		panic("cannot find `child` node, check input")
	}

	// desired path was found, exit condition for recursion.
	if parentNode.name == childNode.name {
		return Paths{{*parentNode}}
	}

	// results will be stored here.
	var paths Paths
	// initial case, where no recursive calls were done.
	if seen == nil {
		seen = make(map[string]bool)
	}
	// mark node as seen to prevent infinity loops where to library points to each other.
	seen[child] = true

	// let's trace parent target library to start library.
	// trace is reversed, because it is implied that parent has many children and many
	// paths not ending on desired child, but there always is a path from child to parent.
	for _, nFrom := range childNode.from {
		logger.Println(nFrom)
		if _, ok = seen[nFrom.name]; ok {
			// infinity loop was found, skip this path.
			const msg = "already seen, parent: %s, child: %s, seen: %v\n"
			logger.Printf(msg, parent, nFrom.name, seen)
			continue
		}
		subPaths := index.FindPaths(parent, nFrom.name, copyMap(seen))
		logger.Println(subPaths)
		if subPaths == nil {
			continue
		}
		for _, path := range subPaths {
			fullPath := make(Path, len(path))
			copy(fullPath, path)
			fullPath = append(fullPath, *childNode)
			paths = append(paths, fullPath)
		}
	}
	logger.Println("======Result begin======")
	logger.Println(parent, child)
	logger.Println("=======Items:=====")
	for _, item := range paths {
		logger.Println(item)
	}
	logger.Println("======Result End======")

	return paths
}
