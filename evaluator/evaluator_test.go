package evaluator

import (
	"bytes"
	"io"
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/lexer"
	"github.com/codecrafters-io/interpreter-starter-go/object"
	"github.com/codecrafters-io/interpreter-starter-go/parser"
)

func testEval(input string, stdout, stderr io.Writer) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	e := NewEvaluator(stdout, stderr)
	return e.Eval(program, env)
}

func testStdout(t *testing.T, stdout bytes.Buffer, expected string) bool {
	result := stdout.String()
	if result != expected {
		t.Errorf("expected stdout %v, got %v", expected, result)
		return false
	}
	return true
}

func testStderr(t *testing.T, stderr bytes.Buffer, expected string) bool {
	result := stderr.String()
	if result != expected {
		t.Errorf("expected stderr %v, got %v", expected, result)
		return false
	}
	return true
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

func testPrintObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.Print)
	if !ok {
		t.Errorf("object is not Print. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value.Inspect() != expected {
		t.Errorf("object has wrong value. got=%q, want=%q", result.Value.Inspect(), expected)
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
	var stdout, stderr bytes.Buffer

	for _, tt := range tests {
		evaluated := testEval(tt.input, &stdout, &stderr)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalNil(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval("nil", &stdout, &stderr)
	testNilObject(t, evaluated)
}

func TestEvalString(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval(`"hello world!"`, &stdout, &stderr)
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
		var stdout, stderr bytes.Buffer
		evaluated := testEval(tt.input, &stdout, &stderr)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestEvalGroupExpression(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval("(10.4)", &stdout, &stderr)
	testNumberObject(t, evaluated, 10.4)

	evaluated = testEval("(true)", &stdout, &stderr)
	testBooleanObject(t, evaluated, true)

	evaluated = testEval("(nil)", &stdout, &stderr)
	testNilObject(t, evaluated)

	evaluated = testEval(`("hello world!")`, &stdout, &stderr)
	testStringObject(t, evaluated, "hello world!")
}

func TestUnaryExpression(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval("-10.4", &stdout, &stderr)
	testNumberObject(t, evaluated, -10.4)

	evaluated = testEval("!true", &stdout, &stderr)
	testBooleanObject(t, evaluated, false)

	evaluated = testEval("!false", &stdout, &stderr)
	testBooleanObject(t, evaluated, true)

	evaluated = testEval("!nil", &stdout, &stderr)
	testBooleanObject(t, evaluated, true)

	evaluated = testEval("-(-10.4)", &stdout, &stderr)
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
		var stdout, stderr bytes.Buffer
		evaluated := testEval(tt.input, &stdout, &stderr)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestStringConcatenation(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval(`"hello" + " " + "world"`, &stdout, &stderr)
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
		var stdout, stderr bytes.Buffer
		evaluated := testEval(tt.input, &stdout, &stderr)
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
		{`false == true`, false},
		{`false != true`, true},
	}

	for _, tt := range tests {
		var stdout, stderr bytes.Buffer
		evaluated := testEval(tt.input, &stdout, &stderr)
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
		{`"foo" * 42`, "Operands must be numbers."},
		{`true / 2`, "Operands must be numbers."},
		{`"foo" * "bar"`, "Operands must be numbers."},
		{`("foo" * "bar")`, "Operands must be numbers."},
		{`false / true`, "Operands must be numbers."},
	}

	for _, tt := range tests {
		var stdout, stderr bytes.Buffer
		evaluated := testEval(tt.input, &stdout, &stderr)
		testErrorObject(t, evaluated, tt.expected)
		testStderr(t, stderr, tt.expected)
	}
}

func TestPrintExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`print true`, "true"},
		{`print false`, "false"},
		{`print "hello world"`, "hello world"},
		{`print 10.4`, "10.4"},
		{`print (10.4 + 10.4)`, "20.8"},
	}

	for _, tt := range tests {
		var stdout, stderr bytes.Buffer
		evaluated := testEval(tt.input, &stdout, &stderr)
		testPrintObject(t, evaluated, tt.expected)
		testStdout(t, stdout, tt.expected+"\n")
	}
}

func TestVarStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"var a = 5; a;", 5},
		{"var a = 5 * 5; a;", 25},
		{"var a = 5; var b = a; b;", 5},
		{"var a = 5; var b = a; var c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		var stdout, stderr bytes.Buffer
		evaluated := testEval(tt.input, &stdout, &stderr)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestVarStatementsError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval("var a = 5; b;", &stdout, &stderr)
	testErrorObject(t, evaluated, "Undefined variable 'b'.")
	testStderr(t, stderr, "Undefined variable 'b'.")
}

func TestAssignStatements(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval("var a = 5; a = 10; a;", &stdout, &stderr)
	testNumberObject(t, evaluated, 10)
}

func TestBlockStatement(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval("{ var x = 10; print x; }", &stdout, &stderr)
	if evaluated != nil {
		t.Errorf("expected nil, got %T (%+v)", evaluated, evaluated)
	}
	testStdout(t, stdout, "10\n")
}

func TestBlockStatementWithScope(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval(
		`
		{
			var hello = "before";
			{
				var hello = "after";
				print hello;
			}
			print hello;
		}
		`,
		&stdout, &stderr,
	)
	if evaluated != nil {
		t.Errorf("expected nil, got %T (%+v)", evaluated, evaluated)
	}
	testStdout(t, stdout, "after\nbefore\n")
}
