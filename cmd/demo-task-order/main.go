package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	resolver "github.com/vizv/pkg/dependency-resolver"
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

func newReaderSource(reader io.Reader, splitter Splitter) <-chan resolver.Dependency {
	ch := make(chan resolver.Dependency)
	scanner := bufio.NewScanner(reader)
	go func() {
		for scanner.Scan() {
			ch <- splitter(scanner.Text())
		}
		close(ch)
	}()
	return ch
}

func newFileSource(filename string, splitter Splitter) <-chan resolver.Dependency {
	if file, err := os.Open(filename); err == nil {
		return newReaderSource(file, splitter)
	} else {
		log.Fatalf("cannot create file source: %v", err)
	}

	return nil
}

func main() {
	splitter := defaultStringSplitter()
	var dependencySource <-chan resolver.Dependency

	switch len(os.Args) {
	case 1:
		log.Println("read from stdin")
		dependencySource = newReaderSource(os.Stdin, splitter)
	case 2:
		filename := os.Args[1]
		log.Printf("read from '%s'\n", filename)
		dependencySource = newFileSource(filename, splitter)
	default:
		args := os.Args[1:]
		log.Fatalln("invalid arguments:", strings.Join(args, " "))
	}

	if sequence, err := resolver.NewResolver(dependencySource).Resolve(); err == nil {
		fmt.Println(sequence)
	} else {
		log.Fatalln("failed to resolve dependency:", err)
	}
}
