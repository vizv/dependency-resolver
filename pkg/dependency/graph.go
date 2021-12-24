package dependency

// Graph represent a directed graph for storing dependency information
type Graph struct {
	nodes  map[string]*Node
	leaves *NodeSet
}

// getOrCreateNode by name
// the node gets created if not exist in graph
func (g *Graph) getOrCreateNode(name string) *Node {
	n, exists := g.nodes[name]

	if !exists {
		n = NewNode(name)
		g.nodes[name] = n
	}

	return n
}

// getOrCreateNodes for a given dependency
// these nodes get created if not exist in graph
func (g *Graph) getOrCreateNodes(d *Dependency) (*Node, *Node) {
	dependant := g.getOrCreateNode(d.Dependant)
	prerequisite := g.getOrCreateNode(d.Prerequisite)

	return dependant, prerequisite
}

// resolve the dependency graph recursively
// update this node and all dependants only if current sequence is larger
// returns error if the dependency graph is not resolvable
func (g *Graph) resolve(n *Node, seq uint) error {
	// resolving a visited node will result in the CircularDependencyError
	if n.visited {
		return NewCircularDependencyError()
	}

	// mark this node as visited and clear it after resolved
	n.visited = true
	defer func() { n.visited = false }()

	// ignore node already been marked with larger sequence number
	if seq <= n.Sequence {
		return nil
	}

	// update sequence number and then resolve all dependants as next sequence
	n.Sequence = seq
	if err := g.resolveAll(n.Dependants, seq+1); err != nil {
		return err
	}

	return nil
}

// resolveAll the nodes in a node set with resolve function
func (g *Graph) resolveAll(s *NodeSet, seq uint) error {
	for leaf := range s.Iterator() {
		if err := g.resolve(leaf, seq); err != nil {
			return err
		}
	}

	return nil
}

// AddDependency to a graph
// also update information stored in the graph
func (g *Graph) AddDependency(d *Dependency) {
	// ignore empty dependency
	if d == nil {
		return
	}

	dep, pre := g.getOrCreateNodes(d)

	// skip updating dependency for single node
	if dep != pre {
		dep.Prerequisites.Add(pre)
		pre.Dependants.Add(dep)

		g.leaves.Delete(dep)
	}

	if pre.Prerequisites.Count() == 0 {
		g.leaves.Add(pre)
	}
}

// Resolve the dependency graph from leaves
func (g *Graph) Resolve() error {
	if err := g.resolveAll(g.leaves, 1); err != nil {
		return err
	}

	return nil
}

// Nodes of this graph
func (g *Graph) Nodes() []*Node {
	slice := make([]*Node, len(g.nodes))

	i := 0
	for _, n := range g.nodes {
		slice[i] = n
		i++
	}

	return slice
}

// NewGraph creates a empty dependency graph, and initialize it
func NewGraph() *Graph {
	return &Graph{
		nodes:  make(map[string]*Node),
		leaves: NewNodeSet(),
	}
}
