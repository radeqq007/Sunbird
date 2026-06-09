package ast_test

import (
	"sunbird/internal/ast"
	"sunbird/internal/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.ExpressionStatement{
				Token: token.Token{Type: token.Let, Literal: "let"},
				Expression: &ast.LetExpression{
					Token: token.Token{Type: token.Let, Literal: "let"},
					Name: &ast.Identifier{
						Token: token.Token{Type: token.Ident, Literal: "myVar"},
						Value: "myVar",
					},
					Value: &ast.Identifier{
						Token: token.Token{Type: token.Ident, Literal: "anotherVar"},
						Value: "anotherVar",
					},
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
