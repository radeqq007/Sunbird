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

	case token.While:
		return p.parseWhileStatement()

	case token.Break:
		return p.parseBreakStatement()

	case token.Continue:
		return p.parseContinueStatement()

	case token.Import:
		return p.parseImportStatement()

	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	stmt := &ast.BreakStatement{Token: p.curToken}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseContinueStatement() *ast.ContinueStatement {
	stmt := &ast.ContinueStatement{Token: p.curToken}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
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

	if p.curTokenIs(token.Let) {
		stmt.Init = p.parseLetExpression()
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

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.curToken}

	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBrace) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseLetExpression() ast.Expression {
	exp := &ast.LetExpression{Token: p.curToken}

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

func (p *Parser) parseConstExpression() ast.Expression {
	exp := &ast.ConstExpression{Token: p.curToken}

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

func (p *Parser) parseImportStatement() *ast.ImportStatement {
	stmt := &ast.ImportStatement{Token: p.curToken}

	if !p.expectPeek(token.String) {
		return nil
	}

	stmt.Path = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.As) {
		p.nextToken()
		if !p.expectPeek(token.Ident) {
			return nil
		}
		stmt.Alias = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken() // skip the ;
	}

	return stmt
}
