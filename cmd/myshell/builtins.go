package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

func typeBuiltin(cmd *Command) bool {
	if len(cmd.Args) == 1 {
		fmt.Fprintln(os.Stderr, "Error: expected argument")
		return true
	}
	if _, ok := builtins[cmd.Args[1]]; ok {
		fmt.Printf("%s is a shell builtin\n", cmd.Args[1])
		return true
	}
	if fp := findInPath(cmd.Args[1]); fp != "" {
		fmt.Printf("%s is %s\n", cmd.Args[1], fp)
		return true
	}
	fmt.Fprintf(os.Stderr, "%s: not found\n", cmd.Args[1])
	return true
}
