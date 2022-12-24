package main

import (
	"github.com/jomy10/nootlang/interpreter"
	"github.com/jomy10/nootlang/parser"
	"os"
)

func main() {
	tokens, err := parser.Tokenize("def add(a, b) { return a + b; }\n noot!(add(1, 2))")
	if err != nil {
		panic(err)
	}
	nodes, err := parser.Parse(tokens)
	if err != nil {
		panic(err)
	}

	interpreter.Interpret(nodes, os.Stdout, os.Stderr, os.Stdin)
}
