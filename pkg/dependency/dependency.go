package dependency

// Dependency represents a relationship between a dependant and a prerequisite
// the Dependant depends on the prerequisite
type Dependency struct {
	Dependant    string
	Prerequisite string
}
