package parser

import (
	"reflect"
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
	source := "c := a + b"
	expected := []Node{
		VarDeclNode{
			"c",
			BinaryExpressionNode{
				VariableNode{"a"},
				"+",
				VariableNode{"b"},
			},
		},
	}
	testParsing(source, expected, t)
}

func TestAssignmentParse(t *testing.T) {
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

func TestPrint(t *testing.T) {
	source := "noot!(5)"
	expected := []Node{
		FunctionCallExprNode{"noot!", []Node{IntegerLiteralNode{5}}},
	}
	testParsing(source, expected, t)
}

func TestFuncDecl(t *testing.T) {
	source := "def test(arg) { argCpy := arg; return argCpy; }"
	expected := []Node{
		FunctionDeclNode{
			"test",
			[]string{"arg"},
			[]Node{
				VarDeclNode{"argCpy", VariableNode{"arg"}},
				ReturnNode{VariableNode{"argCpy"}},
			},
		},
	}
	testParsing(source, expected, t)
}

func TestFuncCallMultiArguments(t *testing.T) {
	source := "call(a, b)"
	expected := []Node{
		FunctionCallExprNode{
			"call",
			[]Node{
				VariableNode{"a"},
				VariableNode{"b"},
			},
		},
	}
	testParsing(source, expected, t)
}

func TestFuncDeclMultiArguments(t *testing.T) {
	source := "def call(a, b) { return a + b; }"
	expected := []Node{
		FunctionDeclNode{
			"call",
			[]string{"a", "b"},
			[]Node{
				ReturnNode{
					BinaryExpressionNode{
						VariableNode{"a"},
						Operator("+"),
						VariableNode{"b"},
					},
				},
			},
		},
	}
	testParsing(source, expected, t)
}

func TestParseList(t *testing.T) {
	tokens, err := Tokenize("(a, b)")
	if err != nil {
		t.Fatal(err.Error())
	}

	iter := newArrayIterator(tokens)
	iter.consume(1)
	list, err := collectList(&iter, ClosedPar)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := [][]*Token{
		{&Token{Ident, "a"}},
		{&Token{Ident, "b"}},
	}

	if len(list) != len(expected) {
		t.Fatal("Unequal lengths")
	}

	for i := 0; i < len(expected); i++ {
		if *expected[i][0] != *list[i][0] {
			t.Fatalf("Unequal element %v an %v", *expected[i][0], *list[i][0])
		}
	}
}

func TestParseListNested(t *testing.T) {
	tokens, err := Tokenize("(add(a, b))")
	if err != nil {
		t.Fatal(err.Error())
	}

	iter := newArrayIterator(tokens)
	iter.consume(1)
	list, err := collectList(&iter, ClosedPar)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := [][]*Token{
		{&Token{Ident, "add"}, &Token{OpenPar, "("}, &Token{Ident, "a"}, &Token{Comma, ","}, &Token{Ident, "b"}},
	}

	if len(list) != len(expected) {
		t.Fatal("Unequal lengths")
	}

	for i := 0; i < len(expected); i++ {
		if len(expected[i]) != len(list[i]) {
			for j := 0; j < len(expected); j++ {
				if *expected[i][j] != *expected[i][j] {
					t.Fatalf("Unequal element %v an %v", *expected[i][j], *list[i][j])
				}
			}
		}
	}
}

func TestParseNil(t *testing.T) {
	source := "a := nil"
	expected := []Node{
		VarDeclNode{"a", NilLiteralNode{}},
	}
	testParsing(source, expected, t)
}

func TestParseString(t *testing.T) {
	source := "a := \"Hello\""
	expected := []Node{
		VarDeclNode{"a", StringLiteralNode{"Hello"}},
	}
	testParsing(source, expected, t)
}

func TestParseFloat(t *testing.T) {
	source := "a := 6.5"
	expected := []Node{
		VarDeclNode{"a", FloatLiteralNode{6.5}},
	}
	testParsing(source, expected, t)
}

// TODO: new boolean operators
func TestParseBool(t *testing.T) {
	source := "a := true == false"
	expected := []Node{
		VarDeclNode{"a", BinaryExpressionNode{
			BoolLiteralNode{true},
			Operator("=="),
			BoolLiteralNode{false},
		}},
	}
	testParsing(source, expected, t)
}

// TODO:
// func TestOperatorPrecedence(t *testing.T) {
// 	source := "a := a + a * a - a == a != a || a && a"
// 	expected := []Node{}
// 	testParsing(source, expected, t)
// }

func TestIfParse(t *testing.T) {
	source := "if true { a := 1; } elsif false { noot!(5); } else { b := 2; }"
	expected := []Node{
		IfNode{
			BoolLiteralNode{true},
			IfNode{
				BoolLiteralNode{false},
				ElseNode{
					[]Node{
						VarDeclNode{
							"b",
							IntegerLiteralNode{2},
						},
					},
				},
				[]Node{
					FunctionCallExprNode{
						"noot!",
						[]Node{IntegerLiteralNode{5}},
					},
				},
			},
			[]Node{
				VarDeclNode{
					"a",
					IntegerLiteralNode{1},
				},
			},
		},
	}
	testParsing(source, expected, t)
}

func TestWhileLoop(t *testing.T) {
	source := `while true { noot!("infinite") }`
	expected := []Node{
		WhileNode{
			BoolLiteralNode{true},
			[]Node{
				FunctionCallExprNode{"noot!", []Node{StringLiteralNode{"infinite"}}},
			},
		},
	}
	testParsing(source, expected, t)
}

func TestArrayLiteral(t *testing.T) {
	source := `a := [5, 6 * 8, getVal()]`
	expected := []Node{
		VarDeclNode{
			"a",
			ArrayLiteralNode{
				[]Node{
					IntegerLiteralNode{5},
					BinaryExpressionNode{IntegerLiteralNode{6}, Operator("*"), IntegerLiteralNode{8}},
					FunctionCallExprNode{"getVal", nil},
				},
			},
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
		t.Fatalf("Expected and nodes have different sizes\n%#v\n%#v\n", expected, nodes)
	}

	for i, node := range expected {
		// if node != nodes[i] {
		if !reflect.DeepEqual(node, nodes[i]) {
			t.Fatalf("Expected %#v\n But got %#v\n", node, nodes[i])
		}
	}
}
