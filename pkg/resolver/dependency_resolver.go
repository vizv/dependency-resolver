package resolver

// DependencyResolver represents a dependency resolver that
type DependencyResolver struct {
	graph *DependencyGraph
}

// Resolve the dependency graph and update sequence for all nodes in it
func (r *DependencyResolver) Resolve() ([]*DependencyNode, error) {
	if err := r.graph.resolveAll(r.graph.leaves(), 1); err != nil {
		return nil, err
	}

	return r.graph.allNodes.ToSlice(), nil
}

// NewDependencyResolver creates a new resolver with information provided from a source
func NewDependencyResolver(dependencySource DependencySource) *DependencyResolver {
	graph := NewDependencyGraph()
	for dependency := range dependencySource {
		graph.addDependency(&dependency)
	}

	return &DependencyResolver{graph}
}
