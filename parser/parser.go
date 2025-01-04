package parser

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/ast"
	"github.com/codecrafters-io/interpreter-starter-go/lexer"
	"github.com/codecrafters-io/interpreter-starter-go/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
	token.EQUAL_EQUAL:   EQUALS,
	token.BANG_EQUAL:    EQUALS,
	token.LESS:          LESSGREATER,
	token.GREATER:       LESSGREATER,
	token.LESS_EQUAL:    LESSGREATER,
	token.GREATER_EQUAL: LESSGREATER,
	token.PLUS:          SUM,
	token.MINUS:         SUM,
	token.SLASH:         PRODUCT,
	token.STAR:          PRODUCT,
	token.LEFT_PAREN:    CALL,
	token.LEFT_BRACE:    INDEX,
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	// If peekToken is not found, we immediately return LOWEST,
	// So we will not consider the token at all
	return LOWEST
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError()
	return false
}

func (p *Parser) peekError() {
	// msg := fmt.Sprintf("[line %d] Error at '%s': Expect expression.", p.peekToken.Line, p.peekToken.Lexeme)
	// p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}

	return program
}

func (p *Parser) CheckErrors(stderr io.Writer) bool {
	if len(p.errors) == 0 {
		return true
	}

	msg := strings.Join(p.errors, "\n")
	fmt.Fprintln(stderr, msg)
	return false
}

func (p *Parser) Errors() []string {
	return p.errors
}

// let 			|			x
// ^   			| 		^
// cur tok	|  		peek tok

// Read the next token and set it as the current token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)

	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.NIL, p.parseNil)
	p.registerPrefix(token.NUMBER, p.parseNumberLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LEFT_PAREN, p.parseGroupExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.PRINT, p.parsePrintStatement)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.STAR, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.LESS, p.parseInfixExpression)
	p.registerInfix(token.GREATER, p.parseInfixExpression)
	p.registerInfix(token.LESS_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.GREATER_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.EQUAL_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.BANG_EQUAL, p.parseInfixExpression)

	// Read two tokens, so curToken and peekToken are both set
	// Sets the peekToken by calling the lexer's NextToken method
	p.nextToken()

	// Sets the curToken by putting peekToken into curToken
	// Sets the peekToken by calling the lexer's NextToken method
	p.nextToken()

	return p
}

func (p *Parser) parseStatement() ast.Statement {
	return p.parseExpressmentStatement()
}

func (p *Parser) parseExpressmentStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) noPrefixParseFnError(t token.Token) {
	msg := fmt.Sprintf("[line %d] Error at '%s': Expect expression.", t.Line, t.Lexeme)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken)
		return nil
	}

	leftExp := prefix()

	// if precedence is less than the peek precedence, we need to parse more
	// e.g. 1 + 2 * 3
	//    +
	// 1     *
	//     2   3

	// if precedence is greater or equal than the peek precedence, we need to stop parsing the group
	// if precedence is greater or equal than the peek precedence, we need to stop parsing the group
	// e.g. 3 * 4 + 5
	//      +
	//   *     5
	// 3   4

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseNil() ast.Expression {
	return &ast.Nil{Token: p.curToken}
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	num, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("could not parse %q as number", p.curToken.Literal))
		return nil
	}

	return &ast.NumberLiteral{Token: p.curToken, Value: num}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseGroupExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RIGHT_PAREN) {
		return nil
	}
	return &ast.GroupExpression{Token: p.curToken, Expression: exp}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Lexeme,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parsePrintStatement() ast.Expression {
	expression := &ast.PrintExpression{
		Token: p.curToken,
	}

	p.nextToken()
	expression.Expression = p.parseExpression(LOWEST)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Lexeme,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()

	expression.Right = p.parseExpression(precedence)
	return expression
}
