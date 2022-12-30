package interpreter

import (
	"errors"
	"fmt"
	"github.com/jomy10/nootlang/corelib"
	"github.com/jomy10/nootlang/parser"
	runtime "github.com/jomy10/nootlang/runtime"
	"io"
	"reflect"
	"strconv"
	"strings"
)

func Interpret(nodes []parser.Node, stdout, stderr io.Writer, stdin io.Reader) error {
	runtime := runtime.NewRuntime(stdout, stderr, stdin)
	corelib.Register(&runtime)

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
		return execFuncCallNode(runtime, node.(parser.FunctionCallExprNode))
	case parser.MethodCallExprNode:
		return execMethodCallNode(runtime, node.(parser.MethodCallExprNode))
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
	case parser.ArrayLiteralNode:
		return execArrayLiteral(runtime, node.(parser.ArrayLiteralNode))
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
	case parser.ArrayIndexNode:
		return execArrayIndexNode(runtime, node.(parser.ArrayIndexNode))
	case parser.ArrayIndexAssignmentNode:
		return nil, execArrayIndexAssignmentNode(runtime, node.(parser.ArrayIndexAssignmentNode))
	case parser.IfNode:
		return nil, execIf(runtime, node.(parser.IfNode))
	case parser.ElseNode:
		return nil, execElse(runtime, node.(parser.ElseNode))
	case parser.WhileNode:
		return nil, execWhile(runtime, node.(parser.WhileNode))
	}
	return nil, errors.New(fmt.Sprintf("Noot error: Invalid node `%#v`", node))
}

func execArrayIndexAssignmentNode(runtime *runtime.Runtime, node parser.ArrayIndexAssignmentNode) error {
	idx, err := ExecNode(runtime, node.Index)
	if err != nil {
		return err
	}
	val, err := ExecNode(runtime, node.Rhs)
	if err != nil {
		return err
	}

	switch idx.(type) {
	case int64:
		// ok
		runtime.SetArrayIndex(node.Array.Name, idx.(int64), val)
		return nil
	default:
		return errors.New("Only integers can be used for array indexing")
	}
}

// Return the value
func execArrayIndexNode(runtime *runtime.Runtime, node parser.ArrayIndexNode) (interface{}, error) {
	array, err := ExecNode(runtime, node.Array)
	if err != nil {
		return nil, err
	}
	switch array.(type) {
	case []interface{}:
		idx, err := ExecNode(runtime, node.Index)
		if err != nil {
			return nil, err
		}
		switch idx.(type) {
		case int64:
			return array.([]interface{})[idx.(int64)], nil
		default:
			return nil, errors.New("Only integer values can be used to index an aray")
		}
	default:
		return nil, errors.New("Cannot index non-array type")
	}
}

func execArrayLiteral(runtime *runtime.Runtime, node parser.ArrayLiteralNode) (interface{}, error) {
	arr := make([]interface{}, len(node.Values))
	for i, element := range node.Values {
		arrElemVal, err := ExecNode(runtime, element)
		if err != nil {
			return nil, err
		}
		arr[i] = arrElemVal
	}
	return arr, nil
}

func execWhile(runtime *runtime.Runtime, node parser.WhileNode) error {
whileLoop:
	for {
		condVal, err := ExecNode(runtime, node.Condition)
		if err != nil {
			return err
		}

		switch condVal.(type) {
		case bool:
			if !(condVal.(bool)) {
				break whileLoop
			}
			for _, node := range node.Body {
				_, err := ExecNode(runtime, node)
				if err != nil {
					return err
				}
			}
		default:
			return errors.New("Condition is not a boolean expression in while loop")
		}

	}

	return nil
}

func execIf(runtime *runtime.Runtime, node parser.IfNode) error {
	val, err := ExecNode(runtime, node.Condition)
	if err != nil {
		return err
	}
	switch val.(type) {
	case bool:
		if val.(bool) {
			for _, bodynode := range node.Body {
				_, err := ExecNode(runtime, bodynode)
				if err != nil {
					return err
				}
			}
			return nil
		} else {
			if node.NextElseBlock != nil {
				_, err := ExecNode(runtime, node.NextElseBlock)
				return err
			} else {
				return nil
			}
		}
	default:
		return errors.New(fmt.Sprintf("%v is not a boolean value", val))
	}
}

func execElse(runtime *runtime.Runtime, node parser.ElseNode) error {
	for _, node := range node.Body {
		_, err := ExecNode(runtime, node)
		if err != nil {
			return err
		}
	}
	return nil
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
				runtime.SetVar(runtime.CurrentScope(), node.ArgumentNames[i], args[i])
			} else {
				// runtime.Vars[node.ArgumentNames[i]] = nil
				runtime.SetVar(runtime.CurrentScope(), node.ArgumentNames[i], nil)
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

func execFuncCallNode(_runtime *runtime.Runtime, node parser.FunctionCallExprNode) (interface{}, error) {
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

	return execFuncCall(_runtime, function, node.Arguments, nil)
}

// In the method call, the value on the left of the method call will be the first
// element in the argument list passed to the native function
func execMethodCallNode(runtime *runtime.Runtime, node parser.MethodCallExprNode) (interface{}, error) {
	calledOnValue, err := ExecNode(runtime, node.CalledOn)
	if err != nil {
		return nil, err
	}
	method := runtime.GetMethod(calledOnValue, node.FunctionCall.FuncName)
	if method == nil {
		return nil, errors.New(fmt.Sprintf("Method %s does not exist on %v\n", node.FunctionCall.FuncName, reflect.TypeOf(calledOnValue)))
	}
	return execFuncCall(runtime, method, node.FunctionCall.Arguments, calledOnValue)
}

// - firstArg: Optional parameter for prepending an argument to the argument list
//	 passed to the function (used in method call).
func execFuncCall(runtime *runtime.Runtime, fn runtime.NativeFunction, callArgs []parser.Node, firstArg interface{}) (interface{}, error) {
	args := []interface{}{}
	if firstArg != nil {
		args = append(args, firstArg)
	}
	for _, argNode := range callArgs {
		val, err := ExecNode(runtime, argNode)
		if err != nil {
			return nil, err
		}
		args = append(args, val)
	}

	return fn(runtime, args)
}

func execVarDecl(runtime *runtime.Runtime, node parser.VarDeclNode) error {
	// if _, exists := runtime.Vars[node.VarName]; exists {
	exists, _ := runtime.VarExists(node.VarName)
	if exists {
		return errors.New(fmt.Sprintf("Variable `%s` is already defined", node.VarName))
	}

	rhs, err := ExecNode(runtime, node.Rhs)
	if err != nil {
		return err
	}
	// runtime.Vars[node.VarName] = rhs
	runtime.SetVar(runtime.CurrentScope(), node.VarName, rhs)
	return nil
}

func execVarAssign(runtime *runtime.Runtime, node parser.VarAssignNode) error {
	// if _, exists := runtime.Vars[node.VarName]; !exists {
	exists, scope := runtime.VarExists(node.VarName)
	if !exists {
		return errors.New(fmt.Sprintf("Variable `%s` is not defined", node.VarName))
	}
	rhs, err := ExecNode(runtime, node.Rhs)
	if err != nil {
		return err
	}

	switch node.Op {
	case parser.Op_Equal:
		runtime.SetVar(scope, node.VarName, rhs)
	case parser.Op_PlusEqual:
		if err := runtime.ApplyToVariable(scope, node.VarName, func(varval interface{}) (interface{}, error) {
			switch varval.(type) {
			case []interface{}:
				return append(varval.([]interface{}), rhs), nil
			default:
				return binaryExpressionResult(varval, rhs, parser.Operator("+"))
			}
		}); err != nil {
			return nil
		}
	case parser.Op_MinEqual:
		if err := runtime.ApplyToVariable(scope, node.VarName, func(varval interface{}) (interface{}, error) {
			return binaryExpressionResult(varval, rhs, parser.Operator("-"))
		}); err != nil {
			return nil
		}
	case parser.Op_TimesEqual:
		if err := runtime.ApplyToVariable(scope, node.VarName, func(varval interface{}) (interface{}, error) {
			return binaryExpressionResult(varval, rhs, parser.Operator("*"))
		}); err != nil {
			return nil
		}
	case parser.Op_DivEqual:
		if err := runtime.ApplyToVariable(scope, node.VarName, func(varval interface{}) (interface{}, error) {
			return binaryExpressionResult(varval, rhs, parser.Operator("/"))
		}); err != nil {
			return nil
		}
	default:
		return errors.New(fmt.Sprintf("Invalid operator %v (interpreter bug)", node.Op))
	}

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
