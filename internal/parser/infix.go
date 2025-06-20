package parser

import (
	"sunbird/internal/ast"
	"sunbird/internal/token"
)

var precedences = map[token.TokenType]int{
	token.Or:       LOGICAL,
	token.And:      LOGICAL,
	token.Eq:       EQUALS,
	token.NotEq:    EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LE:       LESSGREATER,
	token.GE:       LESSGREATER,
	token.Plus:     SUM,
	token.Minus:    SUM,
	token.Slash:    PRODUCT,
	token.Asterisk: PRODUCT,
	token.LParen:   CALL,
	token.Pipe:     PIPE,
	token.LBracket: INDEX,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}
