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

	// String methods
	string_type := reflect.TypeOf("")
	r.Methods[string_type] = make(map[string]runtime.NativeFunction)

	r.Methods[string_type]["concat"] = string__concat
	r.Methods[string_type]["split"] = string__split
	r.Methods[string_type]["len"] = string__len

	// Array methods
	array_type := reflect.TypeOf([]interface{}{})
	r.Methods[array_type] = make(map[string]runtime.NativeFunction)
	r.Methods[array_type]["len"] = array__len
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
			str += " "
		default:
			str += fmt.Sprintf("%v", arg)
			str += " "
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

	// returns a string
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

	// TODO: return as []interface{} !!
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

	return int64(len(lhs)), nil
}

func array__len(runtime *runtime.Runtime, args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("`array.len` expects no arguments")
	}

	lhs, ok := args[0].([]interface{})
	if !ok {
		return nil, errors.New("intepreter error in `arrray.len`")
	}

	return int64(len(lhs)), nil
}
