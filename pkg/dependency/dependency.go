package dependency

// Dependency represents a relationship between a Dependant and a Prerequisite,
// the Dependant depends on the Prerequisite
type Dependency struct {
	Dependant    string
	Prerequisite string
}
