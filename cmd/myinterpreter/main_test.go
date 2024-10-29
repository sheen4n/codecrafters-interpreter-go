package main

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name       string
		command    string
		filename   string
		wantOutput string
		wantErr    []error
		setupFile  func(string) error
	}{
		{
			name:       "missing tokenize command",
			command:    "invalid_command",
			filename:   "test.txt",
			wantOutput: "",
			wantErr:    []error{errors.New("unknown command: invalid_command")},
		},
		{
			name:       "file not found",
			command:    "tokenize",
			filename:   "nonexistent.txt",
			wantOutput: "",
			wantErr:    []error{errors.New("open nonexistent.txt: no such file or directory")},
		},
		{
			name:       "empty file",
			command:    "tokenize",
			filename:   "emptyfile.txt",
			wantOutput: "EOF  null\n",
			wantErr:    nil,
			setupFile: func(filename string) error {
				return os.WriteFile(filename, []byte(""), 0644)
			},
		},
		{
			name:     "symbols",
			command:  "tokenize",
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
			wantErr: nil,
			setupFile: func(filename string) error {
				return os.WriteFile(filename, []byte("((){}*.,+*-;"), 0644)
			},
		},
		{
			name:     "unknown symbols",
			command:  "tokenize",
			filename: "unknown_symbols.txt",
			wantOutput: `COMMA , null
DOT . null
LEFT_PAREN ( null
EOF  null
`,
			wantErr: []error{
				errors.New("[line 1] Error: Unexpected character: $"),
				errors.New("[line 1] Error: Unexpected character: #")},
			setupFile: func(filename string) error {
				return os.WriteFile(filename, []byte(",.$(#"), 0644)
			},
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
			errs := tokenize(tt.command, tt.filename, &stdout, &stderr)

			if len(tt.wantErr) > 0 {
				if len(errs) != len(tt.wantErr) {
					t.Errorf("expected %d errors, got %d", len(tt.wantErr), len(errs))
				}

				for i, err := range errs {
					if err.Error() != tt.wantErr[i].Error() {
						t.Errorf("expected error %v, got %v", tt.wantErr[i], err)
					}
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
