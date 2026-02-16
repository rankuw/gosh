package main

import (
	"fmt"
)

func main() {
	var str string
	fmt.Print("$ ")
	fmt.Scanln(&str, "")
	fmt.Println(str + ": command not found")
}
