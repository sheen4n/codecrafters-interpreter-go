package ast

import (
	"bytes"
	"fmt"

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
}

func (n *Nil) expressionNode()      {}
func (n *Nil) TokenLiteral() string { return "nil" }
func (n *Nil) String() string       { return "nil" }

type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (il *NumberLiteral) expressionNode()      {}
func (il *NumberLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *NumberLiteral) String() string       { return il.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type GroupExpression struct {
	Token      token.Token // the LEFT_PAREN token
	Expression Expression
}

func (ge *GroupExpression) expressionNode()      {}
func (ge *GroupExpression) TokenLiteral() string { return ge.Token.Literal }
func (ge *GroupExpression) String() string {
	return fmt.Sprintf("(group %s)", ge.Expression.String())
}

type PrefixExpression struct {
	Token    token.Token // the prefix token, e.g. -
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(" ")
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token // the operator token, e.g. +
	Operator string
	Left     Expression
	Right    Expression
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {

	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Operator)
	out.WriteString(" " + oe.Left.String() + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}

type PrintExpression struct {
	Token      token.Token // the token.PRINT token
	Expression Expression
}

func (ps *PrintExpression) expressionNode()      {}
func (ps *PrintExpression) TokenLiteral() string { return ps.Token.Literal }
func (ps *PrintExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(print ")
	out.WriteString(ps.Expression.String())
	out.WriteString(")")
	return out.String()
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type VarStatement struct {
	Token token.Token // the token.VAR token
	Name  *Identifier
	Value Expression
}

func (ls *VarStatement) statementNode()       {}
func (ls *VarStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *VarStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}
