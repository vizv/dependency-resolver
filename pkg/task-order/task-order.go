package taskorder

import (
	"fmt"
	"sort"

	mapset "github.com/deckarep/golang-set"
)

type Dependency struct {
	Parent interface{}
	Child  interface{}
}

type Node struct {
	Value    interface{}
	Children mapset.Set
	Level    uint

	visited bool
}

func (n Node) String() string {
	return fmt.Sprintf("Node{value=%v level=%d children=%d}", n.Value, n.Level, n.Children.Cardinality())
}

type resolver struct {
	nodes       map[interface{}]*Node
	parentsMap  map[*Node]*mapset.Set
	parentNodes mapset.Set
	childNodes  mapset.Set
	_leaves     *mapset.Set
	_all        *mapset.Set
}

func (r resolver) getOrCreateNode(value interface{}) *Node {
	if np, okay := r.nodes[value]; okay {
		return np
	} else {
		n := Node{}
		n.Children = mapset.NewSet()
		n.Value = value
		r.nodes[value] = &n
		return &n
	}
}

func (r resolver) addDependency(dependency *Dependency) {
	parent := r.getOrCreateNode(dependency.Parent)
	child := r.getOrCreateNode(dependency.Child)

	parent.Children.Add(child)

	r.parentNodes.Add(parent)
	r.childNodes.Add(child)

	if sp, okay := r.parentsMap[child]; okay {
		(*sp).Add(parent)
	} else {
		s := mapset.NewSet()
		r.parentsMap[child] = &s
		s.Add(parent)
	}
}

func (r resolver) resolve(n *Node, level uint) bool {
	level += 1
	// fmt.Printf("resolve %v[visited: %v, level: %d] at level %d\n", n.Value, n.visited, n.Level, level)

	if n.visited {
		return false
	}
	n.visited = true

	if level <= n.Level {
		n.visited = false
		return true
	}
	n.Level = level
	if parents, okay := r.parentsMap[n]; okay {
		for leaf := range (*parents).Iter() {
			if !r.resolve(leaf.(*Node), level) {
				return false
			}
		}
	}

	n.visited = false
	return true
}

func NewResolver(dependencies <-chan Dependency) resolver {
	resolver := resolver{}
	resolver.nodes = make(map[interface{}]*Node)
	resolver.parentsMap = make(map[*Node]*mapset.Set)
	resolver.parentNodes = mapset.NewSet()
	resolver.childNodes = mapset.NewSet()

	// var dependency Dependency
	for dependency := range dependencies {
		resolver.addDependency(&dependency)
	}

	return resolver
}

func (r resolver) leaves() *mapset.Set {
	if r._leaves == nil {
		leaves := r.childNodes.Difference(r.parentNodes)
		r._leaves = &leaves
	}
	return r._leaves
}

func (r resolver) all() *mapset.Set {
	if r._all == nil {
		all := r.childNodes.Union(r.parentNodes)
		r._all = &all
	}
	return r._all
}

func (r resolver) resetVisited() {
	for leaf := range (*r.all()).Iter() {
		leaf.(*Node).visited = false
	}
}

func (r resolver) Resolve(sequence *[]Node) bool {
	for leaf := range (*r.leaves()).Iter() {
		r.resetVisited()
		if !r.resolve(leaf.(*Node), 0) {
			return false
		}
	}

	*sequence = []Node{}
	for leaf := range (*r.all()).Iter() {
		leafNode, _ := leaf.(*Node)
		*sequence = append(*sequence, *leafNode)
	}
	sort.Slice(*sequence, func(i, j int) bool {
		leftNode := (*sequence)[i]
		rightNode := (*sequence)[j]

		return leftNode.Level < rightNode.Level
	})

	return true
}
