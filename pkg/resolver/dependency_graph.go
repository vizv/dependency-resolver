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
func (r *DependencyGraph) getOrCreateNode(name string) *DependencyNode {
	np, exists := r.nodes[name]

	if !exists {
		np = NewNode(name)
		r.nodes[name] = np
		r.allNodes.Add(np)
	}

	return np
}

// getOrCreateNodes returns dependant node and prerequisite node for a given dependency,
// these nodes get created if not exist in graph
func (r *DependencyGraph) getOrCreateNodes(dependency *Dependency) (*DependencyNode, *DependencyNode) {
	dependant := r.getOrCreateNode(dependency.Dependant)
	prerequisite := r.getOrCreateNode(dependency.Prerequisite)

	return dependant, prerequisite
}

// addToLookupMap adds a dependant with a prerequisite to a map for reverse lookup
func (r *DependencyGraph) addToLookupMap(dependant *DependencyNode, prerequisite *DependencyNode) {
	sp, exists := r.lookupMap[prerequisite]
	if !exists {
		sp = NewSet()
		r.lookupMap[prerequisite] = sp
	}
	sp.Add(dependant)
}

// addDependency add a dependency to a graph,
// and update information stored in the graph
func (r *DependencyGraph) addDependency(dependency *Dependency) {
	dependant, prerequisite := r.getOrCreateNodes(dependency)
	dependant.Prerequisites.Add(prerequisite)

	r.dependantNodes.Add(dependant)
	r.prerequisiteNodes.Add(prerequisite)

	r.addToLookupMap(dependant, prerequisite)
}

// leaves calculates nodes without prerequisite
func (r *DependencyGraph) leaves() *DependencyNodeSet {
	return r.prerequisiteNodes.Difference(r.dependantNodes)
}

func (r *DependencyGraph) resolveAll(nodeSet *DependencyNodeSet, sequence uint) bool {
	for leaf := range nodeSet.Iter() {
		if !r.resolve(leaf, sequence) {
			return false
		}
	}

	return true
}

// resolve the dependency graph recursively,
// update this node and all dependants only if current sequence is larger
func (r *DependencyGraph) resolve(node *DependencyNode, sequence uint) bool {
	if node.visited {
		return false
	}
	node.visited = true
	defer func() { node.visited = false }()

	if sequence <= node.Sequence {
		return true
	}
	node.Sequence = sequence
	if parents, okay := r.lookupMap[node]; okay {
		for leaf := range parents.Iter() {
			if !r.resolve(leaf, sequence+1) {
				return false
			}
		}
	}

	return true
}
