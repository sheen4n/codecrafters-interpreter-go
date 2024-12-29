package evaluator

import (
	"github.com/codecrafters-io/interpreter-starter-go/ast"
	"github.com/codecrafters-io/interpreter-starter-go/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.Boolean:
		return &object.Boolean{Value: node.Value}
	case *ast.Nil:
		return &object.Nil{}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	}
	return nil
}

func evalProgram(stmts []ast.Statement) object.Object {
	if len(stmts) == 0 {
		return nil
	}

	return Eval(stmts[0])
}
