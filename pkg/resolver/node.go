package resolver

import (
	"fmt"
)

// Node of a dependency (directed) graph
type Node struct {
	// Name of this node
	Name string
	// Prerequisites is a set nodes this node depends on
	Prerequisites *NodeSet
	// Level means max depth of dependency chain to reach this node
	Level uint

	visited bool
}

// String function used to pretty print this node
func (n Node) String() string {
	return fmt.Sprintf("%d:%s", n.Level, n.Name)
}

// NewNode creates a node with a name, and initialize it
func NewNode(name string) *Node {
	return &Node{Name: name, Prerequisites: NewSet()}
}
