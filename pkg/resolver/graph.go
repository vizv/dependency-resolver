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

// resolve the dependency
func (r *Graph) resolve(n *Node, sequence uint) uint {
	sequence += 1
	maxSequence := sequence

	if n.visited {
		return 0
	}
	n.visited = true

	if sequence <= n.Sequence {
		n.visited = false
		return sequence
	}
	n.Sequence = sequence
	if parents, okay := r.lookupMap[n]; okay {
		for leaf := range parents.Iter() {
			leafSequence := r.resolve(leaf, sequence)
			if leafSequence == 0 {
				return 0
			}
			if leafSequence > maxSequence {
				maxSequence = leafSequence
			}
		}
	}

	n.visited = false
	return maxSequence
}
