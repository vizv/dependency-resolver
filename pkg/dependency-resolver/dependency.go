package dependencyresolver

// Dependency represents a relationship between a Dependant and a Prerequisite,
// the Dependant depends on the Prerequisite
type Dependency struct {
	Dependant    interface{}
	Prerequisite interface{}
}
