package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	StatusExecutable   = 0
	StatusNotFound     = 1
	StatusUnExecutable = 2
)

func main() {
	paths := os.Getenv("PATH")
	pathsArray := strings.Split(paths, string(os.PathListSeparator))
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")

		input, err := reader.ReadString('\n')

		if err != nil {
			if err.Error() == "EOF" {
				os.Exit(0)
			}

			fmt.Println("Error reading input: ", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		cmd := parts[0]
		args := parts[1:]

	LoopLabel:
		switch cmd {
		case "type":
			if len(args) == 0 {
				continue
			}
			param := args[0]
			if param == "type" || param == "exit" || param == "echo" {
				fmt.Println(param + " is a shell builtin")
			} else {
				for _, path := range pathsArray {
					fullPath := filepath.Join(path, param)
					res := fileExistsAndPermission(fullPath)
					if res == StatusExecutable {
						// fmt.Println(param, "is", fullPath)

						cmd := exec.Command(fullPath, args...)

						if err := cmd.Run(); err != nil {
							fmt.Println("Error", err)
						}
						break LoopLabel
					} else if res == StatusNotFound {
						continue
					}
				}
				fmt.Println(param + ": not found")
			}
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(strings.Join(args, " "))
		default:
			fmt.Println(cmd + ": command not found")
		}

	}
}

func fileExistsAndPermission(path string) int {
	info, err := os.Stat(path)
	if err == nil {
		mode := info.Mode()
		if mode.Perm()&011 != 0 {
			return StatusExecutable
		}
		return StatusUnExecutable
	}
	if errors.Is(err, os.ErrNotExist) {
		return StatusNotFound
	}

	return StatusUnExecutable
}
