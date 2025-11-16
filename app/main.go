package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		commandLn, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		command := commandLn[:len(commandLn)-1]

		cmdParts := strings.Split(command, " ")
		if len(cmdParts) < 2 {
			fmt.Println(command + ": command not found")
			continue
		}

		switch cmdParts[0] {
		case "exit":
			n, _ := strconv.Atoi(cmdParts[1])
			os.Exit(n)
		case "echo":
			args := make([]any, len(cmdParts[1:]))

			for i, v := range cmdParts[1:] {
				args[i] = v
			}

			fmt.Fprintln(os.Stdout, args...)
		default:
			fmt.Println(cmdParts[0] + ": command not found")
		}
	}
}
