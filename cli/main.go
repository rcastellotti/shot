// echo "echo hello bear" | go run cli/main.go
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "shot: cannot read script: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("shot: read %d bytes\n", len(input))

	cmd := exec.Command("/bin/sh", "-s", "--")
	cmd.Args = append(cmd.Args, os.Args[1:]...)
	cmd.Stdin = bytes.NewReader(input)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	err = cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() >= 0 {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "shot: failed to run shell: %v\n", err)
		os.Exit(1)
	}
}
