package interpreter

import (
	"errors"
	"fmt"
	"github.com/jomy10/nootlang/parser"
	runtime "github.com/jomy10/nootlang/runtime"
	"github.com/jomy10/nootlang/stdlib"
	"io"
	"strconv"
	"strings"
)

func Interpret(nodes []parser.Node, stdout, stderr io.Writer, stdin io.Reader) error {
	runtime := runtime.NewRuntime(stdout, stderr, stdin)
	stdlib.Register(&runtime)

	for _, node := range nodes {
		_, err := ExecNode(&runtime, node)
		if err != nil {
			return err
		}
	}

	return nil // program executed without errors
}

// (return 1) Returns the value returned by the expression, or nil of nothing returned
func ExecNode(runtime *runtime.Runtime, node parser.Node) (interface{}, error) {
	// fmt.Printf("Node: %#v\n", node)
	switch node.(type) {
	case parser.VarDeclNode:
		return nil, execVarDecl(runtime, node.(parser.VarDeclNode))
	case parser.VarAssignNode:
		return nil, execVarAssign(runtime, node.(parser.VarAssignNode))
	case parser.FunctionCallExprNode:
		return execFuncCall(runtime, node.(parser.FunctionCallExprNode))
	case parser.IntegerLiteralNode:
		return node.(parser.IntegerLiteralNode).Value, nil
	case parser.NilLiteralNode:
		return nil, nil
	case parser.StringLiteralNode:
		return node.(parser.StringLiteralNode).String, nil
	case parser.FloatLiteralNode:
		return node.(parser.FloatLiteralNode).Value, nil
	case parser.VariableNode:
		return runtime.GetVar(node.(parser.VariableNode).Name)
	case parser.BinaryExpressionNode:
		return execBinaryExpressionNode(runtime, node.(parser.BinaryExpressionNode))
	case parser.FunctionDeclNode:
		return newFunction(runtime, node.(parser.FunctionDeclNode))
	case parser.ReturnNode:
		return ExecNode(runtime, node.(parser.ReturnNode).Expr)
	}
	return nil, errors.New(fmt.Sprintf("Noot error: Invalid node `%#v`", node))
}

func newFunction(_runtime *runtime.Runtime, node parser.FunctionDeclNode) (interface{}, error) {
	_runtime.SetFunc(node.FuncName, func(runtime *runtime.Runtime, args []interface{}) (interface{}, error) {
		// Set scope
		scopeStringBuilder := strings.Builder{}
		scopeStringBuilder.WriteString(runtime.CurrentScope())
		scopeStringBuilder.WriteString("$")
		scopeStringBuilder.WriteString(node.FuncName)
		scope := scopeStringBuilder.String()
		runtime.AddScope(scope)

		// Add variables
		for i := 0; i < len(node.ArgumentNames); i++ {
			if i < len(args) {
				// runtime.Vars[node.ArgumentNames[i]] = args[i]
				runtime.SetVar(node.ArgumentNames[i], args[i])
			} else {
				// runtime.Vars[node.ArgumentNames[i]] = nil
				runtime.SetVar(node.ArgumentNames[i], nil)
			}
		}

		for _, node := range node.Body {
			switch node.(type) {
			case parser.ReturnNode:
				return ExecNode(runtime, node)
			default:
				_, err := ExecNode(runtime, node)
				if err != nil {
					return nil, err
				}
			}
		}

		// Pop scope
		runtime.ExitScope()
		return nil, nil // Function did not return any value
	})

	return nil, nil
}

func execFuncCall(_runtime *runtime.Runtime, node parser.FunctionCallExprNode) (interface{}, error) {
	// function := runtime.Funcs[node.FuncName]
	function := _runtime.GetFunc(node.FuncName)

	if function == nil {
		variable, err := _runtime.GetVar(node.FuncName)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Undeclared function `%s`\n", node.FuncName))
		} else {
			switch variable.(type) {
			case func(*runtime.Runtime, []interface{}) (interface{}, error):
				function = variable.(func(*runtime.Runtime, []interface{}) (interface{}, error))
			default:
				return nil, errors.New(fmt.Sprintf("Undeclared function `%s`\n", node.FuncName))
			}
		}
	}

	args := []interface{}{}
	for _, argNode := range node.Arguments {
		val, err := ExecNode(_runtime, argNode)
		if err != nil {
			return nil, err
		}
		args = append(args, val)
	}

	return function(_runtime, args)
}

func execVarDecl(runtime *runtime.Runtime, node parser.VarDeclNode) error {
	// if _, exists := runtime.Vars[node.VarName]; exists {
	if runtime.VarExists(node.VarName) {
		return errors.New(fmt.Sprintf("Variable `%s` is already defined", node.VarName))
	}

	rhs, err := ExecNode(runtime, node.Rhs)
	if err != nil {
		return err
	}
	// runtime.Vars[node.VarName] = rhs
	runtime.SetVar(node.VarName, rhs)
	return nil
}

func execVarAssign(runtime *runtime.Runtime, node parser.VarAssignNode) error {
	// if _, exists := runtime.Vars[node.VarName]; !exists {
	if !runtime.VarExists(node.VarName) {
		return errors.New(fmt.Sprintf("Variable `%s` is not defined", node.VarName))
	}
	rhs, err := ExecNode(runtime, node.Rhs)
	if err != nil {
		return err
	}
	// runtime.Vars[node.VarName] = rhs
	runtime.SetVar(node.VarName, rhs)
	return nil
}

func execBinaryExpressionNode(runtime *runtime.Runtime, node parser.BinaryExpressionNode) (interface{}, error) {
	left, err := ExecNode(runtime, node.Left)
	if err != nil {
		return nil, err
	}
	right, err := ExecNode(runtime, node.Right)
	if err != nil {
		return nil, err
	}
	return binaryExpressionResult(left, right, node.Operator)
}

func binaryExpressionResult[A any, B any](lhs A, rhs B, op parser.Operator) (interface{}, error) {
	switch any(lhs).(type) {
	case int64:
		rhsInt, ok := any(rhs).(int64)
		if !ok {
			switch any(rhs).(type) {
			case float64:
				rhsInt = int64(any(rhs).(float64))
			default:
				rhsStr := fmt.Sprintf("%v", rhs)
				var err error
				rhsInt, err = strconv.ParseInt(rhsStr, 10, 64)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Cannot convert %s to integer in right hand side of binary expression (%v)", rhsStr, err))
				}
			}
		}

		return binaryOp(any(lhs).(int64), rhsInt, op), nil
	case float64:
		rhsFloat, ok := any(rhs).(float64)
		if !ok {
			switch any(rhs).(type) {
			case int64:
				rhsFloat = float64(any(rhs).(int64))
			default:
				rhsStr := fmt.Sprintf("%v", rhs)
				var err error
				rhsFloat, err = strconv.ParseFloat(rhsStr, 64)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Cannot convert %s to float in right hand side of binary expression (%v)", rhsStr, err))
				}
			}
		}

		return binaryOp(any(lhs).(float64), rhsFloat, op), nil
	case string:
		rhsStr := fmt.Sprintf("%v", rhs)
		switch op {
		case parser.Op_Plus:
			var sb strings.Builder
			sb.WriteString(any(lhs).(string))
			sb.WriteString(rhsStr)
			return sb.String(), nil
		case parser.Op_Min:
			return nil, errors.New("`-` cannot be applied to strings")
		case parser.Op_Mul:
			return nil, errors.New("`-` cannot be applied to strings")
		case parser.Op_Div:
			return nil, errors.New("`-` cannot be applied to strings")
		}
	}

	return nil, errors.New(fmt.Sprintf("Cannot apply `%v` to the given operands", op))
}

// Returns the result of a binary operation on an integer or float
func binaryOp[T int64 | float64](lhs T, rhs T, op parser.Operator) T {
	switch op {
	case parser.Op_Plus:
		return lhs + rhs
	case parser.Op_Min:
		return lhs - rhs
	case parser.Op_Mul:
		return lhs * rhs
	case parser.Op_Div:
		return lhs / rhs
	}

	fmt.Println("Interpreter bug (unreachable)")
	return 0
}
