package resolver

type Graph struct {
	nodes             map[string]*Node
	lookupMap         map[*Node]*NodeSet
	dependantNodes    *NodeSet
	prerequisiteNodes *NodeSet
	allNodes          *NodeSet
}

func (r *Graph) getOrCreateNode(name string) *Node {
	np, exists := r.nodes[name]

	if !exists {
		np = NewNode(name)
		r.nodes[name] = np
		r.allNodes.Add(np)
	}

	return np
}

func (r *Graph) getOrCreateNodes(dependency *Dependency) (*Node, *Node) {
	dependant := r.getOrCreateNode(dependency.Dependant)
	prerequisite := r.getOrCreateNode(dependency.Prerequisite)

	return dependant, prerequisite
}

func (r *Graph) addToLookupMap(dependant *Node, prerequisite *Node) {
	sp, exists := r.lookupMap[prerequisite]
	if !exists {
		sp = NewSet()
		r.lookupMap[prerequisite] = sp
	}
	sp.Add(dependant)
}

func (r *Graph) addDependency(dependency *Dependency) {
	dependant, prerequisite := r.getOrCreateNodes(dependency)
	dependant.Prerequisites.Add(prerequisite)

	r.dependantNodes.Add(dependant)
	r.prerequisiteNodes.Add(prerequisite)

	r.addToLookupMap(dependant, prerequisite)
}

func (r *Graph) leaves() *NodeSet {
	return r.prerequisiteNodes.Difference(r.dependantNodes)
}

func (r *Graph) resetVisited() {
	for leaf := range r.allNodes.Iter() {
		leaf.visited = false
	}
}
