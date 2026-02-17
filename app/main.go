package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
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

		switch cmd {
		case "type":
			if len(args) == 0 {
				continue
			}
			param := args[0]
			fmt.Println(param, "This is the param")
			if param == "type" || param == "exit" || param == "echo" {
				fmt.Println(param + " is a shell builtin")
			} else {
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
