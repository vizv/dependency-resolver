package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/vizv/dependency-resolver/pkg/resolver"
)

type ReaderSource func(reader io.Reader) resolver.Source

func newSplitReaderSource(splitter Splitter) ReaderSource {
	return func(reader io.Reader) resolver.Source {
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
}

func newGraphvizReaderSource() ReaderSource {
	return func(reader io.Reader) resolver.Source {
		ch := make(chan resolver.Dependency)
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
						ch <- resolver.Dependency{Dependant: dependant, Prerequisite: prerequisite}
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

func newFileSource(filename string, readerSource ReaderSource) resolver.Source {
	if file, err := os.Open(filename); err == nil {
		return readerSource(file)
	} else {
		log.Fatalf("cannot create file source: %v", err)
	}

	return nil
}

func main() {
	var dependencySource resolver.Source

	defaultSplitter := defaultStringSplitter()
	defaultReaderSource := newSplitReaderSource(defaultSplitter)

	switch len(os.Args) {
	case 1:
		log.Println("read from stdin")
		dependencySource = defaultReaderSource(os.Stdin)
	case 2:
		filename := os.Args[1]
		ext := filepath.Ext(filename)
		var readerSource ReaderSource
		switch ext {
		case ".in":
			readerSource = defaultReaderSource
		case ".gv":
			readerSource = newGraphvizReaderSource()
		default:
			log.Fatalf("unsupported file: %s", filename)
		}
		log.Printf("read from '%s'\n", filename)
		dependencySource = newFileSource(filename, readerSource)
	default:
		args := os.Args[1:]
		log.Fatalln("invalid arguments:", strings.Join(args, " "))
	}

	if leveledSequence, err := resolver.NewResolver(dependencySource).Resolve(); err == nil {
		for _, sequence := range leveledSequence {
			names := []string{}
			for _, node := range sequence {
				names = append(names, node.Name)
			}
			sort.Strings(names)
			fmt.Println(strings.Join(names, " "))
		}
	} else {
		log.Fatalln("failed to resolve dependency:", err)
	}
}
