package gomodtrace

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"testing"
)

func init() {
	log.SetOutput(io.Discard)
}

// Having this graph of dependencies:
// A──────►B─────►C─────►D
// │       │      ▲      ▲
// │       │      │      │
// └───────►E─────┴──────┘
//         └──►X
// or in vertical view:
// ┌───┐
// │ A │
// └┬─┬┘
//  │┌▽────┐
//  ││  B  │
//  │└┬───┬┘
// ┌▽─▽──┐│
// │  E  ││
// └┬─┬─┬┘│
// ┌▽┐│┌▽─▽─┐
// │X│││ C  │
// └─┘│└┬───┘
// ┌──▽─▽┐
// │  D  │
// └─────┘
// consider that:
// A depends on B and on E
// B depends on C and E
// C depends on D
// E depends on C, D and X
// So, as adjacency list of this graph could be written as:
// A B, A E, B C, B E, C D, E C, E D, E X
// If I want to find all paths from A to D or vise versa, it could be
// A->B->C->D, A->E->D, A->E->C->D, A->B->E->D, A->B->E->C->D. Pay attention
// there is no X in any path.

var defaultAdjList = AdjacencyList{
	AdjacencyListItem{"A", "B"},
	AdjacencyListItem{"A", "E"},
	AdjacencyListItem{"B", "C"},
	AdjacencyListItem{"B", "E"},
	AdjacencyListItem{"C", "D"},
	AdjacencyListItem{"E", "C"},
	AdjacencyListItem{"E", "D"},
	AdjacencyListItem{"E", "X"},
}

func Test_ParseGraph(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  AdjacencyList
	}{
		{
			name: "success",
			input: []string{
				"A B",
				"A E",
				"B C",
				"B E",
				"C D",
				"E C",
				"E D",
				"E X",
			},
			want: defaultAdjList,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseGraph(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseGraph() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_BuildGraphIndex(t *testing.T) {
	tests := []struct {
		name          string
		adjacencyList AdjacencyList
		want          map[string]string
	}{
		{
			name:          "success",
			adjacencyList: defaultAdjList,
			// A──────►B─────►C─────►D
			// │       │      ▲      ▲
			// │       │      │      │
			// └───────►E─────┴──────┘
			//         └──►X
			// consider that:
			// A depends on B and on E
			// B depends on C and E
			// C depends on D
			// E depends on C, D and X
			want: map[string]string{
				"A": "{A from: [] to: [B, E]}",
				"B": "{B from: [A] to: [C, E]}",
				"C": "{C from: [B, E] to: [D]}",
				"D": "{D from: [C, E] to: []}",
				"E": "{E from: [A, B] to: [C, D, X]}",
				"X": "{X from: [E] to: []}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildGraphIndex(tt.adjacencyList)
			if !reflect.DeepEqual(
				fmt.Sprintf("%v", got),
				fmt.Sprintf("%v", tt.want),
			) {
				t.Errorf("BuildGraphIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeIndex_FindPaths(t *testing.T) {
	graph := BuildGraphIndex(defaultAdjList)
	// A──────►B─────►C─────►D
	// │       │      ▲      ▲
	// │       │      │      │
	// └───────►E─────┴──────┘
	//         └──►X
	// If I want to find all paths from A to D or vise versa, it could be
	// A->B->C->D, A->E->D, A->E->C->D, A->B->E->D, A->B->E->C->D.
	tests := []struct {
		name string
		from Library
		to   Library
		seen map[string]bool
		want string
	}{
		{
			name: "success find from A to E",
			from: "A",
			to:   "E",
			seen: nil,
			want: "[[A, E] [A, B, E]]",
		},
		{
			name: "success find from A to C",
			from: "A",
			to:   "C",
			seen: nil,
			want: "[[A, B, C] [A, E, C] [A, B, E, C]]",
		},
		{
			name: "success find from A to D",
			from: "A",
			to:   "D",
			seen: nil,
			want: "[[A, B, C, D] [A, E, C, D] [A, B, E, C, D] [A, E, D] [A, B, E, D]]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := graph.FindPaths(tt.from, tt.to, tt.seen)
			if !reflect.DeepEqual(
				fmt.Sprintf("%v", got),
				fmt.Sprintf("%v", tt.want),
			) {
				t.Errorf("FindPaths() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("success find from A to C on a graph with cycles", func(t *testing.T) {
		adjListWithCycles := AdjacencyList{
			AdjacencyListItem{"A", "B"},
			AdjacencyListItem{"A", "E"},
			AdjacencyListItem{"B", "C"},
			AdjacencyListItem{"B", "E"}, // there is a cycle between B and E
			AdjacencyListItem{"C", "D"},
			AdjacencyListItem{"E", "B"}, // there is a cycle between B and E
			AdjacencyListItem{"E", "C"},
			AdjacencyListItem{"E", "D"},
			AdjacencyListItem{"E", "X"},
		}
		cycledGraphIdx := BuildGraphIndex(adjListWithCycles)
		from := "A"
		to := "C"
		want := "[[A, B, C] [A, E, B, C] [A, E, C] [A, B, E, C]]"
		got := cycledGraphIdx.FindPaths(from, to, nil)
		if !reflect.DeepEqual(
			fmt.Sprintf("%v", got),
			fmt.Sprintf("%v", want),
		) {
			t.Errorf("FindPaths() = %v, want %v", got, want)
		}
	})
}

func TestAdjacencyList_WithOnly(t *testing.T) {
	// A──────►B─────►C─────►D
	// │       │      ▲      ▲
	// │       │      │      │
	// └───────►E─────┴──────┘
	//         └──►X
	t.Run("success with all libraries in graph", func(t *testing.T) {
		al := make(AdjacencyList, len(defaultAdjList))
		copy(al, defaultAdjList)
		got := al.WithOnly([]Library{"A", "B", "C", "D", "E", "X"})
		if !reflect.DeepEqual(got, defaultAdjList) {
			t.Errorf("WithOnly() = %v, want %v", got, defaultAdjList)
		}
	})
	t.Run("success with several libraries in graph", func(t *testing.T) {
		got := defaultAdjList.WithOnly([]Library{"A", "B", "C"})
		want := AdjacencyList{{"A", "B"}, {"B", "C"}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("WithOnly() = %v, want %v", got, want)
		}
	})
}
