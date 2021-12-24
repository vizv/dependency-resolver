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

func newFileSource(filename string, readerSource ReaderSource) dependency.Source {
	if file, err := os.Open(filename); err == nil {
		return readerSource(file)
	} else {
		log.Fatalf("cannot create file source: %v", err)
	}

	return nil
}

func main() {
	var dependencySource dependency.Source

	defaultSplitter := defaultSplitParser()
	defaultReaderSource := newParseReaderSource(defaultSplitter)

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

	if sequence, err := dependency.NewResolver(dependencySource).Resolve(); err == nil {
		sort.Slice(sequence, func(i, j int) bool {
			l, r := sequence[i], sequence[j]
			if l.Sequence == r.Sequence {
				return l.Name < r.Name
			}
			return l.Sequence < r.Sequence
		})
		for _, n := range sequence {
			fmt.Println(n)
		}
	} else {
		log.Fatalln("failed to resolve dependency:", err)
	}
}
