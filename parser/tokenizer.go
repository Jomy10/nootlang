package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Tokens
const (
	Ident     string = "ident"
	Assign    string = ":="
	Equal     string = "="
	Integer   string = "int"
	Print     string = "print"
	EOS       string = "eos" // \n or ;
	OpenPar   string = "opar"
	ClosedPar string = "cpar"
)

type Token struct {
	Type  string
	Value string
}

type Pair struct {
	Key string
	Val *regexp.Regexp
}

// Source code to tokens
func Tokenize(source string) ([]Token, error) {
	re := []Pair{
		{Assign, regexp.MustCompile(`\A(:=)`)},
		{Integer, regexp.MustCompile(`\A\b\d+\b`)},
		{EOS, regexp.MustCompile(`\A(\n)`)},
		{EOS, regexp.MustCompile(`\A(;)`)},
		{Print, regexp.MustCompile(`\A(noot!)`)},
		{OpenPar, regexp.MustCompile(`\A\(`)},
		{ClosedPar, regexp.MustCompile(`\A\)`)},
		{Ident, regexp.MustCompile(`\A(\b\w+\b)`)},
	}

	var tokens []Token
	for source != "" {
		token, err := nextToken(&source, &re)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, *token)
	}

	return tokens, nil
}

// Get the next token
func nextToken(source *string, reg *[]Pair) (*Token, error) {
	for _, pair := range *reg {
		re := pair.Val
		ty := pair.Key
		if idx := re.FindStringIndex(*source); len(idx) != 0 {
			value := (*source)[idx[0]:idx[1]]
			*source = strings.Trim((*source)[idx[1]:], " ")
			return &Token{ty, value}, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Couldn't find token for %s", *source))
}
