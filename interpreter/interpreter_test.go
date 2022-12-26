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
	nodes := nodes("def fn(arg) { noot!(arg); }; fn(56)", t)

	bufStd := new(bytes.Buffer)
	bufErr := new(bytes.Buffer)

	err := Interpret(nodes, bufStd, bufErr, os.Stdin)
	if err != nil {
		t.Fatal(err)
	}

	if bufStd.String() != "56\n" {
		t.Fatal(fmt.Sprintf("got stdout %s", bufStd.String()))
	}
	if bufErr.Len() > 0 {
		t.Fatal(fmt.Sprintf("got stderr %s", bufErr.String()))
	}
}

func TestMultiArgument(t *testing.T) {
	nodes := nodes("def add(a, b) { return a + b }; noot!(add(1, 2))", t)

	bufStd := new(bytes.Buffer)
	bufErr := new(bytes.Buffer)

	err := Interpret(nodes, bufStd, bufErr, os.Stdin)
	if err != nil {
		t.Fatal(err)
	}

	if bufStd.String() != "3\n" {
		t.Fatal(fmt.Sprintf("got stdout %s", bufStd.String()))
	}
	if bufErr.Len() > 0 {
		t.Fatal(fmt.Sprintf("got stderr %s", bufErr.String()))
	}
}

func TestFloatMath(t *testing.T) {
	// Because the first value is a float, the others will automatically be converted to floats
	nodes := nodes("noot!(6.5 + 4 - 0.5)", t)

	fmt.Printf("Nodes: %#v\n", nodes)

	bufStd := new(bytes.Buffer)
	bufErr := new(bytes.Buffer)

	err := Interpret(nodes, bufStd, bufErr, os.Stdin)

	if err != nil {
		t.Fatal(err)
	}

	if bufStd.String() != "10\n" {
		t.Fatal(fmt.Sprintf("Got stdout '%s'", bufStd.String()))
	}
	if bufErr.Len() > 0 {
		t.Fatal(fmt.Sprintf("got stderr %s", bufErr.String()))
	}
}

func TestBoolExpression(t *testing.T) {
	nodes := nodes("noot!(5 != 6)", t)

	bufStd := new(bytes.Buffer)
	bufErr := new(bytes.Buffer)

	err := Interpret(nodes, bufStd, bufErr, os.Stdin)

	if err != nil {
		t.Fatal(err)
	}

	if bufStd.String() != "true\n" {
		t.Fatal(fmt.Sprintf("Got stdout '%s'", bufStd.String()))
	}
	if bufErr.Len() > 0 {
		t.Fatal(fmt.Sprintf("got stderr %s", bufErr.String()))
	}
}

func TestBoolExpression2(t *testing.T) {
	nodes := nodes("noot!(true == false)", t)

	bufStd := new(bytes.Buffer)
	bufErr := new(bytes.Buffer)

	err := Interpret(nodes, bufStd, bufErr, os.Stdin)

	if err != nil {
		t.Fatal(err)
	}

	if bufStd.String() != "false\n" {
		t.Fatal(fmt.Sprintf("Got stdout '%s'", bufStd.String()))
	}
	if bufErr.Len() > 0 {
		t.Fatal(fmt.Sprintf("got stderr %s", bufErr.String()))
	}
}
