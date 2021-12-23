package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	taskorder "github.com/vizv/pkg/task-order"
)

func NewReaderSource(reader io.Reader) <-chan taskorder.Dependency {
	ch := make(chan taskorder.Dependency)
	scanner := bufio.NewScanner(reader)
	go func() {
		for scanner.Scan() {
			tokens := strings.SplitN(scanner.Text(), " ", 2)
			dependency := taskorder.Dependency{Parent: tokens[0], Child: tokens[1]}
			ch <- dependency
		}
		close(ch)
	}()
	return ch
}

func NewFileSource(filename string) <-chan taskorder.Dependency {
	if file, err := os.Open(filename); err == nil {
		return NewReaderSource(file)
	}

	return nil
}

func main() {
	fileDependencySource := NewFileSource("test/test-01.in")
	var sequence []taskorder.Node
	if taskorder.NewResolver(fileDependencySource).Resolve(&sequence) {
		fmt.Println("Done")
	} else {
		fmt.Println("Fail")
	}
	fmt.Println(sequence)
}
