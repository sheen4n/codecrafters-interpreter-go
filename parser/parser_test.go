package parser

import (
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/ast"
	"github.com/codecrafters-io/interpreter-starter-go/lexer"
)

func checkParserErrors(t *testing.T, p *Parser) {
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

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		literal, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}
		if literal.Value != tt.expectedBoolean {
			t.Errorf("literal.Value not %t. got=%t", tt.expectedBoolean, literal.Value)
		}
	}
}

func TestNilExpression(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"nil"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		_, ok = stmt.Expression.(*ast.Nil)
		if !ok {
			t.Fatalf("exp not *ast.Nil. got=%T", stmt.Expression)
		}
	}
}

func TestNumberLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"42.47", 42.47},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		_, ok = stmt.Expression.(*ast.NumberLiteral)
		if !ok {
			t.Fatalf("exp not *ast.Nil. got=%T", stmt.Expression)
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestGroupExpression(t *testing.T) {
	input := `("foo")`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	group, ok := stmt.Expression.(*ast.GroupExpression)
	if !ok {
		t.Fatalf("exp not *ast.GroupExpression. got=%T", stmt.Expression)
	}

	if group.String() != "(group foo)" {
		t.Errorf("group.String() not %q. got=%q", "(group foo)", group.String())
	}
}

func TestUnaryExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-true", "(- true)"},
		{"-42.47", "(- 42.47)"},
		{"!true", "(! true)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		unary, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("exp not *ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if unary.String() != tt.expected {
			t.Errorf("unary.String() not %q. got=%q", tt.expected, unary.String())
		}
	}
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1 + 2", "(+ 1.0 2.0)"},
		{"1 - 2", "(- 1.0 2.0)"},
		{"16 * 38 / 58", "(/ (* 16.0 38.0) 58.0)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		infix, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp not *ast.InfixExpression. got=%T", stmt.Expression)
		}

		if infix.String() != tt.expected {
			t.Errorf("infix.String() not %q. got=%q", tt.expected, infix.String())
		}
	}
}

func TestComparisonExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1 < 2", "(< 1.0 2.0)"},
		{"1 > 2", "(> 1.0 2.0)"},
		{"1 <= 2", "(<= 1.0 2.0)"},
		{"1 >= 2", "(>= 1.0 2.0)"},
		{"83 < 99 < 115", "(< (< 83.0 99.0) 115.0)"},
		{`"baz" == "baz"`, `(== baz baz)`},
		{`"foo" != "bar"`, `(!= foo bar)`},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}
	}
}

func TestPrintExpression(t *testing.T) {
	input := `print "hello world"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	println, ok := stmt.Expression.(*ast.PrintExpression)
	if !ok {
		t.Fatalf("exp not *ast.PrintExpression. got=%T", stmt.Expression)
	}

	if println.String() != "(print hello world)" {
		t.Errorf("println.String() not %q. got=%q", "(print hello world)", println.String())
	}
}

func TestMultiplePrintExpressions(t *testing.T) {
	input := `print "hello world"; print "hello world"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		println, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt not *ast.PrintExpression. got=%T", stmt)
		}

		printExpression, ok := println.Expression.(*ast.PrintExpression)
		if !ok {
			t.Fatalf("stmt not *ast.PrintExpression. got=%T", stmt)
		}

		if printExpression.String() != "(print hello world)" {
			t.Errorf("println.String() not %q. got=%q", "(print hello world)", println.String())
		}
	}
}

func TestVarStatement(t *testing.T) {
	input := `var x = 10`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.VarStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.VarStatement. got=%T", program.Statements[0])
	}

	if stmt.Name.String() != "x" {
		t.Errorf("stmt.Name.String() not %q. got=%q", "x", stmt.Name.String())
	}

	numberLiteral, ok := stmt.Value.(*ast.NumberLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not *ast.NumberLiteral. got=%T", stmt.Value)
	}

	if numberLiteral.Value != 10 {
		t.Errorf("numberLiteral.Value not %f. got=%f", 10.0, numberLiteral.Value)
	}
}

func TestNilVarStatement(t *testing.T) {
	input := `var x;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.VarStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.VarStatement. got=%T", program.Statements[0])
	}

	if stmt.Name.String() != "x" {
		t.Errorf("stmt.Name.String() not %q. got=%q", "x", stmt.Name.String())
	}

	nilLiteral, ok := stmt.Value.(*ast.Nil)
	if !ok {
		t.Fatalf("stmt.Value is not *ast.Nil. got=%T", stmt.Value)
	}

	if nilLiteral.String() != "nil" {
		t.Errorf("nilLiteral.String() not %q. got=%q", "nil", nilLiteral.String())
	}
}

func TestAssignExpression(t *testing.T) {
	input := `var x = 10; x = 20;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[1] is not ast.ExpressionStatement. got=%T", program.Statements[1])
	}

	assignStmt, ok := stmt.Expression.(*ast.AssignExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.AssignExpression. got=%T", stmt.Expression)
	}

	if assignStmt.Name.String() != "x" {
		t.Errorf("assignStmt.Name.String() not %q. got=%q", "x", assignStmt.Name.String())
	}
}

func TestBlockStatement(t *testing.T) {
	input := `{ var x = 10; print x; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	block, ok := program.Statements[0].(*ast.BlockStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.BlockStatement. got=%T", program.Statements[0])
	}

	if len(block.Statements) != 2 {
		t.Fatalf("block statement has not enough statements. got=%d", len(block.Statements))
	}
}

func TestBlockStatementError(t *testing.T) {
	input := `{ var x = 10; `

	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()

	errors := p.Errors()
	if len(errors) == 0 {
		t.Errorf("expected error, got none")
	}

	if errors[0] != "[line 1] Expect '}'." {
		t.Errorf("expected error %q, got %q", "[line 1] Expect '}'.", errors[0])
	}
}

func TestIfStatement(t *testing.T) {
	input := `if (true) print "bar";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T", program.Statements[0])
	}

	condition, ok := stmt.Condition.(*ast.Boolean)
	if !ok {
		t.Fatalf("stmt.Condition is not *ast.Boolean. got=%T", stmt.Condition)
	}

	if condition.Value != true {
		t.Errorf("condition.Value not %t. got=%t", true, condition.Value)
	}

	if stmt.Consequence.String() != "(print bar)" {
		t.Errorf("stmt.Consequence.String() not %q. got=%q", "(print bar)", stmt.Consequence.String())
	}
}

func TestIfBlockStatement(t *testing.T) {
	input := `if (false) {
  print "block body";
}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T", program.Statements[0])
	}

	condition, ok := stmt.Condition.(*ast.Boolean)
	if !ok {
		t.Fatalf("stmt.Condition is not *ast.Boolean. got=%T", stmt.Condition)
	}

	if condition.Value != false {
		t.Errorf("condition.Value not %t. got=%t", false, condition.Value)
	}

	if stmt.Consequence.String() != "{(print block body)}" {
		t.Errorf("stmt.Consequence.String() not %q. got=%q", "(print block body)", stmt.Consequence.String())
	}
}

func TestMultipleIfStatements(t *testing.T) {
	input := `
				if (true) {  }
				if (false) {  }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		if _, ok := stmt.(*ast.IfStatement); !ok {
			t.Fatalf("stmt not *ast.IfStatement. got=%T", stmt)
		}
	}
}

func TestSemicolonAtEndOfStatements(t *testing.T) {
	input := `
	var hello = (65 * 53) - 24;
{
    var foo = "baz" + "10";
    print foo;
}
print hello;
	
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
}

func TestMultiAssignment(t *testing.T) {
	input := `
	var world;
	var baz;
	world = baz = 84 + 33 * 60;
	print world;
	print baz;
	
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 5 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
}

func TestIfElseStatement(t *testing.T) {
	input := `if (true) { print "foo"; } else { print "bar"; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T", program.Statements[0])
	}

	if stmt.Condition.String() != "true" {
		t.Errorf("stmt.Condition.String() not %q. got=%q", "true", stmt.Condition.String())
	}

	if stmt.Consequence.String() != "{(print foo)}" {
		t.Errorf("stmt.Consequence.String() not %q. got=%q", "{(print foo)}", stmt.Consequence.String())
	}

	if stmt.Alternative.String() != "{(print bar)}" {
		t.Errorf("stmt.Alternative.String() not %q. got=%q", "{(print bar)}", stmt.Alternative.String())
	}
}

func TestIfElseStatementWithSemicolon(t *testing.T) {
	input := `if (false) { print "if block"; } else print "else statement";
	if (false) print "if statement"; else {
  	print "else block";
	}
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T", program.Statements[0])
	}

	if stmt.Condition.String() != "false" {
		t.Errorf("stmt.Condition.String() not %q. got=%q", "false", stmt.Condition.String())
	}

	if stmt.Consequence.String() != "{(print if block)}" {
		t.Errorf("stmt.Consequence.String() not %q. got=%q", "{(print if block)}", stmt.Consequence.String())
	}

	if stmt.Alternative.String() != "(print else statement)" {
		t.Errorf("stmt.Alternative.String() not %q. got=%q", "(print else statement)", stmt.Alternative.String())
	}

	stmt2, ok := program.Statements[1].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[1] is not ast.IfStatement. got=%T", program.Statements[1])
	}

	if stmt2.Condition.String() != "false" {
		t.Errorf("stmt2.Condition.String() not %q. got=%q", "false", stmt2.Condition.String())
	}

	if stmt2.Consequence.String() != "(print if statement)" {
		t.Errorf("stmt2.Consequence.String() not %q. got=%q", "(print if statement)", stmt2.Consequence.String())
	}

	if stmt2.Alternative.String() != "{(print else block)}" {
		t.Errorf("stmt2.Alternative.String() not %q. got=%q", "{(print else block)}", stmt2.Alternative.String())
	}
}

func TestWhileStatement(t *testing.T) {
	input := `while (true) print "foo";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.WhileStatement. got=%T", program.Statements[0])
	}

	if stmt.Condition.String() != "true" {
		t.Errorf("stmt.Condition.String() not %q. got=%q", "true", stmt.Condition.String())
	}

	if stmt.Consequence.String() != "(print foo)" {
		t.Errorf("stmt.Consequence.String() not %q. got=%q", "{(print foo)}", stmt.Consequence.String())
	}
}

func TestForStatementWithoutIncrement(t *testing.T) {
	input := `for (var baz = 0; baz < 3;) print baz = baz + 1;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ForStatement. got=%T", program.Statements[0])
	}

	if stmt.Init.String() != "var baz = 0.0;" {
		t.Errorf("stmt.Init.String() not %q. got=%q", "var baz = 0", stmt.Init.String())
	}

	if stmt.Condition.String() != "(< baz 3.0)" {
		t.Errorf("stmt.Condition.String() not %q. got=%q", "(< baz 3.0)", stmt.Condition.String())
	}

	if stmt.Increment != nil {
		t.Errorf("stmt.Increment not nil")
	}

	if stmt.Body.String() != "(print baz = (+ baz 1.0);)" {
		t.Errorf("stmt.Body.String() not %q. got=%q", "(print baz = (+ baz 1.0);)", stmt.Body.String())
	}
}

func TestForStatementWithIncrement(t *testing.T) {
	input := `for (var baz = 0; baz < 3; baz = baz + 1){ print baz; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ForStatement. got=%T", program.Statements[0])
	}

	if stmt.Init.String() != "var baz = 0.0;" {
		t.Errorf("stmt.Init.String() not %q. got=%q", "var baz = 0.0;", stmt.Init.String())
	}

	if stmt.Condition.String() != "(< baz 3.0)" {
		t.Errorf("stmt.Condition.String() not %q. got=%q", "(< baz 3.0)", stmt.Condition.String())
	}

	if stmt.Increment.String() != "baz = (+ baz 1.0);" {
		t.Errorf("stmt.Increment.String() not %q. got=%q", "baz = (+ baz 1.0);", stmt.Increment.String())
	}

	if stmt.Body.String() != "{(print baz)}" {
		t.Errorf("stmt.Body.String() not %q. got=%q", "{(print baz)}", stmt.Body.String())
	}
}

func TestSyntaxError(t *testing.T) {
	tests := []struct {
		input         string
		expectedError string
	}{
		{"(72 +)", "[line 1] Error at ')': Expect expression."},
		{`for (;;) var foo;`, "[line 1] var statement should be in a block."},
		{`for (var a = 1; {}; a = a + 1) {}`, "[line 1] Error at '{': Expect expression."},
		{`for (var a = 1; a < 2; {}) {}`, "[line 1] Empty increment condition."},
		{`for ({}; a < 2; a = a + 1) {}`, "[line 1] Empty initial condition."},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)
		p.ParseProgram()

		errors := p.Errors()
		if len(errors) == 0 {
			t.Errorf("expected error, got none")
			continue
		}

		if errors[0] != tt.expectedError {
			t.Errorf("expected error %q, got %q", tt.expectedError, errors[0])
		}
	}
}
