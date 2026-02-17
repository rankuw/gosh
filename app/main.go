package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
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
					res := fileExistsAndPermission(path + "/" + param)
					if res == 0 {
						fmt.Println(param, "is", path+"/"+param)
						break LoopLabel
					} else if res == 1 {
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
	_, err := os.Stat(path)
	if err == nil {
		return 0
	}
	if errors.Is(err, os.ErrNotExist) {
		return 1
	}

	return 2
}
