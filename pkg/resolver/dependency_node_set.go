package resolver

import "fmt"

type DependencyNodeSet struct {
	mapset map[*DependencyNode]bool
}

func (s *DependencyNodeSet) Size() int {
	return len(s.mapset)
}

func (s *DependencyNodeSet) Add(item *DependencyNode) {
	s.mapset[item] = true
}

func (s *DependencyNodeSet) Del(item *DependencyNode) {
	delete(s.mapset, item)
}

func (s *DependencyNodeSet) Iter() <-chan *DependencyNode {
	ch := make(chan *DependencyNode)
	go func() {
		for item := range s.mapset {
			ch <- item
		}
		close(ch)
	}()

	return ch
}

func (s *DependencyNodeSet) Clone() *DependencyNodeSet {
	set := NewSet()

	for item := range s.Iter() {
		set.Add(item)
	}

	return set
}

func (sl *DependencyNodeSet) Difference(sr *DependencyNodeSet) *DependencyNodeSet {
	set := sl.Clone()

	for item := range sr.Iter() {
		set.Del(item)
	}

	return set
}

func (sl *DependencyNodeSet) Union(sr *DependencyNodeSet) *DependencyNodeSet {
	set := sl.Clone()

	for item := range sr.Iter() {
		set.Add(item)
	}

	return set
}

func (s *DependencyNodeSet) ToSlice() []*DependencyNode {
	slice := make([]*DependencyNode, s.Size())

	i := 0
	for item := range s.mapset {
		slice[i] = item
		i++
	}

	return slice
}

func (s *DependencyNodeSet) String() string {
	return fmt.Sprintf("%v", s.ToSlice())
}

func NewSet() *DependencyNodeSet {
	return &DependencyNodeSet{mapset: make(map[*DependencyNode]bool)}
}
