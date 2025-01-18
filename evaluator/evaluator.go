package evaluator

import (
	"fmt"
	"io"
	"time"

	"github.com/codecrafters-io/interpreter-starter-go/ast"
	"github.com/codecrafters-io/interpreter-starter-go/object"
)

var (
	NIL   = &object.Nil{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

var builtins = map[string]*object.NativeFunction{
	"clock": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError("clock() takes no arguments")
			}
			seconds := float64(time.Now().Unix())
			return &object.Number{Value: seconds}
		},
	},
}

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
		return e.evalProgram(node.Statements, enclosedEnv)
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
	case *ast.InfixExpression:
		if node.Operator == "or" {
			return e.evalOrExpression(node.Left, node.Right, env)
		}
		if node.Operator == "and" {
			return e.evalAndExpression(node.Left, node.Right, env)
		}

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
		io.WriteString(e.stdout, value.Inspect()+"\n")
		return &object.Print{Value: value}
	case *ast.AssignExpression:
		value := e.Eval(node.Value, env)
		if isError(value) {
			return value
		}
		env.Assign(node.Name.Value, value)
		return value
	case *ast.Identifier:
		return e.evalIdentifier(node, env)
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
	case *ast.WhileStatement:
		for isTruthy(e.Eval(node.Condition, env)) {
			e.Eval(node.Consequence, env)
		}
		return nil
	case *ast.ForStatement:
		enclosedEnv := object.NewEnclosedEnvironment(env)
		e.Eval(node.Init, enclosedEnv)
		for isTruthy(e.Eval(node.Condition, enclosedEnv)) {
			e.Eval(node.Body, enclosedEnv)
			e.Eval(node.Increment, enclosedEnv)
		}
		return nil

	case *ast.FunctionLiteral:
		function := &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}

		env.Define(node.Name.Value, function)
		return function
	case *ast.CallExpression:
		function := e.Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := e.evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return e.applyFunction(function, args)
	case *ast.ReturnStatement:
		value := e.Eval(node.ReturnValue, env)
		if isError(value) {
			return value
		}
		return &object.ReturnValue{Value: value}
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
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			io.WriteString(e.stderr, result.Message)
			return result
		case *object.Print:
			continue
		}
	}

	return result
}

func (e *Evaluator) evalAndExpression(left, right ast.Node, env *object.Environment) object.Object {
	leftResult := e.Eval(left, env)
	if !isTruthy(leftResult) {
		return FALSE
	}
	return e.Eval(right, env)
}

func (e *Evaluator) evalOrExpression(left, right ast.Node, env *object.Environment) object.Object {
	leftResult := e.Eval(left, env)
	if isTruthy(leftResult) {
		return leftResult
	}
	rightResult := e.Eval(right, env)
	return rightResult
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

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Define(param.Value, args[paramIdx])
	}
	return env
}

func (e *Evaluator) applyFunction(fn object.Object, args []object.Object) object.Object {

	switch fn := fn.(type) {
	case *object.Function:
		extendEnv := extendFunctionEnv(fn, args)
		returnedValue := e.Eval(fn.Body, extendEnv)
		return unwrapReturnValue(returnedValue)

	case *object.NativeFunction:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func (e *Evaluator) evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("undefined variable: %s", node.Value)
}

func (e *Evaluator) evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, exp := range exps {
		evaluated := e.Eval(exp, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}
