package parser

import (
	"testing"
)

func TestAssignment(t *testing.T) {
	source := "a := 5"
	expected := []Token{
		{Ident, "a"},
		{Assign, ":="},
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

func testTokenizing(source string, expected []Token, t *testing.T) {
	tokens, err := Tokenize(source)

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(expected) != len(tokens) {
		t.Fatalf("Expected and tokens have different sizes %#v | %#v\n", expected, tokens)
	}

	for i, token := range expected {
		if token != tokens[i] {
			t.Fatalf("Expected %#v, but got %#v\n", token, tokens[i])
		}
	}
}