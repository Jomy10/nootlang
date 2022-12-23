// Separate runtime package for recursive imports
package runtime

import "io"

// TODO: scopes
type Runtime struct {
	// Variable names => values
	Vars           map[string]interface{}
	Funcs          map[string]func(*Runtime, []interface{}) (interface{}, error)
	Stdout, Stderr io.Writer
	Stdin          io.Reader
}

func NewRuntime(stdout, stderr io.Writer, stdin io.Reader) Runtime {
	return Runtime{
		Vars:   make(map[string]interface{}),
		Funcs:  make(map[string]func(*Runtime, []interface{}) (interface{}, error)),
		Stdout: stdout,
		Stderr: stderr,
		Stdin:  stdin,
	}
}
