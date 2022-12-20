package interpreter

import (
	"bytes"
	"fmt"
	"github.com/jomy10/nootlang/parser"
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

func TestPrintString(t *testing.T) {
	nodes := nodes("noot!(5)", t)
	runtime := newRuntime()
	bufStd := new(bytes.Buffer)
	bufErr := new(bytes.Buffer)
	ExecNode(&runtime, nodes[0], bufStd, bufErr, os.Stdin)
	if bufStd.String() != "5\n" {
		t.Fatal(fmt.Sprintf("got %s", bufStd.String()))
	}
	if bufErr.Len() > 0 {
		t.Fatal(fmt.Sprintf("Got stderr %s", bufErr.String()))
	}
}
