package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/vizv/dependency-resolver/pkg/dependency"
)

// getDependencySource from commandline arguments
func getDependencySource() dependency.Source {
	defaultParser := defaultSplitParser()
	defaultReaderSource := newParserReaderSource(defaultParser)

	switch len(os.Args) {
	case 1: // no argument
		log.Println("read from stdin")
		return defaultReaderSource(os.Stdin)
	case 2: // 1 file argument
		filename := os.Args[1]
		ext := filepath.Ext(filename)
		var readerSource ReaderSource
		switch ext {
		case ".txt":
			readerSource = defaultReaderSource
		case ".gv":
			readerSource = newGraphvizReaderSource()
		default:
			log.Fatalf("unsupported file: %s", filename)
		}
		log.Printf("read from '%s'\n", filename)
		return newFileSource(filename, readerSource)
	default:
		args := os.Args[1:]
		log.Fatalln("invalid arguments:", strings.Join(args, " "))
	}

	return nil
}

// sortNodes by sequence number
// for the same sequence number, sort by name
// returns a new slice
func sortNodes(nodes []*dependency.Node) []*dependency.Node {
	sorted := make([]*dependency.Node, len(nodes))
	copy(sorted, nodes)

	sort.Slice(sorted, func(i, j int) bool {
		l, r := sorted[i], sorted[j]
		if l.Sequence == r.Sequence {
			return l.Name < r.Name
		}
		return l.Sequence < r.Sequence
	})

	return sorted
}

func main() {
	dependencySource := getDependencySource()
	dependencyResolver := dependency.NewResolver(dependencySource)
	if nodes, err := dependencyResolver.Resolve(); err != nil {
		log.Fatalln("failed to resolve dependency:", err)
	} else {
		for _, n := range sortNodes(nodes) {
			fmt.Println(n)
		}
	}
}
