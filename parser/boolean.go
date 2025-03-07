package parser

import (
	"sunbird/ast"
	"sunbird/token"
)

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}