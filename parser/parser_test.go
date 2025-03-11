package parser

import (
	"fmt"
	"sunbird/ast"
	"sunbird/lexer"
	"testing"
)

func TestVarStatement(t *testing.T) {
  input := `
var x = 5;
var y = 10;
var foobar = 32.1;
  `

  l := lexer.New(input)
  p := New(l)

  program := p.ParseProgram()
	checkParserErrors(t, p)
  if program == nil {
    t.Fatalf("p.ParseProgram() returned nil")
  }

  if len(program.Statements) != 3 {
    t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
  }

  tests := []struct {
    expectedIdentifier string
  } {
    {"x"},
    {"y"},
    {"foobar"},
  }

  for i, tt := range tests {
    stmt := program.Statements[i]
    if !testVarStatement(t, stmt, tt.expectedIdentifier) {
      return
    }
  }
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
  if s.TokenLiteral() != "var" {
    t.Errorf("s.TokenLiteral is not 'var'. got=%q", s.TokenLiteral())
    return false
  }

  varStmt, ok := s.(*ast.VarStatement)
  if !ok {
    t.Errorf("s not *ast.VarStatement. got=%T", s)
    return false
  }

  if varStmt.Name.Value != name {
    t.Errorf("varStmt.Name.Value not '%s'. got=%s", name, varStmt.Name.Value)
    return false
  }
  
  if varStmt.Name.TokenLiteral() != name {
    t.Errorf("s.Name not '%s'. got=%s", name, varStmt.Name)
    return false
  }
  
  return true
}

func TestAssignStatement(t *testing.T) {
  input := "x = 4;"

  l := lexer.New(input)
  p := New(l)
  
  program := p.ParseProgram()
  checkParserErrors(t, p)

  if len(program.Statements) != 1 {
    t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
  }
  
  stmt := program.Statements[0]

  assignStmt, ok := stmt.(*ast.AssignStatement)
  if !ok {
    t.Errorf("stmt not *ast.AssignStmt. got=%T", stmt)
  }

  if assignStmt.TokenLiteral() != "x" {
      t.Errorf("returnStmt.TokenLiteral not 'x', got %q", assignStmt.TokenLiteral())
  }

  if assignStmt.Name.Value != "x" {
    t.Errorf("varStmt.Name.Value not '%s'. got=%s", "x", assignStmt.Name.Value)
  }
  
  if assignStmt.Name.TokenLiteral() != "x" {
    t.Errorf("assignStmt.Name not '%s'. got=%s", "x", assignStmt.Name)
  }

  if assignStmt.Value.String() != "4" {
    t.Errorf("assignStmt.Value.String not '%s'. got=%s", "4", assignStmt.Value.String())
  }

}

func TestReturnStatement(t *testing.T) {
  input := `
return 6;
return 1.24;
return 993322;
`
  l := lexer.New(input)
  p := New(l)

  program := p.ParseProgram()
  checkParserErrors(t, p)

  if len(program.Statements) != 3 {
    t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
  }

  for _, stmt := range program.Statements {
    returnStmt, ok := stmt.(*ast.ReturnStatement)
    if !ok {
      t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
      continue
    }

    if returnStmt.TokenLiteral() != "return" {
        t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
        returnStmt.TokenLiteral())
    }
  }
}

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

func TestIdentifierExpression(t *testing.T) {
  input := "foobar;"

  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()
  checkParserErrors(t, p)

  if len(program.Statements) != 1 {
    t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
  }

  stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
  if !ok {
    t.Fatalf("program.Statements[0] is npt ast.ExpressionStatement. got=%T", program.Statements[0])
  }

  ident, ok := stmt.Expression.(*ast.Identifier)
  if !ok {
    t.Errorf("ident.Value not %s. got=%s ", "foobar", ident.Value)
  }

  if ident.TokenLiteral() != "foobar" {
    t.Errorf("ident.TokenLiteral() not %s. got=%s ", "foobar", ident.TokenLiteral())
  }
}

func TestIntegerLiteralExpression(t *testing.T) {
  input := "5;"
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

  literal, ok := stmt.Expression.(*ast.IntegerLiteral)
  
  if !ok {
    t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
  }
  
  if literal.Value != 5 {
    t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
  }

  if literal.TokenLiteral() != "5" {
    t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
    literal.TokenLiteral())
  }
}


func TestFloatLiteralExpression(t *testing.T) {
  input := "2.4;"
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

  literal, ok := stmt.Expression.(*ast.FloatLiteral)
  
  if !ok {
    t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
  }
  
  if literal.Value != 2.4 {
    t.Errorf("literal.Value not %v. got=%v", 2.4, literal.Value)
  }

  if literal.TokenLiteral() != "2.4" {
    t.Errorf("literal.TokenLiteral not %s. got=%s", "2.4", literal.TokenLiteral())
  }
}

func TestStringLiteralExpression(t *testing.T) {
  input := `"hello world";`
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
    t.Errorf("literal.Value not %s. got=%s", "hello world", literal.Value)
  }

  if literal.TokenLiteral() != "hello world" {
    t.Errorf("literal.TokenLiteral not %s. got=%s", "hello world", literal.TokenLiteral())
  }
}

func TestBooleanLiteralExpression(t *testing.T) {
  input := "true;"
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

  boolean, ok := stmt.Expression.(*ast.Boolean)
  
  if !ok {
    t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
  }
  
  if boolean.Value != true {
    t.Errorf("boolean.Value not %v. got=%v", true, boolean.Value)
  }

  if boolean.TokenLiteral() != "true" {
    t.Errorf("boolean.TokenLiteral not %s. got=%s", "true", boolean.TokenLiteral())
  }
}

func TestParsingPrefixExpressions(t *testing.T) {
  prefixTests := []struct {
    input    string
    operator string
    value    interface{}
    } {
      {"!5", "!",  5},
      {"-14", "-", 14},
      {"!true", "!", true},
      {"!false", "!", false},
  }

  for _, tt := range prefixTests {
    l := lexer.New(tt.input)
    p := New(l)
    program := p.ParseProgram()
    checkParserErrors(t, p)

    if len(program.Statements) != 1 {
      t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

    if !ok {
      t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    exp, ok := stmt.Expression.(*ast.PrefixExpression)
    
    if !ok {
      t.Fatalf("exp.Operator")
    }

    if exp.Operator != tt.operator {
      t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
    }

    if !testLiteralExpression(t, exp.Right, tt.value) {
      return
    }

    
  }
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
  integ, ok := il.(*ast.IntegerLiteral)
  if !ok {
    t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
    return false
  }

  if integ.Value != value {
    t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
    return false
  }

  if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
    t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
    return false
  }

  return true
}

func testFloatLiteral(t *testing.T, fl ast.Expression, value float64) bool {
  float, ok := fl.(*ast.FloatLiteral)
  if !ok {
    t.Errorf("il not *ast.FloatLiteral. got=%T", fl)
    return false
  }

  if float.Value != value {
    t.Errorf("float.Value not %v. got=%v", value, float.Value)
    return false
  }

  if float.TokenLiteral() != fmt.Sprintf("%v", value) {
    t.Errorf("float.TokenLiteral not %v. got=%s", value, float.TokenLiteral())
    return false
  }

  return true
}

func testStringLiteral(t *testing.T, sl ast.Expression, value string) bool {
  str, ok := sl.(*ast.StringLiteral)
  if !ok {
    t.Errorf("il not *ast.StringLiteral. got=%T", sl)
    return false
  }

  if str.Value != value {
    t.Errorf("str.Value not %s. got=%s", value, str.Value)
    return false
  }

  if str.TokenLiteral() != value {
    t.Errorf("str.TokenLiteral not %s. got=%s", value, str.TokenLiteral())
    return false
  }

  return true
}

func testBooleanLiteral(t *testing.T, bl ast.Expression, value bool) bool {
  bo, ok := bl.(*ast.Boolean)
  if !ok {
    t.Errorf("il not *ast.Boolean. got=%T", bl)
    return false
  }

  if bo.Value != value {
    t.Errorf("bo.Value not %v. got=%v", value, bo.Value)
    return false
  }

  if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
    t.Errorf("bo.TokenLiteral not %t. got=%s", value, bo.TokenLiteral())
    return false
  }

  return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input			 string
		leftValue  interface {}
		operator 	 string 
		rightValue interface {}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"3.14 != 4;", 3.14, "!=", 4},
		{"5.21 >= 5.21;", 5.21, ">=", 5.21},
		{"5 <= 5;", 5, "<=", 5},
		{"true || false;", true, "||", false},
		{"true && false;", true, "&&", false},
    {"true == true", true, "==", true},
    {"true != false", true, "!=", false},
    {"false == false", false, "==", false},
    {"\"hello\" + \"world\"", "hello", "+", "world"},
    {"foo + bar", "foo", "+", "bar"},
	}
	
	for _, tt := range infixTests {
    l := lexer.New(tt.input)
    p := New(l)
    program := p.ParseProgram()
   
    checkParserErrors(t, p)
   
    if len(program.Statements) != 1 {
     t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
      1, len(program.Statements))
    }
    
    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    
    if !ok {
      t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
      program.Statements[0])
    }
    
    exp, ok := stmt.Expression.(*ast.InfixExpression)
    
    if !ok {
      t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
    }

    if !testLiteralExpression(t, exp.Left, tt.leftValue) {
      return
    }

    if exp.Operator != tt.operator {
      t.Fatalf("exp.Operator is not '%s'. got=%s",
      tt.operator, exp.Operator)
    }

    if !testLiteralExpression(t, exp.Right, tt.rightValue) {
      return
    }
    }
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	} {
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
    {"true", "true"},
    {"false", "false"},
    {"3 > 5 == false", "((3 > 5) == false)"},
    {"3 < 5 == true", "((3 < 5) == true)"},
    {"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
    {"(5 + 5) * 2", "((5 + 5) * 2)"},
    {"2 / (5 + 5)", "(2 / (5 + 5))"},
    {"-(5 + 5)", "(-(5 + 5))"},
    {"!(true == true)", "(!(true == true))"},
    {
      "a + add(b * c) + d",
      "((a + add((b * c))) + d)",
    },
    {
      "add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
      "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
    },
    {
      "add(a + b + c * d / f + g)",
      "add((((a + b) + ((c * d) / f)) + g))",
    },      
  }

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}


func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
  ident, ok := exp.(*ast.Identifier)
  if !ok {
    t.Errorf("exp not *ast.Identifier. got=%T", exp)
    return false
  }

  if ident.Value != value {
    t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
    return false
  }

  if ident.TokenLiteral() != value {
    t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
    return false
  }

  return true 
}

func testLiteralExpression(
  t *testing.T,
  exp ast.Expression,
  expected interface{},
) bool {
  switch v := expected.(type) {
    case int:
      return testIntegerLiteral(t, exp, int64(v))
    case int64:
      return testIntegerLiteral(t, exp, v)
    
    case float64:
      return testFloatLiteral(t, exp, v)
    
    case string:
      if _, ok := exp.(*ast.Identifier); ok {
        return testIdentifier(t, exp, v)
      }
      return testStringLiteral(t, exp, v)
    
    case bool:
      return testBooleanLiteral(t, exp, v)

    }

    t.Errorf("type of exp not handled. got=%T", exp)
    return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
  opExp, ok := exp.(*ast.InfixExpression)
  if !ok {
    t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
    return false
  }

  if !testLiteralExpression(t, opExp.Left, left) {
    return false
  }

  if opExp.Operator != operator {
    t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
    return false
  }

  if !testLiteralExpression(t, opExp.Right, right) {
    return false
  }

  return true
}

func TestIfExpression(t *testing.T) {
	input := "if x < 0 { x }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", 0) {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
  input := "if x < 0 { x } else { y }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", 0) {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("alternative is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}

  alt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
  if !ok {
		t.Fatalf("exp.Alternative.Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

  if !testIdentifier(t, alt.Expression, "y") {
    return;
  }
}

// func TestIfElseExpression(t *testing.T) {
//   input := "if x < 0 { x } else if x > 0 { y } else { 0 }"

//   l := lexer.New(input)
//   p := New(l)
//   program := p.ParseProgram()
//   checkParserErrors(t, p)

//   if len(program.Statements) != 1 {
//     t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
//   }

//   stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
//   if !ok {
//     t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
//   }

//   exp, ok := stmt.Expression.(*ast.IfExpression)
//   if !ok {
//     t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
//   }

//   if !testInfixExpression(t, exp.Condition, "x", "<", 0) {
//     return
//   }

//   if len(exp.Consequence.Statements) != 1 {
//     t.Errorf("consequence is not 1 statement. got=%d", len(exp.Consequence.Statements))
//   }

//   consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
//   if !ok {
//     t.Fatalf("exp.Consequence.Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
//   }

//   if !testIdentifier(t, consequence.Expression, "x") {
//     return
//   }

//   altExp, ok := exp.Alternative.Statements[0].(*ast.IfExpression)
//   if !ok {
//     t.Fatalf("exp.Alternative.Statements[0] is not ast.IfExpression. got=%T", exp.Alternative.Statements[0])
//   }

//   if !testInfixExpression(t, altExp.Condition, "x", ">", 0) {
//     return
//   }

//   if len(altExp.Consequence.Statements) != 1 {
//     t.Errorf("alternative consequence is not 1 statement. got=%d", len(altExp.Consequence.Statements))
//   }

//   altConsequence, ok := altExp.Consequence.Statements[0].(*ast.ExpressionStatement)
//   if !ok {
//     t.Fatalf("altExp.Consequence.Statements[0] is not ast.ExpressionStatement. got=%T", altExp.Consequence.Statements[0])
//   }

//   if !testIdentifier(t, altConsequence.Expression, "y") {
//     return
//   }

//   finalAlt, ok := altExp.Alternative.(*ast.BlockStatement)
//   if !ok {
//     t.Fatalf("altExp.Alternative is not ast.BlockStatement. got=%T", altExp.Alternative)
//   }

//   if len(finalAlt.Statements) != 1 {
//     t.Errorf("final alternative is not 1 statement. got=%d", len(finalAlt.Statements))
//   }

//   finalStmt, ok := finalAlt.Statements[0].(*ast.ExpressionStatement)
//   if !ok {
//     t.Fatalf("finalAlt.Statements[0] is not ast.ExpressionStatement. got=%T", finalAlt.Statements[0])
//   }

//   if !testIntegerLiteral(t, finalStmt.Expression, 0) {
//     return
//   }
// }

func TestFunctionLiteralParsing(t *testing.T) {
  input := "func (x, y) { x + y; }"

  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()
  checkParserErrors(t, p)

  if len(program.Statements) != 1 {
    t.Fatalf("program.Body does not contain %d statements. got=%d\n", 1, len(program.Statements))
  }

  stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
  if !ok {
    t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
  }

  function, ok := stmt.Expression.(*ast.FunctionLiteral)
  if !ok {
    t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
  }

  if len(function.Parameters) != 2 {
    t.Fatalf("function literal parameters wrong. want 2, got=%d\n", len(function.Parameters))
  }

  testLiteralExpression(t, function.Parameters[0], "x")
  testLiteralExpression(t, function.Parameters[1], "y")
  
  if len(function.Body.Statements) != 1 {
    t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
    len(function.Body.Statements))
  }

  bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
  if !ok {
    t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
  }

  testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
  tests := []struct {
    input          string
    expectedParams []string
  } {
    { input: "func() {};", expectedParams: []string{}},
    { input: "func(a) {};", expectedParams: []string{"a"}},
    { input: "func(a, b, c) {};", expectedParams: []string{"a", "b", "c"}},
  }

  for _, tt := range tests {
    l := lexer.New(tt.input)
    p := New(l)
    program := p.ParseProgram()
    checkParserErrors(t, p)

    stmt := program.Statements[0].(*ast.ExpressionStatement)
    function := stmt.Expression.(*ast.FunctionLiteral)

    if len(function.Parameters) != len(tt.expectedParams) {
      t.Errorf("length parameters wrong. want %d, got=%d\n", len(tt.expectedParams), len(function.Parameters))
    }

    for i, ident := range tt.expectedParams {
      testLiteralExpression(t, function.Parameters[i], ident)
    }
    
  }
}

func TestCallExpressionParsing(t *testing.T) {
  input := "add(1, 2 * 3, 4 + 5);"
  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()
  checkParserErrors(t, p)

  if len(program.Statements) != 1 {
    t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
  }

  stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
  if !ok {
    t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
    program.Statements[0])
  }

  exp, ok := stmt.Expression.(*ast.CallExpression)
  if !ok {
    t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
  }

  if !testIdentifier(t, exp.Function, "add") {
    return
  }

  if len(exp.Arguments) != 3 {
    t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
  }

  testLiteralExpression(t, exp.Arguments[0], 1)
  testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
  testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestParsingArrayLiterals(t *testing.T) {
  input := "[1, 2 * 2, 3 + 3]"

  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()
  checkParserErrors(t, p)

  stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
  if !ok {
    t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
  }
  
  array, ok := stmt.Expression.(*ast.ArrayLiteral)
  if !ok {
    t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
  }

  if len(array.Elements) != 3 {
    t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
  }

  testIntegerLiteral(t, array.Elements[0], 1)
  testInfixExpression(t, array.Elements[1], 2, "*", 2)
  testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
  input := "coolArray[1 + 1]"

  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()
  checkParserErrors(t, p)

  stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
  if !ok {
    t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
  }

  indexExp, ok := stmt.Expression.(*ast.IndexExpression)
  if !ok {
    t.Fatalf("exp not *ast.IndexExpression, got=%T", stmt.Expression)
  }

  if !testIdentifier(t, indexExp.Left, "coolArray") {
    return
  }

  if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
    return
  }
}
