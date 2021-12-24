package main

import (
	"bufio"
	"bytes"
	"io"
	"log"

	"github.com/goccy/go-graphviz"
	"github.com/vizv/dependency-resolver/pkg/dependency"
)

// ReaderSource is a special Source that feeds dependency from an input stream
type ReaderSource func(reader io.Reader) dependency.Source

// newParseReaderSource creates a ParseReaderSource use a given parser
// it parses each dependency with the parser
func newParseReaderSource(splitter Parser) ReaderSource {
	return func(reader io.Reader) dependency.Source {
		ch := make(chan *dependency.Dependency)
		scanner := bufio.NewScanner(reader)
		go func() {
			for scanner.Scan() {
				ch <- splitter(scanner.Text())
			}
			close(ch)
		}()
		return ch
	}
}

// newGraphvizReaderSource creates a GraphvizReaderSource
// it parses the graphviz dot format graph from an input stream
func newGraphvizReaderSource() ReaderSource {
	return func(reader io.Reader) dependency.Source {
		ch := make(chan *dependency.Dependency)
		go func() {
			buf := new(bytes.Buffer)
			if _, err := buf.ReadFrom(reader); err != nil {
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
			close(ch)
		}()
		return ch
	}
}
