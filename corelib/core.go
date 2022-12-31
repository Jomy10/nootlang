package corelib

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/jomy10/nootlang/runtime"
)

// Register (import) the core library in a runtime
func Register(r *runtime.Runtime) {
	r.Funcs["GLOBAL"]["noot!"] = nootLine
	r.Funcs["GLOBAL"]["gets?"] = gets_potentially

	// String methods
	string_type := reflect.TypeOf(reflect.String.String())
	r.Methods[string_type] = make(map[string]runtime.NativeFunction)

	r.Methods[string_type]["concat"] = string__concat
	r.Methods[string_type]["split"] = string__split
}

// `noot!`
func nootLine(runtime *runtime.Runtime, args []interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, errors.New("`noot!` expects at least one argument")
	}

	var str string
	for _, arg := range args {
		switch arg.(type) {
		case string:
			str += arg.(string)
		default:
			str += fmt.Sprintf("%v", arg)
		}
	}

	// noot! is like println
	str += "\n"

	runtime.Stdout.Write([]byte(str))
	return str, nil
}

// string.concat
func string__concat(runtime *runtime.Runtime, args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, errors.New("string.concat expects 1 argument")
	}

	lhs, ok := args[0].(string)
	if !ok {
		return nil, errors.New("interpreter error")
	}
	rhs, ok := args[1].(string)
	if !ok {
		rhs = fmt.Sprintf("%v", rhs)
	}

	return lhs + rhs, nil
}

// string.split
func string__split(runtime *runtime.Runtime, args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, errors.New("string.split expects 1 argument")
	}

	lhs, ok := args[0].(string)
	if !ok {
		return nil, errors.New("interpreter error")
	}

	rhs, ok := args[1].(string)
	if !ok {
		rhs = fmt.Sprintf("%v", rhs)
	}

	return strings.Split(lhs, rhs), nil
}

func string__len(runtime *runtime.Runtime, args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("string.len expects no arguments")
	}

	lhs, ok := args[0].(string)
	if !ok {
		return nil, errors.New("interpreter error")
	}

	return len(lhs), nil
}
