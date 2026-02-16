package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Print("$ ")

	command, err := bufio.NewReader(os.Stdin).ReadString('\n')

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(command[:len(command)-1] + ": command not found")
}
