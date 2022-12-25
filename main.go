package main

import (
	"fmt"
	"github.com/jomy10/nootlang/interpreter"
	"github.com/jomy10/nootlang/parser"
	"os"
)

func main() {
	tokens, err := parser.Tokenize(`noot!("Hello" + " world")`)
	if err != nil {
		panic(err)
	}
	nodes, err := parser.Parse(tokens)
	if err != nil {
		panic(err)
	}

	if err := interpreter.Interpret(nodes, os.Stdout, os.Stderr, os.Stdin); err != nil {
		panic(fmt.Sprintf("[Runtime error] %v\n", err))
	}
}
