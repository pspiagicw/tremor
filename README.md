# Tremor

Tremor is a small statically typed programming language implemented in Go. It takes Tremor source code through a classic compiler pipeline:

`lexer -> parser -> typechecker -> compiler -> fenc bytecode -> fenc VM`

The project is focused on language implementation rather than a large standard library. It already supports a useful subset of the language, a REPL, example programs, diagnostics with source locations, and a test suite around parsing, typechecking, and code generation.

## What the project is

Tremor is best understood as an experimental language and compiler project:

- The front end is implemented in this repository.
- Bytecode emission and execution are delegated to [`fenc`](../fenc) through its emitter and VM APIs.
- The codebase is organized around separate stages, which makes it easy to study or extend one part of the pipeline at a time.

This repository does **not** currently aim to be a production-ready language runtime. Some syntax is parsed and typechecked before it is fully mature at execution time, so the README below distinguishes between "implemented in the frontend" and "safe to rely on end to end".

## Current language surface

The codebase and tests indicate support for the following language features.

### Values and types

- `int`
- `float`
- `string`
- `bool`
- arrays like `[]int`
- hashes like `[string]int`
- functions, including first-class function types such as `fn(int) int`
- `void`
- inferred `let` bindings via `auto`

### Statements and expressions

- `let` declarations with optional type annotations
- assignment to existing variables
- arithmetic on integers and floats
- boolean logic with `and`, `or`, `not`
- comparisons: `==`, `!=`, `<`, `<=`, `>`, `>=`
- prefix negation for numeric values
- string concatenation with `..`
- `if / else / end`
- `return`
- named functions
- lambda expressions
- array and hash literals
- indexing into arrays and hashes
- class declarations in the parser/typechecker/compiler surface

### Built-in functions

The built-ins currently registered in `builtins/builtins.go` are:

- `print(value)`
- `len(value)`
- `str(value)`
- `type(value)`
- `exit()`

## Example

```tm
fn add(a int, b int) int then
    return a + b
end

let result int = add(10, 32)
print(str(result))
```

More examples are available in `examples/`.

## Project structure

- `main.go`: CLI entrypoint. Starts the REPL when no file is passed, otherwise executes a source file.
- `lexer/`: tokenization.
- `parser/`: AST construction and parser diagnostics.
- `ast/`: AST node definitions and source locations.
- `types/`: type model used by the checker and compiler.
- `typechecker/`: semantic analysis, scope management, and type inference/checking.
- `compiler/`: lowers the typed AST into `fenc` bytecode via the emitter.
- `builtins/`: runtime builtin registration plus builtin type information for the checker.
- `batch/`: file execution flow.
- `repl/`: interactive REPL loop.
- `diagnostic/`: rendering of human-readable source diagnostics.

## How execution works

When a file is executed, Tremor does the following:

1. Reads the source file.
2. Lexes and parses it into an AST.
3. Typechecks the AST and records node-to-type information.
4. Compiles the typed AST into `fenc` bytecode.
5. Dumps constants and bytecode instructions in batch mode.
6. Runs the bytecode on the `fenc` VM with Tremor built-ins attached.

The REPL follows the same general path but works one input at a time. If `TREMOR_DEBUG=1` is set, it also prints the AST and bytecode information for each entered expression.

## Prerequisites

- Go `1.24`
- A local checkout of [`fenc`](https://github.com/pspiagicw/fenc)

The module file currently contains:

```go
replace github.com/pspiagicw/fenc => ../fenc
```

That means the repository expects `fenc` to exist as a sibling directory next to this project.

## Build

```bash
go build .
```

This produces the `tremor` binary in the repository root.

If you want to use the provided `makefile`:

```bash
make build
```

## Run

Run the REPL:

```bash
go run .
```

Run a source file:

```bash
go run . examples/functions.tm
```

Run the compiled binary:

```bash
./tremor examples/functions.tm
```

Run all bundled examples:

```bash
make run-tremor
```

## Test

The repository includes tests for:

- lexer behavior
- parser correctness
- diagnostics
- typechecker behavior
- compiler output

Run them with:

```bash
go test ./...
```

Or through the `makefile`:

```bash
make test
```

Note that the `make test` target installs and uses `tparse` for prettier test output.

## Limitations and current caveats

This is the part worth reading before extending the language.

- Tremor depends on a local `../fenc` checkout, so a clean clone of this repository alone is not enough to build.
- Some features are present in the parser and typechecker but are still experimental from a full language-design perspective.
- The examples directory includes files that are clearly exploratory; not every example should be treated as a guaranteed passing integration test.
- `batch` execution currently dumps constants and bytecode before running the VM, which is helpful for development but noisy for end users.
- The REPL keeps compiler/type information alive across iterations in a development-oriented way, so it behaves more like a language workbench than a polished shell.
- The repository contains TODOs around richer built-ins, imports, classes/methods, and stronger runtime coverage.
- One example explicitly notes that recursion is not expected to work yet.

## Why this codebase is interesting

Tremor is a compact compiler project with a clean separation between stages. If you want to experiment with language features, improve diagnostics, or study how a typed AST gets lowered into bytecode, the codebase is small enough to navigate quickly while still covering the full pipeline from source text to VM execution.
