package evaluator

import (
	"fmt"

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
		result := Eval(node.Expression)
		if isError(result) {
			return result
		}
		return result
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.PrintExpression:
		value := Eval(node.Expression)
		if isError(value) {
			return value
		}

		fmt.Println(value.Inspect())

		return NIL
	}
	return nil
}

func nativeToBoolean(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalProgram(stmts []ast.Statement) object.Object {
	if len(stmts) == 0 {
		return nil
	}

	evaluated := Eval(stmts[0])
	if evaluated.Type() == object.ERROR_OBJ {
		return evaluated
	}

	return evaluated
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		if right.Type() != object.NUMBER_OBJ {
			return newError("Operand must be a number.")
		}
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

func evalInfixExpression(operator string, left, right object.Object) object.Object {

	if left.Type() == object.NUMBER_OBJ && right.Type() == object.NUMBER_OBJ {
		return evalNumberInfixExpression(operator, left, right)
	}

	if left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ {
		return evalStringInfixExpression(operator, left, right)
	}

	if (operator == "==" || operator == "!=") && left.Type() != right.Type() {
		return FALSE
	}

	return newError("Operands must be numbers.")
}

func evalNumberInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Number).Value
	rightValue := right.(*object.Number).Value

	switch operator {
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

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftValue + rightValue}
	case "==":
		return nativeToBoolean(leftValue == rightValue)
	case "!=":
		return nativeToBoolean(leftValue != rightValue)
	case "*":
		return newError("Operands must be numbers.")
	}
	return nil
}
