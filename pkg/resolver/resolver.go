package resolver

type Resolver struct {
	nodes             map[string]*Node
	lookupMap         map[*Node]*NodeSet
	dependantNodes    *NodeSet
	prerequisiteNodes *NodeSet
	allNodes          *NodeSet
}

func (r *Resolver) getOrCreateNode(name string) *Node {
	np, exists := r.nodes[name]

	if !exists {
		np = NewNode(name)
		r.nodes[name] = np
		r.allNodes.Add(np)
	}

	return np
}

func (r *Resolver) getOrCreateNodes(dependency *Dependency) (*Node, *Node) {
	dependant := r.getOrCreateNode(dependency.Dependant)
	prerequisite := r.getOrCreateNode(dependency.Prerequisite)

	return dependant, prerequisite
}

func (r *Resolver) addToLookupMap(dependant *Node, prerequisite *Node) {
	sp, exists := r.lookupMap[prerequisite]
	if !exists {
		sp = NewSet()
		r.lookupMap[prerequisite] = sp
	}
	sp.Add(dependant)
}

func (r *Resolver) addDependency(dependency *Dependency) {
	dependant, prerequisite := r.getOrCreateNodes(dependency)
	dependant.Prerequisites.Add(prerequisite)

	r.dependantNodes.Add(dependant)
	r.prerequisiteNodes.Add(prerequisite)

	r.addToLookupMap(dependant, prerequisite)
}

func (r *Resolver) resolve(n *Node, level uint) uint {
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
	if parents, okay := r.lookupMap[n]; okay {
		for leaf := range parents.Iter() {
			leafLevel := r.resolve(leaf, level)
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

func NewResolver(dependencySource Source) *Resolver {
	resolver := &Resolver{
		nodes:             make(map[string]*Node),
		lookupMap:         make(map[*Node]*NodeSet),
		dependantNodes:    NewSet(),
		allNodes:          NewSet(),
		prerequisiteNodes: NewSet(),
	}

	for dependency := range dependencySource {
		resolver.addDependency(&dependency)
	}

	return resolver
}

func (r *Resolver) leaves() *NodeSet {
	return r.prerequisiteNodes.Difference(r.dependantNodes)
}

func (r *Resolver) resetVisited() {
	for leaf := range r.allNodes.Iter() {
		leaf.visited = false
	}
}

func (r *Resolver) Resolve() ([][]*Node, error) {
	maxLevel := uint(0)
	for leaf := range r.leaves().Iter() {
		r.resetVisited()
		leafLevel := r.resolve(leaf, 0)
		if leafLevel == 0 {
			return nil, NewCircularDependencyError()
		}
		if leafLevel > maxLevel {
			maxLevel = leafLevel
		}
	}

	leveledSequence := make([][]*Node, maxLevel)
	for i := uint(0); i < maxLevel; i++ {
		leveledSequence[i] = []*Node{}
	}
	for n := range r.allNodes.Iter() {
		leveledSequence[n.Level-1] = append(leveledSequence[n.Level-1], n)
	}

	return leveledSequence, nil
}
