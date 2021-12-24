# dependency-resolver

* Build demo program: `make build`
* Run tests: `make test`
* Draw graph with Graphviz: `make draw`
* Clean up: `make clean`

## Test file formats

### `*.in`

* one dependency each line
* A depends on B written as `A B`
* see `test/sample.in`

### `*.gz`

Graphviz formatted digraph

* see `test/alpine.gv`
