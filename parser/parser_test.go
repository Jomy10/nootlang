package parser

import (
	"testing"
)

func TestAssign(t *testing.T) {
	source := "a := 5"
	expected := []Node{
		AssignmentNode{"a", "5", "int"},
	}

	testParsing(source, expected, t)
}

func testPrint(t *testing.T) {
	source := "noot!(6)"
	expected := []Node{
		PrintNode{
			LiteralNode{"6", "int"},
		},
	}

	testParsing(source, expected, t)
}

func testParsing(source string, expected []Node, t *testing.T) {
	tokens, err := Tokenize(source)
	if err != nil {
		t.Fatal(err.Error())
	}

	nodes, err := Parse(tokens)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(nodes) != len(expected) {
		t.Fatalf("Expected and nodes have different size %#v | %#v", expected, nodes)
	}

	for i, node := range expected {
		if node != nodes[i] {
			t.Fatalf("Expected %#v, but go %#v\n", node, nodes[i])
		}
	}
}
