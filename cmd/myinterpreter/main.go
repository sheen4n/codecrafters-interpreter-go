package main

import (
	"fmt"
	"io"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/lexer"
	"github.com/codecrafters-io/interpreter-starter-go/token"
)

func tokenize(command, filename string, stdout, stderr io.Writer) bool {
	if command != "tokenize" {
		fmt.Fprintf(stderr, "unknown command: %s\n", command)
		return false
	}

	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(stderr, "error reading file: %v\n", err)
		return false
	}

	l := lexer.New(string(fileContents))

	ok := true
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		if tok.Type == token.ILLEGAL {
			fmt.Fprintf(stderr, "[line %d] Error: Unexpected character: %s\n", tok.Line, string(tok.Lexeme))
			ok = false
		} else {
			fmt.Fprintf(stdout, "%s %s null\n", tok.Type, tok.Lexeme)
		}
	}

	fmt.Fprintln(stdout, "EOF  null")
	return ok
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	ok := tokenize(os.Args[1], os.Args[2], os.Stdout, os.Stderr)
	if !ok {
		os.Exit(65)
	}
	os.Exit(0)
}
