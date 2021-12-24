package dependency

// Resolver represents a dependency resolver
type Resolver struct {
	graph *Graph
}

// Resolve the dependency graph and returns all the nodes
func (r *Resolver) Resolve() ([]*Node, error) {
	if err := r.graph.Resolve(); err != nil {
		return nil, err
	}

	return r.graph.Nodes(), nil
}

// NewResolver creates a new resolver with information provided from a source
func NewResolver(src Source) *Resolver {
	graph := NewGraph()
	for dependency := range src {
		graph.AddDependency(dependency)
	}

	return &Resolver{graph}
}
