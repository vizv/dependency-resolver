package dependency

// Resolver represents a dependency dependency that
type Resolver struct {
	graph *Graph
}

// Resolve the dependency graph and update sequence for all nodes in it
func (r *Resolver) Resolve() ([]*Node, error) {
	if err := r.graph.Resolve(); err != nil {
		return nil, err
	}

	return r.graph.Nodes(), nil
}

// NewResolver creates a new dependency with information provided from a source
func NewResolver(dependencySource Source) *Resolver {
	graph := NewGraph()
	for dependency := range dependencySource {
		graph.AddDependency(&dependency)
	}

	return &Resolver{graph}
}
