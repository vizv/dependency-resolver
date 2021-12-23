package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	taskorder "github.com/vizv/pkg/task-order"
)

type fileDependenciesEnumerator struct {
	scanner bufio.Scanner
	channel chan *taskorder.Dependency
}

func NewFileDependencyEnumerator(filename string) (*fileDependenciesEnumerator, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	enumerator := &fileDependenciesEnumerator{}
	enumerator.scanner = *bufio.NewScanner(file)
	enumerator.channel = make(chan *taskorder.Dependency)

	go func() {
		for enumerator.scanner.Scan() {
			tokens := strings.SplitN(enumerator.scanner.Text(), " ", 2)
			dependency := &taskorder.Dependency{Parent: tokens[0], Child: tokens[1]}
			enumerator.channel <- dependency
		}
		enumerator.channel <- nil
	}()

	return enumerator, nil
}

func (enumerator fileDependenciesEnumerator) NextDependency(dependency *taskorder.Dependency) bool {
	next := <-enumerator.channel
	if next == nil {
		close(enumerator.channel)
		return false
	}
	*dependency = *next
	return true
}

func main() {
	enumerator, _ := NewFileDependencyEnumerator("test/test-01.in")
	if taskorder.NewResolver(enumerator).Resolve() {
		fmt.Println("Done")
	} else {
		fmt.Println("Fail")
	}
	// taskorder.ResolveDependencies(enumerator)
}
