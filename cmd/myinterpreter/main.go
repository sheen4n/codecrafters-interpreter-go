package main

import (
	"fmt"
	"io"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/evaluator"
	"github.com/codecrafters-io/interpreter-starter-go/lexer"
	"github.com/codecrafters-io/interpreter-starter-go/object"
	"github.com/codecrafters-io/interpreter-starter-go/parser"
	"github.com/codecrafters-io/interpreter-starter-go/token"
)

func tokenize(filename string, stdout, stderr io.Writer) bool {

	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(stderr, "error reading file: %v\n", err)
		return false
	}

	l := lexer.New(string(fileContents))

	ok := true
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		if tok.Type == token.UNTERMINATED_STRING {
			fmt.Fprintf(stderr, "[line %d] Error: Unterminated string.\n", tok.Line)
			ok = false
		} else if tok.Type == token.ILLEGAL {
			fmt.Fprintf(stderr, "[line %d] Error: Unexpected character: %s\n", tok.Line, string(tok.Lexeme))
			ok = false
		} else {
			fmt.Fprintf(stdout, "%s %s %s\n", tok.Type, tok.Lexeme, tok.Literal)
		}
	}

	fmt.Fprintln(stdout, "EOF  null")
	return ok
}

func parse(filename string, stdout, stderr io.Writer) bool {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(stderr, "error reading file: %v\n", err)
		return false
	}

	p := parser.New(lexer.New(string(fileContents)))

	program := p.ParseProgram()
	if !p.CheckErrors(stderr) {
		return false
	}

	fmt.Fprintln(stdout, program.String())
	return true
}

func evaluate(filename string, stdout, stderr io.Writer) bool {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(stderr, "error reading file: %v\n", err)
		return false
	}

	p := parser.New(lexer.New(string(fileContents)))

	program := p.ParseProgram()
	if !p.CheckErrors(stderr) {
		os.Exit(65)
		return false
	}

	env := object.NewEnvironment()
	e := evaluator.NewEvaluator(stdout, stderr)

	evaluated := e.Eval(program, env)
	return evaluated == nil || evaluated.Type() != object.ERROR_OBJ
}

func execute(command, filename string, stdout, stderr io.Writer) bool {
	if command == "tokenize" {
		return tokenize(filename, stdout, stderr)
	}

	if command == "parse" {
		return parse(filename, stdout, stderr)
	}

	if command == "evaluate" || command == "run" {
		if !evaluate(filename, stdout, stderr) {
			os.Exit(70)
		}
		return true
	}

	fmt.Fprintf(stderr, "unknown command: %s\n", command)
	return false
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	if !execute(os.Args[1], os.Args[2], os.Stdout, os.Stderr) {
		os.Exit(65)
	}
	os.Exit(0)
}
