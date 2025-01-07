package main

import (
	"fmt"
	"os"
)

var exitCode int

func printLoop() bool {
	fmt.Fprint(os.Stdout, "$ ")
	cmd, err := NewCommand()
	if err != nil {
		return false
	}
	defer cmd.Close()
	if len(cmd.Args) == 0 {
		return false
	}
	return cmd.Execute()
}

func main() {
	for printLoop() {
		os.Stdout.Sync()
	}
	os.Exit(exitCode)
}
