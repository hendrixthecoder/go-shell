package main

import (
	"bufio"
	"fmt"
	"os"
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

			pathVar := os.Getenv("PATH")
			dirs := strings.Split(pathVar, string(os.PathListSeparator))

			found := false
			for _, dir := range dirs {
				full := dir + string(os.PathSeparator) + target

				info, err := os.Stat(full)
				if err != nil {
					continue
				}

				if info.IsDir() {
					continue
				}

				if info.Mode()&0111 != 0 {
					fmt.Printf("%s is %s\n", target, full)
					found = true
					break
				}
			}

			if !found {
				fmt.Printf("%s: not found\n", target)
			}

		default:
			fmt.Println(cmdParts[0] + ": command not found")
		}
	}
}
