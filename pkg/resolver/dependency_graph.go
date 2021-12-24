package resolver

// DependencyGraph represent a directed graph for storing dependency information
type DependencyGraph struct {
	nodes             map[string]*DependencyNode
	lookupMap         map[*DependencyNode]*DependencyNodeSet
	dependantNodes    *DependencyNodeSet
	prerequisiteNodes *DependencyNodeSet
	allNodes          *DependencyNodeSet
}

// getOrCreateNode returns a node by name,
// the node gets created if not exist in graph
func (g *DependencyGraph) getOrCreateNode(name string) *DependencyNode {
	n, exists := g.nodes[name]

	if !exists {
		n = NewNode(name)
		g.nodes[name] = n
		g.allNodes.Add(n)
	}

	return n
}

// getOrCreateNodes returns dependant node and prerequisite node for a given dependency,
// these nodes get created if not exist in graph
func (g *DependencyGraph) getOrCreateNodes(d *Dependency) (*DependencyNode, *DependencyNode) {
	dependant := g.getOrCreateNode(d.Dependant)
	prerequisite := g.getOrCreateNode(d.Prerequisite)

	return dependant, prerequisite
}

// addToLookupMap adds a dependant with a prerequisite to a map for reverse lookup
func (g *DependencyGraph) addToLookupMap(dep *DependencyNode, pre *DependencyNode) {
	s, exists := g.lookupMap[pre]
	if !exists {
		s = NewSet()
		g.lookupMap[pre] = s
	}
	s.Add(dep)
}

// addDependency add a dependency to a graph,
// and update information stored in the graph
func (g *DependencyGraph) addDependency(d *Dependency) {
	dep, pre := g.getOrCreateNodes(d)
	dep.Prerequisites.Add(pre)

	g.dependantNodes.Add(dep)
	g.prerequisiteNodes.Add(pre)

	g.addToLookupMap(dep, pre)
}

// leaves calculates nodes without prerequisite
func (g *DependencyGraph) leaves() *DependencyNodeSet {
	return g.prerequisiteNodes.Difference(g.dependantNodes)
}

func (g *DependencyGraph) resolveAll(s *DependencyNodeSet, seq uint) error {
	for leaf := range s.Iter() {
		if err := g.resolve(leaf, seq); err != nil {
			return err
		}
	}

	return nil
}

// resolve the dependency graph recursively,
// update this node and all dependants only if current sequence is larger
func (g *DependencyGraph) resolve(n *DependencyNode, seq uint) error {
	if n.visited {
		return NewCircularDependencyError()
	}
	n.visited = true
	defer func() { n.visited = false }()

	if seq <= n.Sequence {
		return nil
	}
	n.Sequence = seq
	if pres, exists := g.lookupMap[n]; exists {
		if err := g.resolveAll(pres, seq+1); err != nil {
			return err
		}
	}

	return nil
}
