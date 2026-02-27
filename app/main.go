package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
)

func parseInput(input string) []string {
	var args []string
	var currentArg strings.Builder
	inSingleQuotes := false
	inDoubleQuotes := false
	inArg := false
	escaped := false

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if escaped {
			currentArg.WriteByte(ch)
			escaped = false
			continue
		}

		if ch == '\\' {
			if inSingleQuotes {
				inArg = true
				currentArg.WriteByte(ch)
			} else if inDoubleQuotes {
				if i+1 < len(input) && (input[i+1] == '\\' || input[i+1] == '"' || input[i+1] == '$' || input[i+1] == '`' || input[i+1] == '\n') {
					escaped = true
					inArg = true
				} else {
					inArg = true
					currentArg.WriteByte(ch)
				}
			} else {
				escaped = true
				inArg = true
			}
			continue
		}

		if ch == '\'' && !inDoubleQuotes {
			inSingleQuotes = !inSingleQuotes
			inArg = true
		} else if ch == '"' && !inSingleQuotes {
			inDoubleQuotes = !inDoubleQuotes
			inArg = true
		} else if (ch == ' ' || ch == '\t') && !inSingleQuotes && !inDoubleQuotes {
			if inArg {
				args = append(args, currentArg.String())
				currentArg.Reset()
				inArg = false
			}
		} else {
			inArg = true
			currentArg.WriteByte(ch)
		}
	}
	if inArg {
		args = append(args, currentArg.String())
	}

	return args
}

func main() {

	builtins := map[string]func([]string, io.Writer){
		"exit": func(_ []string, _ io.Writer) {
			os.Exit(0)
		},
		"echo": func(args []string, out io.Writer) {
			fmt.Fprintln(out, strings.Join(args, " "))
		},
		"pwd": func(_ []string, out io.Writer) {
			dir, _ := os.Getwd()
			fmt.Fprintln(out, dir)
		},
		"cd": func(path []string, _ io.Writer) {
			var pathStr string
			if len(path) > 0 && path[0] == "~" {
				home, err := os.UserHomeDir()
				if err != nil {
					return
				}

				pathStr = home
			} else if len(path) > 0 {
				pathStr = filepath.Join(path...)
			} else {
				return
			}

			err := os.Chdir(pathStr)

			if err != nil {
				fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", pathStr)
			}
		},
	}

	builtins["type"] = func(args []string, out io.Writer) {
		if len(args) == 0 {
			return
		}

		command := args[0]
		if _, ok := builtins[command]; ok {
			fmt.Fprintln(out, command+" is a shell builtin")
		} else if path, err := exec.LookPath(command); err == nil {
			fmt.Fprintln(out, command+" is "+path)
		} else {
			fmt.Fprintln(out, command+": not found")
		}
	}
	completer := readline.NewPrefixCompleter(
		readline.PcItem("echo"),
		readline.PcItem("exit"),
	)
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "$ ",
		AutoComplete: completer,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing readline:", err)
		return
	}
	defer rl.Close()

	for {

		input, err := rl.Readline()
		fmt.Println(input)

		if err != nil {
			if err == io.EOF {
				os.Exit(0)
			}

			fmt.Println("Error reading input: ", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		parts := parseInput(input)
		if len(parts) == 0 {
			continue
		}

		var args []string
		var outFile string
		var errFile string
		var appendOut bool
		var appendErr bool

		for i := 0; i < len(parts); i++ {
			if parts[i] == ">" || parts[i] == "1>" {
				if i+1 < len(parts) {
					outFile = parts[i+1]
					appendOut = false
					i++
				}
			} else if parts[i] == ">>" || parts[i] == "1>>" {
				if i+1 < len(parts) {
					outFile = parts[i+1]
					appendOut = true
					i++
				}
			} else if parts[i] == "2>" {
				if i+1 < len(parts) {
					errFile = parts[i+1]
					appendErr = false
					i++
				}
			} else if parts[i] == "2>>" {
				if i+1 < len(parts) {
					errFile = parts[i+1]
					appendErr = true
					i++
				}
			} else {
				args = append(args, parts[i])
			}
		}

		if len(args) == 0 {
			continue
		}

		cmd := args[0]
		callArgs := args[1:]

		var outWriter io.Writer = os.Stdout
		var outFilePtr *os.File
		if outFile != "" {
			var err error
			flags := os.O_WRONLY | os.O_CREATE
			if appendOut {
				flags |= os.O_APPEND
			} else {
				flags |= os.O_TRUNC
			}
			outFilePtr, err = os.OpenFile(outFile, flags, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
				continue
			}
			outWriter = outFilePtr
		}

		var errWriter io.Writer = os.Stderr
		var errFilePtr *os.File
		if errFile != "" {
			var err error
			flags := os.O_WRONLY | os.O_CREATE
			if appendErr {
				flags |= os.O_APPEND
			} else {
				flags |= os.O_TRUNC
			}
			errFilePtr, err = os.OpenFile(errFile, flags, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
				if outFilePtr != nil {
					outFilePtr.Close()
				}
				continue
			}
			errWriter = errFilePtr
		}

		if handler, ok := builtins[cmd]; ok {
			handler(callArgs, outWriter)
		} else if _, err := exec.LookPath(cmd); err == nil {
			command := exec.Command(cmd, callArgs...)
			command.Stdout = outWriter
			command.Stdin = os.Stdin
			command.Stderr = errWriter
			command.Run()
		} else {
			// For unrecognized commands, bash typically outputs command not found to stderr
			fmt.Fprintln(errWriter, cmd+": command not found")
		}

		if outFilePtr != nil {
			outFilePtr.Close()
		}
		if errFilePtr != nil {
			errFilePtr.Close()
		}
	}
}
