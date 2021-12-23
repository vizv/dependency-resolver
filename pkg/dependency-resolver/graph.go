package dependencyresolver

import (
	"fmt"

	mapset "github.com/vizv/pkg/set"
)

// Node of a dependency (directed) graph
type Node struct {
	// Value of this node
	Value interface{}
	// Prerequisites is a set nodes this node depends on
	Prerequisites mapset.Set
	// Level means max depth of dependency chain to reach this node
	Level uint

	visited bool
}

// String function used to pretty print this node
func (n Node) String() string {
	return fmt.Sprintf("{v=%v l=%d}", n.Value, n.Level)
}
