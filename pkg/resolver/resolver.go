package resolver

type Resolver struct {
	graph *Graph
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
		leafLevel := r.graph.resolve(leaf, 0)
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
