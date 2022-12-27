// Separate runtime package for recursive imports
package runtime

import (
	"errors"
	"fmt"
	"io"
)

// TODO: scopes
type Runtime struct {
	Scopes []string
	// Scope names => Variable names => values
	Vars           map[string]map[string]interface{}
	Funcs          map[string]map[string]func(*Runtime, []interface{}) (interface{}, error)
	Stdout, Stderr io.Writer
	Stdin          io.Reader
}

func NewRuntime(stdout, stderr io.Writer, stdin io.Reader) Runtime {
	runtime := Runtime{
		Vars:   make(map[string]map[string]interface{}),
		Funcs:  make(map[string]map[string]func(*Runtime, []interface{}) (interface{}, error)),
		Stdout: stdout,
		Stderr: stderr,
		Stdin:  stdin,
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

func (runtime *Runtime) SetVar(varname string, varval interface{}) {
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

func (runtime *Runtime) VarExists(varname string) bool {
	for i := len(runtime.Scopes) - 1; i >= 0; i-- {
		_, exists := runtime.Vars[runtime.Scopes[i]][varname]
		if exists {
			return true
		}
	}
	return false
}

func (runtime *Runtime) GetFunc(funcname string) func(*Runtime, []interface{}) (interface{}, error) {
	for i := len(runtime.Scopes) - 1; i >= 0; i-- {
		scope := runtime.Scopes[i]
		val, ok := runtime.Funcs[scope][funcname]
		if ok {
			return val
		}
	}

	return nil
}

func (runtime *Runtime) SetFunc(funcname string, fn func(*Runtime, []interface{}) (interface{}, error)) {
	runtime.Funcs[runtime.Scopes[len(runtime.Scopes)-1]][funcname] = fn
}

func (runtime *Runtime) CurrentScope() string {
	return runtime.Scopes[len(runtime.Scopes)-1]
}

// Set the current scope
func (runtime *Runtime) AddScope(scopename string) {
	runtime.Scopes = append(runtime.Scopes, scopename)
	runtime.Vars[scopename] = make(map[string]interface{})
	runtime.Funcs[scopename] = make(map[string]func(*Runtime, []interface{}) (interface{}, error))
}

// Exit the current scope
func (runtime *Runtime) ExitScope() {
	runtime.Scopes = runtime.Scopes[:len(runtime.Scopes)-1]
}
