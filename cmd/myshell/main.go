package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var exitCode = 0

func exitBuiltin(cmd *Command) bool {
	var errno int = 0
	var err error
	if len(cmd.Args) > 1 {
		errno, err = strconv.Atoi(cmd.Args[1])
		if err != nil {
			errno = 1
		}
	}
	exitCode = errno
	return false
}

func echoBuiltin(cmd *Command) bool {
	if len(cmd.Args) > 1 {
		fmt.Fprintf(os.Stdout, "%s\n", strings.Join(cmd.Args[1:], " "))
	}
	return true
}

var builtins = map[string]func(*Command) bool{
	"exit": exitBuiltin,
	"echo": echoBuiltin,
}

type Command struct {
	Args []string
}

func NewCommand() Command {
	var cmd Command
	for {
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading user input: %s\n", err.Error())
			continue
		}
		cmd.Args = strings.Fields(input)
		return cmd
	}
}

func main() {
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
