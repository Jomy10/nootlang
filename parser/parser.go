package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Eos struct{}

func (e Eos) Error() string {
	return "End of statement"
}

// Parse tokens into nodes
func Parse(tokens []Token) ([]Node, error) {
	var currentStatement []Token

	nodes := []Node{}

	// Read until end of statement and parse (or end of input)
	blockLevel := 0 // level of curly brackets
	start := 0
	i := 0
	for true {
		if i != len(tokens) {
			if tokens[i].Type == OpenCurlPar {
				blockLevel += 1
			} else if tokens[i].Type == ClosedCurlPar {
				blockLevel -= 1
			}
		}

		if i == len(tokens) || (tokens[i].Type == EOS && blockLevel == 0) {
			currentStatement = tokens[start:i]
			iter := newArrayIterator(currentStatement)
			stmtNode, err := parseStatement(&iter)
			if err != nil {
				if err.Error() != "Empty statement" {
					return nil, err
				}
			} else {
				nodes = append(nodes, stmtNode)
			}
			start = i + 1
		}

		if i == len(tokens) {
			break
		}

		i += 1
	}

	return nodes, nil
}

func parseStatement(tokenIter Iterator[Token]) (Node, error) {
	firstToken, hasFirst := tokenIter.next()
	if !hasFirst {
		return nil, errors.New("Empty statement")
	}

	switch firstToken.Type {
	case Ident:
		secondToken, hasSecond := tokenIter.peek()
		if !hasSecond {
			return nil, errors.New("Invalid statement: lonely identifier")
		}
		switch secondToken.Type {
		case OpenPar:
			return parseFunctionCall(firstToken.Value, tokenIter)
		case Declare:
			fallthrough
		case Equal:
			_, _ = tokenIter.next() // consume :=/=
			exprNode, err := parseExpression(tokenIter)
			if err != nil {
				return nil, err
			}

			if secondToken.Type == Declare {
				return VarDeclNode{firstToken.Value, exprNode}, nil
			} else {
				return VarAssignNode{firstToken.Value, exprNode}, nil
			}
		default:
			return nil, errors.New(fmt.Sprintf("Node %#v is invalid at current position", secondToken))
		}
	case Return:
		expr, err := parseExpression(tokenIter)
		if err != nil {
			return nil, err
		}
		return ReturnNode{expr}, nil
	case Def:
		return parseFunctionDecl(tokenIter)
	default:
		return nil, errors.New(fmt.Sprintf("Node %#v is invalid at current position", firstToken))
	}
}

func parseExpression(tokenIter Iterator[Token]) (Node, error) {
	firstToken, hasFirst := tokenIter.peek()
	if !hasFirst {
		return nil, errors.New("expected expression")
	}

	switch firstToken.Type {
	case OpenPar:
		return parseBinaryExpression(tokenIter)
	case Bool:
		_, hasSecond := tokenIter.peekN(2)
		// handle lonesome boolean literal
		if !hasSecond {
			tokenIter.consume(1) // consume boolean
			boolean, err := strconv.ParseBool(firstToken.Value)
			if err != nil {
				return nil, err
			}
			return BoolLiteralNode{boolean}, nil
		} else {
			return parseBinaryExpression(tokenIter)
		}
	case Integer:
		_, hasSecond := tokenIter.peekN(2)
		// Handle lonesome integer literal
		if !hasSecond {
			tokenIter.consume(1) // consume integer
			integer, err := strconv.ParseInt(firstToken.Value, 10, 64)
			if err != nil {
				return nil, err
			}
			return IntegerLiteralNode{integer}, nil
		} else {
			return parseBinaryExpression(tokenIter)
		}
	case Float:
		_, hasSecond := tokenIter.peekN(2)
		// Handle lonesome float literal
		if !hasSecond {
			tokenIter.consume(1) // consume float
			float, err := strconv.ParseFloat(firstToken.Value, 64)
			if err != nil {
				return nil, err
			}
			return FloatLiteralNode{float}, nil
		} else {
			return parseBinaryExpression(tokenIter)
		}
	case Ident:
		secondToken, hasSecond := tokenIter.peekN(2)
		// Handle lonesome ident
		if !hasSecond {
			tokenIter.consume(1) // consume integer/ident
			return VariableNode{firstToken.Value}, nil
		}

		switch secondToken.Type {
		case OpenPar:
			// function call
			tokenIter.consume(1) // consume integer/ident
			return parseFunctionCall(firstToken.Value, tokenIter)
		default:
			return parseBinaryExpression(tokenIter)
		}
	case Nil:
		tokenIter.consume(1)
		return NilLiteralNode{}, nil
	case String:
		secondToken, hasSecond := tokenIter.peekN(2)
		if !hasSecond {
			tokenIter.consume(1)
			return StringLiteralNode{parseStringLiteral(firstToken.Value)}, nil
		}

		if secondToken.Type == Plus {
			return parseBinaryExpression(tokenIter)
		} else {
			return nil, errors.New(fmt.Sprintf("Did not expect token %s afte string literal\n", secondToken.Value))
		}
	case Not:
		tokenIter.consume(1)
		node, err := parseExpression(tokenIter)
		if err != nil {
			return nil, err
		}
		return BinaryNotNode{node}, nil
	default:
		return nil, errors.New(fmt.Sprintf("Invalid start of expression `%v`", firstToken))
	}
}

func parseFunctionCall(name string, tokenIter Iterator[Token]) (Node, error) {
	args, err := parseFunctionCallArguments(tokenIter)
	return FunctionCallExprNode{name, args}, err
}

// Parse ( args ,* )
func parseFunctionCallArguments(tokenIter Iterator[Token]) ([]Node, error) {
	openPar, hasOpenPar := tokenIter.next()
	if !hasOpenPar {
		return nil, errors.New("Expected opening parenthesis in function call")
	}
	if openPar.Type != OpenPar {
		return nil, errors.New(fmt.Sprintf("Expected opening parenthesis in function call, but got %s", openPar.Value))
	}

	argList, err := collectList(tokenIter, ClosedPar)
	if err != nil {
		return nil, err
	}

	var args []Node
	for _, arg := range argList {
		argIter := newArrayOfPointerIterator(arg)
		expr, err := parseExpression(&argIter)
		if err != nil {
			return nil, errors.New("Expected expression in argument list")
		}
		args = append(args, expr)
	}

	return args, nil
}

// TODO: not operator (is not binary operator)

func isBinaryOperator(token *Token) bool {
	return token.Type == Star || token.Type == Slash || token.Type == Plus || token.Type == Minus ||
		token.Type == DEqual || token.Type == DNEqual || token.Type == And || token.Type == Or
}

// tokenIter starts at the function's name
func parseFunctionDecl(tokenIter Iterator[Token]) (Node, error) {
	funcNameToken, hasNameToken := tokenIter.next()
	if !hasNameToken {
		return nil, errors.New("Expected function name after `def`")
	}
	if funcNameToken.Type != Ident {
		return nil, errors.New(fmt.Sprintf("Expected function name after `def`, get %s", funcNameToken.Value))
	}

	args, err := parseFunctionDeclArgs(tokenIter)
	if err != nil {
		return nil, err
	}

	body, err := parseBody(tokenIter)
	if err != nil {
		return nil, err
	}

	return FunctionDeclNode{funcNameToken.Value, args, body}, nil
}

// Returns the arguments of a function declaration as strings.
// tokenIter starts at the opened paranthesis
func parseFunctionDeclArgs(tokenIter Iterator[Token]) ([]string, error) {
	openedPar, hasOpenedPar := tokenIter.next()
	if !hasOpenedPar {
		return nil, errors.New("Expected opening bracket after function declaration")
	}
	if openedPar.Type != OpenPar {
		return nil, errors.New(fmt.Sprintf("Expected opening brakcet after function declaration, but got %s", openedPar.Value))
	}

	list, err := collectList(tokenIter, ClosedPar)
	if err != nil {
		return nil, err
	}

	var argNames []string
	for _, arg := range list {
		if len(arg) != 1 {
			return nil, errors.New("Expected comma after argument")
		}
		argNames = append(argNames, arg[0].Value)
	}

	return argNames, err
}

// Parse a block of the form `{` (tokens) `}`
// tokenIter starts at the opening curly bracket
func parseBody(tokenIter Iterator[Token]) ([]Node, error) {
	// Curly bracket level
	curlLevel := 0

	i := 1
	nextToken, hasNext := tokenIter.next()
	for true {
		if !hasNext {
			return nil, errors.New("Expected closing curly bracket to match the opening one, but didn't find one")
		}

		switch nextToken.Type {
		case OpenCurlPar:
			curlLevel += 1
		case ClosedCurlPar:
			curlLevel -= 1
			if curlLevel == 0 {
				subIter := tokenIter.subslice(i - 2) // don't include last curly brace
				return Parse(subIter.collect())
			}
		}

		nextToken, hasNext = tokenIter.peekN(i)
		i += 1
	}

	return nil, errors.New("Parser error")
}

// Collect a list of arguments
func collectList(tokenIter Iterator[Token], closingToken TT) ([][]*Token, error) {
	parLevel := 0
	blockLevel := 0

	var tokenArgs = [][]*Token{}
	idx := 0
	for true {
		nextToken, hasNext := tokenIter.next()

		if !hasNext {
			if closingToken == EOS {
				return tokenArgs, nil
			}
			break
		}

		if nextToken.Type == closingToken && parLevel == 0 && blockLevel == 0 {
			return tokenArgs, nil
		} else if nextToken.Type == OpenPar {
			parLevel += 1
		} else if nextToken.Type == ClosedPar {
			parLevel -= 1
		} else if nextToken.Type == OpenCurlPar {
			blockLevel += 1
		} else if nextToken.Type == ClosedCurlPar {
			blockLevel -= 1
		} else if hasNext && nextToken.Type == Comma && parLevel == 0 && blockLevel == 0 {
			idx += 1
		}
		if nextToken.Type != Comma || parLevel != 0 || blockLevel != 0 {
			if idx == len(tokenArgs) {
				tokenArgs = append(tokenArgs, []*Token{})
			}
			tokenArgs[idx] = append(tokenArgs[idx], nextToken)
		}
	}

	return nil, errors.New("Invalid list")
}

func parseStringLiteral(lit string) string {
	// Remove quotes
	substr := lit[1 : len(lit)-1]
	substr = strings.Replace(substr, `\\`, `\`, -1)
	substr = strings.Replace(substr, `\"`, `"`, -1)
	substr = strings.Replace(substr, `\n`, "\n", -1)
	substr = strings.Replace(substr, `\r`, "\r", -1)
	substr = strings.Replace(substr, `\t`, "\t", -1)
	substr = strings.Replace(substr, `\a`, "\a", -1)
	substr = strings.Replace(substr, `\b`, "\b", -1)
	substr = strings.Replace(substr, `\f`, "\f", -1)
	substr = strings.Replace(substr, `\v`, "\v", -1)
	// TODO: \x, \b, \u
	return substr
}
