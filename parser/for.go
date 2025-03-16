package parser

import (
	"sunbird/ast"
	"sunbird/token"
)

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.curToken}

	p.nextToken()

	if p.curTokenIs(token.VAR) {
		stmt.Init = p.parseVarStatement()
	} else if p.curTokenIs(token.IDENT) {
		stmt.Init = p.parseAssignStatement()
	}

	if !p.curTokenIs(token.SEMICOLON) {
		return nil
	}

	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	p.nextToken()

	if p.curTokenIs(token.IDENT) {
		stmt.Update = p.parseAssignStatement()
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}
