package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
	for {
		fmt.Fprint(os.Stdout, "$ ")
		cmd := NewCommand()
		if len(cmd.Args) == 0 {
			continue
		}
		switch cmd.Args[0] {
		case "exit":
			var exitCode int = 0
			var err error
			if len(cmd.Args) > 1 {
				exitCode, err = strconv.Atoi(cmd.Args[1])
				if err != nil {
					exitCode = 1
				}
			}
			os.Exit(exitCode)
		default:
			fmt.Fprintf(os.Stderr, "%s: command not found\n", cmd.Args[0])
		}
	}
}
