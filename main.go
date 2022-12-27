package main

import (
	// "bytes"
	"fmt"
	"github.com/jomy10/nootlang/interpreter"
	"github.com/jomy10/nootlang/parser"
	"os"
)

func main() {
	tokens, err := parser.Tokenize(`a := [5, 6 + 2, 8]; noot!(a); a[1] = 1; noot!(a)`)
	if err != nil {
		panic(err)
	}
	nodes, err := parser.Parse(tokens)
	if err != nil {
		panic(err)
	}
	// stdout := new(bytes.Buffer)
	if err := interpreter.Interpret(nodes, os.Stdout, os.Stderr, os.Stdin); err != nil {
		panic(fmt.Sprintf("[Runtime error] %v\n", err))
	}
	// fmt.Println(stdout.String())
}
