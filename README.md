# Noot lang

A simple scripting language, created specifically for [NootBot](https://github.com/unitoftime/nootbot).

## Hello World
```noot
// helloWorld.noot

noot!("Hello World")
```

## Description

Nootlang is a simple scripting language mainly developed for [NootBot](https://github.com/unitoftime/nootbot).
The interpreter and parser are intentionally kept simple and readable as it is also
developed as a learning project. Anyone interested in making programming languages or
parsers and interpreters in general should be able to understand the code rather easily.
The interpreter is therefore focussed on readability rather than speed.

## Roadmap

- **Types**
  - [x] integers
  - [x] **strings**
  - [x] **floats**
  - [x] **booleans**
  - [x] functions
  - [ ] anonymous functions
  ```noot
  def myFunc() {}
  def myOtherFunc(fn) { fn() }
  myOtherFunc(myFunc)
  myOtherFunc(|| {})
  ```
  - [ ] structs
  - [ ] interfaces
  - [ ] tuples
  - [ ] type assertion
- **Functions**
  - [x] functions
  - [x] **scopes**
- [ ] modules
- **Statements**
  - [ ] **if/elsif/else**
  - [ ] **match**

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

## Contributing

Contributions are always welcome.

- Additions (e.g. new syntax, etc.) to the language should be discussed first in
an issue befoe submitting a pull request
- Speed improvements to the interpreter, parser or tokenizer will be accepted as
long as they do not compromise on readability. The focus of this project lies in
readability for newcomers.
- However, if you are passionate about making a fast interpreter for nootlang,
this is encouraged. Just make a new folder in this project for the faster interpreter
so that they are separated. Here, readability can be compromised for speed (JIT
compilation, caching, etc.)
- Documentation is highly needed at the moment

## License

Nootlang is licensed under the [MIT license](LICENSE).
