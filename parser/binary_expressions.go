package parser

import "errors"

func parseBinaryExpression(tokenIter Iterator[Token]) (Node, error) {
	// Of the form Node Operator Node Operator Node (...)
	// where Operator is of type *Token
	expression := []interface{}{}
	_ = expression

	// parentheses level
	parLevel := 0

	lhs := []*Token{}
	nextToken, hasNext := tokenIter.next()
	for hasNext {
		switch nextToken.Type {
		case OpenPar:
			parLevel += 1
		case ClosedPar:
			parLevel -= 1
		case Star, Slash, Plus, Minus, DEqual, DNEqual, And, Or:
			if parLevel == 0 {
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

	return __parseBinaryExpression(expression)
}

func __parseBinaryExpression(expr []interface{}) (Node, error) {
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

	// precedences := [][]TT{Minus, Plus, Slash, Star, DEqual, DNEqual, LT, GT, LTE, GTE, And, Or} // Operator precedence in reverse order
	precedences := [][]TT{
		{Or},
		{And},
		{DEqual, DNEqual, LT, GT, LTE, GTE},
		{Plus, Minus}, // Also | and ^
		{Star, Slash}, // Also %, <<, >>, & and &^
	}
	precedenceLevelIdx := 0

	exprIdx := 0

	for precedenceLevelIdx != len(precedences) {
		switch expr[exprIdx].(type) {
		case *Token:
			token := expr[exprIdx].(*Token)
			for precedenceIdx := 0; precedenceIdx < len(precedences[precedenceLevelIdx]); precedenceIdx++ {
				if token.Type == precedences[precedenceLevelIdx][precedenceIdx] {
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
		}

		exprIdx += 1

		if exprIdx == len(expr) {
			exprIdx = 0
			precedenceLevelIdx += 1
		}
	}

	return nil, errors.New("Parser bug?")
}
