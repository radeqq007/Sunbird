package parser

import (
	"sunbird/ast"
	"sunbird/token"
)

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:     p.curToken,
		Function:  function,
		Arguments: p.parseExpressionList(token.RPAREN),
	}

	return exp
}

