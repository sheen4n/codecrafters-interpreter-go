package lexer

import (
	"github.com/codecrafters-io/interpreter-starter-go/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: string(input)}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) curToken() byte {
	return l.input[l.position]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case 0:
		tok = token.New(token.EOF, l.ch, "EOF")
	case '(':
		tok = token.New(token.LEFT_PAREN, l.ch, "LEFT_PAREN")
	case ')':
		tok = token.New(token.RIGHT_PAREN, l.ch, "RIGHT_PAREN")
	case '{':
		tok = token.New(token.LEFT_BRACE, l.ch, "LEFT_BRACE")
	case '}':
		tok = token.New(token.RIGHT_BRACE, l.ch, "RIGHT_BRACE")
	case '.':
		tok = token.New(token.DOT, l.ch, "DOT")
	case '*':
		tok = token.New(token.STAR, l.ch, "STAR")
	case ',':
		tok = token.New(token.COMMA, l.ch, "COMMA")
	case '+':
		tok = token.New(token.PLUS, l.ch, "PLUS")
	case '-':
		tok = token.New(token.MINUS, l.ch, "MINUS")
	case ';':
		tok = token.New(token.SEMICOLON, l.ch, "SEMICOLON")
	default:
		tok = token.New(token.ILLEGAL, l.ch, "ILLEGAL")
	}

	l.readChar()
	return tok
}
