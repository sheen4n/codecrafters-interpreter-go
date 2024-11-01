package lexer

import (
	"fmt"
	"strconv"
	"strings"

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

func (l *Lexer) readString() token.Token {
	startPos := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}

	if l.ch == 0 {
		return token.New(token.UNTERMINATED_STRING, "", "", l.line)
	}

	l.readChar()
	s := l.input[startPos : l.position-1]
	return token.New(token.STRING, fmt.Sprintf(`"%s"`, s), s, l.line)
}

func (l *Lexer) readNumber() token.Token {
	startPos := l.position
	for l.ch >= '0' && l.ch <= '9' {
		l.readChar()
	}

	if l.ch == '.' && l.peekChar() >= '0' && l.peekChar() <= '9' {
		l.readChar()
		for l.ch >= '0' && l.ch <= '9' {
			l.readChar()
		}
	}
	lexeme := l.input[startPos:l.position]

	num, err := strconv.ParseFloat(lexeme, 64)
	if err != nil {
		return token.New(token.ILLEGAL, "", "", l.line)
	}

	literal := strconv.FormatFloat(num, 'f', -1, 64)
	if !strings.Contains(literal, ".") {
		literal += ".0"
	}

	return token.New(token.NUMBER, lexeme, literal, l.line)
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readIdentifier() string {
	startPos := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[startPos:l.position]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()

	switch l.ch {
	case 0:
		tok = token.New(token.EOF, "\x00", "null", l.line)
	case '(':
		tok = token.New(token.LEFT_PAREN, string(l.ch), "null", l.line)
	case ')':
		tok = token.New(token.RIGHT_PAREN, string(l.ch), "null", l.line)
	case '{':
		tok = token.New(token.LEFT_BRACE, string(l.ch), "null", l.line)
	case '}':
		tok = token.New(token.RIGHT_BRACE, string(l.ch), "null", l.line)
	case '.':
		tok = token.New(token.DOT, string(l.ch), "null", l.line)
	case '*':
		tok = token.New(token.STAR, string(l.ch), "null", l.line)
	case ',':
		tok = token.New(token.COMMA, string(l.ch), "null", l.line)
	case '+':
		tok = token.New(token.PLUS, string(l.ch), "null", l.line)
	case '-':
		tok = token.New(token.MINUS, string(l.ch), "null", l.line)
	case ';':
		tok = token.New(token.SEMICOLON, string(l.ch), "null", l.line)
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.LESS_EQUAL, "<=", "null", l.line)
		} else {
			tok = token.New(token.LESS, string(l.ch), "null", l.line)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.GREATER_EQUAL, ">=", "null", l.line)
		} else {
			tok = token.New(token.GREATER, string(l.ch), "null", l.line)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.BANG_EQUAL, "!=", "null", l.line)
		} else {
			tok = token.New(token.BANG, string(l.ch), "null", l.line)
		}
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.EQUAL_EQUAL, "==", "null", l.line)
		} else {
			tok = token.New(token.EQUAL, string(l.ch), "null", l.line)
		}
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}

			return l.NextToken()
		} else {
			tok = token.New(token.SLASH, string(l.ch), "null", l.line)
		}

	case '"':
		l.readChar()
		return l.readString()

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return l.readNumber()

	default:
		if isLetter(l.ch) {
			lexeme := l.readIdentifier()
			return token.New(token.LookupIdent(lexeme), lexeme, "null", l.line)
		}
		tok = token.New(token.ILLEGAL, string(l.ch), "null", l.line)
	}

	l.readChar()
	return tok
}
