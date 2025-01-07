package main

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func GetBuiltins() map[string]func(*Command) bool {
	return map[string]func(*Command) bool{
		"exit": exitBuiltin,
		"echo": echoBuiltin,
		"type": typeBuiltin,
		"pwd":  pwdBuiltin,
		"cd":   cdBuiltin,
	}
}

func exitBuiltin(cmd *Command) bool {
	var errno int
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
		elements := make([]string, 0, len(cmd.Args)-1)
		for _, arg := range cmd.Args[1:] {
			if env, ok := strings.CutPrefix(arg, "$"); ok {
				elements = append(elements, os.Getenv(env))
			} else {
				elements = append(elements, arg)
			}
		}
		cmd.WriteToOut("%s\n", strings.Join(elements, " "))
	} else {
		cmd.WriteToOut("\n")
	}
	return true
}

func typeBuiltin(cmd *Command) bool {
	if len(cmd.Args) == 1 {
		cmd.WriteToErr("Error: expected argument")
	} else if _, ok := GetBuiltins()[cmd.Args[1]]; ok {
		cmd.WriteToOut("%s is a shell builtin\n", cmd.Args[1])
	} else if fp, err := exec.LookPath(cmd.Args[1]); err != nil {
		cmd.WriteToErr("%s: not found\n", cmd.Args[1])
	} else {
		cmd.WriteToOut("%s is %s\n", cmd.Args[1], fp)
	}
	return true
}

func pwdBuiltin(cmd *Command) bool {
	pwd, err := os.Getwd()
	if err != nil {
		cmd.WriteToErr("Error: %s\n", err.Error())
	} else {
		cmd.WriteToOut("%s\n", pwd)
	}
	return true
}

func cdBuiltin(cmd *Command) bool {
	cdPath := cmd.Args[1]
	if cdPath[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			cmd.WriteToErr("Error: HOME not set")
			return true
		}
		cdPath = home + cdPath[1:]
	}
	err := os.Chdir(cdPath)
	if err != nil {
		cmd.WriteToErr("cd: %s: No such file or directory\n", cdPath)
	}
	return true
}
