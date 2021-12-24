package resolver

import "fmt"

type Set interface {
	Size() int
	Add(*Node)
	Del(*Node)
	Iter() <-chan *Node
	Clone() Set
	Difference(s Set) Set
	Union(s Set) Set
	ToSlice() []*Node
	String() string
}

type set struct {
	mapset map[*Node]bool
}

func (s set) Size() int {
	return len(s.mapset)
}

func (s set) Add(item *Node) {
	s.mapset[item] = true
}

func (s set) Del(item *Node) {
	delete(s.mapset, item)
}

func (s set) Iter() <-chan *Node {
	ch := make(chan *Node)
	go func() {
		for item := range s.mapset {
			ch <- item
		}
		close(ch)
	}()

	return ch
}

func (s set) Clone() Set {
	set := NewSet()

	for item := range s.Iter() {
		set.Add(item)
	}

	return set
}

func (sl set) Difference(sr Set) Set {
	set := sl.Clone()

	for item := range sr.Iter() {
		set.Del(item)
	}

	return set
}

func (sl set) Union(sr Set) Set {
	set := sl.Clone()

	for item := range sr.Iter() {
		set.Add(item)
	}

	return set
}

func (s set) ToSlice() []*Node {
	slice := make([]*Node, s.Size())

	i := 0
	for item := range s.mapset {
		slice[i] = item
		i++
	}

	return slice
}

func (s set) String() string {
	return fmt.Sprintf("%v", s.ToSlice())
}

func NewSet() Set {
	set := set{mapset: make(map[*Node]bool)}

	return set
}
