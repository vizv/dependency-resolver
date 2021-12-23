package dependencyresolver

import (
	"sort"

	mapset "github.com/vizv/pkg/set"
)

type Resolver interface {
	Resolve() ([]node, error)
}

type resolver struct {
	nodes       map[interface{}]*node
	parentsMap  map[*node]*mapset.Set
	parentNodes mapset.Set
	childNodes  mapset.Set
	_leaves     *mapset.Set
	_all        *mapset.Set
}

func (r resolver) getOrCreateNode(value interface{}) *node {
	if np, okay := r.nodes[value]; okay {
		return np
	} else {
		n := node{}
		n.Prerequisites = mapset.NewSet()
		n.Value = value
		r.nodes[value] = &n
		return &n
	}
}

func (r resolver) addDependency(dependency *Dependency) {
	// fmt.Println("add dependency", dependency.Parent, dependency.Child)

	parent := r.getOrCreateNode(dependency.Dependant)
	child := r.getOrCreateNode(dependency.Prerequisite)

	parent.Prerequisites.Add(child)

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

func (r resolver) resolve(n *node, level uint) bool {
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
			if !r.resolve(leaf.(*node), level) {
				return false
			}
		}
	}

	n.visited = false
	return true
}

func NewResolver(dependencySource <-chan Dependency) Resolver {
	resolver := resolver{}
	resolver.nodes = make(map[interface{}]*node)
	resolver.parentsMap = make(map[*node]*mapset.Set)
	resolver.parentNodes = mapset.NewSet()
	resolver.childNodes = mapset.NewSet()

	// var dependency Dependency
	for dependency := range dependencySource {
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
		leaf.(*node).visited = false
	}
}

func (r resolver) Resolve() ([]node, error) {
	for leaf := range (*r.leaves()).Iter() {
		r.resetVisited()
		if !r.resolve(leaf.(*node), 0) {
			return nil, NewCircularDependencyError()
		}
	}

	nodeSize := (*r.all()).Size()
	sequence := make([]node, nodeSize)
	for i, n := range (*r.all()).ToSlice() {
		sequence[i] = *n.(*node)
	}

	sort.Slice(sequence, func(i, j int) bool {
		leftNode := (sequence)[i]
		rightNode := (sequence)[j]

		return leftNode.Level < rightNode.Level
	})

	return sequence, nil
}
