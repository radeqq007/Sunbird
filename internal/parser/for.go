package parser

import (
	"sunbird/internal/ast"
	"sunbird/internal/token"
)

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.curToken}

	p.nextToken()

	if p.curTokenIs(token.Var) {
		stmt.Init = p.parseVarStatement()
	} else if p.curTokenIs(token.Ident) {
		stmt.Init = p.parseAssignStatement()
	}

	if !p.curTokenIs(token.Semicolon) {
		return nil
	}

	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.Semicolon) {
		return nil
	}

	p.nextToken()

	if p.curTokenIs(token.Ident) {
		stmt.Update = p.parseAssignStatement()
	}

	if !p.expectPeek(token.LBrace) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}
