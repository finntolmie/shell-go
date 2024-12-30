package main

import (
	"fmt"
	"os"
	"os/exec"
)

var exitCode = 0

var builtins map[string]func(*Command) bool

func execute(cmd *Command) {
	exe := exec.Command(cmd.Args[0], cmd.Args[1:]...)
	exe.Stdout = cmd.Stdout
	exe.Stderr = cmd.Stderr
	err := exe.Run()
	if err != nil {
		fmt.Fprintln(cmd.Stderr, err.Error())
	}
}

func attemptExecute(cmd *Command) bool {
	if fp := findInPath(cmd.Args[0]); fp != "" {
		execute(cmd)
	} else {
		fmt.Fprintf(cmd.Stderr, "%s: command not found\n", cmd.Args[0])
	}
	return true
}

func main() {
	builtins = map[string]func(*Command) bool{
		"exit": exitBuiltin,
		"echo": echoBuiltin,
		"type": typeBuiltin,
		"pwd":  pwdBuiltin,
		"cd":   cdBuiltin,
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
		repl = attemptExecute(&cmd)
	}
	os.Exit(exitCode)
}
