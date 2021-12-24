package resolver

import "fmt"

type NodeSet struct {
	mapset map[*Node]bool
}

func (s *NodeSet) Size() int {
	return len(s.mapset)
}

func (s *NodeSet) Add(item *Node) {
	s.mapset[item] = true
}

func (s *NodeSet) Del(item *Node) {
	delete(s.mapset, item)
}

func (s *NodeSet) Iter() <-chan *Node {
	ch := make(chan *Node)
	go func() {
		for item := range s.mapset {
			ch <- item
		}
		close(ch)
	}()

	return ch
}

func (s *NodeSet) Clone() *NodeSet {
	set := NewSet()

	for item := range s.Iter() {
		set.Add(item)
	}

	return set
}

func (sl *NodeSet) Difference(sr *NodeSet) *NodeSet {
	set := sl.Clone()

	for item := range sr.Iter() {
		set.Del(item)
	}

	return set
}

func (sl *NodeSet) Union(sr *NodeSet) *NodeSet {
	set := sl.Clone()

	for item := range sr.Iter() {
		set.Add(item)
	}

	return set
}

func (s *NodeSet) ToSlice() []*Node {
	slice := make([]*Node, s.Size())

	i := 0
	for item := range s.mapset {
		slice[i] = item
		i++
	}

	return slice
}

func (s *NodeSet) String() string {
	return fmt.Sprintf("%v", s.ToSlice())
}

func NewSet() *NodeSet {
	return &NodeSet{mapset: make(map[*Node]bool)}
}
