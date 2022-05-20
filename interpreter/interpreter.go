package interpreter

import (
	"fmt"
	"github.com/jomy10/nootlang/parser"
	"strconv"
)

type Runtime struct {
	vars map[string]interface{}
}

func newRuntime() Runtime {
	return Runtime{
		vars: make(map[string]interface{}),
	}
}

func Interpret(nodes []parser.Node, stdout chan string, stderr chan string, eop chan int) {
	runtime := newRuntime()

	for _, node := range nodes {
		if !execNode(&runtime, node, stdout, stderr) {
			eop <- 1
			return // error occured
		}
	}

	eop <- 0
}

func execNode(runtime *Runtime, _node parser.Node, stdout, stderr chan string) bool {
	switch _node.(type) {
	case parser.AssignmentNode:
		node := _node.(parser.AssignmentNode)

		// Check if variable with name does not exist
		return execAssignment(runtime, node, stderr)
	case parser.PrintNode:
		node := _node.(parser.PrintNode)
		return execPrint(runtime, node, stdout, stderr)
	}

	// Unreachable state
	return false
}

func evalNode(runtime *Runtime, _node parser.Node, stderr chan string) (interface{}, bool) {
	switch _node.(type) {
	case parser.LiteralNode:
		node := _node.(parser.LiteralNode)

		return node.Value, true
	case parser.IdentifierNode:
		node := _node.(parser.IdentifierNode)

		if runtime.vars[node.Value] != nil {
			return runtime.vars[node.Value], true
		} else {
			stderr <- "RUNTIME ERROR: couldn't evaluate identifier node"
			return nil, false
		}
	default:
		stderr <- "RUNTIME ERROR: couldn't evaluate node"
		return nil, false
	}
}

func execAssignment(runtime *Runtime, node parser.AssignmentNode, stderr chan string) (noErr bool) {
	if _, ok := runtime.vars[node.Name]; ok {
		stderr <- fmt.Sprintf("RUNTIME ERROR: variable %s is already declared", node.Name)
		return false
	}

	switch node.Type {
	case parser.Integer:
		i, err := strconv.Atoi(node.Value)
		if err != nil {
			stderr <- "INTERPRETER ERROR: String conversion failed"
			return false
		}
		runtime.vars[node.Name] = i
		return true
	default:
		stderr <- fmt.Sprintf("RUNTIMER ERROR: unknown type %s", node.Type)
		return false
	}
}

func execPrint(runtime *Runtime, node parser.PrintNode, stdout, stderr chan string) bool {
	evaled, ok := evalNode(runtime, node.Value, stderr)
	if !ok {
		return false
	}
	stdout <- fmt.Sprint(evaled)

	return true
}