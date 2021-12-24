package dependency

// Dependency represents a relationship between a dependant and a prerequisite
// the Dependant depends on the prerequisite
type Dependency struct {
	Dependant    string
	Prerequisite string
}

// NewDependency creates a Dependency with a dependant and a prerequisite
func NewDependency(dep string, pre string) *Dependency {
	return &Dependency{Dependant: dep, Prerequisite: pre}
}
