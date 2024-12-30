package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Command struct {
	Args []string
}

func NewCommand() Command {
	for {
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading user input: %s\n", err.Error())
			continue
		}
		return Command{
			Args: strings.Fields(input),
		}
	}
}
