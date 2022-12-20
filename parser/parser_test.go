package parser

import (
	"testing"
)

func TestDecl(t *testing.T) {
	source := "a := 5"
	expected := []Node{
		VarDeclNode{"a", IntegerLiteralNode{5}},
	}
	testParsing(source, expected, t)
}

func TestBinaryExpression(t *testing.T) {
	source := "a + b"
	expected := []Node{
		BinaryExpressionNode{
			VariableNode{"a"},
			"+",
			VariableNode{"b"},
		},
	}
	testParsing(source, expected, t)
}

func testAssignment(t *testing.T) {
	source := "a := 0; a = 6 - 5"
	expected := []Node{
		VarDeclNode{
			"a",
			IntegerLiteralNode{0},
		},
		VarAssignNode{
			"a",
			BinaryExpressionNode{
				Left:     IntegerLiteralNode{6},
				Operator: "-",
				Right:    IntegerLiteralNode{5},
			},
		},
	}
	testParsing(source, expected, t)
}

func testPrint(t *testing.T) {
	source := "noot!(5)"
	expected := []Node{
		PrintStmtNode{IntegerLiteralNode{5}},
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
		t.Fatalf("Expected and nodes have different sizes\n%#v\n%#v\n", expected, nodes)
	}

	for i, node := range expected {
		if node != nodes[i] {
			t.Fatalf("Expected %#v\n But got %#v\n", node, nodes[i])
		}
	}
}
