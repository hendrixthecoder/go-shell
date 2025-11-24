package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var BUILTIN = []string{"type", "echo", "exit"}

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

		switch cmdParts[0] {
		case "exit":
			// if len(cmdParts) < 2 {
			// 	fmt.Println("exit: requires an argument.")
			// 	continue
			// }

			// n, err := strconv.Atoi(cmdParts[1])
			// if err != nil {
			// 	fmt.Println("exit: invalid exit code.")
			// 	continue
			// }

			os.Exit(0)

		case "echo":
			if len(cmdParts) < 2 {
				fmt.Println()
				continue
			}

			args := make([]any, len(cmdParts[1:]))
			for i, v := range cmdParts[1:] {
				args[i] = v
			}

			fmt.Fprintln(os.Stdout, args...)

		case "type":
			if len(cmdParts) < 2 {
				fmt.Println("type: requires an argument.")
				continue
			}

			target := cmdParts[1]

			if slices.Contains(BUILTIN, target) {
				fmt.Printf("%s is a shell builtin\n", target)
				continue
			}

			path, err := exec.LookPath(target)
			if err != nil {
				fmt.Println(target + ": not found")
				continue
			} else {
				fmt.Printf("%s is %s\n", target, path)
			}
		default:
			_, err := exec.LookPath(cmdParts[0])
			if err != nil {
				fmt.Println(cmdParts[0] + ": not found")
				continue
			} else {
				cmd := exec.Command(cmdParts[0], cmdParts[1:]...)

				// Connect the command's output to the shell's output
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout

				err := cmd.Run()
				if err != nil {
					fmt.Printf("%s: %v\n", cmdParts[0], err)
				}
			}
		}
	}
}
