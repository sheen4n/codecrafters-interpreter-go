package lexer

import (
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/token"
)

func testLexTokens(t *testing.T, input string, expected []token.Token) {
	l := New(input)

	for i, expectedTok := range expected {
		actualTok := l.NextToken()
		if actualTok.Type != expectedTok.Type {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, expectedTok.Type, actualTok.Type)
		}
		if actualTok.Lexeme != expectedTok.Lexeme {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q",
				i, expectedTok.Lexeme, actualTok.Lexeme)
		}
		if actualTok.Literal != expectedTok.Literal {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, expectedTok.Literal, actualTok.Literal)
		}
		if actualTok.Line != expectedTok.Line {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d",
				i, expectedTok.Line, actualTok.Line)
		}
	}
}

func TestNextToken(t *testing.T) {
	input := `=+(){},;==!!=<><=>=/
	
	!`

	expected := []token.Token{
		{Type: token.EQUAL, Lexeme: "=", Literal: "null", Line: 1},
		{Type: token.PLUS, Lexeme: "+", Literal: "null", Line: 1},
		{Type: token.LEFT_PAREN, Lexeme: "(", Literal: "null", Line: 1},
		{Type: token.RIGHT_PAREN, Lexeme: ")", Literal: "null", Line: 1},
		{Type: token.LEFT_BRACE, Lexeme: "{", Literal: "null", Line: 1},
		{Type: token.RIGHT_BRACE, Lexeme: "}", Literal: "null", Line: 1},
		{Type: token.COMMA, Lexeme: ",", Literal: "null", Line: 1},
		{Type: token.SEMICOLON, Lexeme: ";", Literal: "null", Line: 1},
		{Type: token.EQUAL_EQUAL, Lexeme: "==", Literal: "null", Line: 1},
		{Type: token.BANG, Lexeme: "!", Literal: "null", Line: 1},
		{Type: token.BANG_EQUAL, Lexeme: "!=", Literal: "null", Line: 1},
		{Type: token.LESS, Lexeme: "<", Literal: "null", Line: 1},
		{Type: token.GREATER, Lexeme: ">", Literal: "null", Line: 1},
		{Type: token.LESS_EQUAL, Lexeme: "<=", Literal: "null", Line: 1},
		{Type: token.GREATER_EQUAL, Lexeme: ">=", Literal: "null", Line: 1},
		{Type: token.SLASH, Lexeme: "/", Literal: "null", Line: 1},
		{Type: token.BANG, Lexeme: "!", Literal: "null", Line: 3},
		{Type: token.EOF, Lexeme: "\x00", Literal: "null", Line: 3},
	}

	testLexTokens(t, input, expected)
}

func TestLexComments(t *testing.T) {
	input := "=// This is a comment"

	expected := []token.Token{
		{Type: token.EQUAL, Lexeme: "=", Literal: "null", Line: 1},
		{Type: token.EOF, Lexeme: "\x00", Literal: "null", Line: 1},
	}

	testLexTokens(t, input, expected)
}

func TestMultilineError(t *testing.T) {
	input := `# (
)	@`

	expected := []token.Token{
		{Type: token.ILLEGAL, Lexeme: "#", Literal: "null", Line: 1},
		{Type: token.LEFT_PAREN, Lexeme: "(", Literal: "null", Line: 1},
		{Type: token.RIGHT_PAREN, Lexeme: ")", Literal: "null", Line: 2},
		{Type: token.ILLEGAL, Lexeme: "@", Literal: "null", Line: 2},
		{Type: token.EOF, Lexeme: "\x00", Literal: "null", Line: 2},
	}

	testLexTokens(t, input, expected)
}

func TestString(t *testing.T) {
	input := `"hello world"
"foo bar"`

	expected := []token.Token{
		{Type: token.STRING, Lexeme: `"hello world"`, Literal: "hello world", Line: 1},
		{Type: token.STRING, Lexeme: `"foo bar"`, Literal: "foo bar", Line: 2},
		{Type: token.EOF, Lexeme: "\x00", Literal: "null", Line: 2},
	}

	testLexTokens(t, input, expected)
}

func TestUnterminatedString(t *testing.T) {
	input := `"hello world`

	expected := []token.Token{
		{Type: token.UNTERMINATED_STRING, Lexeme: "", Literal: "", Line: 1},
		{Type: token.EOF, Lexeme: "\x00", Literal: "null", Line: 1},
	}

	testLexTokens(t, input, expected)
}

func TestNumberLiterals(t *testing.T) {
	input := `123
123.456
55.0000`

	expected := []token.Token{
		{Type: token.NUMBER, Lexeme: "123", Literal: "123.0", Line: 1},
		{Type: token.NUMBER, Lexeme: "123.456", Literal: "123.456", Line: 2},
		{Type: token.NUMBER, Lexeme: "55.0000", Literal: "55.0", Line: 3},
		{Type: token.EOF, Lexeme: "\x00", Literal: "null", Line: 3},
	}

	testLexTokens(t, input, expected)
}

func TestIdentifiers(t *testing.T) {
	input := `foo bar _hello _123_hello`

	expected := []token.Token{
		{Type: token.IDENTIFIER, Lexeme: "foo", Literal: "null", Line: 1},
		{Type: token.IDENTIFIER, Lexeme: "bar", Literal: "null", Line: 1},
		{Type: token.IDENTIFIER, Lexeme: "_hello", Literal: "null", Line: 1},
		{Type: token.IDENTIFIER, Lexeme: "_123_hello", Literal: "null", Line: 1},
		{Type: token.EOF, Lexeme: "\x00", Literal: "null", Line: 1},
	}

	testLexTokens(t, input, expected)
}

func TestReservedKeywords(t *testing.T) {
	input := `foo bar and class else false for fun if nil or print return super this true var while`

	expected := []token.Token{
		{Type: token.IDENTIFIER, Lexeme: "foo", Literal: "null", Line: 1},
		{Type: token.IDENTIFIER, Lexeme: "bar", Literal: "null", Line: 1},
		{Type: token.AND, Lexeme: "and", Literal: "null", Line: 1},
		{Type: token.CLASS, Lexeme: "class", Literal: "null", Line: 1},
		{Type: token.ELSE, Lexeme: "else", Literal: "null", Line: 1},
		{Type: token.FALSE, Lexeme: "false", Literal: "null", Line: 1},
		{Type: token.FOR, Lexeme: "for", Literal: "null", Line: 1},
		{Type: token.FUN, Lexeme: "fun", Literal: "null", Line: 1},
		{Type: token.IF, Lexeme: "if", Literal: "null", Line: 1},
		{Type: token.NIL, Lexeme: "nil", Literal: "null", Line: 1},
		{Type: token.OR, Lexeme: "or", Literal: "null", Line: 1},
		{Type: token.PRINT, Lexeme: "print", Literal: "null", Line: 1},
		{Type: token.RETURN, Lexeme: "return", Literal: "null", Line: 1},
		{Type: token.SUPER, Lexeme: "super", Literal: "null", Line: 1},
		{Type: token.THIS, Lexeme: "this", Literal: "null", Line: 1},
		{Type: token.TRUE, Lexeme: "true", Literal: "null", Line: 1},
		{Type: token.VAR, Lexeme: "var", Literal: "null", Line: 1},
		{Type: token.WHILE, Lexeme: "while", Literal: "null", Line: 1},
		{Type: token.EOF, Lexeme: "\x00", Literal: "null", Line: 1},
	}

	testLexTokens(t, input, expected)
}
