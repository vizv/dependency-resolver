package main

import (
	"strings"

	"github.com/vizv/dependency-resolver/pkg/dependency"
)

// Parser used to parse a string to a dependency
type Parser func(string) *dependency.Dependency

// DEFAULT_SEP is the default separator
const DEFAULT_SEP = " "

// newSplitParser creates a SplitParser that split a string by separator
func newSplitParser(sep string) Parser {
	return func(dependencyString string) *dependency.Dependency {
		tokens := strings.SplitN(dependencyString, sep, 2)

		switch len(tokens) {
		case 1:
			// single node
			return dependency.NewDependency(tokens[0], tokens[0])
		case 2:
			return dependency.NewDependency(tokens[0], tokens[1])
		default:
			return nil
		}
	}
}

// defaultSplitParser returns the default SplitParser with default separator
func defaultSplitParser() Parser {
	return newSplitParser(DEFAULT_SEP)
}
