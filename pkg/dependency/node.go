package dependency

import (
	"fmt"
)

// Node of a dependency graph
type Node struct {
	// Name of this node
	Name string
	// Prerequisites is a set nodes this node depends on
	Prerequisites *NodeSet
	// Sequence the minimum possible batch number for processing this node
	Sequence uint

	visited bool
}

// String function used to pretty print this node
func (n Node) String() string {
	return fmt.Sprintf("%d:%s", n.Sequence, n.Name)
}

// NewNode creates a node with a name, and initialize it
func NewNode(name string) *Node {
	return &Node{Name: name, Prerequisites: NewNodeSet()}
}
