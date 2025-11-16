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
		if len(cmdParts) > 1 && cmdParts[0] == "exit" {
			n, _ := strconv.Atoi(cmdParts[1])
			os.Exit(n)
		}

		fmt.Println(command + ": command not found")
	}
}
