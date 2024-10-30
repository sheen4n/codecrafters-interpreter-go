package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Name    string
}

const (
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"

	// Operators
	DOT   = "."
	STAR  = "*"
	PLUS  = "+"
	MINUS = "-"

	// Delimiters
	LEFT_PAREN  = "("
	RIGHT_PAREN = ")"
	LEFT_BRACE  = "{"
	RIGHT_BRACE = "}"
	SEMICOLON   = ";"
	COMMA       = ","
)

func New(tokenType TokenType, ch byte, name string) Token {
	return Token{Type: tokenType, Literal: string(ch), Name: name}
}
