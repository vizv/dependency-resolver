package resolver

type Resolver interface {
	Resolve() ([][]Node, error)
}

type resolver struct {
	nodes       map[interface{}]*Node
	parentsMap  map[*Node]*Set
	parentNodes Set
	childNodes  Set
	allNodes    Set
}

func (r resolver) getOrCreateNode(value interface{}) *Node {
	if np, okay := r.nodes[value]; okay {
		return np
	} else {
		n := Node{}
		n.Prerequisites = NewSet()
		n.Value = value
		r.nodes[value] = &n
		return &n
	}
}

func (r resolver) addDependency(dependency *Dependency) {
	parent := r.getOrCreateNode(dependency.Dependant)
	child := r.getOrCreateNode(dependency.Prerequisite)

	parent.Prerequisites.Add(child)

	r.parentNodes.Add(parent)
	r.childNodes.Add(child)
	r.allNodes.Add(parent)
	r.allNodes.Add(child)

	if sp, okay := r.parentsMap[child]; okay {
		(*sp).Add(parent)
	} else {
		s := NewSet()
		r.parentsMap[child] = &s
		s.Add(parent)
	}
}

func (r resolver) resolve(n *Node, level uint) uint {
	level += 1
	maxLevel := level

	if n.visited {
		return 0
	}
	n.visited = true

	if level <= n.Level {
		n.visited = false
		return level
	}
	n.Level = level
	if parents, okay := r.parentsMap[n]; okay {
		for leaf := range (*parents).Iter() {
			leafLevel := r.resolve(leaf.(*Node), level)
			if leafLevel == 0 {
				return 0
			}
			if leafLevel > maxLevel {
				maxLevel = leafLevel
			}
		}
	}

	n.visited = false
	return maxLevel
}

func NewResolver(dependencySource Source) Resolver {
	resolver := resolver{}
	resolver.nodes = make(map[interface{}]*Node)
	resolver.parentsMap = make(map[*Node]*Set)
	resolver.parentNodes = NewSet()
	resolver.allNodes = NewSet()
	resolver.childNodes = NewSet()

	for dependency := range dependencySource {
		resolver.addDependency(&dependency)
	}

	return resolver
}

func (r resolver) leaves() Set {
	return r.childNodes.Difference(r.parentNodes)
}

func (r resolver) resetVisited() {
	for leaf := range r.allNodes.Iter() {
		leaf.(*Node).visited = false
	}
}

func (r resolver) Resolve() ([][]Node, error) {
	maxLevel := uint(0)
	for leaf := range r.leaves().Iter() {
		r.resetVisited()
		leafLevel := r.resolve(leaf.(*Node), 0)
		if leafLevel == 0 {
			return nil, NewCircularDependencyError()
		}
		if leafLevel > maxLevel {
			maxLevel = leafLevel
		}
	}

	leveledSequence := make([][]Node, maxLevel)
	for i := uint(0); i < maxLevel; i++ {
		leveledSequence[i] = []Node{}
	}
	for np := range r.allNodes.Iter() {
		n := *np.(*Node)
		leveledSequence[n.Level-1] = append(leveledSequence[n.Level-1], n)
	}

	return leveledSequence, nil
}
