package resolver

import (
	"fmt"
)

// DependencyNode of a dependency (directed) graph
type DependencyNode struct {
	// Name of this node
	Name string
	// Prerequisites is a set nodes this node depends on
	Prerequisites *DependencyNodeSet
	// Sequence the minimum possible batch number for processing this node
	Sequence uint

	visited bool
}

// String function used to pretty print this node
func (n DependencyNode) String() string {
	return fmt.Sprintf("%d:%s", n.Sequence, n.Name)
}

// NewNode creates a node with a name, and initialize it
func NewNode(name string) *DependencyNode {
	return &DependencyNode{Name: name, Prerequisites: NewSet()}
}
