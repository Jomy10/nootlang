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

// Parse tokens into nodeq
func Parse(tokens []Token) ([]Node, error) {
	var currentStatement []Token

	nodes := []Node{}

	// Read until end of statement and parse (or end of input)
	start := 0
	i := 0
	for true {
		if i == len(tokens) || tokens[i].Type == EOS {
			currentStatement = tokens[start:i]
			fmt.Printf("Statement: %v\n", currentStatement)
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
		secondToken, hasSecond := tokenIter.next()
		if !hasSecond {
			return nil, errors.New("Invalid statement: lonely identifier")
		}
		switch secondToken.Type {
		case Declare:
			fallthrough
		case Equal:
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
	fmt.Println("Parsing expression")

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
			panic("TODO: function calls")
		default:
			return parseBinaryExpression(tokenIter)
		}
	default:
		return nil, errors.New(fmt.Sprintf("Invalid start of expression `%v`", firstToken))
	}
}

// ----------

// Parse tokens into nodes
func ParseOld(tokens []Token) ([]Node, error) {
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

func isBinaryOperator(token *Token) bool {
	return token.Type == Star || token.Type == Slash || token.Type == Plus || token.Type == Minus
}

// Collect all tokens inside of brachets. The parameter tokenIter must start at
// the open paranthesis
func collectBracketsInner(tokenIter Iterator[Token]) (Iterator[Token], error) {
	// _, _ = tokenIter.next()
	_, _ = tokenIter.next()
	// i := 1
	bracketLevel := 1
	var nextToken *Token
	var hasNext bool
	innerTokens := []*Token{}
	for true {
		nextToken, hasNext = tokenIter.next()
		if !hasNext {
			return nil, errors.New("Undetermined opening bracket `(`")
		}

		if nextToken.Type == OpenPar {
			bracketLevel += 1
		} else if nextToken.Type == ClosedPar {
			bracketLevel -= 1
			if bracketLevel == 0 {
				// sub := tokenIter.subslice(i - 1) // exclude closing bracket
				// tokenIter.consume(i)             // also consume closing bracket
				// return sub, nil
				innerIter := newArrayOfPointerIterator(innerTokens)
				return &innerIter, nil
			}
		} else if nextToken.Type == EOS {
			return nil, errors.New("Undetermined opening bracket `(`")
		} else {
			innerTokens = append(innerTokens, nextToken)
		}

		// i += 1
	}

	return nil, errors.New("Parser bug (unreachable)")
}

// Collect until end of expression (or end of input)
func collectExpression(tokenIter Iterator[Token]) Iterator[Token] {
	i := 1
	var nextToken *Token
	var hasNext bool
	for true {
		nextToken, hasNext = tokenIter.peekN(i)
		if !hasNext {
			return tokenIter
		}
		if nextToken.Type == EOS {
			sub := tokenIter.subslice(i)
			tokenIter.consume(i)
			return sub
		}
		i += 1
	}
	fmt.Println("Parser bug: unreachable")
	return nil
}

func parseFunctionArguments(tokenIter Iterator[Token]) ([]Node, error) {
	var nextToken *Token
	var hasNext bool
	currentExpression := []*Token{}
	expressions := []Node{}
	for true {
		nextToken, hasNext = tokenIter.next()
		if !hasNext || nextToken.Type == Comma {
			exprIter := newArrayOfPointerIterator(currentExpression)
			expr, err := parseExpression(&exprIter)
			if err != nil {
				return nil, err
			}
			expressions = append(expressions, expr)
			currentExpression = currentExpression[:0]
		}

		currentExpression = append(currentExpression, nextToken)

		if !hasNext {
			break
		}
	}
	return expressions, nil
}

// func parseExpressionOld(tokenIter Iterator[Token]) (Node, error) {
// 	firstToken, hasToken := tokenIter.peek()
// 	if !hasToken {
// 		return nil, errors.New("Empty expression")
// 	}

// 	switch firstToken.Type {
// 	case OpenPar:
// 		innerTokens, err := collectBracketsInner(tokenIter)
// 		if err != nil {
// 			return nil, err
// 		}
// 		// handle case like (5 + 6) * 8
// 		innerExpr, err := parseExpression(innerTokens)
// 		if err != nil {
// 			return nil, err
// 		}
// 		nextToken, hasToken := tokenIter.peek()
// 		if hasToken {
// 			if isBinaryOperator(nextToken) {
// 				return parseBinaryExpression(&innerExpr, tokenIter)
// 			} else {
// 				return nil, errors.New("Invalid expression, expected operator")
// 			}
// 		}
// 		return innerExpr, err
// 	case Integer:
// 		fallthrough
// 	case Ident:
// 		secondToken, hasToken := tokenIter.peekN(2)
// 		if !hasToken || secondToken.Type == EOS {
// 			_, _ = tokenIter.next() // consume integer/ident
// 			if firstToken.Type == Ident {
// 				return VariableNode{firstToken.Value}, nil
// 			} else if firstToken.Type == Integer {
// 				integer, err := strconv.ParseInt(firstToken.Value, 10, 64)
// 				if err != nil {
// 					return nil, err
// 				}
// 				return IntegerLiteralNode{integer}, nil
// 			}
// 		}

// 		switch secondToken.Type {
// 		case OpenPar:
// 			_, _ = tokenIter.next() // ident
// 			innerIter, err := collectBracketsInner(tokenIter)
// 			if err != nil {
// 				return nil, err
// 			}
// 			node, err := parseFunctionArguments(innerIter)
// 			if err != nil {
// 				return nil, err
// 			}
// 			return FunctionCallExprNode{
// 				firstToken.Value,
// 				node,
// 			}, nil
// 		// case Star, Slash, Plus, Minus:
// 		default:
// 			// Collect until end of expression
// 			newIter := collectExpression(tokenIter)
// 			return parseBinaryExpression(nil, newIter)
// 			// default:
// 			// 	return nil, errors.New(fmt.Sprintf("Invalid token `%s` at current position", secondToken.Value))
// 		}
// 	default:
// 		newIter := collectExpression(tokenIter)
// 		return parseBinaryExpression(nil, newIter)
// 		// default:
// 		// 	return nil, errors.New(fmt.Sprintf("Invalid token `%s` at current position", firstToken.Value))
// 	}
// }

// Parse all situations of an identifier followed by something
func parseIdent(ident *Token, tokenIter Iterator[Token]) (Node, error) {
	nextToken, hasNext := tokenIter.next()
	if !hasNext {
		return VariableNode{ident.Value}, nil
	}

	switch nextToken.Type {
	case Declare:
		// rhs, err := parseNext(tokenIter)
		rhs, err := parseExpression(tokenIter)
		if err != nil {
			return nil, err
		}
		return VarDeclNode{
			ident.Value,
			rhs,
		}, nil
	case Equal:
		// rhs, err := parseNext(tokenIter)
		rhs, err := parseExpression(tokenIter)
		p, e := tokenIter.peek()
		fmt.Printf("Parsed expression successfully; %v. Iterator is now at %v - %v\n", rhs, p, e)
		if err != nil {
			return nil, err
		}
		return VarAssignNode{
			ident.Value,
			rhs,
		}, nil
	// case Plus, Minus, Slash, Star:

	// 	tokenIter.reverse(2)
	// 	return parseBinaryExpressionSequence(tokenIter)
	case EOS, ClosedPar:
		return VariableNode{ident.Value}, nil
	}

	return nil, errors.New(fmt.Sprintf("Token %s is invalid at this location", nextToken.Value))
}

// Parse an integer literal
func parseIntegerLiteral(token *Token, tokenIter Iterator[Token]) (Node, error) {
	return nil, errors.New("Parser bug: parsing integer literal in old way")
	// integer, err := strconv.ParseInt(token.Value, 10, 64)
	// if err != nil {
	// 	return nil, err
	// }
	// nextToken, hasNext := tokenIter.peek()
	// if hasNext {
	// 	// switch nextToken.Type {
	// 	// // case Plus, Minus, Slash, Star:
	// 	// // 	// _, _ = tokenIter.next()
	// 	// // 	tokenIter.reverse(1)
	// 	// // 	return parseBinaryExpressionSequence(tokenIter)
	// 	// default:
	// 	// 	return IntegerLiteralNode{integer}, nil
	// 	// }
	// } else {
	// 	return IntegerLiteralNode{integer}, nil
	// }
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
