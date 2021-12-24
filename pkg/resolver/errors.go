package resolver

import "errors"

// NewCircularDependencyError creates a circular dependency error
func NewCircularDependencyError() error {
	return errors.New("circular dependency detected")
}
