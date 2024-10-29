package main

import (
	"fmt"
	"io"
	"os"
)

func tokenize(command, filename string, stdout, stderr io.Writer) error {
	if command != "tokenize" {
		fmt.Fprintf(stderr, "Unknown command: %s\n", command)
		return fmt.Errorf("unknown command: %s", command)
	}

	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(stderr, "Error reading file: %v\n", err)
		return err
	}

	for _, char := range fileContents {
		switch char {
		case '(':
			fmt.Fprintln(stdout, "LEFT_PAREN ( null") // Updated to use the stdout writer
		case ')':
			fmt.Fprintln(stdout, "RIGHT_PAREN ) null") // Updated to use the stdout writer
		}
	}

	fmt.Fprintln(stdout, "EOF  null") // Updated to use the stdout writer
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	err := tokenize(os.Args[1], os.Args[2], os.Stdout, os.Stderr)
	if err != nil {
		os.Exit(1)
	}
}
