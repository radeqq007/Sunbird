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

	case token.Try:
		return p.parseTryCatchStatement()

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

	if !p.expectPeek(token.Ident) {
		return nil
	}

	stmt.Variable = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.In) {
		return nil
	}

	p.nextToken()

	stmt.Iterable = p.parseExpression(LOWEST)

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

	if p.peekTokenIs(token.Colon) {
		p.nextToken()
		p.nextToken()
		exp.Type = p.parseTypeAnnotation()
	}

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

	if p.peekTokenIs(token.Colon) {
		p.nextToken()
		p.nextToken()
		exp.Type = p.parseTypeAnnotation()
	}

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

func (p *Parser) parseTryCatchStatement() *ast.TryCatchStatement {
	stmt := &ast.TryCatchStatement{Token: p.curToken}

	if !p.expectPeek(token.LBrace) {
		return nil
	}

	stmt.Try = p.parseBlockStatement()

	if !p.expectPeek(token.Catch) {
		return nil
	}

	if !p.expectPeek(token.Ident) {
		return nil
	}

	stmt.Param = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.LBrace) {
		return nil
	}

	stmt.Catch = p.parseBlockStatement()

	if p.peekTokenIs(token.Finally) {
		p.nextToken()

		if !p.expectPeek(token.LBrace) {
			return nil
		}

		stmt.Finally = p.parseBlockStatement()
	}

	return stmt
}
