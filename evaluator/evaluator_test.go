package evaluator

import (
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/lexer"
	"github.com/codecrafters-io/interpreter-starter-go/object"
	"github.com/codecrafters-io/interpreter-starter-go/parser"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func testNilObject(t *testing.T, obj object.Object) bool {
	if obj.Type() != object.NIL_OBJ {
		t.Errorf("object is not nil. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%q, want=%q", result.Value, expected)
		return false
	}
	return true
}

func testNumberObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Number)
	if !ok {
		t.Errorf("object is not Number. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
		return false
	}
	return true
}

func testErrorObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.Error)
	if !ok {
		t.Errorf("object is not Error. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Message != expected {
		t.Errorf("object has wrong value. got=%q, want=%q", result.Message, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalNil(t *testing.T) {
	evaluated := testEval("nil")
	testNilObject(t, evaluated)
}

func TestEvalString(t *testing.T) {
	evaluated := testEval(`"hello world!"`)
	testStringObject(t, evaluated, "hello world!")
}

func TestEvalNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"10.40", 10.4},
		{"10", 10},
		{"10.400", 10.4},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestEvalGroupExpression(t *testing.T) {
	evaluated := testEval("(10.4)")
	testNumberObject(t, evaluated, 10.4)

	evaluated = testEval("(true)")
	testBooleanObject(t, evaluated, true)

	evaluated = testEval("(nil)")
	testNilObject(t, evaluated)

	evaluated = testEval(`("hello world!")`)
	testStringObject(t, evaluated, "hello world!")
}

func TestUnaryExpression(t *testing.T) {
	evaluated := testEval("-10.4")
	testNumberObject(t, evaluated, -10.4)

	evaluated = testEval("!true")
	testBooleanObject(t, evaluated, false)

	evaluated = testEval("!false")
	testBooleanObject(t, evaluated, true)

	evaluated = testEval("!nil")
	testBooleanObject(t, evaluated, true)

	evaluated = testEval("-(-10.4)")
	testNumberObject(t, evaluated, 10.4)
}

func TestEvaluateArithmeticExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"42 / 5", 8.4},
		{"18 * 3 / (3 * 6)", 3},
		{"(10.40 * 2) / 2", 10.4},
		{"10.4 + 10.4", 20.8},
		{"10.4 - 10.4", 0},
		{"10.4 / 10.4", 1},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestStringConcatenation(t *testing.T) {
	evaluated := testEval(`"hello" + " " + "world"`)
	testStringObject(t, evaluated, "hello world")
}

func TestRelationalOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"57 > -5", true},
		{"11 >= 11", true},
		{"(54 - 64) >= -(114 / 57 + 11)", true},
		{"57 > 500", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEqualityOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`"hello" == "world"`, false},
		{`"hello" != "world"`, true},
		{`"hello" == "hello"`, true},
		{`"hello" != "hello"`, false},
		{"61 == 61", true},
		{"61 != 61", false},
		{`61 == "61"`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`-"hello world!"`, "Operand must be a number."},
		{`-true`, "Operand must be a number."},
		{`-false`, "Operand must be a number."},
		{`-("foo" + "bar")	`, "Operand must be a number."},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testErrorObject(t, evaluated, tt.expected)
	}
}
