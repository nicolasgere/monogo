package resolver

import (
	"fmt"
	"testing"

	"github.com/dominikbraun/graph"
)

func TestG(t *testing.T) {
	g := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic())
	_ = g.AddVertex("A")
	_ = g.AddVertex("B")
	_ = g.AddVertex("C")
	_ = g.AddVertex("D")
	_ = g.AddEdge("A", "B")
	_ = g.AddEdge("B", "C")
	err := g.AddEdge("C", "A")
	fmt.Println(err)
	_ = graph.DFS(g, "A", func(value string) bool {
		fmt.Println(value)
		return false
	})
}
