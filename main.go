package main

import (
	// "bytes"
	"fmt"
	"github.com/jomy10/nootlang/interpreter"
	"github.com/jomy10/nootlang/parser"
	"github.com/jomy10/nootlang/runtime"
	"github.com/jomy10/nootlang/stdlib"
	"os"
)

func main() {
	dat, err := os.ReadFile("/Users/jonaseveraert/Documents/projects/Advent-Of-Code-2022/day15/main.noot")
	if err != nil {
		panic(err)
	}

	tokens, err := parser.Tokenize(string(dat))
	if err != nil {
		panic(err)
	}
	nodes, err := parser.Parse(tokens)
	if err != nil {
		panic(err)
	}

	if err := interpreter.Interpret(nodes, os.Stdout, os.Stderr, os.Stdin, []func(*runtime.Runtime){stdlib.Register}); err != nil {
		panic(fmt.Sprintf("[Runtime error] %v\n", err))
	}
}
