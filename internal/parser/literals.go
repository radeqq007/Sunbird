package parser

import (
	"fmt"
	"strconv"
	"sunbird/internal/ast"
	"sunbird/internal/token"
)

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.True)}
}

func (p *Parser) parseNullLiteral() ast.Expression {
	return &ast.NullLiteral{Token: p.curToken}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}

	array.Elements = p.parseExpressionList(token.RBracket)

	return array
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := ast.HashLiteral{Token: p.curToken}

	hash.Pairs = make([]ast.HashPair, 0)

	// Check for an empty hash
	if p.peekTokenIs(token.RBrace) {
		p.nextToken()
		return &hash
	}

	for !p.peekTokenIs(token.RBrace) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.Colon) {
			return nil
		}

		p.nextToken()

		value := p.parseExpression(LOWEST)

		hash.Pairs = append(hash.Pairs, ast.HashPair{Key: key, Value: value})

		// If the next token isn't the end, we MUST see a comma
		if !p.peekTokenIs(token.RBrace) && !p.expectPeek(token.Comma) {
			return nil
		}
	}

	if !p.expectPeek(token.RBrace) {
		return nil
	}

	return &hash
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LParen) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBrace) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RParen) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()

		ident = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RParen) {
		return nil
	}

	return identifiers
}
