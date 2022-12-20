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
	nodes := []Node{}
	iter := newArrayIterator(tokens)

	for iter.hasNext() {
		node, err := parseNext(&iter)
		if err != nil {
			if _, ok := err.(*Eos); ok {
				continue
			}
			return nil, err
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func parseNext(tokenIter Iterator[Token]) (Node, error) {
	token, hasNext := tokenIter.next()
	if !hasNext {
		return nil, errors.New("Unexpectedly found end of input")
	}

	switch token.Type {
	case Ident:
		return parseIdent(token, tokenIter)
	case Integer:
		return parseIntegerLiteral(token, tokenIter)
	case Print:
		return parsePrintStmt(tokenIter)
	case EOS:
		return nil, &Eos{}
	}

	return nil, errors.New(fmt.Sprintf("Parser bug (parsing `%s`)", token.Value))
}

// Parse all situations of an identifier followed by something
func parseIdent(ident *Token, tokenIter Iterator[Token]) (Node, error) {
	nextToken, hasNext := tokenIter.next()
	if !hasNext {
		return VariableNode{ident.Value}, nil
	}

	switch nextToken.Type {
	case Declare:
		rhs, err := parseNext(tokenIter)
		if err != nil {
			return nil, err
		}
		return VarDeclNode{
			ident.Value,
			rhs,
		}, nil
	case Equal:
		rhs, err := parseNext(tokenIter)
		if err != nil {
			return nil, err
		}
		return VarAssignNode{
			ident.Value,
			rhs,
		}, nil
	case Plus, Minus, Slash, Star:
		return parseBinaryExpressionSequence(VariableNode{ident.Value}, Operator(nextToken.Value), tokenIter)
	case EOS, ClosedPar:
		return VariableNode{ident.Value}, nil
	}

	return nil, errors.New(fmt.Sprintf("Token %s is invalid at this location", nextToken.Value))
}

// Parse an integer literal
func parseIntegerLiteral(token *Token, tokenIter Iterator[Token]) (Node, error) {
	integer, err := strconv.ParseInt(token.Value, 10, 64)
	if err != nil {
		return nil, err
	}
	nextToken, hasNext := tokenIter.peek()
	if hasNext {
		switch nextToken.Type {
		case Plus, Minus, Slash, Star:
			_, _ = tokenIter.next()
			return parseBinaryExpressionSequence(IntegerLiteralNode{integer}, Operator(nextToken.Value), tokenIter)
		default:
			return IntegerLiteralNode{integer}, nil
		}
	} else {
		return IntegerLiteralNode{integer}, nil
	}
}

// Parse a binary expression
// - `lhs`: The left side of the expression
// - `operator`: The operator of the expression
// - `rhsIter`: The iterator of the tokens starting after the operator
func parseBinaryExpressionSequence(lhs Node, operator Operator, rhsIter Iterator[Token]) (Node, error) {
	rhs, err := parseNext(rhsIter)
	if err != nil {
		return nil, err
	}

	return BinaryExpressionNode{
		lhs,
		operator,
		rhs,
	}, nil
}

// nodeIter starts after `noot!`
func parsePrintStmt(nodeIter Iterator[Token]) (Node, error) {
	openParen, hasNext := nodeIter.next()
	if !hasNext {
		return nil, errors.New("Expected `(` after noot! statement, but got end of file")
	}
	if openParen.Type != OpenPar {
		return nil, errors.New(fmt.Sprintf("Expected `(` after noot! statement, but got %s", openParen.Value))
	}

	inner, err := parseNext(nodeIter)
	if err != nil {
		return nil, err
	}

	closeParen := nodeIter.prev()
	// if !hasNext {
	// 	return nil, errors.New("Expected `)` to close noot! statement, but got end of file")
	// }
	if closeParen.Type != ClosedPar {
		return nil, errors.New(fmt.Sprintf("Expected `)` to close noot! statement, but got %s", closeParen.Value))
	}

	return PrintStmtNode{inner}, nil
}
