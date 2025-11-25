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

func initialize() []string {
	fmt.Fprint(os.Stdout, "$ ")

	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}

	line = strings.TrimSpace(line)
	return strings.Split(line, " ")
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

func handleEchoBuiltin(parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("echo: requires an argument")
	}

	args := strings.Join(parts[1:], " ")

	hasQuote := strings.Contains(args, "'")
	if hasQuote {
		return strings.ReplaceAll(args, "'", "") + "\n", nil
	}

	return strings.Join(strings.Fields(args), " ") + "\n", nil
}

func handleDefault(parts []string) (string, error) {
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

func execute(parts []string) {
	var (
		out string
		err error
	)

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
		fmt.Println(err.Error() + ".")
		return
	}

	if out != "" {
		fmt.Print(out)
	}
}

func main() {
	for {
		parts := initialize()
		if len(parts) == 0 || parts[0] == "" {
			continue
		}
		execute(parts)
	}
}
