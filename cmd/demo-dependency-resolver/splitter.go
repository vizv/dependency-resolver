package main

import (
	"strings"

	"github.com/vizv/dependency-resolver/pkg/dependency"
)

type Splitter func(string) dependency.Dependency

const DEFAULT_SEP = " "

func newStringSplitter(sep string) Splitter {
	return func(dependencyString string) dependency.Dependency {
		tokens := strings.SplitN(dependencyString, sep, 2)

		return dependency.Dependency{Dependant: tokens[0], Prerequisite: tokens[1]}
	}
}

func defaultStringSplitter() Splitter {
	return newStringSplitter(DEFAULT_SEP)
}
