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
// Helpers
//

var hasPrecedingDelimiter = func(idx int, line string, delim string) bool {
	return idx > 0 && string(line[idx-1]) == delim
}

var hasSucceedingDelimiter = func(idx int, line string, delim string) bool {
	return idx+1 < len(line) && string(line[idx+1]) == delim
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
		// ----------------------------------------
		// ESCAPING MODE
		// ----------------------------------------
		case escaping:
			// In escaping mode, the character is always literal
			buf = append(buf, r)
			escaping = false
			continue

		// ----------------------------------------
		// BACKSLASH HANDLING
		// ----------------------------------------
		case r == '\\':
			// Inside ANY quotes → literal backslash
			if inSingle || inDouble {
				buf = append(buf, r)
				continue
			}

			// Outside quotes → escape next char
			escaping = true
			continue

		// ----------------------------------------
		// SINGLE QUOTES
		// ----------------------------------------
		case r == '\'' && !inDouble:
			inSingle = !inSingle
			continue

		// ----------------------------------------
		// DOUBLE QUOTES
		// ----------------------------------------
		case r == '"' && !inSingle:
			inDouble = !inDouble
			continue

		// ----------------------------------------
		// SPACE HANDLING
		// ----------------------------------------
		case r == ' ' && !inSingle && !inDouble:
			if len(buf) > 0 {
				words = append(words, string(buf))
				buf = nil
			}
			continue
		}

		// Normal character
		buf = append(buf, r)
	}

	// Flush last word
	if len(buf) > 0 {
		words = append(words, string(buf))
	}

	return words
}

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

func handleEchoBuiltin(line string) (string, error) {
	parts := parseQuotedLine(line)
	if len(parts) < 2 {
		return "", errors.New("echo: requires an argument")
	}

	return strings.Join(parts[1:], " ") + "\n", nil
}

func handleDefault(line string) (string, error) {
	parts := parseQuotedLine(line)
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

func execute(line string) {
	var out string
	var err error

	parts := strings.Split(line, " ")

	switch parts[0] {
	case "exit":
		os.Exit(0)

	case "echo":
		out, err = handleEchoBuiltin(line)

	case "type":
		out, err = handleTypeBuiltin(parts)

	default:
		out, err = handleDefault(line)
	}

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if out != "" {
		fmt.Print(out)
	}
}

func main() {
	for {
		line := initialize()
		if line == "" {
			continue
		}
		execute(line)
	}
}
