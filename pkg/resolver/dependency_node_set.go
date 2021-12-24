package resolver

import "fmt"

// DependencyNodeSet is a set for dependency nodes
type DependencyNodeSet struct {
	mapset map[*DependencyNode]bool
}

// Count of nodes in this set
func (s *DependencyNodeSet) Count() int {
	return len(s.mapset)
}

// Add a node to the set
func (s *DependencyNodeSet) Add(item *DependencyNode) {
	s.mapset[item] = true
}

// Delete a node from the set
func (s *DependencyNodeSet) Delete(item *DependencyNode) {
	delete(s.mapset, item)
}

// Iterator used to iterate through the set
func (s *DependencyNodeSet) Iterator() <-chan *DependencyNode {
	ch := make(chan *DependencyNode)
	go func() {
		for item := range s.mapset {
			ch <- item
		}
		close(ch)
	}()

	return ch
}

// Clone the set
func (s *DependencyNodeSet) Clone() *DependencyNodeSet {
	set := NewDependencyNodeSet()

	for item := range s.Iterator() {
		set.Add(item)
	}

	return set
}

// Exclude nodes from another set
// returns a new set
func (sl *DependencyNodeSet) Exclude(sr *DependencyNodeSet) *DependencyNodeSet {
	set := sl.Clone()

	for item := range sr.Iterator() {
		set.Delete(item)
	}

	return set
}

// Union nodes with another set
// returns a new set
func (sl *DependencyNodeSet) Union(sr *DependencyNodeSet) *DependencyNodeSet {
	set := sl.Clone()

	for item := range sr.Iterator() {
		set.Add(item)
	}

	return set
}

// ToSlice function converts current set to slice
func (s *DependencyNodeSet) ToSlice() []*DependencyNode {
	slice := make([]*DependencyNode, s.Count())

	i := 0
	for item := range s.mapset {
		slice[i] = item
		i++
	}

	return slice
}

// String function used to pretty print this set
func (s *DependencyNodeSet) String() string {
	return fmt.Sprintf("%v", s.ToSlice())
}

// NewDependencyNodeSet creates an empty set, and initialize it
func NewDependencyNodeSet() *DependencyNodeSet {
	return &DependencyNodeSet{mapset: make(map[*DependencyNode]bool)}
}
