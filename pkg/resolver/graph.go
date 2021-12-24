package resolver

// Graph represent a directed graph for storing dependency information
type Graph struct {
	nodes             map[string]*Node
	lookupMap         map[*Node]*NodeSet
	dependantNodes    *NodeSet
	prerequisiteNodes *NodeSet
	allNodes          *NodeSet
}

// getOrCreateNode returns a node by name,
// the node gets created if not exist in Graph
func (r *Graph) getOrCreateNode(name string) *Node {
	np, exists := r.nodes[name]

	if !exists {
		np = NewNode(name)
		r.nodes[name] = np
		r.allNodes.Add(np)
	}

	return np
}

// getOrCreateNodes returns dependant node and prerequisite node for a given Dependency,
// these nodes get created if not exist in Graph
func (r *Graph) getOrCreateNodes(dependency *Dependency) (*Node, *Node) {
	dependant := r.getOrCreateNode(dependency.Dependant)
	prerequisite := r.getOrCreateNode(dependency.Prerequisite)

	return dependant, prerequisite
}

// addToLookupMap adds a dependant with a prerequisite to a map for reverse lookup
func (r *Graph) addToLookupMap(dependant *Node, prerequisite *Node) {
	sp, exists := r.lookupMap[prerequisite]
	if !exists {
		sp = NewSet()
		r.lookupMap[prerequisite] = sp
	}
	sp.Add(dependant)
}

// addDependency add a Dependency to a Graph, and update information stored in the Graph
func (r *Graph) addDependency(dependency *Dependency) {
	dependant, prerequisite := r.getOrCreateNodes(dependency)
	dependant.Prerequisites.Add(prerequisite)

	r.dependantNodes.Add(dependant)
	r.prerequisiteNodes.Add(prerequisite)

	r.addToLookupMap(dependant, prerequisite)
}

// leaves calculates nodes without prerequisite
func (r *Graph) leaves() *NodeSet {
	return r.prerequisiteNodes.Difference(r.dependantNodes)
}

func (r *Graph) resolveAll(nodeSet *NodeSet, sequence uint) bool {
	for leaf := range nodeSet.Iter() {
		if !r.resolve(leaf, sequence) {
			return false
		}
	}

	return true
}

// resolve the dependency
func (r *Graph) resolve(node *Node, sequence uint) bool {
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
