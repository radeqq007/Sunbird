package parser

import (
	"testing"
	"vex-programming-language/ast"
	"vex-programming-language/lexer"
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

func testReturnStatement(t *testing.T) {
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