package parser

import (
	"testing"
)

func TestAssignment(t *testing.T) {
	source := "a := 5"
	expected := []Token{
		{Ident, "a"},
		{Declare, ":="},
		{Integer, "5"},
	}

	testTokenizing(source, expected, t)
}

func TestNoot(t *testing.T) {
	source := "noot!(6);"
	expected := []Token{
		{Ident, "noot!"},
		{OpenPar, "("},
		{Integer, "6"},
		{ClosedPar, ")"},
		{EOS, ";"},
	}

	testTokenizing(source, expected, t)
}

func TestOperators(t *testing.T) {
	source := "*+-/ 6 + 4 == != || && !true > < >= <="
	expected := []Token{
		{Star, "*"},
		{Plus, "+"},
		{Minus, "-"},
		{Slash, "/"},
		{Integer, "6"},
		{Plus, "+"},
		{Integer, "4"},
		{DEqual, "=="},
		{DNEqual, "!="},
		{Or, "||"},
		{And, "&&"},
		{Not, "!"},
		{Bool, "true"},
		{GT, ">"},
		{LT, "<"},
		{GTE, ">="},
		{LTE, "<="},
	}

	testTokenizing(source, expected, t)
}

func TestFunction(t *testing.T) {
	source := "def f(arg1, arg2) { return arg1 }"
	expected := []Token{
		{Def, "def"},
		{Ident, "f"},
		{OpenPar, "("},
		{Ident, "arg1"},
		{Comma, ","},
		{Ident, "arg2"},
		{ClosedPar, ")"},
		{OpenCurlPar, "{"},
		{Return, "return"},
		{Ident, "arg1"},
		{ClosedCurlPar, "}"},
	}

	testTokenizing(source, expected, t)
}

func TestNil(t *testing.T) {
	source := "nil"
	expected := []Token{{Nil, "nil"}}

	testTokenizing(source, expected, t)
}

func TestStringToken(t *testing.T) {
	source := "\"Hello \\\" World\""
	expected := []Token{{String, source}}

	testTokenizing(source, expected, t)
}

func TestFloatToken(t *testing.T) {
	source := "1. 4.56"
	expected := []Token{{Float, "1."}, {Float, "4.56"}}

	testTokenizing(source, expected, t)
}

func TestBoolToken(t *testing.T) {
	source := "true false"
	expected := []Token{{Bool, "true"}, {Bool, "false"}}
	testTokenizing(source, expected, t)
}

func TestIfElseToken(t *testing.T) {
	source := "if true { } elsif { } else {}"
	expected := []Token{
		{If, "if"}, {Bool, "true"}, {OpenCurlPar, "{"}, {ClosedCurlPar, "}"},
		{Elsif, "elsif"}, {OpenCurlPar, "{"}, {ClosedCurlPar, "}"},
		{Else, "else"}, {OpenCurlPar, "{"}, {ClosedCurlPar, "}"},
	}
	testTokenizing(source, expected, t)
}

func TestLoopingTokens(t *testing.T) {
	source := "while"
	expected := []Token{{While, "while"}}
	testTokenizing(source, expected, t)
}

func TestSquareBracket(t *testing.T) {
	source := "[]"
	expected := []Token{{OpenSquarePar, "["}, {ClosedSquarePar, "]"}}
	testTokenizing(source, expected, t)
}

func TestAssignmentOperatorsTokens(t *testing.T) {
	source := "+= -= *= /="
	expected := []Token{{PlusEqual, "+="}, {MinEqual, "-="}, {StarEqual, "*="}, {SlashEqual, "/="}}
	testTokenizing(source, expected, t)
}

func TestDotToken(t *testing.T) {
	source := "."
	expected := []Token{{Dot, "."}}
	testTokenizing(source, expected, t)
}

func testTokenizing(source string, expected []Token, t *testing.T) {
	tokens, err := Tokenize(source)

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(expected) != len(tokens) {
		t.Fatalf("Expected and tokens have different sizes\n%#v\n%#v\n", expected, tokens)
	}

	for i, token := range expected {
		if token != tokens[i] {
			t.Fatalf("Expected %#v, but got %#v\n", token, tokens[i])
		}
	}
}
