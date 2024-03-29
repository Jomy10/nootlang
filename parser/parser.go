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
	arrayLevel := 0 // level of square brackets
	start := 0
	i := 0
	for true {
		if i != len(tokens) {
			if tokens[i].Type == OpenCurlPar {
				blockLevel += 1
			} else if tokens[i].Type == ClosedCurlPar {
				blockLevel -= 1
			} else if tokens[i].Type == OpenSquarePar {
				arrayLevel += 1
			} else if tokens[i].Type == ClosedSquarePar {
				arrayLevel -= 1
			}
		}

		if i == len(tokens) || (tokens[i].Type == EOS && blockLevel == 0 && arrayLevel == 0) {
			currentStatement = tokens[start:i]
			iter := newArrayIterator(currentStatement)
			stmtNode, err := parseStatement(&iter)
			if err != nil {
				if err.Error() != "Empty statement" {
					return nil, err
				}
			} else if stmtNode != nil { // nil check to exclude comments
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
		case Declare, Equal, PlusEqual, MinEqual, StarEqual, SlashEqual:
			_, _ = tokenIter.next() // consume :=/=
			exprNode, err := parseExpression(tokenIter)
			if err != nil {
				return nil, err
			}

			if secondToken.Type == Declare {
				return VarDeclNode{firstToken.Value, exprNode}, nil
			} else {
				return VarAssignNode{firstToken.Value, Operator(secondToken.Value), exprNode}, nil
			}
		case OpenSquarePar:
			idxNode, err := parseArrayIndex(tokenIter)
			if err != nil {
				return nil, err
			}

			nextToken, hasNext := tokenIter.next()
			if !hasNext {
				return nil, errors.New("Cannot use array index expression as statement")
			}

			if nextToken.Type == Equal {
				rhs, err := parseExpression(tokenIter)
				if err != nil {
					return nil, err
				}
				return ArrayIndexAssignmentNode{
					VariableNode{"a"},
					idxNode,
					rhs,
				}, nil
			} else {
				return nil, errors.New(fmt.Sprintf("Unexpected token %s", nextToken.Value))
			}
		case Dot:
			tokenIter.consume(1) // consume dot
			return parseMethodCall(VariableNode{firstToken.Value}, tokenIter)
		default:
			return nil, errors.New(fmt.Sprintf("%#v is invalid at current position", secondToken))
		}
	case String, Integer, Float, Bool:
		// Literals follew by a dot are valid in statements
		secondToken, hasSecond := tokenIter.next()
		if !hasSecond {
			return nil, errors.New("Literals cannot be used as statements")
		}
		if secondToken.Type == Dot {
			firstTokenIter := newArrayOfPointerIterator([]*Token{firstToken})
			lhsexpr, err := parseExpression(&firstTokenIter)
			if err != nil {
				return nil, err
			}
			return parseMethodCall(lhsexpr, tokenIter)
		} else {
			return nil, errors.New(fmt.Sprintf("Invalid token %v", secondToken))
		}
	case Return:
		expr, err := parseExpression(tokenIter)
		if err != nil {
			return nil, err
		}
		return ReturnNode{expr}, nil
	case Def:
		return parseFunctionDecl(tokenIter)
	case If:
		tokenIter.reverse(1)
		return parseIf(tokenIter)
	case While:
		return parseWhile(tokenIter)
	case Comment:
		return nil, nil // Currently ignored
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
		secondToken, hasSecond := tokenIter.peekN(2)
		// handle lonesome boolean literal
		if !hasSecond {
			tokenIter.consume(1) // consume boolean
			boolean, err := strconv.ParseBool(firstToken.Value)
			if err != nil {
				return nil, err
			}
			return BoolLiteralNode{boolean}, nil
		} else {
			if secondToken.Type == Dot {
				tokenIter.consume(2) // (bool).
				boolean, err := strconv.ParseBool(firstToken.Value)
				if err != nil {
					return nil, err
				}
				return parseMethodCall(BoolLiteralNode{boolean}, tokenIter)
			} else {
				return parseBinaryExpression(tokenIter)
			}
		}
	case Integer:
		secondToken, hasSecond := tokenIter.peekN(2)
		// Handle lonesome integer literal
		if !hasSecond {
			tokenIter.consume(1) // consume integer
			integer, err := strconv.ParseInt(firstToken.Value, 10, 64)
			if err != nil {
				return nil, err
			}
			return IntegerLiteralNode{integer}, nil
		} else {
			if secondToken.Type == Dot {
				tokenIter.consume(2) // consume (int).
				integer, err := strconv.ParseInt(firstToken.Value, 10, 64)
				if err != nil {
					return nil, err
				}
				return parseMethodCall(IntegerLiteralNode{integer}, tokenIter)
			} else {
				return parseBinaryExpression(tokenIter)
			}
		}
	case Float:
		secondToken, hasSecond := tokenIter.peekN(2)
		// Handle lonesome float literal
		if !hasSecond {
			tokenIter.consume(1) // consume float
			float, err := strconv.ParseFloat(firstToken.Value, 64)
			if err != nil {
				return nil, err
			}
			return FloatLiteralNode{float}, nil
		} else {
			if secondToken.Type == Dot {
				tokenIter.consume(2) // consume (int).
				float, err := strconv.ParseFloat(firstToken.Value, 64)
				if err != nil {
					return nil, err
				}
				return parseMethodCall(FloatLiteralNode{float}, tokenIter)
			} else {
				return parseBinaryExpression(tokenIter)
			}
		}
	case Ident:
		secondToken, hasSecond := tokenIter.peekN(2)
		// Handle lonesome ident
		if !hasSecond {
			tokenIter.consume(1) // consume ident
			return VariableNode{firstToken.Value}, nil
		}

		switch secondToken.Type {
		case OpenPar:
			// function call
			tokenIter.consume(1) // consume ident
			return parseFunctionCall(firstToken.Value, tokenIter)
		case OpenSquarePar:
			tokenIter.consume(1) // consume ident
			innerIndexExpression, err := parseArrayIndex(tokenIter)
			if err != nil {
				return nil, err
			}
			return ArrayIndexNode{
				VariableNode{firstToken.Value},
				innerIndexExpression,
			}, nil
		case Dot:
			tokenIter.consume(2) // consume ident and dot
			return parseMethodCall(VariableNode{firstToken.Value}, tokenIter)
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
			if secondToken.Type == Dot {
				tokenIter.consume(2) // consume (str).
				str := parseStringLiteral(firstToken.Value)
				return parseMethodCall(str, tokenIter)
			} else {
				return nil, errors.New(fmt.Sprintf("Did not expect token %s afte string literal\n", secondToken.Value))
			}
		}
	case Not:
		tokenIter.consume(1)
		node, err := parseExpression(tokenIter)
		if err != nil {
			return nil, err
		}
		return BinaryNotNode{node}, nil
	case OpenSquarePar: // start of array initialization
		return parseArrayLiteral(tokenIter)
	default:
		return nil, errors.New(fmt.Sprintf("Invalid start of expression `%v`", firstToken))
	}
}

// tokenIter is at the method name (ident)
func parseMethodCall(calledOn Node, tokenIter Iterator[Token]) (Node, error) {
	funcNameToken, hasNext := tokenIter.next()
	if !hasNext || funcNameToken.Type != Ident {
		return nil, errors.New("Expected function name after `.`")
	}
	funcCallNode, err := parseFunctionCall(funcNameToken.Value, tokenIter)
	if err != nil {
		return nil, err
	}
	return MethodCallExprNode{
		calledOn,
		funcCallNode,
	}, nil
}

// tokenIter is at [
func parseArrayIndex(tokenIter Iterator[Token]) (Node, error) {
	tokenIter.consume(1) // [
	arrayLevel := 1

	var expression []*Token
	for {
		nextToken, hasNext := tokenIter.next()
		if !hasNext {
			return nil, errors.New("Expected ] to close array index")
		}

		if nextToken.Type == ClosedSquarePar {
			arrayLevel -= 1
		} else if nextToken.Type == OpenSquarePar {
			arrayLevel += 1
		}

		if arrayLevel == 0 {
			exprIter := newArrayOfPointerIterator(expression)
			return parseExpression(&exprIter)
		} else {
			expression = append(expression, nextToken)
		}
	}
}

// tokenIter starts at [
func parseArrayLiteral(tokenIter Iterator[Token]) (Node, error) {
	tokenIter.consume(1) // consume [

	list, err := collectList(tokenIter, ClosedSquarePar)
	if err != nil {
		return nil, err
	}

	expressions := make([]Node, len(list))
	for i, exprTokens := range list {
		exprTokenIter := newArrayOfPointerIterator(exprTokens)
		expr, err := parseExpression(&exprTokenIter)
		if err != nil {
			return nil, err
		}
		expressions[i] = expr
	}

	return ArrayLiteralNode{
		expressions,
	}, nil
}

// MethodCallExprNodeenIter starts at the condition of the while loop
func parseWhile(tokenIter Iterator[Token]) (Node, error) {
	// Collect the while loop's condition
	var condition []*Token
	for {
		nextToken, hasNext := tokenIter.peek()

		if !hasNext {
			return nil, errors.New("Expected opening curly bracket after while condition")
		}

		if nextToken.Type == OpenCurlPar {
			break
		}

		tokenIter.consume(1) // consume if not {
		condition = append(condition, nextToken)
	}

	if len(condition) == 0 {
		return nil, errors.New("While loop has empty condition")
	}

	conditionIter := newArrayOfPointerIterator(condition)
	expr, err := parseExpression(&conditionIter)
	if err != nil {
		return nil, err
	}

	body, err := parseBody(tokenIter)
	if err != nil {
		return nil, err
	}

	return WhileNode{
		expr,
		body,
	}, nil
}

// Parse if or elsif (starting at if or elsi)
func parseIf(tokenIter Iterator[Token]) (Node, error) {
	ifToken, _ := tokenIter.next() // if / elsif
	var condition []*Token
	if ifToken.Type == If || ifToken.Type == Elsif {
		nextToken, hasNext := tokenIter.peek()
		for {
			if !hasNext {
				return nil, errors.New("If exepcted an opening curly bracket")
			}
			if nextToken.Type == OpenCurlPar {
				break
			}

			condition = append(condition, nextToken)
			// nextToken, hasNext = tokenIter.next()
			tokenIter.consume(1)
			nextToken, hasNext = tokenIter.peek()
		}
	}

	var conditionExpr Node
	if ifToken.Type == If || ifToken.Type == Elsif {
		conditionIter := newArrayOfPointerIterator(condition)
		var err error
		conditionExpr, err = parseExpression(&conditionIter)
		if err != nil {
			return nil, err
		}
	}

	bodyExpr, err := parseBody(tokenIter)
	if err != nil {
		return nil, err
	}

	nextToken, hasNext := tokenIter.peek()
	var elseBlock Node
	if hasNext && (nextToken.Type == Else || nextToken.Type == Elsif) {
		ifNode, err := parseIf(tokenIter)
		if err != nil {
			return nil, err
		}
		elseBlock = ifNode
	} else {
		elseBlock = nil
	}

	if ifToken.Type == Else {
		return ElseNode{bodyExpr}, nil
	} else {
		return IfNode{
			conditionExpr,
			elseBlock,
			bodyExpr,
		}, nil
	}
}

// tokenIter is at the opening bracket of the function call
func parseFunctionCall(name string, tokenIter Iterator[Token]) (FunctionCallExprNode, error) {
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
			return nil, err
		}
		args = append(args, expr)
	}

	return args, nil
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
				tokenIter.consume(i - 1)
				return Parse(subIter.collect())
			}
		}

		nextToken, hasNext = tokenIter.peekN(i)
		i += 1
	}

	return nil, errors.New("Parser error")
}

// Collect a list of arguments between brackets
func collectList(tokenIter Iterator[Token], closingToken TT) ([][]*Token, error) {
	parLevel := 0
	blockLevel := 0 // {}
	arrayLevel := 0 // []

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

		if nextToken.Type == closingToken && parLevel == 0 && blockLevel == 0 && arrayLevel == 0 {
			return tokenArgs, nil
		} else if nextToken.Type == OpenPar {
			parLevel += 1
		} else if nextToken.Type == ClosedPar {
			parLevel -= 1
		} else if nextToken.Type == OpenCurlPar {
			blockLevel += 1
		} else if nextToken.Type == ClosedCurlPar {
			blockLevel -= 1
		} else if nextToken.Type == OpenSquarePar {
			arrayLevel += 1
		} else if nextToken.Type == ClosedSquarePar {
			arrayLevel -= 1
		} else if nextToken.Type == Comma && parLevel == 0 && blockLevel == 0 && arrayLevel == 0 {
			idx += 1
		}
		if nextToken.Type != Comma || parLevel != 0 || blockLevel != 0 || arrayLevel != 0 {
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
