package object

import (
	"bytes"
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/ast"
)

type ObjectType string

const (
	BOOLEAN_OBJ         ObjectType = "BOOLEAN"
	NIL_OBJ             ObjectType = "NIL"
	STRING_OBJ          ObjectType = "STRING"
	NUMBER_OBJ          ObjectType = "NUMBER"
	ERROR_OBJ           ObjectType = "ERROR"
	BUILTIN_OBJ         ObjectType = "BUILTIN"
	PRINT_OBJ           ObjectType = "PRINT"
	NATIVE_FUNCTION_OBJ            = "NATIVE_FUNCTION"
	FUNCTION_OBJ                   = "FUNCTION"
	RETURN_VALUE_OBJ               = "RETURN_VALUE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Nil struct{}

func (n *Nil) Type() ObjectType { return NIL_OBJ }
func (n *Nil) Inspect() string  { return "nil" }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

type Number struct {
	Value float64
}

func (n *Number) Type() ObjectType { return NUMBER_OBJ }
func (n *Number) Inspect() string  { return fmt.Sprintf("%g", n.Value) }

type Error struct{ Message string }

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return e.Message }

type Print struct {
	Value Object
}

func (p *Print) Type() ObjectType { return PRINT_OBJ }
func (p *Print) Inspect() string  { return "" }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	out.WriteString("<fn foo>")

	return out.String()
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
