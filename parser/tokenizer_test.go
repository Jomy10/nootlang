package parser

import (
	"fmt"
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
		{Print, "noot!"},
		{OpenPar, "("},
		{Integer, "6"},
		{ClosedPar, ")"},
		{EOS, ";"},
	}

	testTokenizing(source, expected, t)
}

func testOperators(t *testing.T) {
	source := "*+-/ 6 + 4"
	fmt.Println("Testing", source)
	expected := []Token{
		{Star, "*"},
		{Plus, "+"},
		{Minus, "-"},
		{Slash, "+"},
		{Integer, "6"},
		{Plus, "+"},
		{Integer, "4"},
	}

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
