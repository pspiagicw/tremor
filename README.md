# `tremor`

`tremor` is a compiled, typed language inspired from Go,Java which uses similar syntax to Lua.


### Features

The following features are currently present in the language:

- Types including int, float, string, bool, arrays and hashes.
- Function declaration and calling including return types and argument types.
- Support for built-in functions including print, len.

### Implementation
`tremor` depends on another library `fenc` for bytecode compilation and VM support.

`tremor` only implements a lexer, parser, a typechecker and a compiler.

The compiler simply makes the appropriate `fenc` API calls.

### Future Plans
Future plans include:

- Add robust documentation for explaining the inner workings.
- Implementation of classes and methods.
- Support for imports.
- Better bytecode execution including C support.

