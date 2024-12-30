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
		fmt.Fprintf(cmd.Stdout, "%s\n", strings.Join(cmd.Args[1:], " "))
	}
	return true
}

func typeBuiltin(cmd *Command) bool {
	if len(cmd.Args) == 1 {
		fmt.Fprintln(cmd.Stderr, "Error: expected argument")
		return true
	}
	if _, ok := builtins[cmd.Args[1]]; ok {
		fmt.Fprintf(cmd.Stdout, "%s is a shell builtin\n", cmd.Args[1])
		return true
	}
	if fp := findInPath(cmd.Args[1]); fp != "" {
		fmt.Fprintf(cmd.Stdout, "%s is %s\n", cmd.Args[1], fp)
		return true
	}
	fmt.Fprintf(cmd.Stderr, "%s: not found\n", cmd.Args[1])
	return true
}

func pwdBuiltin(cmd *Command) bool {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(cmd.Stderr, "Error: %s\n", err.Error())
	} else {
		fmt.Fprintln(cmd.Stdout, pwd)
	}
	return true
}

func cdBuiltin(cmd *Command) bool {
	cdPath := cmd.Args[1]
	if cdPath[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(cmd.Stderr, "Error: HOME not set")
			return true
		}
		cdPath = home + cdPath[1:]
	}
	err := os.Chdir(cdPath)
	if err != nil {
		fmt.Fprintf(cmd.Stderr, "cd: %s: No such file or directory\n", cdPath)
	}
	return true
}
