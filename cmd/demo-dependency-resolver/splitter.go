package main

import (
	"strings"

	"github.com/vizv/dependency-resolver/pkg/resolver"
)

type Splitter func(string) resolver.Dependency

const DEFAULT_SEP = " "

func newStringSplitter(sep string) Splitter {
	return func(dependencyString string) resolver.Dependency {
		tokens := strings.SplitN(dependencyString, sep, 2)

		return resolver.Dependency{Dependant: tokens[0], Prerequisite: tokens[1]}
	}
}

func defaultStringSplitter() Splitter {
	return newStringSplitter(DEFAULT_SEP)
}
