package parser

import (
	"errors"
	"fmt"
)

func parseBinaryExpression(tokenIter Iterator[Token]) (Node, error) {
	fmt.Println("Parsing binary expression")

	// Of the form Node Operator Node Operator Node (...)
	// where Operator is of type *Token
	expression := []interface{}{}
	_ = expression

	// parentheses level
	parLevel := 0

	lhs := []*Token{}
	nextToken, hasNext := tokenIter.next()
	for hasNext {
		fmt.Printf("Next token: %v\n", nextToken)
		switch nextToken.Type {
		case OpenPar:
			parLevel += 1
		case ClosedPar:
			parLevel -= 1
		case Star, Slash, Plus, Minus:
			if parLevel == 0 {
				fmt.Printf("LHS: %v\n", lhs)
				lhsIter := newArrayOfPointerIterator(lhs)
				exprNode, err := parseExpression(&lhsIter)
				if err != nil {
					return nil, err
				}
				expression = append(expression, exprNode, nextToken)
				lhs = lhs[:0]
			} else {
				lhs = append(lhs, nextToken)
			}
		default:
			lhs = append(lhs, nextToken)
		}

		nextToken, hasNext = tokenIter.next()
	}

	lhsIter := newArrayOfPointerIterator(lhs)
	exprNode, err := parseExpression(&lhsIter)
	if err != nil {
		return nil, err
	}
	expression = append(expression, exprNode)

	fmt.Printf("Expression: %v\n", expression)

	return __parseBinaryExpression(expression)
}

func __parseBinaryExpression(expr []interface{}) (Node, error) {
	fmt.Printf("Binary expression: %v\n", expr)
	if len(expr) == 1 {
		switch expr[0].(type) {
		case Node:
			return expr[0].(Node), nil
		case *Token:
			return nil, errors.New("Invalid: expected literal or identifier before and after operator")
		default:
			return nil, errors.New("Invalid: expected literal or identifier before and after operator")
		}
	} else if len(expr) == 0 {
		return nil, errors.New("Empty expression")
	}

	precedence := []TT{Minus, Plus, Slash, Star} // Operator precedence in reverse order
	precedenceIdx := 0

	exprIdx := 0

	for precedenceIdx != len(precedence) {
		switch expr[exprIdx].(type) {
		case *Token:
			token := expr[exprIdx].(*Token)
			if token.Type == precedence[precedenceIdx] {
				lhs, err := __parseBinaryExpression(expr[:exprIdx])
				if err != nil {
					return nil, err
				}
				rhs, err := __parseBinaryExpression(expr[exprIdx+1:])
				if err != nil {
					return nil, err
				}
				return BinaryExpressionNode{
					lhs,
					Operator(token.Value),
					rhs,
				}, nil
			}
		}

		exprIdx += 1

		if exprIdx == len(expr) {
			exprIdx = 0
			precedenceIdx += 1
		}
	}

	return nil, errors.New("Parser bug?")
}

// func parseBinaryExpression(tokenIter Iterator[Token]) (Node, error) {
// 	// Of the form Node Operator Node Operator Node (...)
// 	// where Operator is of type *Token
// 	expression := []interface{}{}

// 	// Parentheses level
// 	parLevel := 0

// 	nextToken, hasNext := tokenIter.next()
// 	lhs := []*Token{}
// 	for hasNext && nextToken.Type != EOS {
// 		fmt.Printf("Next in bin expr: %v\n", nextToken)
// 		// Update parenthesis level
// 		if nextToken.Type == OpenPar {
// 			parLevel += 1
// 		} else if nextToken.Type == ClosedPar {
// 			parLevel -= 1
// 		} else if isBinaryOperator(nextToken) && parLevel == 0 {
// 			fmt.Println("Found operator")
// 			// if len(lhs) != 1 { // when starts with operator (e.g. in expression (5 + 6) * 9)
// 			// 	lhsIter := newArrayOfPointerIterator(lhs)
// 			// 	expr, err := parseExpression(&lhsIter)
// 			// 	if err != nil {
// 			// 		return nil, err
// 			// 	}
// 			// 	expression = append(expression, expr)
// 			// }
// 			// // Append the operator
// 			// expression = append(expression, nextToken)
// 			// // Reset lhs
// 			// lhs = lhs[:0]
// 			if err := addLhs(lhs, expression); err != nil {
// 				return nil, err
// 			}
// 		} else {
// 			lhs = append(lhs, nextToken)
// 		}

// 		nextToken, hasNext = tokenIter.next()
// 	}
// 	err := addLhs(lhs, expression)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// return __parseBinaryExpression(expression)
// 	return nil, nil
// }

// func addLhs(lhs []*Token, expression []interface{}) error {
// 	// if len(lhs) == 1 {
// 	// 	return errors.New(fmt.Sprintf("Parser bug: lhs is 1 in binary expression (%v)", *lhs[0]))
// 	// }
// 	for _, p := range lhs {
// 		fmt.Printf("lhs: %v\n", p)
// 	}
// 	// fmt.Printf("Iterating expression over: %v\n", lhs)
// 	lhsIter := newArrayOfPointerIterator(lhs)
// 	expr, err := parseExpression(&lhsIter)
// 	if err != nil {
// 		return err
// 	}
// 	expression = append(expression, expr)
// 	// Append the operator
// 	expression = append(expression, nextToken)
// 	// Reset lhs
// 	lhs = lhs[:0]

// 	return nil
// }

// ------------

// tokenIter only contains items relevant for the binary parseBinaryExpression
// - `lhs` is an optional parameter
func parseBinaryExpressionOld(lhsNode *Node, tokenIter Iterator[Token]) (Node, error) {
	// Of the form Node Operator Node Operator Node (...)
	// where Operator is of type *Token
	expression := []interface{}{}
	if lhsNode != nil {
		expression = append(expression, *lhsNode)
	}
	// parenthesis level
	parLevel := 0
	nextToken, hasNext := tokenIter.next()
	lhs := []*Token{}
	for hasNext && nextToken.Type != EOS {
		if nextToken.Type == OpenPar {
			parLevel += 1
			lhs = append(lhs, nextToken)
		} else if nextToken.Type == ClosedPar {
			parLevel -= 1
			lhs = append(lhs, nextToken)
		} else if isBinaryOperator(nextToken) && parLevel == 0 {
			lhs = append(lhs, &Token{EOS, ""})
			if len(lhs) != 1 { // when starts with operator (e.g. in expression (5 + 6) * 9)
				lhsIter := newArrayOfPointerIterator(lhs)
				expr, err := parseExpression(&lhsIter)
				if err != nil {
					return nil, err
				}
				expression = append(expression, expr)
			}
			// Append the operator
			expression = append(expression, nextToken)
			lhs = lhs[:0]
			// fmt.Printf("Operator and lhs done %v\n", expr)
		} else {
			lhs = append(lhs, nextToken)
		}

		nextToken, hasNext = tokenIter.next()
	}
	lhs = append(lhs, &Token{EOS, ""})
	lhsIter := newArrayOfPointerIterator(lhs)
	expr, err := parseExpression(&lhsIter)
	if err != nil {
		return nil, err
	}
	expression = append(expression, expr)

	return __parseBinaryExpression(expression)
}

func __parseBinaryExpressionOld(expr []interface{}) (Node, error) {
	if len(expr) == 1 {
		switch expr[0].(type) {
		case Token:
			return nil, errors.New("Invalid")
		}
		return expr[0].(Node), nil
	} else if len(expr) == 0 {
		return nil, errors.New("Lhss of expression is empty")
	}

	precedence := []TT{Minus, Plus, Slash, Star} // Operator precedence in reverse order
	precedenceIdx := 0

	exprIdx := 0

	for precedenceIdx != len(precedence) {
		switch expr[exprIdx].(type) {
		case *Token:
			token := expr[exprIdx].(*Token)
			if token.Type == precedence[precedenceIdx] {
				lhs, err := __parseBinaryExpression(expr[:exprIdx])
				if err != nil {
					return nil, err
				}
				rhs, err := __parseBinaryExpression(expr[exprIdx+1:])
				if err != nil {
					return nil, err
				}
				return BinaryExpressionNode{
					lhs,
					Operator(token.Value),
					rhs,
				}, nil
			}
		}

		exprIdx += 1
		if exprIdx == len(expr) {
			exprIdx = 0
			precedenceIdx += 1
		}
	}

	return nil, errors.New("Invalid binary expression")
}

// func parseBinaryExpressionSequence(tokenIter Iterator[Token]) (Node, error) {
// 	// Collect until end of expression
// 	bracketLevel := 0
// 	nextToken, hasMore := tokenIter.next()
// 	binaryExpressionTokens := []*Token{}
// 	for hasMore {
// 		if nextToken.Type == OpenPar {
// 			bracketLevel += 1
// 		} else if nextToken.Type == ClosedPar {
// 			bracketLevel -= 1
// 		} else if nextToken.Type == EOS /*&& bracketLevel == 0*/ {
// 			break
// 		}
// 		binaryExpressionTokens = append(binaryExpressionTokens, nextToken)

// 		nextToken, hasMore = tokenIter.next()
// 	}

// 	fmt.Printf("Collected: %v", binaryExpressionTokens)
// 	return __parseBinaryExpression(binaryExpressionTokens)
// }

// // Loop through the operator precedences in order.
// // In the code from left to right.
// // When an operator is found, then we call this function on the lhs and rhs
// // of that operator and return a BinaryExpressionNode with those
// func __parseBinaryExpression(tokens []*Token) (Node, error) {
// 	// precedence := []TT{Star, Slash, Plus, Minus}
// 	precedence := []TT{Minus, Plus, Slash, Star} // operator precedence in reversed order
// 	currentPrec := 0
// 	codePtr := 0
// 	openLevel := 0 // level of open brackets
// 	for true {
// 		if tokens[codePtr].Type == OpenPar {
// 			fmt.Println("OpenLevel += 1")
// 			openLevel += 1
// 		} else if tokens[codePtr].Type == ClosedPar {
// 			openLevel -= 1
// 		} else if tokens[codePtr].Type == precedence[currentPrec] && openLevel == 0 {
// 			fmt.Printf("Found operator %s\n", tokens[codePtr].Value)
// 			lhs := tokens[:codePtr]
// 			rhs := tokens[codePtr+1:]
// 			lhsNode, err := __parseBinaryExpression(lhs)
// 			if err != nil {
// 				return nil, err
// 			}
// 			rhsNode, err := __parseBinaryExpression(rhs)
// 			if err != nil {
// 				return nil, err
// 			}

// 			return BinaryExpressionNode{
// 				lhsNode,
// 				Operator(tokens[codePtr].Value),
// 				rhsNode,
// 			}, nil
// 		}

// 		codePtr += 1
// 		if codePtr == len(tokens) {
// 			codePtr = 0
// 			currentPrec += 1
// 		}
// 		if currentPrec == len(precedence) {
// 			break
// 		}
// 	}

// 	// Did not return => no binary expression left
// 	tokensIter := newArrayOfPointerIterator(tokens)
// 	fmt.Printf("Calling parseNext in binary expression parser with iterator over {")
// 	for _, token := range tokens {
// 		fmt.Printf("%v", token.Value)
// 	}
// 	fmt.Printf("}\n")
// 	return parseNext(&tokensIter)
// }

// TODO: /\A.*\*.*$/

// // Parse a (possible) sequence of binary operator expresssions
// func parseBinaryExpressionSequence(lhs *Token, operator *Token, rhsIter Iterator[Token]) (Node, error) {
// 	// Collect until end of expression
// 	bracketLevel := 0
// 	nextToken, hasMore := rhsIter.next()
// 	var binaryExpressionTokens []*Token = []*Token{lhs, operator}
// 	for hasMore {
// 		if nextToken.Type == OpenPar {
// 			bracketLevel += 1
// 		} else if nextToken.Type == ClosedPar {
// 			if bracketLevel == 0 {
// 				break
// 			} else {
// 				bracketLevel -= 1
// 			}
// 		} else if nextToken.Type == EOS {
// 			break
// 		}
// 		binaryExpressionTokens = append(binaryExpressionTokens, nextToken)
// 		nextToken, hasMore = rhsIter.next()
// 	}

// 	// return parseBinaryExpression(lhs, operator, rhsIter)
// 	return parseBinaryExpressionSequence__OperatorPrecedence(binaryExpressionTokens)
// }

// // Parse a binary expression in the order of operator precedence
// func parseBinaryExpressionSequence__OperatorPrecedence(tokens []*Token) (Node, error) {
// 	// Contains the operators in the order they are parsed
// 	precedence := []TT{Star, Slash, Plus, Minus}
// 	operatorPtr := 0
// 	codePtr := 0
// 	bracketLevel := 0
// 	var result BinaryExpressionNode
// 	var resultLevel *BinaryExpressionNode = nil
// 	for operatorPtr != len(precedence) {
// 		if tokens[codePtr].Type == OpenPar {
// 			bracketLevel += 1
// 		} else if tokens[codePtr].Type == ClosedPar {
// 			bracketLevel -= 1
// 		} else if tokens[codePtr].Type == precedence[operatorPtr] && bracketLevel == 0 {
// 			// found first occurance of current operator
// 			if tokens[codePtr-1].Type == ClosedPar {
// 				innerBracketLevel := 1
// 				i := codePtr - 2
// 				for true { // tokens[i].Type != OpenPar && innerBracketLevel != 0 {
// 					if tokens[i].Type == ClosedPar {
// 						innerBracketLevel += 1
// 					} else if tokens[i].Type == OpenPar {
// 						innerBracketLevel -= 1
// 					}

// 					if tokens[i].Type == OpenPar && innerBracketLevel == 0 {

// 					}

// 					i -= 1
// 				}

// 				innerTokens := tokens[i+1 : codePtr-1]
// 				innerTokensIter := newArrayOfPointerIterator(innerTokens)
// 				innerNodes, err := parseNext(&innerTokensIter)
// 				if err != nil {
// 					return nil, err
// 				}

// 			} else {
// 				lhs := []*Token{tokens[codePtr-1]}
// 				lhsIter := newArrayOfPointerIterator(lhs)
// 				lhsNode, err := parseNext(&lhsIter)
// 				if err != nil {
// 					return nil, err
// 				}
// 				rhsIter := newArrayOfPointerIterator(tokens[codePtr+1:])
// 				exprNode, err := __parseBinaryExpression(lhsNode, Operator(tokens[codePtr].Value), &rhsIter)
// 				if err != nil {
// 					return nil, err
// 				}
// 				if resultLevel == nil {
// 					result = exprNode
// 					resultLevel = &result
// 				}
// 			}
// 		}
// 		codePtr += 1
// 		if codePtr == len(tokens) {
// 			codePtr = 0
// 			operatorPtr += 1
// 		}
// 	}
// 	return result, nil
// }

// func __addToResult(result *BinaryExpressionNode, resultLevel *BinaryExpressionNode, expr *BinaryExpressionNode) {
// 	if resultLevel == nil {
// 		*result = *expr
// 		resultLevel = result
// 	} else {
// 		resultLevel.Right = expr
// 	}
// }

// // Parse a binary expression
// // - `lhs`: The left side of the expression
// // - `operator`: The operator of the expression
// // - `rhsIter`: The iterator of the tokens starting after the operator
// func __parseBinaryExpression(lhs Node, operator Operator, rhsIter Iterator[Token]) (BinaryExpressionNode, error) {

// 	// rhs, err := parseNext(rhsIter)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	return BinaryExpressionNode{
// 		lhs,
// 		operator,
// 		rhs,
// 	}, nil
// }
