// Separate runtime package for recursive imports
package runtime

import (
	"errors"
	"fmt"
	"io"
	"reflect"
)

type NativeFunction = func(*Runtime, []interface{}) (interface{}, error)

// TODO: scopes
type Runtime struct {
	Scopes []string
	// Scope names => Variable names => values
	Vars           map[string]map[string]interface{}
	Funcs          map[string]map[string]func(*Runtime, []interface{}) (interface{}, error)
	Methods        map[reflect.Type]map[string]func(*Runtime, []interface{}) (interface{}, error)
	Stdout, Stderr io.Writer
	Stdin          io.Reader
}

func NewRuntime(stdout, stderr io.Writer, stdin io.Reader) Runtime {
	runtime := Runtime{
		Vars:    make(map[string]map[string]interface{}),
		Funcs:   make(map[string]map[string]func(*Runtime, []interface{}) (interface{}, error)),
		Methods: make(map[reflect.Type]map[string]NativeFunction),
		Stdout:  stdout,
		Stderr:  stderr,
		Stdin:   stdin,
	}
	runtime.Vars["GLOBAL"] = make(map[string]interface{})
	runtime.Funcs["GLOBAL"] = make(map[string]func(*Runtime, []interface{}) (interface{}, error))
	runtime.Scopes = []string{"GLOBAL"}
	return runtime
}

func (runtime *Runtime) GetVar(varname string) (interface{}, error) {
	for i := len(runtime.Scopes) - 1; i >= 0; i-- {
		scope := runtime.Scopes[i]
		val, ok := runtime.Vars[scope][varname]
		if ok {
			return val, nil
		}

		val, ok = runtime.Funcs[scope][varname]
		if ok {
			return val, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Variable %s is not declared", varname))
}

func (runtime *Runtime) SetVar(scopename string, varname string, varval interface{}) {
	runtime.Vars[runtime.Scopes[len(runtime.Scopes)-1]][varname] = varval
}

func (runtime *Runtime) SetArrayIndex(varname string, index int64, varval interface{}) {
	for i := len(runtime.Scopes) - 1; i >= 0; i-- {
		scope := runtime.Scopes[i]
		_, ok := runtime.Vars[scope][varname]
		if ok {
			runtime.Vars[scope][varname].([]interface{})[index] = varval
		}
	}
}

func (runtime *Runtime) ApplyToVariable(scopename string, varname string, operation func(interface{}) (interface{}, error)) error {
	val, err := operation(runtime.Vars[scopename][varname])
	if err != nil {
		return err
	}
	runtime.Vars[scopename][varname] = val
	return nil
}

func (runtime *Runtime) VarExists(varname string) (bool, string) {
	for i := len(runtime.Scopes) - 1; i >= 0; i-- {
		_, exists := runtime.Vars[runtime.Scopes[i]][varname]
		if exists {
			return true, runtime.Scopes[i]
		}
	}
	return false, ""
}

func (runtime *Runtime) GetFunc(funcname string) NativeFunction {
	for i := len(runtime.Scopes) - 1; i >= 0; i-- {
		scope := runtime.Scopes[i]
		val, ok := runtime.Funcs[scope][funcname]
		if ok {
			return val
		}
	}

	return nil
}

func (runtime *Runtime) SetFunc(funcname string, fn NativeFunction) {
	runtime.Funcs[runtime.Scopes[len(runtime.Scopes)-1]][funcname] = fn
}

func (runtime *Runtime) CurrentScope() string {
	return runtime.Scopes[len(runtime.Scopes)-1]
}

// Set the current scope
func (runtime *Runtime) AddScope(scopename string) {
	runtime.Scopes = append(runtime.Scopes, scopename)
	runtime.Vars[scopename] = make(map[string]interface{})
	runtime.Funcs[scopename] = make(map[string]NativeFunction)
}

// Exit the current scope
func (runtime *Runtime) ExitScope() {
	runtime.Scopes = runtime.Scopes[:len(runtime.Scopes)-1]
}

func (runtime *Runtime) GetMethod(calledOnValue interface{}, methodname string) NativeFunction {
	methodMap, hasType := runtime.Methods[reflect.TypeOf(calledOnValue)]
	if !hasType {
		return nil
	}

	return methodMap[methodname]
}

func (runtime *Runtime) SetMethod(onType reflect.Type, methodname string, method NativeFunction) {
	methodMap, hasType := runtime.Methods[onType]
	if !hasType {
		runtime.Methods[onType] = make(map[string]NativeFunction)
	}
	methodMap[methodname] = method
}
