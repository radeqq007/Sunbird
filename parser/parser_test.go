package parser

import (
	"fmt"
	"testing"
	"vex/ast"
	"vex/lexer"
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

func TestParsingPrefixExpressions(t *testing.T) {
  prefixTests := []struct {
    input    string
    operator string
    intValue    int64
    } {
      {"!5", "!",  5},
      {"-14", "-", 14},
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

    if !testIntegerLiteral(t, exp.Right, tt.intValue) {
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

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input			 string
		leftValue  int64
		operator 	 string 
		rightValue string
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"5 >= 5;", 5, ">=", 5},
		{"5 <= 5;", 5, "<=", 5}
	}
	
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements)) 
		}

		stmt, ok := program.Statements[0].(*ast.InfixExpression)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression))
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
			tt.operator, exp.Operator)
		}
	
		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}
