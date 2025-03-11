package parser

import (
	"sunbird/ast"
	"sunbird/token"
)

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
	stmt := &ast.AssignStatement{ Token: p.curToken }

	if !p.curTokenIs(token.IDENT) {
		return nil
	}

	stmt.Name =  &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
