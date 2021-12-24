package dependency

// Graph represent a directed graph for storing dependency information
type Graph struct {
	nodes     map[string]*Node
	lookupMap map[*Node]*NodeSet

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

// addToLookupMap of a dependant with a prerequisite for reverse lookup
func (g *Graph) addToLookupMap(dep *Node, pre *Node) {
	s, exists := g.lookupMap[pre]
	if !exists {
		s = NewNodeSet()
		g.lookupMap[pre] = s
	}
	s.Add(dep)
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
	if deps, exists := g.lookupMap[n]; exists {
		if err := g.resolveAll(deps, seq+1); err != nil {
			return err
		}
	}

	return nil
}

// AddDependency to a graph
// also update information stored in the graph
func (g *Graph) AddDependency(d *Dependency) {
	dep, pre := g.getOrCreateNodes(d)
	dep.Prerequisites.Add(pre)

	g.leaves.Delete(dep)
	if pre.Prerequisites.Count() == 0 {
		g.leaves.Add(pre)
	}

	g.addToLookupMap(dep, pre)
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
		nodes:     make(map[string]*Node),
		lookupMap: make(map[*Node]*NodeSet),
		leaves:    NewNodeSet(),
	}
}
