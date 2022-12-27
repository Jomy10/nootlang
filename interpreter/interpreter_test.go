package interpreter

import (
	"bytes"
	"fmt"
	"github.com/jomy10/nootlang/parser"
	"github.com/jomy10/nootlang/runtime"
	"github.com/jomy10/nootlang/stdlib"
	"os"
	"testing"
)

func nodes(source string, t *testing.T) []parser.Node {
	tokens, err := parser.Tokenize(source)
	if err != nil {
		t.Fatal(err.Error())
	}
	nodes, err := parser.Parse(tokens)
	if err != nil {
		t.Fatal(err.Error())
	}
	return nodes
}

func TestPrint(t *testing.T) {
	nodes := nodes("noot!(5)", t)
	bufStd := new(bytes.Buffer)
	bufErr := new(bytes.Buffer)
	runtime := runtime.NewRuntime(bufStd, bufErr, os.Stdin)
	stdlib.Register(&runtime)
	n, err := ExecNode(&runtime, nodes[0])
	if err != nil {
		t.Fatal(err)
	}
	if n != "5\n" {
		t.Fatal()
	}

	if bufStd.String() != "5\n" {
		t.Fatal(fmt.Sprintf("got %s", bufStd.String()))
	}
	if bufErr.Len() > 0 {
		t.Fatal(fmt.Sprintf("Got stderr %s", bufErr.String()))
	}
}

func TestFunction(t *testing.T) {
	testWithOutput("def fn(arg) { noot!(arg); }; fn(56)", "56\n", t)
}

func TestMultiArgument(t *testing.T) {
	testWithOutput("def add(a, b) { return a + b }; noot!(add(1, 2))", "3\n", t)
}

func TestFloatMath(t *testing.T) {
	// 4 will be converted to float because others are floats
	testWithOutput("noot!(6.5 + 4 - 0.5)", "10\n", t)
}

func TestBoolExpression(t *testing.T) {
	testWithOutput("noot!(5 != 6)", "true\n", t)
}

func TestBoolExpression2(t *testing.T) {
	testWithOutput("noot!(!true)", "false\n", t)
}

func TestIf(t *testing.T) {
	testWithOutput("if true { noot!(\"works\")}", "works\n", t)
}

func TestElsif(t *testing.T) {
	testWithOutput(`if 1 == 2 { noot!("wrong") } elsif true { noot!("correct") }`, "correct\n", t)
}

func TestElse(t *testing.T) {
	testWithOutput(`if 2.0 != 2.0 { noot!("wrong 1") } elsif !true { noot!("wrong 2") } else { noot!("correct") }`, "correct\n", t)
}

// Test interpreter and check its stdout
func testWithOutput(source string, expectedStdout string, t *testing.T) {
	nodes := nodes(source, t)

	bufStd := new(bytes.Buffer)
	bufErr := new(bytes.Buffer)

	err := Interpret(nodes, bufStd, bufErr, os.Stdin)

	if err != nil {
		t.Fatal(err)
	}

	if bufStd.String() != expectedStdout {
		t.Fatal(fmt.Sprintf("Got stdout '%s', but expected %v", bufStd.String(), expectedStdout))
	}
	if bufErr.Len() > 0 {
		t.Fatal(fmt.Sprintf("got stderr %s", bufErr.String()))
	}
}
