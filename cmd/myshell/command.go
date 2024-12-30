package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func ExtractArgs(input string) []string {
	var args []string
	var sb strings.Builder
	singleQuote := false
	doubleQuote := false
	escaping := false
	for i, c := range input {
		if escaping {
			sb.WriteRune(c)
			escaping = false
			continue
		}
		if !singleQuote && !doubleQuote && unicode.IsSpace(c) {
			if sb.Len() > 0 {
				args = append(args, sb.String())
				sb.Reset()
			}
			continue
		}
		switch {
		case c == '\'' && !doubleQuote:
			singleQuote = !singleQuote
		case c == '"' && !singleQuote:
			doubleQuote = !doubleQuote
		case c == '\\' && doubleQuote:
			// edge case im not handling it
			if i == len(input)-1 {
				sb.WriteRune(c)
				continue
			}
			peek := input[i+1]
			if peek == '\\' || peek == '"' || peek == '$' || peek == '\n' {
				escaping = true
			} else {
				sb.WriteRune(c)
			}
		case c == '\\' && !doubleQuote && !singleQuote:
			escaping = true
		default:
			sb.WriteRune(c)
		}
	}
	if sb.Len() > 0 {
		args = append(args, sb.String())
	}
	return args
}

type Command struct {
	Args   []string
	Stdout *os.File
	Stderr *os.File
}

func NewCommand() Command {
	for {
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading user input: %s\n", err.Error())
			continue
		}
		args := ExtractArgs(strings.TrimSpace(input))
		return Command{
			Args:   args,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
	}
}
