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
		wantErr    error
		setupFile  func(string) error
	}{
		{
			name:       "missing tokenize command",
			command:    "invalid_command",
			filename:   "test.txt",
			wantOutput: "",
			wantErr:    errors.New("unknown command: invalid_command"),
		},
		{
			name:       "file not found",
			command:    "tokenize",
			filename:   "nonexistent.txt",
			wantOutput: "Error reading file:",
			wantErr:    errors.New("open nonexistent.txt: no such file or directory"),
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
			name:     "paranthesis",
			command:  "tokenize",
			filename: "parens.txt",
			wantOutput: `LEFT_PAREN ( null
LEFT_PAREN ( null
RIGHT_PAREN ) null
EOF  null
`,
			wantErr: nil,
			setupFile: func(filename string) error {
				return os.WriteFile(filename, []byte("(()"), 0644)
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
			err := tokenize(tt.command, tt.filename, &stdout, &stderr)

			if tt.wantErr != nil && err == nil {
				t.Errorf("expected error %v, got nil", tt.wantErr)
			} else if tt.wantErr == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			} else if tt.wantErr != nil && err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}

			// Check output
			output := stdout.String()
			if tt.wantErr == nil && tt.wantOutput != "" && output != tt.wantOutput {
				t.Errorf("expected output %q, got %q", tt.wantOutput, output)
			}
		})
	}
}
