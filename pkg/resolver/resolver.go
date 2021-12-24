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

func (r *Resolver) Resolve() ([]*Node, error) {
	maxLevel := uint(0)
	for leaf := range r.graph.leaves().Iter() {
		leafSequence := r.graph.resolve(leaf, 0)
		if leafSequence == 0 {
			return nil, NewCircularDependencyError()
		}
		if leafSequence > maxLevel {
			maxLevel = leafSequence
		}
	}

	sequence := r.graph.allNodes.ToSlice()

	return sequence, nil
}
