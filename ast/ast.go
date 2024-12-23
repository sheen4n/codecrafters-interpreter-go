package ast

import (
	"bytes"

	"github.com/codecrafters-io/interpreter-starter-go/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

// Statement is a node that represents a statement in the program.
// Statements do not produce a value.
type Statement interface {
	Node
	statementNode()
}

// Expression is a node that represents an expression in the program.
// Expressions produce a value.
type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

// TokenLiteral returns the literal value of the token that represents the program.
// Program is the root node of the AST.
// It implements the Node interface.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// String returns a string representation of the program.
// It implements the Node interface.
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Lexeme }

type Nil struct {
	Token token.Token
}

func (n *Nil) expressionNode()      {}
func (n *Nil) TokenLiteral() string { return n.Token.Literal }
func (n *Nil) String() string       { return n.Token.Lexeme }
