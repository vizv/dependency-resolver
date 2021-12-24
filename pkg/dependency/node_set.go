package dependency

import "fmt"

// NodeSet represents a set of dependency nodes
type NodeSet struct {
	mapset map[*Node]bool
}

// Count of nodes in this set
func (s *NodeSet) Count() int {
	return len(s.mapset)
}

// Add a node to the set
func (s *NodeSet) Add(item *Node) {
	s.mapset[item] = true
}

// Delete a node from the set
func (s *NodeSet) Delete(item *Node) {
	delete(s.mapset, item)
}

// Iterator used to iterate through the set
func (s *NodeSet) Iterator() <-chan *Node {
	ch := make(chan *Node)
	go func() {
		for item := range s.mapset {
			ch <- item
		}
		close(ch)
	}()

	return ch
}

// Clone the set
func (s *NodeSet) Clone() *NodeSet {
	set := NewNodeSet()

	for item := range s.Iterator() {
		set.Add(item)
	}

	return set
}

// Exclude nodes from another set
// returns a new set
func (sl *NodeSet) Exclude(sr *NodeSet) *NodeSet {
	set := sl.Clone()

	for item := range sr.Iterator() {
		set.Delete(item)
	}

	return set
}

// Union nodes with another set
// returns a new set
func (sl *NodeSet) Union(sr *NodeSet) *NodeSet {
	set := sl.Clone()

	for item := range sr.Iterator() {
		set.Add(item)
	}

	return set
}

// ToSlice function converts current set to slice
func (s *NodeSet) ToSlice() []*Node {
	slice := make([]*Node, s.Count())

	i := 0
	for item := range s.mapset {
		slice[i] = item
		i++
	}

	return slice
}

// String function used to pretty print this set
func (s *NodeSet) String() string {
	return fmt.Sprintf("%v", s.ToSlice())
}

// NewNodeSet creates an empty set, and initialize it
func NewNodeSet() *NodeSet {
	return &NodeSet{mapset: make(map[*Node]bool)}
}
