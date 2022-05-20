package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jomy10/nootlang/interpreter"
	"github.com/jomy10/nootlang/parser"
)

func main() {
	fileData, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(fmt.Sprintf("Couldn't read source file %s\n", err.Error()))
	}
	tokens, err := parser.Tokenize(string(fileData))
	if err != nil {
		panic(err)
	}

	nodes, err := parser.Parse(tokens)
	if err != nil {
		panic(err)
	}

	stdout := make(chan string)
	stderr := make(chan string)
	eop := make(chan int, 1)

	defer close(stdout)
	defer close(stderr)
	defer close(eop)

	go interpreter.Interpret(nodes, stdout, stderr, eop)
	go printStdout(stdout)
	go printStderr(stderr)

	exitCode := <-eop
	time.Sleep(2 * time.Millisecond)
	fmt.Printf("Program exited with exit code %d\n", exitCode)
}

func printStdout(stdout chan string) {
	fmt.Println(<-stdout)
}

func printStderr(stderr chan string) {
	fmt.Fprintf(os.Stderr, <-stderr)
}
