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

func (g *DependencyGraph) resolveAll(s *DependencyNodeSet, seq uint) bool {
	for leaf := range s.Iter() {
		if !g.resolve(leaf, seq) {
			return false
		}
	}

	return true
}

// resolve the dependency graph recursively,
// update this node and all dependants only if current sequence is larger
func (g *DependencyGraph) resolve(n *DependencyNode, seq uint) bool {
	if n.visited {
		return false
	}
	n.visited = true
	defer func() { n.visited = false }()

	if seq <= n.Sequence {
		return true
	}
	n.Sequence = seq
	if pres, okay := g.lookupMap[n]; okay {
		if !g.resolveAll(pres, seq+1) {
			return false
		}
	}

	return true
}
