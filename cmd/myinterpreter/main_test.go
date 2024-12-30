package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name       string
		filename   string
		wantOutput string
		wantErr    string
		setupFile  func(string) error
	}{
		{
			name:       "file not found",
			filename:   "nonexistent.txt",
			wantOutput: "",
			wantErr:    "error reading file: open nonexistent.txt: no such file or directory\n",
		},
		{
			name:       "empty file",
			filename:   "emptyfile.txt",
			wantOutput: "EOF  null\n",
			wantErr:    "",
			setupFile: func(filename string) error {
				return os.WriteFile(filename, []byte(""), 0644)
			},
		},
		{
			name:     "symbols",
			filename: "symbols.txt",
			wantOutput: `LEFT_PAREN ( null
LEFT_PAREN ( null
RIGHT_PAREN ) null
LEFT_BRACE { null
RIGHT_BRACE } null
STAR * null
DOT . null
COMMA , null
PLUS + null
STAR * null
MINUS - null
SEMICOLON ; null
EOF  null
`,
			wantErr: "",
			setupFile: func(filename string) error {
				return os.WriteFile(filename, []byte("((){}*.,+*-;"), 0644)
			},
		},
		{
			name:     "unknown symbols",
			filename: "unknown_symbols.txt",
			wantOutput: `COMMA , null
DOT . null
LEFT_PAREN ( null
EOF  null
`,
			wantErr: `[line 1] Error: Unexpected character: $
[line 2] Error: Unexpected character: #
`,
			setupFile: func(filename string) error {
				return os.WriteFile(filename, []byte(",.$(\n#"), 0644)
			},
		},
		{
			name:     "group",
			filename: "group.txt",
			wantOutput: `LEFT_PAREN ( null
STRING "foo" foo
RIGHT_PAREN ) null
EOF  null
`,
			wantErr:   "",
			setupFile: func(filename string) error { return os.WriteFile(filename, []byte(`("foo")`), 0644) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up any necessary files
			if tt.setupFile != nil {
				defer os.Remove(tt.filename)
				if err := tt.setupFile(tt.filename); err != nil {
					t.Fatalf("failed to set up file: %v", err)
				}
			}

			var stdout, stderr bytes.Buffer
			ok := tokenize(tt.filename, &stdout, &stderr)

			// Check error
			errOutput := stderr.String()
			if tt.wantErr != "" && errOutput != tt.wantErr {
				t.Errorf("expected error output\n%v, got\n%v", tt.wantErr, errOutput)
				if ok {
					t.Errorf("expected error, got success")
				}
			}

			// Check output
			output := stdout.String()
			if tt.wantOutput != "" && output != tt.wantOutput {
				t.Errorf("expected output\n%v, got\n%v, \nexpected: \n%q\n,output: \n%q", tt.wantOutput, output, tt.wantOutput, output)
			}
		})
	}
}

func TestRunUnknownCommand(t *testing.T) {
	// Arrange
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	command := "unknown"
	filename := "test.txt"

	// Act
	ok := run(command, filename, stdout, stderr)

	// Assert
	expectedError := "unknown command: unknown\n"
	if stderr.String() != expectedError {
		t.Errorf("Expected error message %q, got %q", expectedError, stderr.String())
	}
	if ok {
		t.Error("Expected run to return false for unknown command")
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		filename   string
		wantOutput string
		wantErr    string
		setupFile  func(string) error
	}{
		{
			name:       "parse boolean",
			filename:   "boolean.txt",
			wantOutput: "true",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte("true"), 0644) },
		},
		{
			name:       "parse nil",
			filename:   "nil.txt",
			wantOutput: "nil",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte("nil"), 0644) },
		},
		{
			name:       "parse number",
			filename:   "number.txt",
			wantOutput: "42.47",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte("42.47"), 0644) },
		},
		{
			name:       "parse integer",
			filename:   "integer.txt",
			wantOutput: "35.0",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte("35"), 0644) },
		},
		{
			name:       "parse string",
			filename:   "string.txt",
			wantOutput: "hello world",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte(`"hello world"`), 0644) },
		},
		{
			name:       "parse group",
			filename:   "group.txt",
			wantOutput: "(group foo)",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte(`("foo")`), 0644) },
		},
		{
			name:       "parse unary",
			filename:   "unary.txt",
			wantOutput: "(- true)",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte(`-true`), 0644) },
		},
		{
			name:       "parse bang",
			filename:   "bang.txt",
			wantOutput: "(! true)",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte(`!true`), 0644) },
		},
		{
			name:       "parse infix",
			filename:   "infix.txt",
			wantOutput: "(/ (* 16.0 38.0) 58.0)",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte(`16 * 38 / 58`), 0644) },
		},
		{
			name:       "parse comparison",
			filename:   "comparison.txt",
			wantOutput: "(< (< 83.0 99.0) 115.0)",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte(`83 < 99 < 115`), 0644) },
		},
		{
			name:       "parse equality",
			filename:   "equality.txt",
			wantOutput: `(!= foo bar)`,
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte(`"foo" != "bar"`), 0644) },
		},
		{
			name:       "parse syntax error",
			filename:   "syntax_error.txt",
			wantOutput: "",
			wantErr:    "[line 1] Error at ')': Expect expression.",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte(`(72 + )`), 0644) },
		},
		{
			name:       "parse no syntax error",
			filename:   "no_syntax_error.txt",
			wantOutput: `(!= baz hello)`,
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte(`"baz"!="hello"`), 0644) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up any necessary files
			if tt.setupFile != nil {
				defer os.Remove(tt.filename)
				if err := tt.setupFile(tt.filename); err != nil {
					t.Fatalf("failed to set up file: %v", err)
				}
			}

			var stdout, stderr bytes.Buffer
			ok := parse(tt.filename, &stdout, &stderr)

			// Check error
			errOutput := stderr.String()

			if tt.wantErr != "" {
				if strings.TrimSpace(errOutput) != strings.TrimSpace(tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, errOutput)
				}
				if ok {
					t.Errorf("expected parse to return false for syntax error")
				}
			} else {
				if errOutput != "" {
					t.Errorf("expected no error, got %v", errOutput)
				}
			}

			// Check output
			output := stdout.String()
			if tt.wantOutput != "" && strings.TrimSpace(output) != tt.wantOutput {
				t.Errorf("expected output %v, got %v", tt.wantOutput, output)
			}
		})
	}
}

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name       string
		filename   string
		wantOutput string
		wantErr    string
		setupFile  func(string) error
	}{
		{
			name:       "evaluate boolean",
			filename:   "boolean.txt",
			wantOutput: "true",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte("true"), 0644) },
		},
		{
			name:       "evaluate nil",
			filename:   "nil.txt",
			wantOutput: "nil",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte("nil"), 0644) },
		},
		{
			name:       "evaluate string",
			filename:   "string.txt",
			wantOutput: "hello world!",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte(`"hello world!"`), 0644) },
		},
		{
			name:       "evaluate number",
			filename:   "number.txt",
			wantOutput: "10.4",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte("10.4"), 0644) },
		},
		{
			name:       "evaluate group",
			filename:   "group.txt",
			wantOutput: "10.4",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte("(10.4)"), 0644) },
		},
		{
			name:       "evaluate unary",
			filename:   "unary.txt",
			wantOutput: "false",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte("!true"), 0644) },
		},
		{
			name:       "evaluate arithmetic",
			filename:   "arithmetic.txt",
			wantOutput: "20.8",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte("10.4 + 10.4"), 0644) },
		},
		{
			name:       "evaluate string concatenation",
			filename:   "string_concatenation.txt",
			wantOutput: "hello world",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte(`"hello" + " " + "world"`), 0644) },
		},
		{
			name:       "evaluate relational operators",
			filename:   "relational_operators.txt",
			wantOutput: "true",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte(`57 > -5`), 0644) },
		},
		{
			name:       "evaluate equality operators",
			filename:   "equality_operators.txt",
			wantOutput: "true",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte(`"hello" == "hello"`), 0644) },
		},
		{
			name:       "evaluate number equality",
			filename:   "number_equality.txt",
			wantOutput: "false",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte(`61 == "61"`), 0644) },
		},
		{
			name:       "evaluate error",
			filename:   "error.txt",
			wantOutput: "Operand must be a number.",
			wantErr:    "",
			setupFile:  func(filename string) error { return os.WriteFile(filename, []byte(`-true`), 0644) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up any necessary files
			if tt.setupFile != nil {
				defer os.Remove(tt.filename)
				if err := tt.setupFile(tt.filename); err != nil {
					t.Fatalf("failed to set up file: %v", err)
				}
			}

			var stdout, stderr bytes.Buffer
			ok := evaluate(tt.filename, &stdout, &stderr)

			// Check error
			errOutput := stderr.String()

			if tt.wantErr != "" {
				if strings.TrimSpace(errOutput) != strings.TrimSpace(tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, errOutput)
				}
				if ok {
					t.Errorf("expected evaluate to return false for syntax error")
				}
			} else {
				if errOutput != "" {
					t.Errorf("expected no error, got %v", errOutput)
				}
			}

			// Check output
			output := stdout.String()
			if tt.wantOutput != "" && strings.TrimSpace(output) != tt.wantOutput {
				t.Errorf("expected output %v, got %v", tt.wantOutput, output)
			}
		})
	}
}
