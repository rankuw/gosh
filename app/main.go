package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	builtins := map[string]func([]string){
		"exit": func(_ []string) {
			os.Exit(0)
		},
		"echo": func(args []string) {
			fmt.Println(strings.Join(args, " "))
		},
		"pwd": func(_ []string) {
			dir, _ := os.Getwd()
			fmt.Println(dir)
		},
		"cd": func(path []string) {
			var pathStr string
			if path[0] == "~" {
				home, err := os.UserHomeDir()
				if err != nil {

				}

				pathStr = home
			} else {
				pathStr = filepath.Join(path...)
			}

			err := os.Chdir(pathStr)

			if err != nil {
				fmt.Printf("cd: %s: No such file or directory\n", pathStr)
			}
		},
	}

	builtins["type"] = func(args []string) {
		if len(args) == 0 {
			return
		}

		command := args[0]
		if _, ok := builtins[command]; ok {
			fmt.Println(command + " is a shell builtin")
		} else if path, err := exec.LookPath(command); err == nil {
			fmt.Println(command + " is " + path)
		} else {
			fmt.Println(command + ": not found")
		}
	}
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

		if handler, ok := builtins[cmd]; ok {
			handler(args)
		} else if _, err := exec.LookPath(cmd); err == nil {
			command := exec.Command(cmd, args...)
			command.Stdout = os.Stdout
			command.Stdin = os.Stdin
			command.Stderr = os.Stderr
			command.Run()
		} else {
			fmt.Println(cmd + ": command not found")
		}
	}
}
