package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Name    string
	Line    int
}

const (
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"

	// Operators
	DOT           = "."
	STAR          = "*"
	PLUS          = "+"
	MINUS         = "-"
	EQUAL         = "="
	EQUAL_EQUAL   = "=="
	BANG          = "!"
	BANG_EQUAL    = "!="
	LESS          = "<"
	GREATER       = ">"
	LESS_EQUAL    = "<="
	GREATER_EQUAL = ">="
	SLASH         = "/"

	// Delimiters
	LEFT_PAREN  = "("
	RIGHT_PAREN = ")"
	LEFT_BRACE  = "{"
	RIGHT_BRACE = "}"
	SEMICOLON   = ";"
	COMMA       = ","
)

func New(tokenType TokenType, ch string, name string, line int) Token {
	return Token{Type: tokenType, Literal: ch, Name: name, Line: line}
}
