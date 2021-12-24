package dependency

// Resolver represents a dependency dependency that
type Resolver struct {
	graph *Graph
}

// Resolve the dependency graph and update sequence for all nodes in it
func (r *Resolver) Resolve() ([]*Node, error) {
	if err := r.graph.resolveAll(r.graph.leaves(), 1); err != nil {
		return nil, err
	}

	return r.graph.allNodes.ToSlice(), nil
}

// NewResolver creates a new dependency with information provided from a source
func NewResolver(dependencySource Source) *Resolver {
	graph := NewGraph()
	for dependency := range dependencySource {
		graph.addDependency(&dependency)
	}

	return &Resolver{graph}
}
