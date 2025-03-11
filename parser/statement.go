package parser

import (
	"sunbird/ast"
	"sunbird/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()

	case token.RETURN:
		return p.parseReturnStatement()

	case token.IDENT:
		if p.peekTokenIs(token.ASSIGN) {
				return p.parseAssignStatement()
		}
		return p.parseExpressionStatement()
	
	default:
		return p.parseExpressionStatement()
	}
}