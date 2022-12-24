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
	Ident   TT = iota
	Declare    // :=
	Equal      // =
	Integer    // \d
	// End Of Statement
	EOS           // \n or ;
	Plus          // +
	Minus         // -
	Slash         // /
	Star          // *
	OpenPar       // (
	ClosedPar     // )
	OpenCurlPar   // {
	ClosedCurlPar // }
	Comma         // ,
	Def           // def
	Return        // return
	Nil           // nil
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
	re := []Pair{
		{Declare, regexp.MustCompile(`\A(:=)`)},
		{Equal, regexp.MustCompile(`\A(=)`)},
		{Plus, regexp.MustCompile(`\A\+`)},
		{Minus, regexp.MustCompile(`\A-`)},
		{Star, regexp.MustCompile(`\A\*`)},
		{Slash, regexp.MustCompile(`\A/`)},
		{Integer, regexp.MustCompile(`\A\b\d+\b`)},
		{EOS, regexp.MustCompile(`\A(\n|;)`)},
		{OpenPar, regexp.MustCompile(`\A\(`)},
		{ClosedPar, regexp.MustCompile(`\A\)`)},
		{OpenCurlPar, regexp.MustCompile(`\A{`)},
		{ClosedCurlPar, regexp.MustCompile(`\A}`)},
		{Comma, regexp.MustCompile(`\A(,)`)},
		{Def, regexp.MustCompile(`\A(def)`)},
		{Return, regexp.MustCompile(`\A(return)`)},
		{Nil, regexp.MustCompile(`\A(nil)`)},
		{Ident, regexp.MustCompile(`\A(\w|!)+`)},
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
