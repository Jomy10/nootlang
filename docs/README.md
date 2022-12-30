# Nootlang

Nootlang is a simple scripting language.

```noot
noot!("Hello world")
```

## Basics

### Variables

To declare a new variable:
```
my_var := 5
```

To re-assign a variable:
```
my_var = 6
```

This distinction between declaration and assignment is important in scopes:

```
var1 := 5
def my_func() {
  noot!(var1) # Outputs 5 (global scope)
  var1 := 6
  noot!(var1) # Outputs 6 (function scope)
}

my_func()
noot!(var1) # Output 5 (global scope)
```