package main

import (
	"fmt"
	"os"
)

var exitCode = 0

var builtins map[string]func(*Command) bool

func main() {
	builtins = map[string]func(*Command) bool{
		"exit": exitBuiltin,
		"echo": echoBuiltin,
		"type": typeBuiltin,
	}
	var repl bool = true
	for repl {
		fmt.Fprint(os.Stdout, "$ ")
		cmd := NewCommand()
		if len(cmd.Args) == 0 {
			continue
		}
		if builtin, ok := builtins[cmd.Args[0]]; ok {
			repl = builtin(&cmd)
			continue
		}
		fmt.Fprintf(os.Stderr, "%s: command not found\n", cmd.Args[0])
	}
	os.Exit(exitCode)
}
