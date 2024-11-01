package token

type TokenType string

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
}

const (
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"

	// Operators
	DOT           = "DOT"
	STAR          = "STAR"
	PLUS          = "PLUS"
	MINUS         = "MINUS"
	EQUAL         = "EQUAL"
	EQUAL_EQUAL   = "EQUAL_EQUAL"
	BANG          = "BANG"
	BANG_EQUAL    = "BANG_EQUAL"
	LESS          = "LESS"
	GREATER       = "GREATER"
	LESS_EQUAL    = "LESS_EQUAL"
	GREATER_EQUAL = "GREATER_EQUAL"
	SLASH         = "SLASH"

	// Delimiters
	LEFT_PAREN  = "LEFT_PAREN"
	RIGHT_PAREN = "RIGHT_PAREN"
	LEFT_BRACE  = "LEFT_BRACE"
	RIGHT_BRACE = "RIGHT_BRACE"
	SEMICOLON   = "SEMICOLON"
	COMMA       = "COMMA"
)

func New(tokenType TokenType, ch string, line int) Token {
	return Token{Type: tokenType, Lexeme: ch, Line: line}
}
