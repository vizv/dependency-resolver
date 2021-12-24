package main

import (
	"bufio"
	"bytes"
	"io"
	"log"

	"github.com/goccy/go-graphviz"
	"github.com/vizv/dependency-resolver/pkg/dependency"
)

// Reader feeds dependency from an input stream
type Reader func(io.Reader, chan *dependency.Dependency)

// newParserReader creates a ParserReader
// it parses each dependency with the parser line-by-line
func newParserReader(parserFn Parser) Reader {
	return func(r io.Reader, ch chan *dependency.Dependency) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			ch <- parserFn(scanner.Text())
		}
	}
}

// newGraphvizReader creates a GraphvizReader
// it parses the graphviz dot format graph
func newGraphvizReader() Reader {
	return func(r io.Reader, ch chan *dependency.Dependency) {
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(r); err != nil {
			log.Fatalf("IO error: %v", err)
		}

		if graph, err := graphviz.ParseBytes(buf.Bytes()); err != nil {
			log.Fatalf("graphviz parsing error: %v", err)
		} else {
			node := graph.FirstNode()
			for node != nil {
				dependant := node.Name()
				edge := graph.FirstOut(node)
				for edge != nil {
					prerequisite := edge.Node().Name()
					ch <- dependency.NewDependency(dependant, prerequisite)

					edge = graph.NextOut(edge)
				}

				node = graph.NextNode(node)
			}
		}
	}
}
