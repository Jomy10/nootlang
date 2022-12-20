package interpreter

import (
	"errors"
	"fmt"
	"github.com/jomy10/nootlang/parser"
	"io"
	"strconv"
	"strings"
)

type Runtime struct {
	// Variable names => values
	vars map[string]interface{}
}

func NR() Runtime {
	return newRuntime()
}

func newRuntime() Runtime {
	return Runtime{
		vars: make(map[string]interface{}),
	}
}

func Interpret(nodes []parser.Node, stdout, stderr io.Writer, stdin io.Reader) error {
	runtime := newRuntime()

	for _, node := range nodes {
		_, err := ExecNode(&runtime, node, stdout, stderr, stdin)
		if err != nil {
			fmt.Printf("GOT ERROR %v\n", err)
			return err
		}
	}

	return nil // program executed without errors
}

// (return 1) Returns the value returned by the expression, or nil of nothing returned
func ExecNode(runtime *Runtime, node parser.Node, stdout, stderr io.Writer, stdin io.Reader) (interface{}, error) {
	// fmt.Printf("Node: %#v\n", node)
	switch node.(type) {
	case parser.VarDeclNode:
		return nil, execVarDecl(runtime, node.(parser.VarDeclNode), stdout, stderr, stdin)
	case parser.VarAssignNode:
		return nil, execVarAssign(runtime, node.(parser.VarAssignNode), stdout, stderr, stdin)
	case parser.PrintStmtNode:
		return execPrintStmt(runtime, node.(parser.PrintStmtNode), stdout, stderr, stdin)
	case parser.IntegerLiteralNode:
		return node.(parser.IntegerLiteralNode).Value, nil
	case parser.VariableNode:
		return getVariable(runtime, node.(parser.VariableNode))
	case parser.BinaryExpressionNode:
		return execBinaryExpressionNode(runtime, node.(parser.BinaryExpressionNode), stdout, stderr, stdin)
	}
	return nil, errors.New(fmt.Sprintf("Noot error: Invalid node `%#v`", node))
}

func execVarDecl(runtime *Runtime, node parser.VarDeclNode, stdout, stderr io.Writer, stdin io.Reader) error {
	rhs, err := ExecNode(runtime, node.Rhs, stdout, stderr, stdin)
	if err != nil {
		return err
	}
	runtime.vars[node.VarName] = rhs
	return nil
}

func execVarAssign(runtime *Runtime, node parser.VarAssignNode, stdout, stderr io.Writer, stdin io.Reader) error {
	rhs, err := ExecNode(runtime, node.Rhs, stdout, stderr, stdin)
	if err != nil {
		return err
	}
	runtime.vars[node.VarName] = rhs
	return nil
}

func execPrintStmt(runtime *Runtime, node parser.PrintStmtNode, stdout, stderr io.Writer, stdin io.Reader) (interface{}, error) {
	inner, err := ExecNode(runtime, node.Inner, stdout, stderr, stdin)
	if err != nil {
		return nil, err
	}
	str := fmt.Sprintf("%v\n", inner)
	stdout.Write([]byte(str))
	return str, nil
}

func getVariable(runtime *Runtime, node parser.VariableNode) (interface{}, error) {
	val, ok := runtime.vars[node.Name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Variable %s is not declared", node.Name))
	}
	return val, nil
}

func execBinaryExpressionNode(runtime *Runtime, node parser.BinaryExpressionNode, stdout, stderr io.Writer, stdin io.Reader) (interface{}, error) {
	left, err := ExecNode(runtime, node.Left, stdout, stderr, stdin)
	if err != nil {
		return nil, err
	}
	right, err := ExecNode(runtime, node.Right, stdout, stderr, stdin)
	if err != nil {
		return nil, err
	}
	return binaryExpressionResult(left, right, node.Operator)
}

func binaryExpressionResult[A any, B any](lhs A, rhs B, op parser.Operator) (A, error) {
	switch any(lhs).(type) {
	case int64:
		rhsInt, ok := any(rhs).(int64)
		if !ok {
			rhsStr := fmt.Sprintf("%v", rhs)
			var err error
			rhsInt, err = strconv.ParseInt(rhsStr, 10, 64)
			if err != nil {
				return toA[A](nil), err
			}
		}

		switch op {
		case parser.Op_Plus:
			return toA[A](any(lhs).(int64) + rhsInt), nil
		case parser.Op_Min:
			return toA[A](any(lhs).(int64) - rhsInt), nil
		case parser.Op_Mul:
			return toA[A](any(lhs).(int64) * rhsInt), nil
		case parser.Op_Div:
			return toA[A](any(lhs).(int64) / rhsInt), nil
		}
	case string:
		rhsStr := fmt.Sprintf("%v", rhs)
		switch op {
		case parser.Op_Plus:
			var sb strings.Builder
			sb.WriteString(any(lhs).(string))
			sb.WriteString(rhsStr)
			return toA[A](sb.String()), nil
		case parser.Op_Min:
			return toA[A](nil), errors.New("`-` cannot be applied to strings")
		case parser.Op_Mul:
			return toA[A](nil), errors.New("`-` cannot be applied to strings")
		case parser.Op_Div:
			return toA[A](nil), errors.New("`-` cannot be applied to strings")
		}
	}

	return toA[A](nil), errors.New(fmt.Sprintf("Cannot apply `%v` to the given operands", op))
}

func toA[A any](a any) A {
	return a.(A)
}
