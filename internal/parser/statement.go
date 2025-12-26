package parser

import (
	"sunbird/internal/ast"
	"sunbird/internal/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.Var:
		return p.parseVarStatement()

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
