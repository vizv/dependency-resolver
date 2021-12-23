package taskorder

import (
	"fmt"

	mapset "github.com/deckarep/golang-set"
)

type Dependency struct {
	Parent interface{}
	Child  interface{}
}

type DependenciesEnumerator interface {
	NextDependency(dependency *Dependency) bool
}

type node struct {
	identifier interface{}
	visited    bool
	level      uint
}

type resolver struct {
	nodes       map[interface{}]*node
	parentMap   map[*node]*mapset.Set
	parentNodes mapset.Set
	childNodes  mapset.Set
	_leaves     *mapset.Set
	_all        *mapset.Set
}

func (r resolver) getOrCreateNode(identifier interface{}) *node {
	if np, okay := r.nodes[identifier]; okay {
		return np
	} else {
		n := node{}
		n.identifier = identifier
		r.nodes[identifier] = &n
		return &n
	}
}

func (r resolver) addDependency(dependency *Dependency) {
	parent := r.getOrCreateNode(dependency.Parent)
	child := r.getOrCreateNode(dependency.Child)

	r.parentNodes.Add(parent)
	r.childNodes.Add(child)

	if sp, okay := r.parentMap[child]; okay {
		(*sp).Add(parent)
	} else {
		s := mapset.NewSet()
		r.parentMap[child] = &s
		s.Add(parent)
	}
}

func (r resolver) resolve(n *node, level uint) bool {
	level += 1
	fmt.Printf("resolve %v[visited: %v, level: %d] at level %d\n", n.identifier, n.visited, n.level, level)

	if n.visited {
		return false
	}
	n.visited = true

	if level <= n.level {
		n.visited = false
		return true
	}
	n.level = level
	if parents, okay := r.parentMap[n]; okay {
		for leaf := range (*parents).Iter() {
			if !r.resolve(leaf.(*node), level) {
				return false
			}
		}
	}

	n.visited = false
	return true
}

func NewResolver(dependencies DependenciesEnumerator) resolver {
	resolver := resolver{}
	resolver.nodes = make(map[interface{}]*node)
	resolver.parentMap = make(map[*node]*mapset.Set)
	resolver.parentNodes = mapset.NewSet()
	resolver.childNodes = mapset.NewSet()

	var dependency Dependency
	for dependencies.NextDependency(&dependency) {
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

func (r resolver) Resolve() bool {
	for leaf := range (*r.leaves()).Iter() {
		r.resetVisited()
		if !r.resolve(leaf.(*node), 0) {
			return false
		}
	}

	for leaf := range (*r.all()).Iter() {
		fmt.Println(leaf)
	}

	return true
}
