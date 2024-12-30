package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func isStdoutRedirect(arg string) bool {
	return arg == ">>" || arg == ">" || arg == "1>" || arg == "1>>"
}

func isStderrRedirect(arg string) bool {
	return arg == "2>" || arg == "2>>"
}

func (cmd *Command) HandleRedirects() {
	newArgs := make([]string, 0, len(cmd.Args))
	for i := 0; i < len(cmd.Args); i++ {
		arg := cmd.Args[i]
		if i == len(cmd.Args)-1 {
			newArgs = append(newArgs, arg)
			break
		}
		stdoutRedirect := isStdoutRedirect(arg)
		stderrRedirect := isStderrRedirect(arg)
		if !stdoutRedirect && !stderrRedirect {
			newArgs = append(newArgs, arg)
			continue
		}
		flags := os.O_WRONLY | os.O_CREATE
		if strings.Contains(arg, ">>") {
			flags |= os.O_APPEND
		} else {
			flags |= os.O_TRUNC
		}
		redirectPath := cmd.Args[i+1]
		i++

		file, err := os.OpenFile(redirectPath, flags, 0644)
		if err != nil {
			fmt.Fprintf(cmd.Stderr, "Error: %s\n", err)
			continue
		}

		if stdoutRedirect {
			cmd.Stdout = file
		} else {
			cmd.Stderr = file
		}
	}
	cmd.Args = newArgs
}

func (cmd *Command) cleanupRedirects() {
	if cmd.Stdout != os.Stdout {
		cmd.Stdout.Close()
	}
	if cmd.Stderr != os.Stderr {
		cmd.Stderr.Close()
	}
}

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
