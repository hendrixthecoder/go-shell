package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var BUILTIN = []string{"type", "echo", "exit"}

func initialize() string {
	fmt.Fprint(os.Stdout, "$ ")
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}
	return strings.TrimSpace(line)
}

//
// Line parser
//

func parseLine(line string) (string, []string) {
	var words []string
	var buf []rune

	inSingle := false
	inDouble := false
	escaping := false

	for _, r := range line {
		switch {

		case escaping:
			if inDouble {
				switch r {
				case '"', '\\':
					buf = append(buf, r)
				default:
					buf = append(buf, '\\', r)
				}
			} else {
				buf = append(buf, r)
			}
			escaping = false
			continue

		case r == '\\':
			if inSingle {
				buf = append(buf, r)
			} else {
				escaping = true
			}
			continue

		case r == '\'' && !inDouble:
			inSingle = !inSingle
			continue

		case r == '"' && !inSingle:
			inDouble = !inDouble
			continue

		case r == ' ' && !inSingle && !inDouble:
			if len(buf) > 0 {
				words = append(words, string(buf))
				buf = nil
			}
			continue
		}

		buf = append(buf, r)
	}

	if len(buf) > 0 {
		words = append(words, string(buf))
	}

	if len(words) == 0 {
		return "", []string{}
	}

	return words[0], words[1:]
}

//
// Builtins
//

func handleTypeBuiltin(args []string) {
	if len(args) < 2 {
		fmt.Println("type: missing argument")
		return
	}

	name := args[1]

	for _, b := range BUILTIN {
		if b == name {
			fmt.Printf("%s is a shell builtin\n", name)
			return
		}
	}

	path, err := exec.LookPath(name)
	if err == nil {
		fmt.Printf("%s is %s\n", name, path)
		return
	}

	fmt.Printf("%s: not found\n", name)
}

func handleEchoBuiltin(parts []string) error {
	if len(parts) < 2 {
		return errors.New("echo: requires an argument")
	}

	// match expected behavior: no trailing newline for redirect tests
	fmt.Fprintln(os.Stdout, strings.Join(parts[1:], " "))
	return nil
}

func handleDefault(parts []string) error {
	if len(parts) == 0 {
		return fmt.Errorf("invalid command")
	}

	path, err := exec.LookPath(parts[0])
	if err != nil {
		return fmt.Errorf("%s: not found", parts[0])
	}

	cmd := exec.Command(path, parts[1:]...)
	cmd.Args[0] = parts[0] // critical: program sees correct argv[0]

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	// swallow "exit status X"
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return nil
	}

	return err
}

//
// Execute
//

func execute(line string) {
	command, arguments := parseLine(line)
	if command == "" {
		return
	}

	originalStdout := os.Stdout

	var outputFile *os.File
	for i, arg := range arguments {
		if (arg == ">" || arg == "1>") && i+1 < len(arguments) {
			file, err := os.Create(arguments[i+1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
				return
			}

			outputFile = file
			arguments = arguments[:i] // remove redirection
			break
		}
	}

	if outputFile != nil {
		os.Stdout = outputFile
	}

	parts := append([]string{command}, arguments...)
	var err error

	switch command {
	case "exit":
		os.Exit(0)

	case "echo":
		err = handleEchoBuiltin(parts)

	case "type":
		handleTypeBuiltin(parts)

	default:
		err = handleDefault(parts)
	}

	if outputFile != nil {
		outputFile.Close()
		os.Stdout = originalStdout
	}

	if err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}

//
// Main
//

func main() {
	for {
		line := initialize()
		if line == "" {
			continue
		}
		execute(line)
	}
}
