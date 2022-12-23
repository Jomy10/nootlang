# Noot lang

A simple scripting language, created specifically for [NootBot](https://github.com/unitoftime/nootbot).


## Hello World
```noot
// helloWorld.noot

noot!("Hello World")
```

## Native Function Interface

Native functions (like the ones in the [standard library](/stdlib)), are of the
following signature:

```go
func(*runtime.Runtime, args []interface{}) (interface{}, error)
```

- The first argument passed to any function is always the runtime, this contains
all the variables and functions available.
- The second is an array of all the arguments passed to this function

- The function can return a value as its first return type, or nil if it does not
return a value
- If a runtime error occurs during execution, the function should return a
descriptive error as its second argument
