package ast_test

import (
	"testing"

	"github.com/radeqq007/sunbird/internal/ast"
	"github.com/radeqq007/sunbird/internal/token"
)

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		program  *ast.Program
		expected string
	}{
		{
			name: "const declaration",
			program: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token: token.Token{Type: token.DoubleColon, Literal: "::"},
						Expression: &ast.DeclarationExpression{
							Token:   token.Token{Type: token.DoubleColon, Literal: "::"},
							Name: &ast.Identifier{
								Token: token.Token{Type: token.Ident, Literal: "myVar"},
								Value: "myVar",
							},
							IsConst: true,
							Value: &ast.Identifier{
								Token: token.Token{Type: token.Ident, Literal: "anotherVar"},
								Value: "anotherVar",
							},
						},
					},
				},
			},
			expected: "myVar :: anotherVar;",
		},
		{
			name: "mutable declaration",
			program: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token: token.Token{Type: token.ColonAssign, Literal: ":="},
						Expression: &ast.DeclarationExpression{
							Token:   token.Token{Type: token.ColonAssign, Literal: ":="},
							Name: &ast.Identifier{
								Token: token.Token{Type: token.Ident, Literal: "myVar"},
								Value: "myVar",
							},
							IsConst: false,
							Value: &ast.Identifier{
								Token: token.Token{Type: token.Ident, Literal: "anotherVar"},
								Value: "anotherVar",
							},
						},
					},
				},
			},
			expected: "myVar := anotherVar;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.program.String(); got != tt.expected {
				t.Errorf("program.String() wrong. got=%q, want=%q", got, tt.expected)
			}
		})
	}
}
