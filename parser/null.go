package parser

import "sunbird/ast"

func (p *Parser) parseNullLiteral() ast.Expression {
	return &ast.NullLiteral{ Token: p.curToken }
}