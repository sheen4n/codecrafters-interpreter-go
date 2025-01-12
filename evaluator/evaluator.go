package evaluator

import (
	"fmt"
	"io"

	"github.com/codecrafters-io/interpreter-starter-go/ast"
	"github.com/codecrafters-io/interpreter-starter-go/object"
)

var (
	NIL   = &object.Nil{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

type Evaluator struct {
	stdout io.Writer
	stderr io.Writer
}

func NewEvaluator(stdout, stderr io.Writer) *Evaluator {
	return &Evaluator{stdout: stdout, stderr: stderr}
}

func (e *Evaluator) Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node.Statements, env)
	case *ast.BlockStatement:
		enclosedEnv := object.NewEnclosedEnvironment(env)
		e.evalProgram(node.Statements, enclosedEnv)
		return nil
	case *ast.ExpressionStatement:
		return e.Eval(node.Expression, env)
	case *ast.Boolean:
		return nativeToBoolean(node.Value)
	case *ast.Nil:
		return NIL
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}
	case *ast.GroupExpression:
		result := e.Eval(node.Expression, env)
		if isError(result) {
			return result
		}
		return result
	case *ast.PrefixExpression:
		right := e.Eval(node.Right, env)
		return evalPrefixExpression(node.Operator, right)
		// TODO: fix this to not eval too early if short-circuit by OR
	case *ast.InfixExpression:
		left := e.Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := e.Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.PrintExpression:
		value := e.Eval(node.Expression, env)
		if isError(value) {
			return value
		}

		return &object.Print{Value: value}
	case *ast.AssignExpression:
		value := e.Eval(node.Value, env)
		if isError(value) {
			return value
		}
		env.Assign(node.Name.Value, value)
		return value
	case *ast.Identifier:
		obj, ok := env.Get(node.Value)
		if !ok {
			return newError("Undefined variable '%s'.", node.Value)
		}
		return obj
	case *ast.VarStatement:
		value := e.Eval(node.Value, env)
		if isError(value) {
			return value
		}
		env.Define(node.Name.Value, value)
		return nil
	case *ast.IfStatement:
		condition := e.Eval(node.Condition, env)
		if isError(condition) {
			return condition
		}
		if isTruthy(condition) {
			return e.Eval(node.Consequence, env)
		}
		if node.Alternative != nil {
			return e.Eval(node.Alternative, env)
		}
	}
	return nil
}

func isTruthy(obj object.Object) bool {
	if obj == NIL || obj == FALSE {
		return false
	}
	return true
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

func (e *Evaluator) evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	if len(stmts) == 0 {
		return nil
	}

	var result object.Object
	for _, stmt := range stmts {
		result = e.Eval(stmt, env)
		switch result := result.(type) {
		case *object.Error:
			io.WriteString(e.stderr, result.Message)
			return result
		case *object.Print:
			io.WriteString(e.stdout, result.Value.Inspect()+"\n")
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

	if operator == "or" {
		if isTruthy(left) {
			return left
		}
		if isTruthy(right) {
			return right
		}
		return FALSE
	}

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
