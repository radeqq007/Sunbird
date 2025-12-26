package parser

import (
	"sunbird/internal/ast"
	"sunbird/internal/token"
)

func (p *Parser) parsePropertyExpression(left ast.Expression) ast.Expression {
	exp := &ast.PropertyExpression{
		Token:  p.curToken,
		Object: left,
	}

	if !p.expectPeek(token.Ident) {
		return nil
	}

	exp.Property = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	return exp
}
