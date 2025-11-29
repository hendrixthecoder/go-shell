package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
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
// Quoted parser
//

func parseQuotedLine(line string) []string {
	var words []string
	var buf []rune

	inSingle := false
	inDouble := false
	escaping := false

	for _, r := range line {
		switch {

		// ESCAPING MODE
		case escaping:
			// If inside double quotes, only escape " and \
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

		// BACKSLASH
		case r == '\\':
			if inSingle {
				// literal inside single quotes
				buf = append(buf, r)
			} else {
				escaping = true
			}

			continue

		// SINGLE QUOTES
		case r == '\'' && !inDouble:
			inSingle = !inSingle
			continue

		// DOUBLE QUOTES
		case r == '"' && !inSingle:
			inDouble = !inDouble
			continue

		// WHITESPACE SPLITTING
		case r == ' ' && !inSingle && !inDouble:
			if len(buf) > 0 {
				words = append(words, string(buf))
				buf = nil
			}
			continue
		}

		// normal character
		buf = append(buf, r)
	}

	if len(buf) > 0 {
		words = append(words, string(buf))
	}

	return words
}

//
// Builtins
//

func handleTypeBuiltin(parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("type: requires an argument")
	}

	target := parts[1]

	if slices.Contains(BUILTIN, target) {
		return fmt.Sprintf("%s is a shell builtin\n", target), nil
	}

	path, err := exec.LookPath(target)
	if err != nil {
		return "", fmt.Errorf("%s: not found", target)
	}

	return fmt.Sprintf("%s is %s\n", target, path), nil
}

func handleEchoBuiltin(parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("echo: requires an argument")
	}
	return strings.Join(parts[1:], " ") + "\n", nil
}

func handleDefault(parts []string) (string, error) {
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid command")
	}

	_, err := exec.LookPath(parts[0])
	if err != nil {
		return "", fmt.Errorf("%s: not found", parts[0])
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s: %v", parts[0], err)
	}

	return "", nil
}

//
// Execute
//

func execute(line string) {
	parts := parseQuotedLine(line)
	if len(parts) == 0 {
		return
	}

	var out string
	var err error

	switch parts[0] {
	case "exit":
		os.Exit(0)

	case "echo":
		out, err = handleEchoBuiltin(parts)

	case "type":
		out, err = handleTypeBuiltin(parts)

	default:
		out, err = handleDefault(parts)
	}

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if out != "" {
		fmt.Print(out)
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
