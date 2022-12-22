package parser

import (
	"errors"
	"fmt"
	"strconv"
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
	start := 0
	i := 0
	for true {
		if i == len(tokens) || tokens[i].Type == EOS {
			currentStatement = tokens[start:i]
			iter := newArrayIterator(currentStatement)
			stmtNode, err := parseStatement(&iter)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, stmtNode)
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
	case Integer:
		// fallthrough
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
	default:
		return nil, errors.New(fmt.Sprintf("Invalid start of expression `%v`", firstToken))
	}
}

func parseFunctionCall(name string, tokenIter Iterator[Token]) (Node, error) {
	args, err := parseFunctionArguments(tokenIter)
	return FunctionCallExprNode{name, args}, err
}

// Parse ( args ,* )
func parseFunctionArguments(tokenIter Iterator[Token]) ([]Node, error) {
	argTokens := []*Token{}
	args := []Node{}

	nextToken, hasNext := tokenIter.next()
	for hasNext && nextToken.Type != ClosedPar {
		if nextToken.Type == Comma {
			argsIter := newArrayOfPointerIterator(argTokens)
			expr, err := parseExpression(&argsIter)
			if err != nil {
				return nil, err
			}
			args = append(args, expr)
		} else {
			argTokens = append(argTokens, nextToken)
		}

		nextToken, hasNext = tokenIter.next()
	}
	// collect last argument
	argsIter := newArrayOfPointerIterator(argTokens)
	expr, err := parseExpression(&argsIter)
	if err != nil {
		return nil, err
	}
	args = append(args, expr)

	return args, nil
}

func isBinaryOperator(token *Token) bool {
	return token.Type == Star || token.Type == Slash || token.Type == Plus || token.Type == Minus
}
