package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

var ErrInvalidShebang = errors.New("not a valid shebang line")

// parseShebang parses a shebang line and returns the interpreter name followed by any arguments
//
// The line must start with #! (leading/trailing whitespace).
// If it does not, ErrInvalidShebang is returned.
//
// Returns nil for empty lines, invalid shebangs, or lines that don't start with an interpreter path.
func parseShebang(line string) ([]string, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, ErrInvalidShebang
	}

	if !strings.HasPrefix(line, "#!") {
		return nil, ErrInvalidShebang
	}

	lineWithoutShebang := line[2:]
	if lineWithoutShebang == "" {
		return nil, ErrInvalidShebang
	}

	parts := strings.Fields(lineWithoutShebang)
	if len(parts) == 0 {
		return nil, ErrInvalidShebang
	}

	prog := parts[0]

	result := make([]string, 0, len(parts))
	result = append(result, prog)

	if len(parts) > 1 {
		result = append(result, parts[1:]...)
	}

	return result, nil
}

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "shot: cannot read script: %v\n", err)
		os.Exit(1)
	}
	log.Printf("shot: read %v bytes\n", len(input))

	sheBangLine := strings.SplitN(string(input), "\n", 2)[0]
	parsedShebang, err := parseShebang(sheBangLine)
	if err != nil {
		fmt.Fprintf(os.Stderr, "shot: cannot parse shebang: %q: %v\n", sheBangLine, err)
		os.Exit(1)
	}

	cmd := exec.Command(parsedShebang[0], parsedShebang[1:]...)
	cmd.Args = append(cmd.Args, os.Args[1:]...)
	cmd.Stdin = bytes.NewReader(input)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	err = cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() >= 0 {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "shot: %v\n", err)
		os.Exit(1)
	}
}
