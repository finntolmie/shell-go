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
	for _, c := range input {
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
		case c == '\\' && !doubleQuote && !singleQuote:
			escaping = true
		case c == '\'' && !doubleQuote:
			singleQuote = !singleQuote
		case c == '"' && !singleQuote:
			doubleQuote = !doubleQuote
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
