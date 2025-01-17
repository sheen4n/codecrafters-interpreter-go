package token

type TokenType string

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal string
	Line    int
}

const (
	EOF = "EOF"

	// Errors
	ILLEGAL             = "ILLEGAL"
	UNTERMINATED_STRING = "UNTERMINATED_STRING"

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

	// Keywords
	STRING     = "STRING"
	NUMBER     = "NUMBER"
	IDENTIFIER = "IDENTIFIER"

	// Reserved Keywords
	AND      = "AND"
	CLASS    = "CLASS"
	ELSE     = "ELSE"
	FALSE    = "FALSE"
	FOR      = "FOR"
	FUNCTION = "FUN"
	IF       = "IF"
	NIL      = "NIL"
	OR       = "OR"
	PRINT    = "PRINT"
	RETURN   = "RETURN"
	SUPER    = "SUPER"
	THIS     = "THIS"
	TRUE     = "TRUE"
	VAR      = "VAR"
	WHILE    = "WHILE"
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUNCTION,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

func New(tokenType TokenType, lexeme, literal string, line int) Token {
	return Token{Type: tokenType, Lexeme: lexeme, Line: line, Literal: literal}
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENTIFIER
}
