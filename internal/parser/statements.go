package parser

import (
	"sunbird/internal/ast"
	"sunbird/internal/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.Return:
		return p.parseReturnStatement()

	case token.Ident:
		return p.parseExpressionStatement()

	case token.For:
		return p.parseForStatement()

	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	// Skip semicolon if present
	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.curToken}

	p.nextToken()

	if p.curTokenIs(token.Var) {
		stmt.Init = p.parseVarExpression()
	}

	if !p.expectPeek(token.Semicolon) {
		return nil
	}

	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.Semicolon) {
		return nil
	}

	p.nextToken()

	if p.curTokenIs(token.Ident) {
		ident := p.parseIdentifier()
		if !p.expectPeek(token.Assign) {
			return nil
		}
		stmt.Update = p.parseAssignExpression(ident)
	}

	if !p.expectPeek(token.LBrace) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

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
