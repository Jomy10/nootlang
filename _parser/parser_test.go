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

func TestPrint(t *testing.T) {
	source := "noot!(6)"
	expected := []Node{
		PrintNode{
			LiteralNode{"6", "int"},
		},
	}

	testParsing(source, expected, t)
}

func TestPrintVar(t *testing.T) {
	source := "a := 7; noot!(a)"
	expected := []Node{
		AssignmentNode{"a", "7", "int"},
		PrintNode{
			IdentifierNode{"a"},
		},
	}

	testParsing(source, expected, t)
}

func TestAddition(t *testing.T) {
	source := "a + b"
	expected := []Node{
		AdditionNode{
			IdentifierNode{"a"},
			IdentifierNode{"b"},
		},
	}

	testParsing(source, expected, t)
}

func TestAdditionInteger(t *testing.T) {
	source := "1 + 4"
	expected := []Node{
		AdditionNode{
			LiteralNode{"1", "int"},
			LiteralNode{"4", "int"},
		},
	}

	testParsing(source, expected, t)
}

func TestAdditionMix(t *testing.T) {
	source := "a + 4"
	expected := []Node{
		AdditionNode{
			IdentifierNode{"a"},
			LiteralNode{"4", "int"},
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
