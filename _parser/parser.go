package parser

import (
	"errors"
	"fmt"
)

type Node interface{}
type OperatorNode interface {
	GetLeft() Node
	GetRight() Node
}

// Type
const (
	Type_String  string = "string"
	Type_Integer string = "int"
)

// An assignment node `a := 1`
type AssignmentNode struct {
	Name  string
	Value string
	Type  string
}

// noot!(value)
type PrintNode struct {
	Value Node
}

type LiteralNode struct {
	Value string
	Type  string
}

// e.g. variable
type IdentifierNode struct {
	Value string
}

type AdditionNode struct {
	Left  Node
	Right Node
}

func (n AdditionNode) GetLeft() Node  { return n.Left }
func (n AdditionNode) GetRight() Node { return n.Right }

type SubtractNode struct {
	Left  Node
	Right Node
}

func (n SubtractNode) GetLeft() Node  { return n.Left }
func (n SubtractNode) GetRight() Node { return n.Right }

type MultiplyNode struct {
	Left  Node
	Right Node
}

func (n MultiplyNode) GetLeft() Node  { return n.Left }
func (n MultiplyNode) GetRight() Node { return n.Right }

type DivideNode struct {
	Left  Node
	Right Node
}

func (n DivideNode) GetLeft() Node  { return n.Left }
func (n DivideNode) GetRight() Node { return n.Right }

// Tokens to nodes
func Parse(tokens []Token) ([]Node, error) {
	nodes := []Node{}

	iter := newIterator(tokens)

	for iter.hasNext() {
		node, err := parseNext(&iter)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// Parse the next statement
func parseNext(iter *Iterator[Token]) (Node, error) {
	stmt := getStmt(iter)
	return parseStatement(stmt, iter)
}

// Read the tokens until the next EOS
func getStmt(iter *Iterator[Token]) []*Token {
	var stmt []*Token
	for func() string {
		if token := iter.peek(); token != nil {
			return token.Type
		} else {
			return EOS
		}
	}() != EOS {
		token, _ := iter.next()
		stmt = append(stmt, token)
	}
	_, _ = iter.next() // EOS

	return stmt
}

// Usually separated by ; or \n, but ome need extra parsing, like for loops (with {})
func parseStatement(stmt []*Token, iter *Iterator[Token]) (Node, error) {
	switch stmt[0].Type {
	case Ident:
		return parseIdent(stmt, iter)
	case Print:
		return parsePrint(stmt, iter)
	case Integer:
		return parseLiteral(stmt, iter)
	}

	return nil, errors.New("ERROR: Token type not found")
}

// Parse the case where a statement starts with and identifier
func parseIdent(stmt []*Token, iter *Iterator[Token]) (Node, error) {
	ident := stmt[0]
	if len(stmt) == 1 {
		return IdentifierNode{ident.Value}, nil
	}
	switch stmt[1].Type {
	case Assign: // :=
		if len(stmt) == 3 {
			val := stmt[2]
			return AssignmentNode{ident.Value, val.Value, val.Type}, nil
		} else {
			for _, p := range stmt {
				fmt.Println(p)
			}
			return nil, errors.New(fmt.Sprintf("ERROR: Statements are not yet supported in assignments (parsing %#v)\n", stmt))
			// TODO: parseStatement(stmt[2:])
		}
	case Plus, Minus, Slash, Star:
		node, err := parseStatement(stmt[2:], iter)
		if err != nil {
			return nil, err
		}

		switch stmt[1].Type {
		case Plus:
			return AdditionNode{
				IdentifierNode{ident.Value},
				node,
			}, nil
		case Minus:
			return SubtractNode{
				IdentifierNode{ident.Value},
				node,
			}, nil
		case Slash:
			return DivideNode{
				IdentifierNode{ident.Value},
				node,
			}, nil
		case Star:
			return MultiplyNode{
				IdentifierNode{ident.Value},
				node,
			}, nil
		}
	}

	return nil, errors.New("ERROR: Unknown syntax after ident")
}

func parsePrint(stmt []*Token, iter *Iterator[Token]) (Node, error) {
	if stmt[1].Type != OpenPar {
		return nil, errors.New("ERROR: Expected open parenthesis after noot!")
	}
	if stmt[len(stmt)-1].Type != ClosedPar {
		return nil, errors.New("ERROR: Expected closed parenthesis following open parenthesis of noot!")
	}

	innerStmt := stmt[2 : len(stmt)-1]
	innerNode, err := parseStatement(innerStmt, iter)
	if err != nil {
		return nil, err
	}

	return PrintNode{innerNode}, nil
}

func parseLiteral(stmt []*Token, iter *Iterator[Token]) (Node, error) {
	if len(stmt) != 1 {
		if stmt[1].Type == Plus {
			node, err := parseStatement(stmt[2:], iter)
			if err != nil {
				return nil, err
			}

			return AdditionNode{
				Left:  LiteralNode{stmt[0].Value, stmt[0].Type},
				Right: node,
			}, nil
		}
	}

	return LiteralNode{stmt[0].Value, stmt[0].Type}, nil
}
