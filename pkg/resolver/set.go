package resolver

import "fmt"

type Set interface {
	Size() int
	Add(interface{})
	Del(interface{})
	Iter() <-chan interface{}
	Clone() Set
	Difference(s Set) Set
	Union(s Set) Set
	ToSlice() []interface{}
	String() string
}

type set struct {
	mapset map[interface{}]bool
}

func (s set) Size() int {
	return len(s.mapset)
}

func (s set) Add(item interface{}) {
	s.mapset[item] = true
}

func (s set) Del(item interface{}) {
	delete(s.mapset, item)
}

func (s set) Iter() <-chan interface{} {
	ch := make(chan interface{})
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

func (s set) ToSlice() []interface{} {
	slice := make([]interface{}, s.Size())

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
	set := set{mapset: make(map[interface{}]bool)}

	return set
}
