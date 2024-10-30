package main

import (
	"fmt"
	"io"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/lexer"
	"github.com/codecrafters-io/interpreter-starter-go/token"
)

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
	l := lexer.New(fileContents)

	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		if tok.Type == token.ILLEGAL {
			err := fmt.Errorf("[line 1] Error: Unexpected character: %s", string(tok.Literal))
			errors = append(errors, err)
		} else {
			fmt.Fprintf(stdout, "%s %s null\n", tok.Name, tok.Literal)
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
