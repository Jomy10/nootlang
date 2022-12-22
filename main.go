package main

import (
	// "bytes"
	// "fmt"
	"github.com/jomy10/nootlang/interpreter"
	"github.com/jomy10/nootlang/parser"
	"os"
)

func main() {
	tokens, err := parser.Tokenize("myVar := 5\n myVar = myVar + 6 * 2\n noot!(myVar)")
	if err != nil {
		panic(err)
	}
	nodes, err := parser.Parse(tokens)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("Tree: %#v\n", nodes)
	interpreter.Interpret(nodes, os.Stdout, os.Stderr, os.Stdin)
	// runtime := interpreter.NR()
	// buf := new(bytes.Buffer)
	// _, _ = interpreter.ExecNode(&runtime, nodes[0], os.Stdout, os.Stderr, os.Stdin)
	// fmt.Printf("buf %s", buf.String())
}
