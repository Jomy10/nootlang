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
	case parser.BoolLiteralNode:
		return node.(parser.BoolLiteralNode).Value, nil
	case parser.VariableNode:
		return runtime.GetVar(node.(parser.VariableNode).Name)
	case parser.BinaryExpressionNode:
		return execBinaryExpressionNode(runtime, node.(parser.BinaryExpressionNode))
	case parser.FunctionDeclNode:
		return newFunction(runtime, node.(parser.FunctionDeclNode))
	case parser.ReturnNode:
		return ExecNode(runtime, node.(parser.ReturnNode).Expr)
	case parser.BinaryNotNode:
		return execBinaryNotExpressionNode(runtime, node.(parser.BinaryNotNode))
	}
	return nil, errors.New(fmt.Sprintf("Noot error: Invalid node `%#v`", node))
}

func execBinaryNotExpressionNode(runtime *runtime.Runtime, node parser.BinaryNotNode) (interface{}, error) {
	val, err := ExecNode(runtime, node.Expr)
	if err != nil {
		return nil, err
	}
	switch val.(type) {
	case bool:
		return !(val.(bool)), nil
	default:
		return nil, errors.New(fmt.Sprintf("Cannot apply `!` to %v\n", val))
	}
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

func binaryExpressionResult(lhs interface{}, rhs interface{}, op parser.Operator) (interface{}, error) {
	switch lhs.(type) {
	case int64:
		switch rhs.(type) {
		case int64:
			return binaryOp(lhs.(int64), rhs.(int64), op)
		case float64:
			return binaryOp(float64(lhs.(int64)), rhs.(float64), op)
		case string:
			rhsInt, err := strconv.ParseInt(rhs.(string), 10, 64)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Cannot convert %s to integer in right hand side of binary expression (%v)", rhs, err))
			}
			return binaryOp(lhs.(int64), rhsInt, op)
		default:
			return nil, errors.New(fmt.Sprintf("Cannot apply binary operator to %v\n", rhs))
		}
	case float64:
		switch rhs.(type) {
		case int64:
			return binaryOp(lhs.(float64), float64(rhs.(int64)), op)
		case float64:
			return binaryOp(lhs.(float64), rhs.(float64), op)
		case string:
			rhsFloat, err := strconv.ParseFloat(rhs.(string), 64)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Cannot convert %s to floating point in right hand side of binary expression (%v)", rhs, err))
			}
			return binaryOp(lhs.(float64), rhsFloat, op)
		default:
			return nil, errors.New(fmt.Sprintf("Cannot apply binary operator to %v\n", rhs))
		}
	default:
		return nil, errors.New(fmt.Sprintf("Cannot apply binary operator to %v", rhs))
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
	case bool:
		switch rhs.(type) {
		case int64:
			return binaryOpBool(lhs.(bool), rhs.(int64) == 0, op)
		case float64:
			return binaryOpBool(lhs.(bool), rhs.(float64) == 0, op)
		case string:
			rhsBool, err := strconv.ParseBool(rhs.(string))
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Cannot convert %s to boolean in right hand side of boolean expression (%v)", rhs, err))
			}
			return binaryOpBool(lhs.(bool), rhsBool, op)
		case bool:
			return binaryOpBool(lhs.(bool), rhs.(bool), op)
		}
	}

	return nil, errors.New(fmt.Sprintf("Cannot apply `%v` to the given operands", op))
}

// Returns the result of a binary operation on an integer or float
func binaryOp[T int64 | float64](lhs T, rhs T, op parser.Operator) (interface{}, error) {
	// fmt.Printf("%v %s %v\n", lhs, op, rhs)
	switch op {
	case parser.Op_Plus:
		return lhs + rhs, nil
	case parser.Op_Min:
		return lhs - rhs, nil
	case parser.Op_Mul:
		return lhs * rhs, nil
	case parser.Op_Div:
		return lhs / rhs, nil
	case parser.Op_CompEqual:
		return lhs == rhs, nil
	case parser.Op_CompNEqual:
		return lhs != rhs, nil
	case parser.Op_LT:
		return lhs < rhs, nil
	case parser.Op_GT:
		return lhs > rhs, nil
	case parser.Op_LTE:
		return lhs <= rhs, nil
	case parser.Op_GTE:
		return lhs >= rhs, nil
	case parser.Op_Or:
		return nil, errors.New("Cannot apply `||` on integers and floats")
	case parser.Op_And:
		return nil, errors.New("Cannot apply `&&` on integers and floats")
	default:
		return 0, errors.New("Interpreter bug (unreachable)")
	}
}

func binaryOpBool(lhs bool, rhs bool, op parser.Operator) (bool, error) {
	switch op {
	case parser.Op_Plus, parser.Op_Min, parser.Op_Mul, parser.Op_Div,
		parser.Op_LT, parser.Op_GT, parser.Op_LTE, parser.Op_GTE:
		return false, errors.New(fmt.Sprintf("Cannot apply `%s` to boolean", op))
	case parser.Op_CompEqual:
		return lhs == rhs, nil
	case parser.Op_CompNEqual:
		return lhs != rhs, nil
	case parser.Op_Or:
		return lhs || rhs, nil
	case parser.Op_And:
		return lhs && rhs, nil
	default:
		return false, errors.New("Interpreter bug (unreachable; unhandled operator in boolean binary operator)")
	}
}
