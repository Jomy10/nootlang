package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jomy10/nootlang/interpreter"
	"github.com/jomy10/nootlang/parser"
	"github.com/jomy10/nootlang/runtime"
	"github.com/jomy10/nootlang/stdlib"
)

// TODO: accept expressions
func main() {
	fmt.Println(`                               $$\     $$\           $$\       
                               $$ |    \__|          $$ |      
$$$$$$$\   $$$$$$\   $$$$$$\ $$$$$$\   $$\  $$$$$$$\ $$$$$$$\  
$$  __$$\ $$  __$$\ $$  __$$\\_$$  _|  $$ |$$  _____|$$  __$$\ 
$$ |  $$ |$$ /  $$ |$$ /  $$ | $$ |    $$ |\$$$$$$\  $$ |  $$ |
$$ |  $$ |$$ |  $$ |$$ |  $$ | $$ |$$\ $$ | \____$$\ $$ |  $$ |
$$ |  $$ |\$$$$$$  |\$$$$$$  | \$$$$  |$$ |$$$$$$$  |$$ |  $$ |
\__|  \__| \______/  \______/   \____/ \__|\_______/ \__|  \__|
		`)

	// Start runtime
	runtime := runtime.NewRuntime(os.Stdout, os.Stderr, os.Stdin)
	stdlib.Register(&runtime)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("$ The nootlang interactive shell v0.0.1")
	for {
		fmt.Print("$ ")

		scanner.Scan()
		source := scanner.Text()

		tokens, err := parser.Tokenize(source)
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			continue
		}
		nodes, err := parser.Parse(tokens)
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			continue
		}

		for _, node := range nodes {
			val, err := interpreter.ExecNode(&runtime, node)
			if err != nil {
				os.Stderr.WriteString(err.Error() + "\n")
				continue
			}
			fmt.Printf("> %v\n", val)
		}
	}
}
