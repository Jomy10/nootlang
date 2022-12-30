package corelib

import (
	"errors"
	"fmt"
	"github.com/jomy10/nootlang/runtime"
	"reflect"
)

// Register (import) the core library in a runtime
func Register(r *runtime.Runtime) {
	r.Funcs["GLOBAL"]["noot!"] = nootLine

	string_type := reflect.TypeOf(reflect.String.String())
	r.Methods[string_type] = make(map[string]runtime.NativeFunction)
	r.Methods[string_type]["concat"] = string__concat
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
