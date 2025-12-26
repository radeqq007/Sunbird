package parser

import (
	"sunbird/internal/ast"
	"sunbird/internal/token"
)

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := ast.HashLiteral{Token: p.curToken}

	hash.Pairs = make([]ast.HashPair, 0)

	for !p.curTokenIs(token.RBrace) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.Colon) {
			return nil
		}

		p.nextToken()

		value := p.parseExpression(LOWEST)

		hash.Pairs = append(hash.Pairs, ast.HashPair{Key: key, Value: value})

		if !p.peekTokenIs(token.RBrace) && !p.expectPeek(token.Comma) {
			return nil
		}
	}

	if !p.expectPeek(token.RBrace) {
		return nil
	}

	return &hash
}
