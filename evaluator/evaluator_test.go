package evaluator

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/lexer"
	"github.com/codecrafters-io/interpreter-starter-go/object"
	"github.com/codecrafters-io/interpreter-starter-go/parser"
)

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testEval(t *testing.T, input string, stdout, stderr io.Writer) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	checkParserErrors(t, p)
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
		evaluated := testEval(t, tt.input, &stdout, &stderr)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalNil(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval(t, "nil", &stdout, &stderr)
	testNilObject(t, evaluated)
}

func TestEvalString(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval(t, `"hello world!"`, &stdout, &stderr)
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
		evaluated := testEval(t, tt.input, &stdout, &stderr)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestEvalGroupExpression(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval(t, "(10.4)", &stdout, &stderr)
	testNumberObject(t, evaluated, 10.4)

	evaluated = testEval(t, "(true)", &stdout, &stderr)
	testBooleanObject(t, evaluated, true)

	evaluated = testEval(t, "(nil)", &stdout, &stderr)
	testNilObject(t, evaluated)

	evaluated = testEval(t, `("hello world!")`, &stdout, &stderr)
	testStringObject(t, evaluated, "hello world!")
}

func TestUnaryExpression(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval(t, "-10.4", &stdout, &stderr)
	testNumberObject(t, evaluated, -10.4)

	evaluated = testEval(t, "!true", &stdout, &stderr)
	testBooleanObject(t, evaluated, false)

	evaluated = testEval(t, "!false", &stdout, &stderr)
	testBooleanObject(t, evaluated, true)

	evaluated = testEval(t, "!nil", &stdout, &stderr)
	testBooleanObject(t, evaluated, true)

	evaluated = testEval(t, "-(-10.4)", &stdout, &stderr)
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
		evaluated := testEval(t, tt.input, &stdout, &stderr)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestStringConcatenation(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval(t, `"hello" + " " + "world"`, &stdout, &stderr)
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
		evaluated := testEval(t, tt.input, &stdout, &stderr)
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
		evaluated := testEval(t, tt.input, &stdout, &stderr)
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
		evaluated := testEval(t, tt.input, &stdout, &stderr)
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
		evaluated := testEval(t, tt.input, &stdout, &stderr)
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
		evaluated := testEval(t, tt.input, &stdout, &stderr)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestVarStatementsError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval(t, "var a = 5; b;", &stdout, &stderr)
	testErrorObject(t, evaluated, "undefined variable: b")
	testStderr(t, stderr, "undefined variable: b")
}

func TestAssignStatements(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval(t, "var a = 5; a = 10; a;", &stdout, &stderr)
	testNumberObject(t, evaluated, 10)
}

func TestAssignByEquality(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval(t, "var age = 50; var condition = age >= 18; condition;", &stdout, &stderr)
	testBooleanObject(t, evaluated, true)
}

func TestBlockStatement(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval(t, "{ var x = 10; print x; }", &stdout, &stderr)
	if evaluated != nil {
		t.Errorf("expected nil, got %T (%+v)", evaluated, evaluated)
	}
	testStdout(t, stdout, "10\n")
}

func TestBlockStatementWithScope(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval(t,
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

func TestIfCondition(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`if (true) print "bar";`, "bar\n"},
		{`if (false) print "bar";`, ""},
		{`if (true) { print "block body"; }`, "block body\n"},
		{`var a = false; if (a = true) { print (a == true); }`, "true\n"},
		{`
				if (true) { print "eligible for voting: true"; }
				if (false) { print "eligible for voting: false"; }`, "eligible for voting: true\n"},
		{`var stage = "unknown";
		var age = 50;
		if (age < 18) { stage = "child"; }
		if (age >= 18) { stage = "adult"; }
		print stage;

		var isAdult = age >= 18;
		if (isAdult) { print "eligible for voting: true"; }
		if (!isAdult) { print "eligible for voting: false"; }`, "adult\neligible for voting: true\n"},
	}

	for _, tt := range tests {
		var stdout, stderr bytes.Buffer
		testEval(t, tt.input, &stdout, &stderr)
		testStdout(t, stdout, tt.expected)
	}
}

func TestIfElseCondition(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`if (true) print "if branch"; else print "else branch";`, "if branch\n"},
		{`
						var age = 21;
						if (age > 18) print "adult"; else print "child";`, "adult\n"},
		{`if (false) { print "if block"; } else print "else statement";
	if (false) print "if statement"; else {
  	print "else block";
	}`, "else statement\nelse block\n"},
		{
			`var celsius = 67;
		var fahrenheit = 0;
		var isHot = false;

		{
		  fahrenheit = celsius * 9 / 5 + 32;
		  print celsius; print fahrenheit;

		  if (celsius > 30) {
		    isHot = true;
		    print "It's a hot day. Stay hydrated!";
		  } else {
		    print "It's cold today. Wear a jacket!";
		  }

		  if (isHot) { print "Remember to use sunscreen!"; }
				}`, "67\n152.6\nIt's a hot day. Stay hydrated!\nRemember to use sunscreen!\n"},
	}

	for _, tt := range tests {
		var stdout, stderr bytes.Buffer
		testEval(t, tt.input, &stdout, &stderr)
		testStdout(t, stdout, tt.expected)
	}
}

func TestIfElseIfStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`if (true) print "if branch"; else if (false) print "else-if branch";`, "if branch\n"},
		{`if (true) {
				print "hello";
			} else if (true) print "hello";

			if (true) print "hello"; else if (true) {
				print "hello";
			}`, "hello\nhello\n"},

		{`
			var age = 88;
			var stage = "unknown";
			if (age < 18) { stage = "child"; }
			else if (age >= 18) { stage = "adult"; }
			else if (age >= 65) { stage = "senior"; }
			else if (age >= 100) { stage = "centenarian"; }
			print stage;
			`, "adult\n"},
		{
			`var age = 67;
			var isAdult = age >= 18;
			if (isAdult) { print "eligible for voting: true"; }
			else { print "eligible for voting: false"; }

			if (age < 16) { print "eligible for driving: false"; }
			else if (age < 18) { print "eligible for driving: learner's permit"; }
			else { print "eligible for driving: full license"; }

			if (age < 21) { print "eligible for drinking (US): false"; }
			else { print "eligible for drinking (US): true"; }`, "eligible for voting: true\neligible for driving: full license\neligible for drinking (US): true\n"},
	}

	for _, tt := range tests {
		var stdout, stderr bytes.Buffer
		testEval(t, tt.input, &stdout, &stderr)
		testStdout(t, stdout, tt.expected)
	}
}

func TestOrExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`if (false or "ok") print "baz";
if (nil or "ok") print "baz";

if (false or false) print "world";
if (true or "world") print "world";

if (24 or "bar") print "bar";
if ("bar" or "bar") print "bar";`, "baz\nbaz\nworld\nbar\nbar\n"},

		{`
			print 41 or true;
print false or 41;
print false or false or true;

print false or false;
print false or false or false;
print true or true or true or true;
			`, "41\n41\ntrue\nfalse\nfalse\ntrue\n",
		}, {
			`var stage = "unknown";
var age = 23;
if (age < 18) { stage = "child"; }
if (age >= 18) { stage = "adult"; }
print stage;

var isAdult = age >= 18;
if (isAdult) { print "eligible for voting: true"; }
if (!isAdult) { print "eligible for voting: false"; }`, "adult\neligible for voting: true\n"},
	}

	for _, tt := range tests {
		var stdout, stderr bytes.Buffer
		testEval(t, tt.input, &stdout, &stderr)
		testStdout(t, stdout, tt.expected)
	}
}

func TestAndExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`if (false and "bad") print "foo";`, ""},
		{`if (nil and "bad") print "foo";`, ""},
		{`if (true and "hello") print "hello";`, "hello\n"},
		{`if (97 and "baz") print "baz";`, "baz\n"},
		{`if ("baz" and "baz") print "baz";`, "baz\n"},
		{`if ("" and "bar") print "bar";`, "bar\n"},

		{`print false and 1;`, "false\n"},
		{`print true and 1;`, "1\n"},
		{`print 23 and "hello" and false;`, "false\n"},
		{`print 23 and true;`, "true\n"},
		{`print 23 and "hello" and 23;`, "23\n"},
	}

	for _, tt := range tests {
		var stdout, stderr bytes.Buffer
		testEval(t, tt.input, &stdout, &stderr)
		testStdout(t, stdout, tt.expected)
	}
}

func TestWhileStatement(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval(t, `var baz = 0; while (baz < 3) print baz = baz + 1;`, &stdout, &stderr)
	if evaluated != nil {
		t.Errorf("expected nil, got %T (%+v)", evaluated, evaluated)
	}
	testStdout(t, stdout, "1\n2\n3\n")
}

func TestForStatement(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval(t, `for (var baz = 0; baz < 3; baz = baz + 1) print baz;`, &stdout, &stderr)
	if evaluated != nil {
		t.Errorf("expected nil, got %T (%+v)", evaluated, evaluated)
	}
	testStdout(t, stdout, "0\n1\n2\n")
}

func TestForStatementWithoutIncrement(t *testing.T) {
	var stdout, stderr bytes.Buffer
	evaluated := testEval(t, `for (var baz = 0; baz < 3;) print baz = baz + 1;`, &stdout, &stderr)
	if evaluated != nil {
		t.Errorf("expected nil, got %T (%+v)", evaluated, evaluated)
	}
	testStdout(t, stdout, "1\n2\n3\n")
}

func TestClockFunction(t *testing.T) {
	var stdout, stderr bytes.Buffer
	input := `print clock() + 10;`

	testEval(t, input, &stdout, &stderr)

	// Since clock returns current time, we can only verify the output is a number
	output := stdout.String()
	_, err := strconv.ParseFloat(strings.TrimSpace(output), 64)
	if err != nil {
		t.Errorf("clock() did not return a valid number: %s", output)
	}
}
