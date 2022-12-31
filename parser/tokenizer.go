package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Token Type
type TT = int

// Tokens
const (
	Ident      TT = iota
	Declare       // :=
	Equal         // =
	PlusEqual     // +=
	MinEqual      // -=
	StarEqual     // *=
	SlashEqual    // /=
	Float         // \d.\d
	Integer       // \d
	String        // ".* \" "
	Bool          // true false

	// Binary operators
	And     // &&
	Or      // ||
	DEqual  // ==
	DNEqual // !=
	LT      // <
	GT      // >
	LTE     // <=
	GTE     // >=
	Plus    // +
	Minus   // -
	Slash   // /
	Star    // *

	Not // !

	// End Of Statement
	EOS             // \n or ;
	OpenPar         // (
	ClosedPar       // )
	OpenCurlPar     // {
	ClosedCurlPar   // }
	OpenSquarePar   // [ // TODO
	ClosedSquarePar // ] // TODO
	Comma           // ,
	Def             // def
	Return          // return
	Nil             // nil
	If              // if
	Else            // else
	Elsif           // elsif
	While           // while
	Dot             // .
)

// A single token
type Token struct {
	Type  TT
	Value string
}

// Pair of token type and its regex definition
type Pair struct {
	Type  TT
	Regex *regexp.Regexp
}

// Source code to tokens
func Tokenize(source string) ([]Token, error) {
	// Token regex definitions
	// The first in the list is the one that is matched first
	re := []Pair{
		{Declare, regexp.MustCompile(`\A(:=)`)},
		{DEqual, regexp.MustCompile(`\A(==)`)},
		{PlusEqual, regexp.MustCompile(`\A(\+=)`)},
		{MinEqual, regexp.MustCompile(`\A(-=)`)},
		{StarEqual, regexp.MustCompile(`\A(\*=)`)},
		{SlashEqual, regexp.MustCompile(`\A(/=)`)},
		{LTE, regexp.MustCompile(`\A(<=)`)},
		{GTE, regexp.MustCompile(`\A(>=)`)},
		{LT, regexp.MustCompile(`\A(<)`)},
		{GT, regexp.MustCompile(`\A(>)`)},
		{Equal, regexp.MustCompile(`\A(=)`)},
		{Plus, regexp.MustCompile(`\A\+`)},
		{Minus, regexp.MustCompile(`\A-`)},
		{Star, regexp.MustCompile(`\A\*`)},
		{Slash, regexp.MustCompile(`\A/`)},
		{Float, regexp.MustCompile(`\A\d+\.\d*`)},
		{Integer, regexp.MustCompile(`\A\b\d+\b`)},
		{String, regexp.MustCompile(`\A"[^"\\]*(\\.[^"\\]*)*"`)},
		{Bool, regexp.MustCompile(`\A(true|false)`)},
		{DNEqual, regexp.MustCompile(`\A(!=)`)},
		{And, regexp.MustCompile(`\A(&&)`)},
		{Or, regexp.MustCompile(`\A(\|\|)`)},
		{Not, regexp.MustCompile(`\A(!)`)},
		{EOS, regexp.MustCompile(`\A(\n|;)`)},
		{OpenPar, regexp.MustCompile(`\A\(`)},
		{ClosedPar, regexp.MustCompile(`\A\)`)},
		{OpenCurlPar, regexp.MustCompile(`\A{`)},
		{ClosedCurlPar, regexp.MustCompile(`\A}`)},
		{OpenSquarePar, regexp.MustCompile(`\A\[`)},
		{ClosedSquarePar, regexp.MustCompile(`\A\]`)},
		{Comma, regexp.MustCompile(`\A(,)`)},
		{Dot, regexp.MustCompile(`\A(\.)`)},
		{Def, regexp.MustCompile(`\A(def)`)},
		{Return, regexp.MustCompile(`\A(return)`)},
		{Nil, regexp.MustCompile(`\A(nil)`)},
		{If, regexp.MustCompile(`\A(if)`)},
		{Else, regexp.MustCompile(`\A(else)`)},
		{Elsif, regexp.MustCompile(`\A(elsif)`)},
		{While, regexp.MustCompile(`\A(while)`)},
		{Ident, regexp.MustCompile(`\A(\w|!|\?)+`)},
	}

	source = strings.TrimSpace(source)
	source = strings.Replace(source, "\t", "", -1)
	// source = strings.Replace(source, "\n", "", -1)

	// Collect tokens
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
// - `source`: the remaining part of the source that needs to be tokenized
// - `reg`: the tokens and their regex definitions
func nextToken(source *string, reg *[]Pair) (*Token, error) {
	for _, pair := range *reg {
		re := pair.Regex
		ty := pair.Type
		if idx := re.FindStringIndex(*source); len(idx) != 0 {
			if idx[1] == 0 {
				continue
			}
			value := (*source)[idx[0]:idx[1]]
			*source = strings.Trim((*source)[idx[1]:], " ")
			return &Token{ty, value}, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Couldn't find token for `%s`", *source))
}
