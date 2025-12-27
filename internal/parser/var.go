package parser

import (
	"sunbird/internal/ast"
	"sunbird/internal/token"
)

func (p *Parser) parseVarExpression() ast.Expression {
	exp := &ast.VarExpression{Token: p.curToken}

	if !p.expectPeek(token.Ident) {
		return nil
	}

	exp.Name = p.parseIdentifier()

	if !p.expectPeek(token.Assign) {
		return nil
	}

	p.nextToken()
	exp.Value = p.parseExpression(LOWEST)

	return exp
}

func (p *Parser) validateDeclarationTarget(exp ast.Expression) bool {
	switch exp.(type) {
	case *ast.Identifier, *ast.PropertyExpression, *ast.IndexExpression:
		return true
	}
	return false
}
