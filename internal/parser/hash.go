package parser

import (
	"sunbird/internal/ast"
	"sunbird/internal/token"
)

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
