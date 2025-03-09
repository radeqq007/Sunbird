package parser

import (
	"sunbird/ast"
	"sunbird/token"
)

var precedences = map[token.TokenType]int {
  token.EQ: EQUALS,
  token.NOT_EQ: EQUALS,
  token.LT: LESSGREATER,
  token.GT: LESSGREATER,
  token.LE: LESSGREATER,
  token.GE: LESSGREATER,
  token.PLUS: SUM,
  token.MINUS: SUM,
  token.SLASH: PRODUCT,
  token.ASTERISK: PRODUCT,
  token.LPAREN: CALL,
  token.LBRACKET: INDEX,
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