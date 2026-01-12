package main

import (
	"os"

	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/tremor/batch"
	"github.com/pspiagicw/tremor/repl"
)

func main() {
	if len(os.Args) == 1 {
		repl.StartREPL()
	}
	if len(os.Args) != 2 {
		goreland.LogFatal("Expected 1 arguments program")
	}

	program := os.Args[1]

	batch.ExecFile(program)
}

// Add other builtins

// TODO: print
// TODO: len
// TODO: push
// TODO: pop
// TODO: sqrt
// TODO: exp
// TODO: min
// TODO: min/max
// TODO: string
// TODO: lowercase
// TODO: uppercase
// TODO: count
// TODO: type (most important)
// TODO: exit
