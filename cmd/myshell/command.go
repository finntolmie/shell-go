package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"unicode"
)

type Command struct {
	raw     string
	Args    []string
	Stdouts []*os.File
	Stderrs []*os.File
}

func NewCommand() (*Command, error) {
	var cmd Command
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading user input: %s\n", err.Error())
		return nil, err
	}
	cmd.raw = input
	cmd.parseArgs()
	cmd.handleRedirects()
	return &cmd, nil
}

func (cmd *Command) Execute() bool {
	if builtin, ok := GetBuiltins()[cmd.Args[0]]; ok {
		return builtin(cmd)
	}
	if _, err := exec.LookPath(cmd.Args[0]); err != nil {
		cmd.WriteToErr("%s: command not found\n", cmd.Args[0])
	} else {
		cmd.tryExecute()
	}
	return true
}

func (cmd *Command) tryExecute() {
	var wg sync.WaitGroup
	exe := exec.Command(cmd.Args[0], cmd.Args[1:]...)
	exe.Stdin = os.Stdin
	stdoutPipe, err := exe.StdoutPipe()
	if err != nil {
		return
	}
	stderrPipe, err := exe.StderrPipe()
	if err != nil {
		return
	}
	exe.Start()
	wg.Add(2)
	go streamOutput(&wg, stdoutPipe, cmd.WriteToOut)
	go streamOutput(&wg, stderrPipe, cmd.WriteToErr)
	wg.Wait()
}

func streamOutput(wg *sync.WaitGroup, pipe io.ReadCloser, writeFunc func(format string, a ...interface{})) {
	defer wg.Done()
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		writeFunc("%s\n", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		writeFunc("Error reading output: %s\n", err.Error())
	}
}

func (cmd *Command) handleRedirects() {
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
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			continue
		}

		if stdoutRedirect {
			cmd.Stdouts = append(cmd.Stdouts, file)
		} else {
			cmd.Stderrs = append(cmd.Stderrs, file)
		}
	}
	cmd.Args = newArgs
	if len(cmd.Stdouts) == 0 {
		cmd.Stdouts = append(cmd.Stdouts, os.Stdout)
	}
	if len(cmd.Stderrs) == 0 {
		cmd.Stderrs = append(cmd.Stderrs, os.Stderr)
	}
}

func (cmd *Command) parseArgs() {
	var sb strings.Builder
	singleQuote := false
	doubleQuote := false
	escaping := false
	for i, c := range cmd.raw {
		if escaping {
			sb.WriteRune(c)
			escaping = false
			continue
		}
		if !singleQuote && !doubleQuote && unicode.IsSpace(c) {
			if sb.Len() > 0 {
				cmd.Args = append(cmd.Args, sb.String())
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
			if i == len(cmd.raw)-1 {
				sb.WriteRune(c)
				continue
			}
			peek := cmd.raw[i+1]
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
		cmd.Args = append(cmd.Args, sb.String())
	}
}

func (cmd *Command) WriteToOut(format string, a ...any) {
	for _, file := range cmd.Stdouts {
		fmt.Fprintf(file, format, a...)
	}
}

func (cmd *Command) WriteToErr(format string, a ...any) {
	for _, file := range cmd.Stderrs {
		fmt.Fprintf(file, format, a...)
	}
}

func (cmd *Command) Close() {
	for _, file := range cmd.Stdouts {
		if file != os.Stdout {
			file.Close()
		}
	}
	for _, file := range cmd.Stderrs {
		if file != os.Stderr {
			file.Close()
		}
	}
}

func isStdoutRedirect(arg string) bool {
	return arg == ">>" || arg == ">" || arg == "1>" || arg == "1>>"
}

func isStderrRedirect(arg string) bool {
	return arg == "2>" || arg == "2>>"
}
