package resolver

// DependencyGraph represent a directed graph for storing dependency information
type DependencyGraph struct {
	nodes             map[string]*DependencyNode
	lookupMap         map[*DependencyNode]*DependencyNodeSet
	dependantNodes    *DependencyNodeSet
	prerequisiteNodes *DependencyNodeSet
	allNodes          *DependencyNodeSet
}

// getOrCreateNode by name
// the node gets created if not exist in graph
func (g *DependencyGraph) getOrCreateNode(name string) *DependencyNode {
	n, exists := g.nodes[name]

	if !exists {
		n = NewDependencyNode(name)
		g.nodes[name] = n
		g.allNodes.Add(n)
	}

	return n
}

// getOrCreateNodes for a given dependency
// these nodes get created if not exist in graph
func (g *DependencyGraph) getOrCreateNodes(d *Dependency) (*DependencyNode, *DependencyNode) {
	dependant := g.getOrCreateNode(d.Dependant)
	prerequisite := g.getOrCreateNode(d.Prerequisite)

	return dependant, prerequisite
}

// addToLookupMap of a dependant with a prerequisite for reverse lookup
func (g *DependencyGraph) addToLookupMap(dep *DependencyNode, pre *DependencyNode) {
	s, exists := g.lookupMap[pre]
	if !exists {
		s = NewDependencyNodeSet()
		g.lookupMap[pre] = s
	}
	s.Add(dep)
}

// addDependency to a graph
// also update information stored in the graph
func (g *DependencyGraph) addDependency(d *Dependency) {
	dep, pre := g.getOrCreateNodes(d)
	dep.Prerequisites.Add(pre)

	g.dependantNodes.Add(dep)
	g.prerequisiteNodes.Add(pre)

	g.addToLookupMap(dep, pre)
}

// leaves are nodes without prerequisite
func (g *DependencyGraph) leaves() *DependencyNodeSet {
	return g.prerequisiteNodes.Exclude(g.dependantNodes)
}

// resolveAll the nodes in a node set with resolve function
func (g *DependencyGraph) resolveAll(s *DependencyNodeSet, seq uint) error {
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
func (g *DependencyGraph) resolve(n *DependencyNode, seq uint) error {
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

// NewDependencyGraph creates a empty dependency graph, and initialize it
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		nodes:             make(map[string]*DependencyNode),
		lookupMap:         make(map[*DependencyNode]*DependencyNodeSet),
		dependantNodes:    NewDependencyNodeSet(),
		allNodes:          NewDependencyNodeSet(),
		prerequisiteNodes: NewDependencyNodeSet(),
	}
}
