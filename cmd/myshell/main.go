package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput() string {
	for {
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading user input: %s\n", err.Error())
			continue
		}
		return strings.TrimSpace(input)
	}
}

func main() {
	fmt.Fprint(os.Stdout, "$ ")
	input := cleanInput()
	switch input {
	default:
		fmt.Fprintf(os.Stderr, "%s: command not found\n", input)
	}
}
