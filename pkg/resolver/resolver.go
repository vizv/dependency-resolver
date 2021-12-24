package resolver

type Resolver struct {
	graph *Graph
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
	if parents, okay := r.graph.lookupMap[n]; okay {
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
	graph := &Graph{
		nodes:             make(map[string]*Node),
		lookupMap:         make(map[*Node]*NodeSet),
		dependantNodes:    NewSet(),
		allNodes:          NewSet(),
		prerequisiteNodes: NewSet(),
	}

	for dependency := range dependencySource {
		graph.addDependency(&dependency)
	}

	return &Resolver{graph}
}

func (r *Resolver) Resolve() ([][]*Node, error) {
	maxLevel := uint(0)
	for leaf := range r.graph.leaves().Iter() {
		r.graph.resetVisited()
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
	for n := range r.graph.allNodes.Iter() {
		leveledSequence[n.Level-1] = append(leveledSequence[n.Level-1], n)
	}

	return leveledSequence, nil
}
