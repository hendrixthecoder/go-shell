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

	line = strings.TrimSpace(line)
	return line
}

func parseQuotedLine(line string, delimiter string) []string {
	word := ""
	words := []string{}
	inQuote := false

	hasPrecedingDelimiter := func(currIdx int, line string, lineLen int, delimiter string) bool {
		return currIdx+1 < lineLen && string(line[currIdx-1]) == delimiter
	}

	hasSucceedingDelimiter := func(currIdx int, line string, lineLen int, delimiter string) bool {
		return currIdx+1 < lineLen && string(line[currIdx+1]) == delimiter
	}

	for idx, char := range line {
		if string(char) == delimiter {
			// So we can continue capturing characters.
			if hasPrecedingDelimiter(idx, line, len(line), delimiter) || hasSucceedingDelimiter(idx, line, len(line), delimiter) {
				continue
			}

			inQuote = !inQuote

			if word != "" {
				words = append(words, word)
				word = ""
			}
		} else if string(char) == " " && !inQuote {
			// Since we are not in quotes and the char is an empty space
			// it can be ignored as it is a separation.
			if word != "" {
				words = append(words, word)
				word = ""
				inQuote = false
			}
		} else {
			word += string(char)
		}

		// Cases where there's a word at the last char not appended.
		if idx == len(line)-1 && word != "" {
			words = append(words, word)
		}
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
	if len(line) < 2 {
		return "", errors.New("echo: requires an argument")
	}

	parts := parseQuotedLine(line, "'")

	return strings.Join(parts[1:], " ") + "\n", nil
}

func handleDefault(line string) (string, error) {
	parts := parseQuotedLine(line, "'")
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

	// Command printed its own output, so nothing to return
	return "", nil
}

func execute(line string) {
	var (
		out string
		err error
	)

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
		fmt.Println(err.Error() + ".")
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
