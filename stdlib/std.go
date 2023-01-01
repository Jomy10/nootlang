package stdlib

import (
	"errors"
	"github.com/jomy10/nootlang/runtime"
	"os"
	"reflect"
	"regexp"
)

func Register(r *runtime.Runtime) {
	r.Funcs["GLOBAL"]["read_to_string"] = read_to_string

	// string methods
	string_type := reflect.TypeOf("")
	// > regex methods
	_, ok := r.Methods[string_type]
	if !ok {
		r.Methods[string_type] = make(map[string]runtime.NativeFunction)
	}
	r.Methods[string_type]["match_indices"] = string__match_indices
	r.Methods[string_type]["submatch"] = string__submatch
}

func read_to_string(runtime *runtime.Runtime, args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("`read_to_string` expects a file name as an argument")
	}

	fileName, ok := args[0].(string)
	if !ok {
		return nil, errors.New("`read_to_string` expects a string argument")
	}

	dat, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return string(dat), nil
}

// NOTE: name might change
func string__match_indices(runtime *runtime.Runtime, args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, errors.New("`string.match_indices` expects one argument")
	}

	str, ok := args[0].(string)
	if !ok {
		return nil, errors.New("interpreter error")
	}

	var regex *regexp.Regexp

	switch args[1].(type) {
	case string:
		var err error
		regex, err = regexp.Compile(args[1].(string))
		if err != nil {
			return nil, err
		}
	case *regexp.Regexp:
		regex = args[1].(*regexp.Regexp)
	default:
		return nil, errors.New("Invalid argument type in `string.match`")
	}

	// TODO: change to an array of tuples
	result := regex.FindAllStringIndex(str, -1)
	arr := make([]interface{}, len(result))
	for i, element := range result {
		arr[i] = make([]interface{}, len(element))
		for j, e := range element {
			arr[i].([]interface{})[j] = e
		}
	}
	return arr, nil
}

func string__submatch(runtime *runtime.Runtime, args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, errors.New("`string.re_find_index` expects one argument")
	}

	str, ok := args[0].(string)
	if !ok {
		return nil, errors.New("interpreter error")
	}

	var regex *regexp.Regexp

	switch args[1].(type) {
	case string:
		var err error
		regex, err = regexp.Compile(args[1].(string))
		if err != nil {
			return nil, err
		}
	case *regexp.Regexp:
		regex = args[1].(*regexp.Regexp)
	default:
		return nil, errors.New("Invalid argument type in `string.match`")
	}

	result := regex.FindAllStringSubmatch(str, -1)
	arr := make([]interface{}, len(result))
	for i, element := range result {
		arr[i] = make([]interface{}, len(element))
		for j, e := range element {
			arr[i].([]interface{})[j] = e
		}
	}
	return arr, nil
}
