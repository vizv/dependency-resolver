package dependencyresolver

import (
	"fmt"

	mapset "github.com/vizv/pkg/set"
)

// node of a dependency (directed) graph
type node struct {
	// Value of this node
	Value interface{}
	// Prerequisites is a set nodes this node depends on
	Prerequisites mapset.Set
	// Level means max depth of dependency chain to reach this node
	Level uint

	visited bool
}

// String function used to pretty print this node
func (n node) String() string {
	return fmt.Sprintf("{v=%v l=%d}", n.Value, n.Level)
}
