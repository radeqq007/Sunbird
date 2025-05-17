package parser

import (
	"sunbird/internal/ast"
	"sunbird/internal/token"
)

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	// Skip semicolon if present
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
