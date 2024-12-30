package evaluator

import (
	"github.com/codecrafters-io/interpreter-starter-go/ast"
	"github.com/codecrafters-io/interpreter-starter-go/object"
)

var (
	NIL   = &object.Nil{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.Boolean:
		return nativeToBoolean(node.Value)
	case *ast.Nil:
		return NIL
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}
	case *ast.GroupExpression:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node)
	case *ast.InfixExpression:
		return evalInfixExpression(node)
	}
	return nil
}

func nativeToBoolean(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// func newError(format string, a ...interface{}) *object.Error {
// 	return &object.Error{Message: fmt.Sprintf(format, a...)}
// }

func evalProgram(stmts []ast.Statement) object.Object {
	if len(stmts) == 0 {
		return nil
	}

	return Eval(stmts[0])
}

func evalPrefixExpression(node *ast.PrefixExpression) object.Object {
	right := Eval(node.Right)
	switch node.Operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	}
	return nil
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NIL:
		return TRUE
	}
	return FALSE
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.NUMBER_OBJ {
		return NIL
	}
	value := right.(*object.Number).Value
	return &object.Number{Value: -value}
}

func evalInfixExpression(node *ast.InfixExpression) object.Object {
	left := Eval(node.Left)
	right := Eval(node.Right)

	if left.Type() == object.NUMBER_OBJ && right.Type() == object.NUMBER_OBJ {
		return evalNumberInfixExpression(node, left, right)
	}

	if left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ {
		return evalStringInfixExpression(node, left, right)
	}

	if (node.Operator == "==" || node.Operator == "!=") && left.Type() != right.Type() {
		return FALSE
	}

	return nil
}

func evalNumberInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	leftValue := left.(*object.Number).Value
	rightValue := right.(*object.Number).Value

	switch node.Operator {
	case "+":
		return &object.Number{Value: leftValue + rightValue}
	case "-":
		return &object.Number{Value: leftValue - rightValue}
	case "*":
		return &object.Number{Value: leftValue * rightValue}
	case "/":
		return &object.Number{Value: leftValue / rightValue}
	case ">":
		return nativeToBoolean(leftValue > rightValue)
	case ">=":
		return nativeToBoolean(leftValue >= rightValue)
	case "<":
		return nativeToBoolean(leftValue < rightValue)
	case "<=":
		return nativeToBoolean(leftValue <= rightValue)
	case "==":
		return nativeToBoolean(leftValue == rightValue)
	case "!=":
		return nativeToBoolean(leftValue != rightValue)
	}
	return nil
}

func evalStringInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	switch node.Operator {
	case "+":
		return &object.String{Value: leftValue + rightValue}
	case "==":
		return nativeToBoolean(leftValue == rightValue)
	case "!=":
		return nativeToBoolean(leftValue != rightValue)
	}
	return nil
}
