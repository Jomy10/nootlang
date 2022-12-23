package stdlib

import (
	"errors"
	"fmt"
	"github.com/jomy10/nootlang/runtime"
)

func Register(runtime *runtime.Runtime) {
	runtime.Funcs["noot!"] = nootLine
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
