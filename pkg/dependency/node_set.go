package dependency

// NodeSet represents a set of dependency nodes
type NodeSet struct {
	mapset map[*Node]bool
}

// Count of nodes in this set
func (s *NodeSet) Count() int {
	return len(s.mapset)
}

// Add a node to the set
func (s *NodeSet) Add(n *Node) {
	s.mapset[n] = true
}

// Delete a node from the set
func (s *NodeSet) Delete(n *Node) {
	delete(s.mapset, n)
}

// Iterator used to iterate through the set
func (s *NodeSet) Iterator() <-chan *Node {
	ch := make(chan *Node)
	go func() {
		for n := range s.mapset {
			ch <- n
		}
		close(ch)
	}()

	return ch
}

// Clone the set
func (s *NodeSet) Clone() *NodeSet {
	set := NewNodeSet()

	for n := range s.Iterator() {
		set.Add(n)
	}

	return set
}

// Exclude nodes from another set
// returns a new set
func (sl *NodeSet) Exclude(sr *NodeSet) *NodeSet {
	set := sl.Clone()

	for n := range sr.Iterator() {
		set.Delete(n)
	}

	return set
}

// NewNodeSet creates an empty set, and initialize it
func NewNodeSet() *NodeSet {
	return &NodeSet{mapset: make(map[*Node]bool)}
}
