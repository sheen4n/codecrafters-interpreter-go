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

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.Boolean:
		return nativeToBoolean(node.Value)
	case *ast.Nil:
		return NIL
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}
	case *ast.GroupExpression:
		result := Eval(node.Expression, env)
		if isError(result) {
			return result
		}
		return result
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.PrintExpression:
		value := Eval(node.Expression, env)
		if isError(value) {
			return value
		}

		return &object.Print{Value: value}
	case *ast.AssignExpression:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}
		env.Set(node.Name.Value, value)
		return value
	case *ast.Identifier:
		obj, ok := env.Get(node.Value)
		if !ok {
			return newError("Undefined variable '%s'.", node.Value)
		}
		return obj
	case *ast.VarStatement:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}
		env.Set(node.Name.Value, value)
		return nil

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

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	if len(stmts) == 0 {
		return nil
	}

	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt, env)
		switch result := result.(type) {
		case *object.Error:
			return result
		case *object.Print:
			fmt.Println(result.Value.Inspect())
		}
	}

	return result
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

	if operator == "==" || operator == "!=" {
		if left.Type() != right.Type() {
			return FALSE
		}

		if operator == "==" {
			return nativeToBoolean(left.Inspect() == right.Inspect())
		}

		return nativeToBoolean(left.Inspect() != right.Inspect())
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
