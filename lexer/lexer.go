package lexer

import (
	"github.com/codecrafters-io/interpreter-starter-go/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
}

func New(input string) *Lexer {
	l := &Lexer{input: string(input), line: 1}
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

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
		}
		l.readChar()
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()

	switch l.ch {
	case 0:
		tok = token.New(token.EOF, "\x00", "EOF", l.line)
	case '(':
		tok = token.New(token.LEFT_PAREN, string(l.ch), "LEFT_PAREN", l.line)
	case ')':
		tok = token.New(token.RIGHT_PAREN, string(l.ch), "RIGHT_PAREN", l.line)
	case '{':
		tok = token.New(token.LEFT_BRACE, string(l.ch), "LEFT_BRACE", l.line)
	case '}':
		tok = token.New(token.RIGHT_BRACE, string(l.ch), "RIGHT_BRACE", l.line)
	case '.':
		tok = token.New(token.DOT, string(l.ch), "DOT", l.line)
	case '*':
		tok = token.New(token.STAR, string(l.ch), "STAR", l.line)
	case ',':
		tok = token.New(token.COMMA, string(l.ch), "COMMA", l.line)
	case '+':
		tok = token.New(token.PLUS, string(l.ch), "PLUS", l.line)
	case '-':
		tok = token.New(token.MINUS, string(l.ch), "MINUS", l.line)
	case ';':
		tok = token.New(token.SEMICOLON, string(l.ch), "SEMICOLON", l.line)
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.LESS_EQUAL, "<=", "LESS_EQUAL", l.line)
		} else {
			tok = token.New(token.LESS, string(l.ch), "LESS", l.line)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.GREATER_EQUAL, ">=", "GREATER_EQUAL", l.line)
		} else {
			tok = token.New(token.GREATER, string(l.ch), "GREATER", l.line)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.BANG_EQUAL, "!=", "BANG_EQUAL", l.line)
		} else {
			tok = token.New(token.BANG, string(l.ch), "BANG", l.line)
		}
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.EQUAL_EQUAL, "==", "EQUAL_EQUAL", l.line)
		} else {
			tok = token.New(token.EQUAL, string(l.ch), "EQUAL", l.line)
		}
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}

			return l.NextToken()
		} else {
			tok = token.New(token.SLASH, string(l.ch), "SLASH", l.line)
		}

	default:
		tok = token.New(token.ILLEGAL, string(l.ch), "ILLEGAL", l.line)
	}

	l.readChar()
	return tok
}
