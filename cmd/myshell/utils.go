package main

import (
	"fmt"
	"os"
	"strings"
)

func findInPath(target string) string {
	paths := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
	for _, p := range paths {
		fullPath := fmt.Sprintf("%s/%s", p, target)
		fs, err := os.Stat(fullPath)
		if err != nil {
			continue
		}
		if !fs.IsDir() && (fs.Mode().Perm()&0100) != 0 {
			return fullPath
		}
	}
	return ""
}
