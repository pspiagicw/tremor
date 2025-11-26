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
