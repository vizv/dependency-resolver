package main

import (
	"io"
	"log"
	"os"

	"github.com/vizv/dependency-resolver/pkg/dependency"
)

// ReaderSource is a special Source that feeds dependency from an input stream
type ReaderSource func(io.Reader) dependency.Source

// newReaderSource creates a ReaderSource with a Reader
func newReaderSource(readerFn Reader) ReaderSource {
	return func(r io.Reader) dependency.Source {
		ch := make(chan *dependency.Dependency)
		go func() {
			readerFn(r, ch)
			close(ch)
		}()
		return ch
	}
}

// newParserReaderSource creates a ParseReaderSource use a given parser
func newParserReaderSource(parserFn Parser) ReaderSource {
	return newReaderSource(newParserReader(parserFn))
}

// newGraphvizReaderSource creates a GraphvizReaderSource
func newGraphvizReaderSource() ReaderSource {
	return newReaderSource(newGraphvizReader())
}

// newFileSource reads file as input and feeds dependency using ReaderSource
func newFileSource(filename string, sourceFn ReaderSource) dependency.Source {
	if file, err := os.Open(filename); err == nil {
		return sourceFn(file)
	} else {
		log.Fatalf("cannot create file source: %v", err)
	}

	return nil
}
