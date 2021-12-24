package resolver

type Resolver struct {
	graph *DependencyGraph
}

func NewResolver(dependencySource DependencySource) *Resolver {
	graph := &DependencyGraph{
		nodes:             make(map[string]*DependencyNode),
		lookupMap:         make(map[*DependencyNode]*DependencyNodeSet),
		dependantNodes:    NewSet(),
		allNodes:          NewSet(),
		prerequisiteNodes: NewSet(),
	}

	for dependency := range dependencySource {
		graph.addDependency(&dependency)
	}

	return &Resolver{graph}
}

// Resolve the dependency graph and update sequence for all nodes
func (r *Resolver) Resolve() ([]*DependencyNode, error) {
	if !r.graph.resolveAll(r.graph.leaves(), 1) {
		return nil, NewCircularDependencyError()
	}

	return r.graph.allNodes.ToSlice(), nil
}
