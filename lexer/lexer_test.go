package lexer

import (
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;==!!=<><=>=/
	
	!`

	tests := []struct {
		expectedType   token.TokenType
		expectedLexeme string
	}{
		{token.EQUAL, "="},
		{token.PLUS, "+"},
		{token.LEFT_PAREN, "("},
		{token.RIGHT_PAREN, ")"},
		{token.LEFT_BRACE, "{"},
		{token.RIGHT_BRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EQUAL_EQUAL, "=="},
		{token.BANG, "!"},
		{token.BANG_EQUAL, "!="},
		{token.LESS, "<"},
		{token.GREATER, ">"},
		{token.LESS_EQUAL, "<="},
		{token.GREATER_EQUAL, ">="},
		{token.SLASH, "/"},
		{token.BANG, "!"},
		{token.EOF, "\x00"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Lexeme != tt.expectedLexeme {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q",
				i, tt.expectedLexeme, tok.Lexeme)
		}
	}
}

func TestLexComments(t *testing.T) {
	input := "=// This is a comment"

	tests := []struct {
		expectedType   token.TokenType
		expectedLexeme string
	}{
		{token.EQUAL, "="},
		{token.EOF, "\x00"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Lexeme != tt.expectedLexeme {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q",
				i, tt.expectedLexeme, tok.Lexeme)
		}
	}
}

func TestMultilineError(t *testing.T) {
	input := `# (
)	@`

	tests := []struct {
		expectedType   token.TokenType
		expectedLexeme string
		expectedLine   int
	}{
		{token.ILLEGAL, "#", 1},
		{token.LEFT_PAREN, "(", 1},
		{token.RIGHT_PAREN, ")", 2},
		{token.ILLEGAL, "@", 2},
		{token.EOF, "\x00", 2},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Lexeme != tt.expectedLexeme {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q",
				i, tt.expectedLexeme, tok.Lexeme)
		}
		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d",
				i, tt.expectedLine, tok.Line)
		}
	}
}

func TestString(t *testing.T) {
	input := `"hello world"
"foo bar"`

	tests := []struct {
		expectedType    token.TokenType
		expectedLexeme  string
		expectedLiteral string
	}{
		{token.STRING, `"hello world"`, "hello world"},
		{token.STRING, `"foo bar"`, "foo bar"},
		{token.EOF, "\x00", "null"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Lexeme != tt.expectedLexeme {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q",
				i, tt.expectedLexeme, tok.Lexeme)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestUnterminatedString(t *testing.T) {
	input := `"hello world`

	tests := []struct {
		expectedType    token.TokenType
		expectedLexeme  string
		expectedLiteral string
		expectedLine    int
	}{
		{token.UNTERMINATED_STRING, "", "", 1},
		{token.EOF, "\x00", "null", 1},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Lexeme != tt.expectedLexeme {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q",
				i, tt.expectedLexeme, tok.Lexeme)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNumberLiterals(t *testing.T) {
	input := `123
123.456
55.0000`

	tests := []struct {
		expectedType    token.TokenType
		expectedLexeme  string
		expectedLiteral string
	}{
		{token.NUMBER, "123", "123.0"},
		{token.NUMBER, "123.456", "123.456"},
		{token.NUMBER, "55.0000", "55.0"},
		{token.EOF, "\x00", "null"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Lexeme != tt.expectedLexeme {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q",
				i, tt.expectedLexeme, tok.Lexeme)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
