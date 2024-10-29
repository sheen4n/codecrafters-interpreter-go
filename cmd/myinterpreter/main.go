package main

import (
	"fmt"
	"io"
	"os"
)

var TOKENS map[byte]string = map[byte]string{
	'(': "LEFT_PAREN",
	')': "RIGHT_PAREN",
	'{': "LEFT_BRACE",
	'}': "RIGHT_BRACE",
	'.': "DOT",
	'*': "STAR",
	',': "COMMA",
	'+': "PLUS",
	'-': "MINUS",
	';': "SEMICOLON",
}

func tokenize(command, filename string, stdout, stderr io.Writer) []error {
	if command != "tokenize" {
		fmt.Fprintf(stderr, "Unknown command: %s\n", command)
		return []error{fmt.Errorf("unknown command: %s", command)}
	}

	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(stderr, "Error reading file: %v\n", err)
		return []error{err}
	}

	errors := []error{}
	for _, char := range fileContents {
		token, ok := TOKENS[char]

		if ok {
			fmt.Fprintf(stdout, "%s %s null\n", token, string(char))
		} else {
			err := fmt.Errorf("[line 1] Error: Unexpected character: %s", string(char))
			errors = append(errors, err)
		}

	}

	fmt.Fprintln(stdout, "EOF  null") // Updated to use the stdout writer
	return errors
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	errs := tokenize(os.Args[1], os.Args[2], os.Stdout, os.Stderr)
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(65)
	}
}
